package gamemap

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"server/internal/game/block"
)

type BaseMapGenerator struct {
	config GeneratorConfig
	rng    *rand.Rand
	gradientCache map[string][2]float64
}

func NewBaseMapGenerator(config GeneratorConfig) *BaseMapGenerator {
	var rng *rand.Rand
	if config.Seed != 0 {
		rng = rand.New(rand.NewSource(config.Seed))
	} else {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	return &BaseMapGenerator{
		config: config,
		rng:    rng,
		gradientCache: make(map[string][2]float64),
	}
}

func (g *BaseMapGenerator) Name() string {
	return "base"
}

func (g *BaseMapGenerator) Generate(size Size, players []Player, config ...GeneratorConfig) (Map, error) {
	if len(config) > 0 {
		g.config = config[0]
		if config[0].Seed != 0 {
			g.rng = rand.New(rand.NewSource(config[0].Seed))
		} else {
			g.rng = rand.New(rand.NewSource(rand.Int63()))
		}
	}

	mapId := fmt.Sprintf("generated-%d", g.config.Seed)
	info := Info{
		Id:   mapId,
		Name: "Generated Map",
		Desc: "Procedurally generated map",
	}

	gameMap := NewEmptyBaseMap(size, info)
	if gameMap == nil {
		return nil, errors.New("failed to create empty map")
	}

	// Clear gradient cache before generating content for new randomness
	g.gradientCache = make(map[string][2]float64)
	
	if err := g.generateContent(gameMap, players); err != nil {
		return nil, err
	}

	return gameMap, nil
}

func (g *BaseMapGenerator) Configure(config GeneratorConfig) {
	g.config = config
}

func (g *BaseMapGenerator) generateContent(gameMap *BaseMap, players []Player) error {
	size := gameMap.Size()

	playerPositions := g.generatePlayerStartPositions(size, players)

	for _, player := range players {
		if player.IsActive && player.Index < len(playerPositions) {
			pos := playerPositions[player.Index]

			// Only place the King, no soldiers - the king will generate units over time
			kingBlock := block.NewBlock(block.KingName, 1, player.Owner)
			if err := gameMap.SetBlock(pos, kingBlock); err != nil {
				return err
			}
		}
	}

	if err := g.generateTerrain(gameMap, playerPositions); err != nil {
		return err
	}

	// Fill all remaining empty spaces with blank blocks
	if err := g.fillEmptySpaces(gameMap); err != nil {
		return err
	}

	return nil
}

func (g *BaseMapGenerator) generatePlayerStartPositions(size Size, players []Player) []Pos {
	activePlayerCount := 0
	for _, p := range players {
		if p.IsActive {
			activePlayerCount++
		}
	}

	positions := make([]Pos, len(players))

	switch activePlayerCount {
	case 2:
		positions[0] = Pos{X: 3, Y: 3}
		positions[1] = Pos{X: size.Width - 2, Y: size.Height - 2}
	case 3:
		positions[0] = Pos{X: 3, Y: 3}
		positions[1] = Pos{X: size.Width - 2, Y: 3}
		positions[2] = Pos{X: size.Width / 2, Y: size.Height - 2}
	case 4:
		positions[0] = Pos{X: 3, Y: 3}
		positions[1] = Pos{X: size.Width - 2, Y: 3}
		positions[2] = Pos{X: 3, Y: size.Height - 2}
		positions[3] = Pos{X: size.Width - 2, Y: size.Height - 2}
	default:
		for i := 0; i < activePlayerCount; i++ {
			angle := float64(i) * 2 * math.Pi / float64(activePlayerCount)
			centerX := float64(size.Width) / 2
			centerY := float64(size.Height) / 2
			radius := float64(min(size.Width, size.Height)) / 3

			x := centerX + radius*math.Cos(angle)
			y := centerY + radius*math.Sin(angle)

			positions[i] = Pos{
				X: uint16(max(1, min(int(size.Width), int(x)))),
				Y: uint16(max(1, min(int(size.Height), int(y)))),
			}
		}
	}

	return positions
}

func (g *BaseMapGenerator) generateTerrain(gameMap *BaseMap, playerPositions []Pos) error {
	size := gameMap.Size()

	noiseMap := g.generatePerlinNoise(int(size.Width), int(size.Height))

	// Adjust thresholds to generate strategic amount of mountains
	// Default density 0.7 should give ~15-20% mountain coverage
	// Higher threshold = fewer mountains
	// After Perlin noise normalization fix, we need higher thresholds
	mountainThreshold := 0.7 - (g.config.MountainDensity * 0.2)  // With density 0.7, threshold = 0.56
	castleThreshold := 0.5 - (g.config.CastleDensity * 0.15)
	

	maxAttempts := 100
	baseCastleCount := int(size.Width * size.Height / 50)
	targetCastleCount := int(float64(baseCastleCount) * (0.5 + g.config.CastleDensity))

	var castlePositions []Pos
	mountainCount := 0
	validMapFound := false

	for attempt := 0; attempt < maxAttempts; attempt++ {
		castlePositions = nil
		mountainCount = 0

		// First pass: place terrain based on noise
		for y := uint16(1); y <= size.Height; y++ {
			for x := uint16(1); x <= size.Width; x++ {
				pos := Pos{X: x, Y: y}

				existing, _ := gameMap.Block(pos)
				if existing != nil {
					continue
				}

				noise := noiseMap[y-1][x-1]

				if noise > mountainThreshold {
					// Don't place mountains too close to kings
					tooCloseToKing := false
					for _, playerPos := range playerPositions {
						dx := int(pos.X) - int(playerPos.X)
						dy := int(pos.Y) - int(playerPos.Y)
						dist := math.Sqrt(float64(dx*dx + dy*dy))
						if dist < 3.0 { // Keep 3 tile radius around kings clear
							tooCloseToKing = true
							break
						}
					}
					
					if !tooCloseToKing {
						mountain := block.NewBlock(block.MountainName, 0, 0)
						gameMap.SetBlock(pos, mountain)
						mountainCount++
					}
				} else if noise > castleThreshold && noise <= mountainThreshold {
					if len(castlePositions) < targetCastleCount &&
						g.canPlaceCastle(pos, castlePositions, g.config.MinCastleDistance) {

						// Neutral castles start with some units (10-30)
						castleNum := g.rng.Intn(20) + 10
						castle := block.NewBlock(block.CastleName, block.Num(castleNum), 0)
						gameMap.SetBlock(pos, castle)
						castlePositions = append(castlePositions, pos)
					}
				}
			}
		}

		// Calculate mountain percentage
		totalTiles := int(size.Width * size.Height)
		mountainPercentage := float64(mountainCount) / float64(totalTiles) * 100.0
		
		if g.validateMap(gameMap, castlePositions) {
			fmt.Printf("Map generated: %d mountains (%.1f%%), %d castles\n", mountainCount, mountainPercentage, len(castlePositions))
			validMapFound = true
			break
		} else {
			fmt.Printf("Validation failed: attempt %d, %d mountains (%.1f%%), %d castles\n", attempt+1, mountainCount, mountainPercentage, len(castlePositions))
		}

		// Only clear terrain if we're not on the last attempt
		if attempt < maxAttempts - 1 {
			g.clearTerrain(gameMap)
		}
	}

	// If no valid map was found after all attempts, generate a simple fallback map
	if !validMapFound {
		fmt.Printf("Warning: Could not generate valid map after %d attempts, using last attempt\n", maxAttempts)
	}

	return nil
}

func (g *BaseMapGenerator) fillEmptySpaces(gameMap *BaseMap) error {
	size := gameMap.Size()
	for y := uint16(1); y <= size.Height; y++ {
		for x := uint16(1); x <= size.Width; x++ {
			pos := Pos{X: x, Y: y}
			existing, _ := gameMap.Block(pos)
			if existing == nil {
				blankBlock := block.NewBlock(block.BlankName, 0, 0)
				if err := gameMap.SetBlock(pos, blankBlock); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (g *BaseMapGenerator) generatePerlinNoise(width, height int) [][]float64 {
	noise := make([][]float64, height)
	for i := range noise {
		noise[i] = make([]float64, width)
	}

	scale := 0.1
	octaves := 4
	persistence := 0.5
	lacunarity := 2.0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			amplitude := 1.0
			frequency := scale
			noiseValue := 0.0
			maxValue := 0.0

			for i := 0; i < octaves; i++ {
				noiseValue += g.perlin(float64(x)*frequency, float64(y)*frequency) * amplitude
				maxValue += amplitude
				amplitude *= persistence
				frequency *= lacunarity
			}

			// Normalize to [0, 1] range (Perlin noise typically produces values in [-1, 1])
			noise[y][x] = (noiseValue/maxValue + 1.0) / 2.0
		}
	}

	return noise
}

func (g *BaseMapGenerator) perlin(x, y float64) float64 {
	x0 := int(math.Floor(x))
	x1 := x0 + 1
	y0 := int(math.Floor(y))
	y1 := y0 + 1

	sx := x - float64(x0)
	sy := y - float64(y0)

	n0 := g.dotGridGradient(x0, y0, x, y)
	n1 := g.dotGridGradient(x1, y0, x, y)
	ix0 := g.interpolate(n0, n1, sx)

	n0 = g.dotGridGradient(x0, y1, x, y)
	n1 = g.dotGridGradient(x1, y1, x, y)
	ix1 := g.interpolate(n0, n1, sx)

	return g.interpolate(ix0, ix1, sy)
}

func (g *BaseMapGenerator) dotGridGradient(ix, iy int, x, y float64) float64 {
	gradient := g.randomGradient(ix, iy)
	dx := x - float64(ix)
	dy := y - float64(iy)
	return dx*gradient[0] + dy*gradient[1]
}

func (g *BaseMapGenerator) randomGradient(ix, iy int) [2]float64 {
	// Create a cache key for this gradient
	key := fmt.Sprintf("%d,%d", ix, iy)
	
	// Check if we already have this gradient cached
	if gradient, exists := g.gradientCache[key]; exists {
		return gradient
	}
	
	// Generate a random angle using the generator's RNG
	angle := g.rng.Float64() * 2 * math.Pi
	gradient := [2]float64{math.Cos(angle), math.Sin(angle)}
	
	// Cache the gradient for consistency
	g.gradientCache[key] = gradient
	
	return gradient
}

func (g *BaseMapGenerator) interpolate(a0, a1, w float64) float64 {
	return (a1-a0)*((w*(w*6.0-15.0)+10.0)*w*w*w) + a0
}

func (g *BaseMapGenerator) canPlaceCastle(pos Pos, existing []Pos, minDistance int) bool {
	for _, castle := range existing {
		dx := int(pos.X) - int(castle.X)
		dy := int(pos.Y) - int(castle.Y)
		distance := int(math.Sqrt(float64(dx*dx + dy*dy)))
		if distance < minDistance {
			return false
		}
	}
	return true
}

func (g *BaseMapGenerator) validateMap(gameMap *BaseMap, castlePositions []Pos) bool {
	// Ensure there's at least a path between kings
	size := gameMap.Size()
	
	// Find king positions
	var kingPositions []Pos
	for y := uint16(1); y <= size.Height; y++ {
		for x := uint16(1); x <= size.Width; x++ {
			pos := Pos{X: x, Y: y}
			b, _ := gameMap.Block(pos)
			if b != nil && b.Meta().Name == block.KingName {
				kingPositions = append(kingPositions, pos)
			}
		}
	}
	
	// Must have exactly as many kings as active players
	if len(kingPositions) < 2 {
		return false
	}
	
	// Check that kings can reach each other (not completely blocked)
	for i := 1; i < len(kingPositions); i++ {
		if !g.isReachable(gameMap, kingPositions[0], kingPositions[i]) {
			return false
		}
	}
	
	// If we have castles, ensure they're reachable too
	if len(castlePositions) > 0 {
		for _, castle := range castlePositions {
			reachable := false
			for _, king := range kingPositions {
				if g.isReachable(gameMap, king, castle) {
					reachable = true
					break
				}
			}
			if !reachable {
				return false
			}
		}
	}
	
	// Check overall map connectivity - ensure most of the map is accessible
	// Use BFS from first king to count reachable tiles
	reachableTiles := g.bfsReachability(gameMap, kingPositions[0])
	totalTiles := int(size.Width * size.Height)
	
	// Count mountains to get non-walkable tiles
	mountainCount := 0
	for y := uint16(1); y <= size.Height; y++ {
		for x := uint16(1); x <= size.Width; x++ {
			pos := Pos{X: x, Y: y}
			b, _ := gameMap.Block(pos)
			if b != nil && b.Meta().Name == block.MountainName {
				mountainCount++
			}
		}
	}
	
	// The reachable area should be at least 70% of walkable tiles
	walkableTiles := totalTiles - mountainCount
	reachablePercentage := float64(reachableTiles) / float64(walkableTiles)
	
	if reachablePercentage < 0.7 {
		fmt.Printf("Map connectivity too low: %.1f%% of walkable tiles are reachable\n", reachablePercentage * 100)
		return false
	}
	
	return true
}

func (g *BaseMapGenerator) bfsReachability(gameMap *BaseMap, start Pos) int {
	size := gameMap.Size()
	visited := make([][]bool, size.Height)
	for i := range visited {
		visited[i] = make([]bool, size.Width)
	}

	queue := []Pos{start}
	visited[start.Y-1][start.X-1] = true
	count := 1

	directions := []Pos{
		{X: 0, Y: 1},
		{X: 0, Y: size.Height - 1},
		{X: 1, Y: 0},
		{X: size.Width - 1, Y: 0},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, dir := range directions {
			var nextX, nextY uint16
			if dir.X == size.Width-1 {
				if current.X == 1 {
					continue
				}
				nextX = current.X - 1
			} else {
				nextX = current.X + dir.X
			}

			if dir.Y == size.Height-1 {
				if current.Y == 1 {
					continue
				}
				nextY = current.Y - 1
			} else {
				nextY = current.Y + dir.Y
			}

			next := Pos{X: nextX, Y: nextY}

			if !size.IsPosValid(next) || visited[next.Y-1][next.X-1] {
				continue
			}

			b, _ := gameMap.Block(next)
			if b != nil && b.Meta().Name == block.MountainName {
				continue
			}

			visited[next.Y-1][next.X-1] = true
			queue = append(queue, next)
			count++
		}
	}

	return count
}

func (g *BaseMapGenerator) isReachable(gameMap *BaseMap, start, end Pos) bool {
	size := gameMap.Size()
	visited := make([][]bool, size.Height)
	for i := range visited {
		visited[i] = make([]bool, size.Width)
	}

	queue := []Pos{start}
	visited[start.Y-1][start.X-1] = true

	directions := []Pos{
		{X: 0, Y: 1},
		{X: 0, Y: size.Height - 1},
		{X: 1, Y: 0},
		{X: size.Width - 1, Y: 0},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.X == end.X && current.Y == end.Y {
			return true
		}

		for _, dir := range directions {
			var nextX, nextY uint16
			if dir.X == size.Width-1 {
				if current.X == 1 {
					continue
				}
				nextX = current.X - 1
			} else {
				nextX = current.X + dir.X
			}

			if dir.Y == size.Height-1 {
				if current.Y == 1 {
					continue
				}
				nextY = current.Y - 1
			} else {
				nextY = current.Y + dir.Y
			}

			next := Pos{X: nextX, Y: nextY}

			if !size.IsPosValid(next) || visited[next.Y-1][next.X-1] {
				continue
			}

			b, _ := gameMap.Block(next)
			if b != nil && b.Meta().Name == block.MountainName {
				continue
			}

			visited[next.Y-1][next.X-1] = true
			queue = append(queue, next)
		}
	}

	return false
}

func (g *BaseMapGenerator) clearTerrain(gameMap *BaseMap) {
	size := gameMap.Size()
	for y := uint16(1); y <= size.Height; y++ {
		for x := uint16(1); x <= size.Width; x++ {
			pos := Pos{X: x, Y: y}
			b, _ := gameMap.Block(pos)
			if b != nil {
				meta := b.Meta()
				if meta.Name == block.MountainName || meta.Name == block.CastleName {
					gameMap.SetBlock(pos, nil)
				}
			}
		}
	}
}

func init() {
	RegisterGenerator("base", func(size Size, players []Player, config ...GeneratorConfig) (Map, error) {
		var cfg GeneratorConfig
		if len(config) > 0 {
			cfg = config[0]
		} else {
			cfg = DefaultGeneratorConfig()
		}

		generator := NewBaseMapGenerator(cfg)
		return generator.Generate(size, players)
	})
}
