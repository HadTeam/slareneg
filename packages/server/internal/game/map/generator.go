package gamemap

import (
	"errors"
	"log/slog"
	"server/internal/game/block"
)

type MapGenerator interface {
	Name() string
	Generate(size Size, info Info, players []Player) (Map, error)
}

type Player struct {
	Index    int
	Owner    block.Owner
	IsActive bool
}

type GeneratorConfig struct {
	MountainDensity   float64 // 0.0-1.0, 默认 0.5
	CastleDensity     float64 // 0.0-1.0, 默认 0.5
	MinCastleDistance int     // 城堡间最小距离
}

func DefaultGeneratorConfig() GeneratorConfig {
	return GeneratorConfig{
		MountainDensity:   0.5,
		CastleDensity:     0.5,
		MinCastleDistance: 5,
	}
}

type GeneratorFunc func(size Size, info Info, players []Player, config ...GeneratorConfig) (Map, error)

type generatorEntry struct {
	name      string
	generator GeneratorFunc
}

var generators = make(map[string]generatorEntry)
var generatorNames []string

func RegisterGenerator(name string, generator GeneratorFunc) {
	entry := generatorEntry{
		name:      name,
		generator: generator,
	}

	generators[name] = entry
	generatorNames = append(generatorNames, name)

	slog.Debug("Registered map generator", "name", name)
}

func GetAllGeneratorNames() []string {
	result := make([]string, len(generatorNames))
	copy(result, generatorNames)
	return result
}

func GeneratorExists(name string) bool {
	_, exists := generators[name]
	return exists
}

func GenerateMap(generatorName string, size Size, info Info, players []Player, config ...GeneratorConfig) (Map, error) {
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

	return entry.generator(size, info, players, config...)
}
