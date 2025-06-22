package game

import (
	"server/internal/game/block"
	gamemap "server/internal/game/map"
	"testing"
)

func createTestCore() *BaseCore {
	size := gamemap.Size{Width: 3, Height: 3}
	blocks := make([][]block.Block, size.Height)
	for i := range blocks {
		blocks[i] = make([]block.Block, size.Width)
		for j := range blocks[i] {
			blocks[i][j] = block.NewBlock(block.BlankName, 0, 0)
		}
	}

	soldierBlock := block.NewBlock(block.SoldierName, 10, 0)
	blocks[1][1] = soldierBlock

	mapInfo := gamemap.Info{
		Id:   "test_map",
		Name: "Test Map",
		Desc: "test",
	}

	testMap := gamemap.NewBaseMap(blocks, size, mapInfo)

	mapManager := createTestMapManager()
	core := NewBaseCore("test_game", TestMode, mapManager)
	core._map = testMap

	core.players = append(core.players, Player{Id: "player1", Name: "Player 1", Status: PlayerStatusInGame, Moves: 2})
	core.players = append(core.players, Player{Id: "player2", Name: "Player 2", Status: PlayerStatusInGame, Moves: 2})
	core.status = StatusInProgress

	return core
}

func TestMoveValidation(t *testing.T) {
	t.Run("valid_move", func(t *testing.T) {
		core := createTestCore()

		move := Move{
			Pos:     gamemap.Pos{X: 2, Y: 2},
			Towards: MoveTowardsRight,
			Num:     5,
		}

		err := core.Move("player1", move)
		if err != nil {
			t.Errorf("Expected valid move to succeed, got error: %v", err)
		}
	})

	t.Run("invalid_game_state", func(t *testing.T) {
		core := createTestCore()
		core.status = StatusWaiting

		move := Move{
			Pos:     gamemap.Pos{X: 2, Y: 2},
			Towards: MoveTowardsRight,
			Num:     5,
		}

		err := core.Move("player1", move)
		if err == nil {
			t.Error("Expected error for move in waiting state")
		}
	})

	t.Run("invalid_player_state", func(t *testing.T) {
		core := createTestCore()

		core.players[0].Status = PlayerStatusWaiting

		move := Move{
			Pos:     gamemap.Pos{X: 2, Y: 2},
			Towards: MoveTowardsRight,
			Num:     5,
		}

		err := core.Move("player1", move)
		if err == nil {
			t.Error("Expected error for move with invalid player state")
		}

		core.players[0].Status = PlayerStatusInGame
	})

	t.Run("no_moves_left", func(t *testing.T) {
		core := createTestCore()

		core.players[0].Moves = 0

		move := Move{
			Pos:     gamemap.Pos{X: 2, Y: 2},
			Towards: MoveTowardsRight,
			Num:     5,
		}

		err := core.Move("player1", move)
		if err == nil {
			t.Error("Expected error for move with no moves left")
		}

		core.players[0].Moves = 2
	})

	t.Run("invalid_position", func(t *testing.T) {
		core := createTestCore()

		move := Move{
			Pos:     gamemap.Pos{X: 10, Y: 10},
			Towards: MoveTowardsRight,
			Num:     5,
		}

		err := core.Move("player1", move)
		if err == nil {
			t.Error("Expected error for invalid position")
		}
	})

	t.Run("invalid_destination", func(t *testing.T) {
		core := createTestCore()

		move := Move{
			Pos:     gamemap.Pos{X: 3, Y: 2},
			Towards: MoveTowardsRight,
			Num:     5,
		}

		err := core.Move("player1", move)
		if err == nil {
			t.Error("Expected error for invalid destination")
		}
	})

	t.Run("wrong_owner", func(t *testing.T) {
		core := createTestCore()

		move := Move{
			Pos:     gamemap.Pos{X: 2, Y: 2},
			Towards: MoveTowardsRight,
			Num:     5,
		}

		err := core.Move("player2", move)
		if err == nil {
			t.Error("Expected error for moving wrong player's piece")
		}
	})

	t.Run("insufficient_troops", func(t *testing.T) {
		core := createTestCore()

		move := Move{
			Pos:     gamemap.Pos{X: 2, Y: 2},
			Towards: MoveTowardsRight,
			Num:     100,
		}

		err := core.Move("player1", move)
		if err == nil {
			t.Error("Expected error for insufficient troops")
		}
	})

	t.Run("special_move_numbers", func(t *testing.T) {
		core := createTestCore()

		move := Move{
			Pos:     gamemap.Pos{X: 2, Y: 2},
			Towards: MoveTowardsRight,
			Num:     0,
		}

		err := core.Move("player1", move)
		if err != nil {
			t.Errorf("Expected move with Num=0 to succeed, got error: %v", err)
		}

		core = createTestCore()
		core._map.SetBlock(gamemap.Pos{X: 2, Y: 2}, block.NewBlock(block.SoldierName, 10, 0))

		move = Move{
			Pos:     gamemap.Pos{X: 2, Y: 2},
			Towards: MoveTowardsRight,
			Num:     1,
		}

		err = core.Move("player1", move)
		if err != nil {
			t.Errorf("Expected move with Num=1 to succeed, got error: %v", err)
		}
	})
}

func TestMoveExecution(t *testing.T) {
	size := gamemap.Size{Width: 3, Height: 1}
	blocks := make([][]block.Block, size.Height)
	for i := range blocks {
		blocks[i] = make([]block.Block, size.Width)
		for j := range blocks[i] {
			blocks[i][j] = block.NewBlock(block.BlankName, 0, 0)
		}
	}

	blocks[0][0] = block.NewBlock(block.SoldierName, 10, 0)
	blocks[0][1] = block.NewBlock(block.BlankName, 0, 0)
	blocks[0][2] = block.NewBlock(block.BlankName, 0, 0)

	mapInfo := gamemap.Info{Id: "test_map", Name: "Test Map", Desc: "test"}
	testMap := gamemap.NewBaseMap(blocks, size, mapInfo)

	mapManager := createTestMapManager()
	core := NewBaseCore("test_game", TestMode, mapManager)
	core._map = testMap
	core.players = append(core.players, Player{Id: "player1", Name: "Player 1", Status: PlayerStatusInGame, Moves: 2})
	core.status = StatusInProgress

	move := Move{
		Pos:     gamemap.Pos{X: 1, Y: 1},
		Towards: MoveTowardsRight,
		Num:     5,
	}

	err := core.Move("player1", move)
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	fromBlock, err := core._map.Block(gamemap.Pos{X: 1, Y: 1})
	if err != nil {
		t.Fatalf("Failed to get from block: %v", err)
	}
	if fromBlock.Num() != 5 {
		t.Errorf("Expected from block to have 5 troops, got %d", fromBlock.Num())
	}

	toBlock, err := core._map.Block(gamemap.Pos{X: 2, Y: 1})
	if err != nil {
		t.Fatalf("Failed to get to block: %v", err)
	}
	if toBlock.Num() != 5 {
		t.Errorf("Expected to block to have 5 troops, got %d", toBlock.Num())
	}

	if toBlock.Owner() != block.Owner(0) {
		t.Errorf("Expected to block to be owned by player 0, got owner %d", toBlock.Owner())
	}
}
