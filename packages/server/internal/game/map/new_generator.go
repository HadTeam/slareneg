package gamemap

import (
	"errors"
	"fmt"
	"math"
	"server/internal/game/block"
)

// NewMapGenerator represents a new map generator with configuration parameters.
type NewMapGenerator struct {
	config               GeneratorConfig // Configuration parameters
	precomputedGradients [][]Vec2        // Precomputed gradient vectors
	gradientMapSize      int             // Size of the gradient vector table
}

func init() {
	RegisterGenerator("new", func(size Size, players []Player, config ...GeneratorConfig) (Map, error) {
		var cfg GeneratorConfig
		if len(config) > 0 {
			cfg = config[0]
		} else {
			cfg = DefaultGeneratorConfig()
		}

		generator := NewNewMapGenerator(cfg)
		return generator.Generate(size, players)
	})
}

// NewNewMapGenerator creates a new map generator with the given configuration.
func NewNewMapGenerator(config GeneratorConfig) *NewMapGenerator {
	generator := &NewMapGenerator{
		config: config,
	}
	// Initialize precomputed gradient vector table during initialization
	generator.initPrecomputedGradients(100)
	return generator
}

// initPrecomputedGradients precomputes the gradient vector table.
func (g *NewMapGenerator) initPrecomputedGradients(size int) {
	g.gradientMapSize = size
	g.precomputedGradients = make([][]Vec2, size)
	for i := range g.precomputedGradients {
		g.precomputedGradients[i] = make([]Vec2, size)
		for j := range g.precomputedGradients[i] {
			g.precomputedGradients[i][j] = getGradient(i, j)
		}
	}
}

// getGradient gets the precomputed gradient vector, or calculates it in real-time if out of range.
func (g *NewMapGenerator) getGradient(x, y int) Vec2 {
	// Map coordinates to the range of the precomputed table
	mapX := ((x % g.gradientMapSize) + g.gradientMapSize) % g.gradientMapSize
	mapY := ((y % g.gradientMapSize) + g.gradientMapSize) % g.gradientMapSize
	return g.precomputedGradients[mapX][mapY]
}

// Name returns the name of the generator.
func (g *NewMapGenerator) Name() string {
	return "new"
}

// Generate creates a new map with the specified size and players.
func (g *NewMapGenerator) Generate(size Size, players []Player) (Map, error) {
	playerCount := len(players)

	// Validate map size
	switch {
	case size.Width < 10 || size.Height < 10:
		return nil, errors.New("invalid map size: must be at least 10x10")
	case size.Width*size.Height < uint16(playerCount)*10:
		return nil, errors.New("invalid map size: too small for number of players")
	case size.Width > 100 || size.Height > 100:
		return nil, errors.New("invalid map size: must not exceed 100x100")
	}

	mapId := fmt.Sprintf("generated-%d", g.config.Seed)
	info := Info{
		Id:   mapId,
		Name: "Generated Map",
		Desc: "Procedurally generated map",
	}

	// Create a new empty base map
	result := NewEmptyBaseMap(size, info)
	if result == nil {
		return nil, errors.New("failed to create empty base map")
	}

	// Generate Perlin noise terrain
	noiseMap := generatePerlinNoise(int(size.Width), int(size.Height))

	// Generate player starting positions
	playerPositions := g.generatePlayerStartPositions(size, players)

	// Determine block type for each position
	for y := uint16(0); y < size.Height; y++ {
		for x := uint16(0); x < size.Width; x++ {
			pos := Pos{X: x, Y: y}

			// Check if this is a player starting position
			isPlayerStart := false
			for i, player := range players {
				if player.IsActive && i < len(playerPositions) && playerPositions[i] == pos {
					// Place player's king
					kingBlock := block.NewBlock(block.KingName, 1, player.Owner)
					if err := result.SetBlock(pos, kingBlock); err != nil {
						return nil, err
					}
					isPlayerStart = true
					break
				}
			}

			if !isPlayerStart {
				// Determine terrain type based on noise value
				noiseValue := noiseMap[y][x]
				blockToPlace := g.getBlockTypeFromNoise(noiseValue, pos)

				if err := result.SetBlock(pos, blockToPlace); err != nil {
					return nil, err
				}
			}
		}
	}

	return result, nil
}

