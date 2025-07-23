package gamemap

import (
	"server/internal/game/block"
	"testing"
)

func TestNewMapGenerator_Generate(t *testing.T) {
	// 创建测试配置
	config := GeneratorConfig{
		MountainDensity:   0.2,
		CastleDensity:     0.1,
		MinCastleDistance: 3,
		Seed:              42, // 固定种子确保可重现
	}

	// 创建玩家
	players := []Player{
		{Index: 0, Owner: 1, IsActive: true},
		{Index: 1, Owner: 2, IsActive: true},
	}

	// 创建地图大小
	size := Size{Width: 20, Height: 20}

	// 生成地图
	generator := NewNewMapGenerator(config)
	gameMap, err := generator.Generate(size, players)
	if err != nil {
		t.Fatalf("生成地图失败: %v", err)
	}

	// 验证地图基本属性
	if gameMap.Size().Width != size.Width || gameMap.Size().Height != size.Height {
		t.Errorf("地图大小不匹配，期望 %v，得到 %v", size, gameMap.Size())
	}

	// 验证玩家起始位置有王
	kingCount := 0
	castleCount := 0
	mountainCount := 0
	for y := uint16(0); y < size.Height; y++ {
		for x := uint16(0); x < size.Width; x++ {
			pos := Pos{X: x, Y: y}
			b, err := gameMap.Block(pos)
			if err != nil {
				continue
			}
			switch {
			case b.Meta().Name == block.KingName:
				kingCount++
			case b.Meta().Name == block.CastleName:
				castleCount++
			case b.Meta().Name == block.MountainName:
				mountainCount++
			}
		}
	}

	if kingCount != len(players) {
		t.Errorf("期望 %d 个王，实际找到 %d 个", len(players), kingCount)
	}
	t.Log("地图生成测试通过。与期望的大小", size.Width, "*", size.Height, "匹配。王的数量:", kingCount, "；城堡的数量：", castleCount, "；山脉的数量：", mountainCount, "。")
}

func TestNewMapGenerator_SmallMapValidation(t *testing.T) {
	config := DefaultGeneratorConfig()
	generator := NewNewMapGenerator(config)

	// 测试太小的地图
	size := Size{Width: 5, Height: 5}
	players := []Player{{Index: 0, Owner: 1, IsActive: true}}

	_, err := generator.Generate(size, players)
	if err == nil {
		t.Error("期望小地图生成失败，但成功了")
	}
}

func TestNewMapGenerator_Name(t *testing.T) {
	config := DefaultGeneratorConfig()
	generator := NewNewMapGenerator(config)

	if generator.Name() != "new" {
		t.Errorf("期望生成器名称为 'new'，得到 '%s'", generator.Name())
	}
}
