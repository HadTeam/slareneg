package game

import (
	"server/internal/game/block"
	gamemap "server/internal/game/map"
	"testing"
)

func TestMapDataStructures(t *testing.T) {
	size := gamemap.Size{Width: 10, Height: 10}
	info := gamemap.Info{
		Id:   "test_map",
		Name: "Test Map",
		Desc: "A test map for unit testing",
	}

	t.Run("empty_map_creation", func(t *testing.T) {
		emptyMap := gamemap.NewEmptyBaseMap(size, info)
		if emptyMap == nil {
			t.Error("Expected empty map to be created successfully")
		}

		if emptyMap.Size() != size {
			t.Errorf("Expected map size %v, got %v", size, emptyMap.Size())
		}

		if emptyMap.Info() != info {
			t.Errorf("Expected map info %v, got %v", info, emptyMap.Info())
		}
	})

	t.Run("map_with_blocks", func(t *testing.T) {
		blocks := make(gamemap.Blocks, size.Height)
		for i := range blocks {
			blocks[i] = make([]block.Block, size.Width)
			for j := range blocks[i] {
				blocks[i][j] = block.NewBlock(block.BlankName, 0, 0)
			}
		}

		testBlock := block.NewBlock(block.SoldierName, 5, 1)
		blocks[5][5] = testBlock

		testMap := gamemap.NewBaseMap(blocks, size, info)
		if testMap == nil {
			t.Error("Expected map with blocks to be created successfully")
		}

		retrievedBlock, err := testMap.Block(gamemap.Pos{X: 6, Y: 6})
		if err != nil {
			t.Errorf("Failed to get block: %v", err)
		}

		if retrievedBlock.Num() != 5 {
			t.Errorf("Expected block to have 5 units, got %d", retrievedBlock.Num())
		}

		if retrievedBlock.Owner() != 1 {
			t.Errorf("Expected block owner to be 1, got %d", retrievedBlock.Owner())
		}
	})

	t.Run("invalid_position_access", func(t *testing.T) {
		emptyMap := gamemap.NewEmptyBaseMap(size, info)

		_, err := emptyMap.Block(gamemap.Pos{X: 0, Y: 5})
		if err == nil {
			t.Error("Expected error for X=0 position")
		}

		_, err = emptyMap.Block(gamemap.Pos{X: 5, Y: 0})
		if err == nil {
			t.Error("Expected error for Y=0 position")
		}

		_, err = emptyMap.Block(gamemap.Pos{X: 11, Y: 5})
		if err == nil {
			t.Error("Expected error for X=11 position (out of bounds)")
		}

		_, err = emptyMap.Block(gamemap.Pos{X: 5, Y: 11})
		if err == nil {
			t.Error("Expected error for Y=11 position (out of bounds)")
		}
	})

	t.Run("pos_validation", func(t *testing.T) {
		size := gamemap.Size{Width: 10, Height: 10}

		validPos := gamemap.Pos{X: 5, Y: 5}
		if !size.IsPosValid(validPos) {
			t.Error("Expected valid position to be valid")
		}

		validBoundaryPos1 := gamemap.Pos{X: 1, Y: 1}
		if !size.IsPosValid(validBoundaryPos1) {
			t.Error("Expected boundary position (1,1) to be valid")
		}

		validBoundaryPos2 := gamemap.Pos{X: 10, Y: 10}
		if !size.IsPosValid(validBoundaryPos2) {
			t.Error("Expected boundary position (10,10) to be valid")
		}

		invalidPos1 := gamemap.Pos{X: 0, Y: 5}
		if size.IsPosValid(invalidPos1) {
			t.Error("Expected X=0 to be invalid")
		}

		invalidPos2 := gamemap.Pos{X: 5, Y: 0}
		if size.IsPosValid(invalidPos2) {
			t.Error("Expected Y=0 to be invalid")
		}

		invalidPos3 := gamemap.Pos{X: 11, Y: 5}
		if size.IsPosValid(invalidPos3) {
			t.Error("Expected X=11 to be invalid for width=10")
		}

		invalidPos4 := gamemap.Pos{X: 5, Y: 11}
		if size.IsPosValid(invalidPos4) {
			t.Error("Expected Y=11 to be invalid for height=10")
		}
	})
}

func TestMapFogOfWar(t *testing.T) {
	size := gamemap.Size{Width: 3, Height: 3}
	info := gamemap.Info{Id: "fog_test", Name: "Fog Test", Desc: "Test fog of war"}

	blocks := make(gamemap.Blocks, size.Height)
	for i := range blocks {
		blocks[i] = make([]block.Block, size.Width)
		for j := range blocks[i] {
			blocks[i][j] = block.NewBlock(block.BlankName, 0, 0)
		}
	}

	playerBlock := block.NewBlock(block.SoldierName, 10, 1)
	enemyBlock := block.NewBlock(block.SoldierName, 5, 2)
	blocks[1][1] = playerBlock
	blocks[2][2] = enemyBlock

	testMap := gamemap.NewBaseMap(blocks, size, info)

	sight := make(gamemap.Sight, size.Height)
	for i := range sight {
		sight[i] = make([]bool, size.Width)
	}
	sight[1][1] = true

	err := testMap.Fog([]block.Owner{1}, sight)
	if err != nil {
		t.Errorf("Failed to apply fog: %v", err)
	}
}

func TestMapRoundOperations(t *testing.T) {
	size := gamemap.Size{Width: 2, Height: 2}
	info := gamemap.Info{Id: "round_test", Name: "Round Test", Desc: "Test round operations"}

	emptyMap := gamemap.NewEmptyBaseMap(size, info)

	emptyMap.RoundStart(1)
	emptyMap.RoundEnd(1)
}
