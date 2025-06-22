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

## 扩展

要添加新的生成器，只需实现 `MapGenerator` 接口：

```go
type MyGenerator struct{}

func (g *MyGenerator) Name() string {
    return "my_generator"
}

func (g *MyGenerator) Generate(size Size, info Info, players []Player) (Map, error) {
    // 实现生成逻辑
    return gameMap, nil
}

// 注册生成器
gamemap.RegisterGenerator("my_generator", func(size Size, info Info, players []Player) (Map, error) {
    generator := &MyGenerator{}
    return generator.Generate(size, info, players)
})
``` 