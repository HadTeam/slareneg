# 地图生成器系统

本包提供了一个插拔式的地图生成器系统，支持通过配置调整各种方块的生成倍率。

## 基本使用

### 使用默认配置生成地图

```go
players := []gamemap.Player{
    {Index: 0, Owner: block.Owner(0), IsActive: true},
    {Index: 1, Owner: block.Owner(1), IsActive: true},
}

// 使用默认配置
gameMap, err := gamemap.GenerateMap("base", size, info, players)
```

### 使用自定义配置

```go
// 创建自定义配置
config := gamemap.GeneratorConfig{
    MountainDensity:   0.7,  // 山脉密度 (0.0-1.0, 默认0.5)
    CastleDensity:     0.8,  // 城堡密度 (0.0-1.0, 默认0.5)
    MinCastleDistance: 4,    // 城堡间最小距离4格
}

// 直接传递配置参数
gameMap, err := gamemap.GenerateMap("base", size, info, players, config)
```

## 常用配置示例

```go
// 高山地形 - 更多山脉和城堡
highMountainConfig := gamemap.GeneratorConfig{
    MountainDensity:   0.8,
    CastleDensity:     0.7,
    MinCastleDistance: 4,
}

// 开阔地形 - 适合快速游戏
openTerrainConfig := gamemap.GeneratorConfig{
    MountainDensity:   0.2,
    CastleDensity:     0.3,
    MinCastleDistance: 6,
}

// 丰富资源 - 大量城堡
richResourceConfig := gamemap.GeneratorConfig{
    MountainDensity:   0.5,
    CastleDensity:     0.9,
    MinCastleDistance: 3,
}

## 配置参数说明

### 密度系统
所有密度参数都是 0.0-1.0 的浮点数，以 0.5 为默认档位：
- `0.0`: 最少生成
- `0.5`: 默认生成量
- `1.0`: 最多生成

### 具体参数
- `MountainDensity`: 山脉密度，影响地图中山脉的生成概率
- `CastleDensity`: 城堡密度，影响中性城堡的数量
- `MinCastleDistance`: 城堡间最小距离（格子数）

## 地图管理器

推荐使用地图管理器来统一管理地图的获取：

```go
// 创建地图管理器
manager := gamemap.NewMapManager()

// 生成地图 ID
config := gamemap.GeneratorConfig{
    MountainDensity: 0.7,
    CastleDensity:   0.8,
    Seed:           12345, // 指定种子确保可重现
}
mapId := manager.GenerateMapId("base", size, len(players), config)
// 结果：base-20x20-2-m0.7-c0.8-d5-s12345

// 获取地图（会自动缓存）
gameMap, err := manager.GetMap(mapId, players)
```

### Map ID 格式

Map ID 采用以下格式：`{generator}-{size}-{playerCount}-{config}`

- `generator`: 生成器名称（如 "base"）
- `size`: 地图尺寸（如 "20x20"）  
- `playerCount`: 玩家数量（如 "2"）
- `config`: 配置字符串（如 "m0.7-c0.8-d5-s12345"）

### 扩展地图提供者

系统支持多种地图源，通过前缀路由：

```go
// 生成器提供者（默认注册）
"base-20x20-2-m0.5-c0.5-d5-s0"  // 使用 base 生成器

// 数据库提供者
"db-map123"                       // 从数据库加载

// 文件提供者  
"file-custom_map"                 // 从文件加载

// 自定义提供者
type MyProvider struct{}

func (p *MyProvider) CanHandle(prefix string) bool {
    return prefix == "my"
}

func (p *MyProvider) GetMap(mapId string, players []Player) (Map, error) {
    // 实现自定义逻辑
    return gameMap, nil
}

// 注册提供者
manager.RegisterProvider(&MyProvider{})
```

## 种子系统

指定种子可确保相同参数生成相同地图：

```go
config1 := gamemap.GeneratorConfig{
    MountainDensity: 0.5,
    CastleDensity:   0.5,
    Seed:           12345,
}

config2 := gamemap.GeneratorConfig{
    MountainDensity: 0.5,
    CastleDensity:   0.5,
    Seed:           12345,
}

// 两次生成的地图完全相同
map1, _ := gamemap.GenerateMap("base", size, info, players, config1)
map2, _ := gamemap.GenerateMap("base", size, info, players, config2)
``` 