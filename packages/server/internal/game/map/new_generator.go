package gamemap

import (
	"errors"
	"fmt"
	"math"
)

const MinMapSizeForSinglePlayer = 10

type NewMapGenerator struct {
	config GeneratorConfig // 配置参数
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

// NewNewMapGenerator 参数 	config 用于初始化生成器配置。
func NewNewMapGenerator(config GeneratorConfig) *NewMapGenerator {
	return &NewMapGenerator{
		config: config,
	}
}

// Name returns the name of the generator.
func (g *NewMapGenerator) Name() string {
	return "new"
}

func (g *NewMapGenerator) Generate(size Size, players []Player) (Map, error) {
	playerCount := len(players)

	// 验证地图大小
	switch {
	case size.Width < 10 || size.Height < 10:
		return nil, errors.New("invalid map size: must be at least 10x10")
	case size.Width*size.Height < uint16(playerCount)*MinMapSizeForSinglePlayer:
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

	// 创建一个新空基础地图
	result := NewEmptyBaseMap(size, info)
	if result == nil {
		return nil, errors.New("failed to create empty base map")
	}
	// noise := Perlin(0, 0, int(size.Width), int(size.Height), 10, 0.1)
	return result, nil
}

func Perlin(startX, startY, endX, endY int, step float64) [][]float64 {
	// 初始化结果切片
	xCount := int((float64(endX)-float64(startX))/step) + 1
	yCount := int((float64(endY)-float64(startY))/step) + 1
	results := make([][]float64, xCount)
	for i := range results {
		results[i] = make([]float64, yCount)
	}

	for i := 0; i < xCount; i++ {
		for j := 0; j < yCount; j++ {
			x := float64(startX) + float64(i)*step
			y := float64(startY) + float64(j)*step

			// 边界检查，防止超出区间
			if x > float64(endX) || y > float64(endY) {
				continue
			}

			// 左下格点
			floorX := math.Floor(x)
			floorY := math.Floor(y)
			p00 := GridPoint{Point2: Point2{X: floorX, Y: floorY}}
			p00.Gradient = getGradient(int(p00.X), int(p00.Y))
			// 右下格点
			p10 := GridPoint{Point2: Point2{X: floorX + 1, Y: floorY}}
			p10.Gradient = getGradient(int(p10.X), int(p10.Y))
			// 左上格点
			p01 := GridPoint{Point2: Point2{X: floorX, Y: floorY + 1}}
			p01.Gradient = getGradient(int(p01.X), int(p01.Y))
			// 右上格点
			p11 := GridPoint{Point2: Point2{X: floorX + 1, Y: floorY + 1}}
			p11.Gradient = getGradient(int(p11.X), int(p11.Y))
			currentPoint := Point2{X: x, Y: y}
			results[i][j] = PerlinSinglePoint(currentPoint, p00, p01, p10, p11)
		}
	}
	return results
}

func PerlinSinglePoint(currentPoint Point2, gp00, gp01, gp10, gp11 GridPoint) float64 {
	dot00 := gp00.Gradient.Dot(gp00.VectorTo(currentPoint))
	dot10 := gp10.Gradient.Dot(gp10.VectorTo(currentPoint))
	dot01 := gp01.Gradient.Dot(gp01.VectorTo(currentPoint))
	dot11 := gp11.Gradient.Dot(gp11.VectorTo(currentPoint))
	dx := currentPoint.X - gp00.X
	ro := lerp(dot00, dot10, fade(dx))
	r1 := lerp(dot01, dot11, fade(dx))
	dy := currentPoint.Y - gp00.Y
	return lerp(ro, r1, fade(dy))
}

// fade 函数用于平滑插值。
func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

// Vec2 表示一个二维向量。
type Vec2 struct {
	X float64
	Y float64
}

// Dot 返回两个二维向量的点乘结果。
func (v Vec2) Dot(other Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Point2 表示一个二维点。
type Point2 struct {
	X float64
	Y float64
}

// VectorTo 返回从当前点到目标点的二维向量。
func (p Point2) VectorTo(target Point2) Vec2 {
	return Vec2{
		X: target.X - p.X,
		Y: target.Y - p.Y,
	}
}

// GridPoint 继承自 Point2，表示格点并包含梯度向量
type GridPoint struct {
	Point2
	Gradient Vec2
}

// 预生成扰动表（Permutation Table），用于哈希格点坐标
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

// 16个均匀分布的单位向量作为梯度方向
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

// 插值函数：线性插值
func lerp(a, b, t float64) float64 {
	return a*(1-t) + b*t
}

// 生成格点梯度向量
func getGradient(x, y int) Vec2 {
	return gradients[permutation[(permutation[x&255]+y)&255]%16]
}
