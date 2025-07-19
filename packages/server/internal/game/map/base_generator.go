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

			kingBlock := block.NewBlock(block.KingName, 1, player.Owner)
			if err := gameMap.SetBlock(pos, kingBlock); err != nil {
				return err
			}

			directions := []Pos{
				{X: pos.X + 1, Y: pos.Y},
				{X: pos.X - 1, Y: pos.Y},
				{X: pos.X, Y: pos.Y + 1},
				{X: pos.X, Y: pos.Y - 1},
			}

			for _, dir := range directions {
				if size.IsPosValid(dir) {
					soldierBlock := block.NewBlock(block.SoldierName, 1, player.Owner)
					gameMap.SetBlock(dir, soldierBlock)
				}
			}
		}
	}

	if err := g.generateTerrain(gameMap); err != nil {
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

func (g *BaseMapGenerator) generateTerrain(gameMap *BaseMap) error {
	size := gameMap.Size()

	noiseMap := g.generatePerlinNoise(int(size.Width), int(size.Height))

	mountainThreshold := 0.8 - (g.config.MountainDensity * 0.6)
	castleThreshold := 0.6 - (g.config.CastleDensity * 0.4)

	maxAttempts := 100
	baseCastleCount := int(size.Width * size.Height / 50)
	targetCastleCount := int(float64(baseCastleCount) * (0.5 + g.config.CastleDensity))

	var castlePositions []Pos

	for attempt := 0; attempt < maxAttempts; attempt++ {
		castlePositions = nil

		for y := uint16(1); y <= size.Height; y++ {
			for x := uint16(1); x <= size.Width; x++ {
				pos := Pos{X: x, Y: y}

				existing, _ := gameMap.Block(pos)
				if existing != nil {
					continue
				}

				noise := noiseMap[y-1][x-1]

				if noise > mountainThreshold {
					mountain := block.NewBlock(block.MountainName, 0, 0)
					gameMap.SetBlock(pos, mountain)
				} else if noise > castleThreshold && noise <= mountainThreshold {
					if len(castlePositions) < targetCastleCount &&
						g.canPlaceCastle(pos, castlePositions, g.config.MinCastleDistance) {

						castleNum := g.rng.Intn(20) + 10
						castle := block.NewBlock(block.CastleName, block.Num(castleNum), 0)
						gameMap.SetBlock(pos, castle)
						castlePositions = append(castlePositions, pos)
					}
				}
			}
		}

		if g.validateMap(gameMap, castlePositions) {
			break
		}

		g.clearTerrain(gameMap)
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

			noise[y][x] = noiseValue / maxValue
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
	w := 32
	s := w / 2
	a := uint32(ix)
	b := uint32(iy)
	a *= 3284157443
	b ^= a<<uint32(s) | a>>uint32(w-s)
	b *= 1911520717
	a ^= b<<uint32(s) | b>>uint32(w-s)
	a *= 2048419325
	random := float64(a) * (3.14159265 / float64(^uint32(0)>>1))
	return [2]float64{math.Cos(random), math.Sin(random)}
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
	if len(castlePositions) < 2 {
		return false
	}

	size := gameMap.Size()
	totalCells := int(size.Width * size.Height)

	reachableCells := g.bfsReachability(gameMap, castlePositions[0])
	unreachableRatio := float64(totalCells-reachableCells) / float64(totalCells)

	if unreachableRatio > 0.1 {
		return false
	}

	for i := 1; i < len(castlePositions); i++ {
		if !g.isReachable(gameMap, castlePositions[0], castlePositions[i]) {
			return false
		}
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
