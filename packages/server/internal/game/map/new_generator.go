package gamemap

import (
	"errors"
	"fmt"
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

// Perlin 将给定范围自动切分为多个网格，对每个网格调用 PerlinSingleGrid，并汇总所有噪声结果。
// 参数说明：
//
//	startX, startY: 总区域起始点的整数坐标（左下角）
//	endX, endY: 总区域结束点的整数坐标（右上角）
//	gridSize: 每个网格的边长（整数，单位与坐标一致）
//	step: 采样步长（float64）
//
// 返回值：
//
//	指向二维 float64 切片的指针，包含整个区域的 Perlin 噪声值
func Perlin(startX, startY, endX, endY, gridSize int, step float64) *[][]float64 {
	// 校验范围能否被 gridSize 整除
	if (endX-startX)%gridSize != 0 || (endY-startY)%gridSize != 0 {
		panic("给定范围不能被 gridSize 整除")
	}

	xCount := int((float64(endX) - float64(startX)) / step)
	yCount := int((float64(endY) - float64(startY)) / step)
	results := make([][]float64, xCount)
	for i := range results {
		results[i] = make([]float64, yCount)
	}
	for gx := startX; gx < endX; gx += gridSize {
		for gy := startY; gy < endY; gy += gridSize {
			gridEndX := gx + gridSize
			gridEndY := gy + gridSize
			gridNoise := PerlinSingleGrid(gx, gy, gridEndX, gridEndY, step)
			// 将当前网格的噪声结果汇总到总结果中
			for xi, x := 0, float64(gx); x <= float64(gridEndX); x, xi = x+step, xi+1 {
				globalXi := int((x - float64(startX)) / step)
				if globalXi >= xCount {
					break
				}
				for yi, y := 0, float64(gy); y <= float64(gridEndY); y, yi = y+step, yi+1 {
					globalYi := int((y - float64(startY)) / step)
					if globalYi >= yCount {
						break
					}
					(*gridNoise)[xi][yi] = clamp((*gridNoise)[xi][yi], -1, 1) // 可选：归一化
					results[globalXi][globalYi] = (*gridNoise)[xi][yi]
				}
			}
		}
	}
	return &results
}

// clamp 用于将值限制在指定区间
func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// PerlinSingleGrid 生成一个二维网格的 Perlin 噪声值切片。
// 参数说明：
//
//	startX, startY: 网格起始点的整数坐标（左下角）
//	endX, endY: 网格结束点的整数坐标（右上角）
//	step: 网格采样的步长（float64），决定采样精度
//
// 返回值：
//
//	指向二维 float64 切片的指针，每个元素为对应采样点的 Perlin 噪声值
func PerlinSingleGrid(startX, startY, endX, endY int, step float64) *[][]float64 {
	// 校验范围能否被 step 整除
	xRange := float64(endX - startX)
	yRange := float64(endY - startY)
	if xRange/step != float64(int(xRange/step)) || yRange/step != float64(int(yRange/step)) {
		panic("给定范围不能被 step 整除")
	}

	// 初始化结果切片
	xCount := int((float64(endX)-float64(startX))/step) + 1
	yCount := int((float64(endY)-float64(startY))/step) + 1
	results := make([][]float64, xCount)
	for i := range results {
		results[i] = make([]float64, yCount)
	}
	// 左下格点
	p00 := GridPoint{Point2: Point2{X: float64(startX), Y: float64(startY)}}
	p00.Gradient = getGradient(startX, startY)
	// 右下格点
	p10 := GridPoint{Point2: Point2{X: float64(endX), Y: float64(startY)}}
	p10.Gradient = getGradient(endX, startY)
	// 左上格点
	p01 := GridPoint{Point2: Point2{X: float64(startX), Y: float64(endY)}}
	p01.Gradient = getGradient(startX, endY)
	// 右上格点
	p11 := GridPoint{Point2: Point2{X: float64(endX), Y: float64(endY)}}
	p11.Gradient = getGradient(endX, endY)

	for xi, x := 0, float64(startX); x <= float64(endX); x, xi = x+step, xi+1 {
		for yi, y := 0, float64(startY); y <= float64(endY); y, yi = y+step, yi+1 {
			currentPoint := Point2{X: x, Y: y}
			dx := x - float64(startX)
			dy := y - float64(startY)
			// 计算四个角点的点乘结果
			dot00 := p00.Gradient.Dot(p00.VectorTo(currentPoint))
			dot10 := p10.Gradient.Dot(p10.VectorTo(currentPoint))
			dot01 := p01.Gradient.Dot(p01.VectorTo(currentPoint))
			dot11 := p11.Gradient.Dot(p11.VectorTo(currentPoint))
			// 计算插值
			ro := lerp(dot00, dot10, fade(dx))
			r1 := lerp(dot01, dot11, fade(dx))
			// 最终结果
			results[xi][yi] = lerp(ro, r1, fade(dy))
		}
	}
	return &results
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
var permutation = [256]uint8{
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
	return gradients[((permutation[x&255]+uint8(y))&255)%16]
}
