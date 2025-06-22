package gamemap

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// MapManager 地图管理器接口
type MapManager interface {
	GetMap(mapId string, players []Player) (Map, error)
	GenerateMapId(generator string, size Size, playerCount int, config GeneratorConfig) string
}

// MapProvider 地图提供者接口
type MapProvider interface {
	CanHandle(prefix string) bool
	GetMap(mapId string, players []Player) (Map, error)
}

// DefaultMapManager 默认地图管理器实现
type DefaultMapManager struct {
	providers []MapProvider
	cache     map[string]Map
}

func NewMapManager() *DefaultMapManager {
	manager := &DefaultMapManager{
		providers: make([]MapProvider, 0),
		cache:     make(map[string]Map),
	}

	manager.RegisterProvider(&GeneratorProvider{})

	return manager
}

func (m *DefaultMapManager) RegisterProvider(provider MapProvider) {
	m.providers = append(m.providers, provider)
}

func (m *DefaultMapManager) GenerateMapId(generator string, size Size, playerCount int, config GeneratorConfig) string {
	return fmt.Sprintf("%s-%s-%d-%s", generator, size.String(), playerCount, config.String())
}

func (m *DefaultMapManager) GetMap(mapId string, players []Player) (Map, error) {
	if cachedMap, exists := m.cache[mapId]; exists {
		return cachedMap, nil
	}

	parts := strings.Split(mapId, "-")
	if len(parts) < 1 {
		return nil, errors.New("invalid map id format")
	}

	prefix := parts[0]

	for _, provider := range m.providers {
		if provider.CanHandle(prefix) {
			gameMap, err := provider.GetMap(mapId, players)
			if err != nil {
				return nil, err
			}

			m.cache[mapId] = gameMap
			return gameMap, nil
		}
	}

	return nil, fmt.Errorf("no provider found for map id: %s", mapId)
}

// GeneratorProvider 生成器提供者
type GeneratorProvider struct{}

func (p *GeneratorProvider) CanHandle(prefix string) bool {
	return GeneratorExists(prefix)
}

func (p *GeneratorProvider) GetMap(mapId string, players []Player) (Map, error) {
	parts := strings.Split(mapId, "-")
	if len(parts) < 4 {
		return nil, errors.New("invalid generator map id format")
	}

	generator := parts[0]
	sizeStr := parts[1]
	_ = parts[2]
	configStr := strings.Join(parts[3:], "-")

	size, err := parseSize(sizeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse size: %w", err)
	}

	config, err := parseGeneratorConfig(configStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return GenerateMap(generator, size, players, config)
}

// DatabaseProvider 数据库提供者示例
type DatabaseProvider struct {
	// db connection, etc.
}

func (p *DatabaseProvider) CanHandle(prefix string) bool {
	return prefix == "db" || prefix == "saved"
}

func (p *DatabaseProvider) GetMap(mapId string, players []Player) (Map, error) {
	// TODO: 从数据库加载地图
	return nil, errors.New("database provider not implemented")
}

// FileProvider 文件提供者示例
type FileProvider struct {
	basePath string
}

func (p *FileProvider) CanHandle(prefix string) bool {
	return prefix == "file" || prefix == "custom"
}

func (p *FileProvider) GetMap(mapId string, players []Player) (Map, error) {
	// TODO: 从文件加载地图
	return nil, errors.New("file provider not implemented")
}

// 辅助函数
func parseSize(sizeStr string) (Size, error) {
	parts := strings.Split(sizeStr, "x")
	if len(parts) != 2 {
		return Size{}, errors.New("invalid size format, expected WxH")
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return Size{}, err
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return Size{}, err
	}

	return Size{Width: uint16(width), Height: uint16(height)}, nil
}

func parseGeneratorConfig(configStr string) (GeneratorConfig, error) {
	config := DefaultGeneratorConfig()

	parts := strings.Split(configStr, "-")
	for _, part := range parts {
		if len(part) < 2 {
			continue
		}

		prefix := part[0]
		valueStr := part[1:]

		switch prefix {
		case 'm':
			if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
				config.MountainDensity = value
			}
		case 'c':
			if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
				config.CastleDensity = value
			}
		case 'd':
			if value, err := strconv.Atoi(valueStr); err == nil {
				config.MinCastleDistance = value
			}
		case 's':
			if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
				config.Seed = value
			}
		}
	}

	return config, nil
}