// getBlockTypeFromNoise determines the block type based on noise value and configuration thresholds.
func (g *NewMapGenerator) getBlockTypeFromNoise(noiseValue float64, _ Pos) block.Block {
	// Determine whether to generate mountains based on configured mountain density threshold
	mountainThreshold := 1.0 - g.config.MountainDensity

	// Determine whether to generate castles based on configured castle density threshold
	castleThreshold := 1.0 - g.config.CastleDensity*0.1 // Castles should be relatively rare

	switch {
	case noiseValue > mountainThreshold:
		// Generate mountains in high noise value areas
		return block.NewBlock(block.MountainName, 0, 0)
	case noiseValue > castleThreshold:
		// Generate castles in medium-high noise value areas
		return block.NewBlock(block.CastleName, 0, 0)
	default:
		// Generate blank terrain in other areas
		return block.NewBlock(block.BlankName, 0, 0)
	}
}

// generatePlayerStartPositions generates starting positions for players.
func (g *NewMapGenerator) generatePlayerStartPositions(size Size, players []Player) []Pos {
	positions := make([]Pos, len(players))

	// Simple uniform distribution strategy: distribute players on map edges
	switch len(players) {
	case 1:
		positions[0] = Pos{X: size.Width / 2, Y: size.Height / 2}
	case 2:
		positions[0] = Pos{X: 1, Y: size.Height / 2}
		positions[1] = Pos{X: size.Width - 2, Y: size.Height / 2}
	case 3:
		positions[0] = Pos{X: 1, Y: 1}
		positions[1] = Pos{X: size.Width - 2, Y: 1}
		positions[2] = Pos{X: size.Width / 2, Y: size.Height - 2}
	case 4:
		positions[0] = Pos{X: 1, Y: 1}
		positions[1] = Pos{X: size.Width - 2, Y: 1}
		positions[2] = Pos{X: 1, Y: size.Height - 2}
		positions[3] = Pos{X: size.Width - 2, Y: size.Height - 2}
	default:
		// For more players, use circular distribution
		centerX := float64(size.Width) / 2
		centerY := float64(size.Height) / 2
		radius := math.Min(centerX, centerY) * 0.8

		for i := range positions {
			angle := 2 * math.Pi * float64(i) / float64(len(players))
			x := centerX + radius*math.Cos(angle)
			y := centerY + radius*math.Sin(angle)

			// Ensure position is within map bounds
			x = math.Max(1, math.Min(float64(size.Width-2), x))
			y = math.Max(1, math.Min(float64(size.Height-2), y))

			positions[i] = Pos{X: uint16(x), Y: uint16(y)}
		}
	}

	return positions
}

// generatePerlinNoise generates a Perlin noise map of the specified size.
func generatePerlinNoise(width, height int) [][]float64 {
	noiseMap := make([][]float64, height)
	for i := range noiseMap {
		noiseMap[i] = make([]float64, width)
	}

	// Generate noise values for each position
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Use a smaller scale factor for better terrain variation
			scale := 0.1
			point := Point2{X: float64(x) * scale, Y: float64(y) * scale}

			// Generate layered noise
			noiseValue := generateLayeredPerlinNoise(point, 4, 0.5, 2.0)
			// Normalize to [0, 1] range
			noiseMap[y][x] = (noiseValue + 1.0) * 0.5
		}
	}

	return noiseMap
}

// generateLayeredPerlinNoise generates multi-layered overlapping noise.
func generateLayeredPerlinNoise(point Point2, octaves int, persistence, lacunarity float64) float64 {
	var total float64
	var amplitude = 1.0
	var frequency = 1.0
	var maxValue = 0.0 // Used for normalization

	for i := 0; i < octaves; i++ {
		n := PerlinSinglePoint(Point2{X: point.X * frequency, Y: point.Y * frequency})
		total += n * amplitude
		maxValue += amplitude

		amplitude *= persistence
		frequency *= lacunarity
	}

	// Normalize result to stabilize range to [-1, 1]
	return total / maxValue
}

