package gamemap

import (
	"errors"
	"fmt"
	"log/slog"
	"server/internal/game/block"
)

// MapGenerator 接口定义了地图生成器的基本行为。
// Name 方法返回生成器的名称，Generate 方法用于生成地图。
type MapGenerator interface {
	Name() string
	Generate(size Size, info Info, players []Player) (Map, error)
}

// Player 结构体表示一个玩家的信息。
// 包含玩家索引、所属者和是否激活状态。
type Player struct {
	Index    int         // 玩家索引
	Owner    block.Owner // 玩家所属者
	IsActive bool        // 玩家是否处于激活状态
}

// GeneratorConfig 结构体用于配置地图生成器参数。
// 包含山脉密度、城堡密度、城堡最小距离和随机种子。
type GeneratorConfig struct {
	MountainDensity   float64 // 0.0-1.0, 默认 0.5，山脉密度
	CastleDensity     float64 // 0.0-1.0, 默认 0.5，城堡密度
	MinCastleDistance int     // 城堡间最小距离
	Seed              int64   // 随机种子，0 表示使用随机种子
}

// DefaultGeneratorConfig 返回默认的地图生成器配置。
// 包含默认的山脉密度、城堡密度、最小城堡距离和随机种子。
func DefaultGeneratorConfig() GeneratorConfig {
	return GeneratorConfig{
		MountainDensity:   0.7,  // Increased from 0.5 to generate more mountains
		CastleDensity:     0.5,
		MinCastleDistance: 5,
		Seed:              0,
	}
}

// String 返回 GeneratorConfig 的字符串表示。
// 格式为 m{山脉密度}-c{城堡密度}-d{最小距离}-s{种子}。
func (c GeneratorConfig) String() string {
	return fmt.Sprintf("m%.1f-c%.1f-d%d-s%d",
		c.MountainDensity, c.CastleDensity, c.MinCastleDistance, c.Seed)
}

// GeneratorFunc 类型定义了地图生成器的函数签名。
// 接受地图大小、玩家列表和可选的生成器配置，返回生成的地图或错误信息。
type GeneratorFunc func(size Size, players []Player, config ...GeneratorConfig) (Map, error)

type generatorEntry struct {
	name      string
	generator GeneratorFunc
}

var generators = make(map[string]generatorEntry)
var generatorNames []string

// RegisterGenerator 函数用于注册新的地图生成器。
// 将生成器名称和对应的生成器函数添加到生成器映射中。
func RegisterGenerator(name string, generator GeneratorFunc) {
	entry := generatorEntry{
		name:      name,
		generator: generator,
	}

	generators[name] = entry
	generatorNames = append(generatorNames, name)

	slog.Debug("Registered map generator", "name", name)
}

// GetAllGeneratorNames 函数返回所有已注册的地图生成器名称的副本。
// 用于在游戏中显示可用的地图生成器选项。
func GetAllGeneratorNames() []string {
	result := make([]string, len(generatorNames))
	copy(result, generatorNames)
	return result
}

// GeneratorExists 函数检查指定名称的地图生成器是否已注册。
// 返回布尔值表示生成器是否存在。
func GeneratorExists(name string) bool {
	_, exists := generators[name]
	return exists
}

// GenerateMap 函数用于生成地图。
// 根据指定的生成器名称、地图大小和玩家列表，调用相应的生成器函数生成地图。
// 如果指定的生成器名称未知，则返回错误信息。
func GenerateMap(generatorName string, size Size, players []Player, config ...GeneratorConfig) (Map, error) {
	entry, exists := generators[generatorName]
	if !exists {
		slog.Warn("Unknown map generator", "name", generatorName, "available", GetAllGeneratorNames())
		if len(generators) > 0 {
			for _, g := range generators {
				entry = g
				break
			}
		} else {
			return nil, errors.New("no map generators available")
		}
	}

	return entry.generator(size, players, config...)
}