// PerlinSinglePoint computes the Perlin noise value for a single point,
// internally auto-calculating the four grid points.
func PerlinSinglePoint(currentPoint Point2) float64 {
	// Calculate grid point coordinates
	floorX := int(math.Floor(currentPoint.X))
	floorY := int(math.Floor(currentPoint.Y))

	// Construct four grid points
	p00 := GridPoint{Point2: Point2{X: float64(floorX), Y: float64(floorY)}}
	p00.Gradient = getGradient(floorX, floorY)

	p10 := GridPoint{Point2: Point2{X: float64(floorX + 1), Y: float64(floorY)}}
	p10.Gradient = getGradient(floorX+1, floorY)

	p01 := GridPoint{Point2: Point2{X: float64(floorX), Y: float64(floorY + 1)}}
	p01.Gradient = getGradient(floorX, floorY+1)

	p11 := GridPoint{Point2: Point2{X: float64(floorX + 1), Y: float64(floorY + 1)}}
	p11.Gradient = getGradient(floorX+1, floorY+1)

	// Calculate dot products
	dot00 := p00.Gradient.Dot(p00.VectorTo(currentPoint))
	dot10 := p10.Gradient.Dot(p10.VectorTo(currentPoint))
	dot01 := p01.Gradient.Dot(p01.VectorTo(currentPoint))
	dot11 := p11.Gradient.Dot(p11.VectorTo(currentPoint))

	// Bilinear interpolation
	dx := currentPoint.X - p00.X
	ro := interpolation(dot00, dot10, fade(dx))
	r1 := interpolation(dot01, dot11, fade(dx))
	dy := currentPoint.Y - p00.Y
	return interpolation(ro, r1, fade(dy))
}

// fade function for smooth interpolation.
func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

// Vec2 represents a two-dimensional vector.
type Vec2 struct {
	X float64
	Y float64
}

// Dot returns the dot product of two two-dimensional vectors.
func (v Vec2) Dot(other Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Point2 represents a two-dimensional point.
type Point2 struct {
	X float64
	Y float64
}

// VectorTo returns the two-dimensional vector from the current point to the target point.
func (p Point2) VectorTo(target Point2) Vec2 {
	return Vec2{
		X: target.X - p.X,
		Y: target.Y - p.Y,
	}
}

// GridPoint inherits from Point2, represents a grid point and contains a gradient vector.
type GridPoint struct {
	Point2
	Gradient Vec2
}

// Precomputed permutation table for hashing grid point coordinates
var permutation = [256]int{
	151, 160, 137, 91, 90, 15, 131, 13, 201, 95, 96, 53, 194, 233, 7, 225,
	140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23, 190, 6, 148,
	247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32,
	57, 177, 33, 88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175,
	74, 165, 71, 134, 139, 48, 27, 166, 77, 146, 158, 231, 83, 111, 229, 122,
	60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244, 102, 143, 54,
	65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169,
	200, 196, 135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64,
	52, 217, 226, 250, 124, 123, 5, 202, 38, 147, 118, 126, 255, 82, 85, 212,
	207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42, 223, 183, 170, 213,
	119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104,
	218, 246, 97, 228, 251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241,
	81, 51, 145, 235, 249, 14, 239, 107, 49, 192, 214, 31, 181, 199, 106, 157,
	184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254, 138, 236, 205, 93,
	222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
}

// 16 uniformly distributed unit vectors as gradient directions
var gradients = []Vec2{
	{X: 1, Y: 0},
	{X: 0.9239, Y: 0.3827},
	{X: 0.7071, Y: 0.7071},
	{X: 0.3827, Y: 0.9239},
	{X: 0, Y: 1},
	{X: -0.3827, Y: 0.9239},
	{X: -0.7071, Y: 0.7071},
	{X: -0.9239, Y: 0.3827},
	{X: -1, Y: 0},
	{X: -0.9239, Y: -0.3827},
	{X: -0.7071, Y: -0.7071},
	{X: -0.3827, Y: -0.9239},
	{X: 0, Y: -1},
	{X: 0.3827, Y: -0.9239},
	{X: 0.7071, Y: -0.7071},
	{X: 0.9239, Y: -0.3827},
}

// interpolation function: linear interpolation
func interpolation(a, b, t float64) float64 {
	return a*(1-t) + b*t
}

// getGradient generates gradient vectors for grid points
func getGradient(x, y int) Vec2 {
	return gradients[permutation[(permutation[x&255]+y)&255]%16]
}
