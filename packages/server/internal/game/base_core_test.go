package game

import (
	"fmt"
	gamemap "server/internal/game/map"
	"server/internal/queue"
	"testing"
	"time"
)

func createTestMapManager() gamemap.MapManager {
	return gamemap.NewMapManager()
}

func TestBaseCore_NewBaseCore(t *testing.T) {
	gameId := "test-game-id"
	mode := TestMode
	mapManager := createTestMapManager()

	core := NewBaseCore(gameId, mode, mapManager)

	if core == nil {
		t.Error("Expected non-nil BaseCore")
	}
	if core.Status() != StatusWaiting {
		t.Errorf("Expected status %v, got %v", StatusWaiting, core.Status())
	}
	if len(core.Players()) != 0 {
		t.Errorf("Expected 0 players, got %d", len(core.Players()))
	}
	if core.TurnNumber() != 0 {
		t.Errorf("Expected turn number 0, got %d", core.TurnNumber())
	}
}

func TestBaseCore_BasicOperations(t *testing.T) {
	t.Run("create_core", func(t *testing.T) {
		mapManager := createTestMapManager()
		core := NewBaseCore("test-game-1", TestMode, mapManager)
		if core == nil {
			t.Fatal("Failed to create BaseCore")
		}

		if core.Status() != StatusWaiting {
			t.Errorf("Expected status waiting, got %v", core.Status())
		}
	})

	t.Run("player_management", func(t *testing.T) {
		mapManager := createTestMapManager()
		core := NewBaseCore("test-game-2", TestMode, mapManager)

		player := Player{Id: "player1", Name: "Player One"}
		err := core.Join(player)
		if err != nil {
			t.Errorf("Failed to join player: %v", err)
		}

		players := core.Players()
		if len(players) != 1 {
			t.Errorf("Expected 1 player, got %d", len(players))
		}

		err = core.Leave("player1")
		if err != nil {
			t.Errorf("Failed to leave player: %v", err)
		}

		players = core.Players()
		if len(players) != 0 {
			t.Errorf("Expected 0 players after leave, got %d", len(players))
		}
	})
}

func TestBaseCore_EdgeCases(t *testing.T) {
	t.Run("empty_game_id", func(t *testing.T) {
		mapManager := createTestMapManager()
		core := NewBaseCore("", TestMode, mapManager)
		if core == nil {
			t.Error("Should handle empty game ID gracefully")
		}
	})

	t.Run("nil_mode", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when creating core with nil mode")
			}
		}()

		mapManager := createTestMapManager()
		var nilMode GameMode
		NewBaseCore("test", nilMode, mapManager)
	})
}

func TestBaseCore_PlayerJoin(t *testing.T) {
	mapManager := createTestMapManager()
	core := NewBaseCore("test-game", TestMode, mapManager)

	t.Run("successful_join", func(t *testing.T) {
		player := Player{
			Id:   "player1",
			Name: "Player One",
		}

		err := core.Join(player)
		if err != nil {
			t.Fatalf("Expected successful join, got error: %v", err)
		}

		if len(core.Players()) != 1 {
			t.Errorf("Expected 1 player, got %d", len(core.Players()))
		}

		if core.Players()[0].Status != PlayerStatusWaiting {
			t.Errorf("Expected player status %s, got %s", PlayerStatusWaiting, core.Players()[0].Status)
		}
	})

	t.Run("duplicate_player_join", func(t *testing.T) {
		player := Player{
			Id:   "player1",
			Name: "Player One Again",
		}

		err := core.Join(player)
		if err == nil {
			t.Error("Expected error for duplicate player join, got nil")
		}

		if len(core.Players()) != 1 {
			t.Errorf("Expected 1 player after duplicate join, got %d", len(core.Players()))
		}
	})

	t.Run("join_after_game_started", func(t *testing.T) {
		player2 := Player{
			Id:   "player2",
			Name: "Player Two",
		}
		core.Join(player2)
		core.Start()

		player3 := Player{
			Id:   "player3",
			Name: "Player Three",
		}

		err := core.Join(player3)
		if err == nil {
			t.Error("Expected error for join after game started, got nil")
		}

		core.Stop()
	})
}

func TestBaseCore_PlayerLeave(t *testing.T) {
	mapManager := createTestMapManager()
	core := NewBaseCore("test-game", TestMode, mapManager)

	t.Run("leave_before_game_start", func(t *testing.T) {
		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		err := core.Leave("player1")
		if err != nil {
			t.Fatalf("Expected successful leave, got error: %v", err)
		}

		if len(core.Players()) != 1 {
			t.Errorf("Expected 1 player after leave, got %d", len(core.Players()))
		}

		if core.Players()[0].Id != "player2" {
			t.Errorf("Expected remaining player to be player2, got %s", core.Players()[0].Id)
		}
	})

	t.Run("leave_during_game", func(t *testing.T) {
		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)
		core.Start()

		err := core.Leave("player1")
		if err != nil {
			t.Fatalf("Expected successful leave, got error: %v", err)
		}

		if len(core.Players()) != 2 {
			t.Errorf("Expected 2 players (disconnected), got %d", len(core.Players()))
		}

		var leftPlayer *Player
		for i := range core.Players() {
			if core.Players()[i].Id == "player1" {
				leftPlayer = &core.Players()[i]
				break
			}
		}

		if leftPlayer == nil {
			t.Error("Left player not found")
		} else if leftPlayer.Status != PlayerStatusDisconnected {
			t.Errorf("Expected left player status %s, got %s", PlayerStatusDisconnected, leftPlayer.Status)
		}

		core.Stop()
	})

	t.Run("leave_nonexistent_player", func(t *testing.T) {
		err := core.Leave("nonexistent")
		if err == nil {
			t.Error("Expected error for leaving nonexistent player, got nil")
		}
	})
}

func TestBaseCore_GetPlayer(t *testing.T) {
	mapManager := createTestMapManager()
	core := NewBaseCore("test-game", TestMode, mapManager)

	player := Player{Id: "player1", Name: "Player One"}
	core.Join(player)

	t.Run("get_existing_player", func(t *testing.T) {
		retrieved, err := core.GetPlayer("player1")
		if err != nil {
			t.Fatalf("Expected successful get, got error: %v", err)
		}

		if retrieved.Id != "player1" {
			t.Errorf("Expected player ID player1, got %s", retrieved.Id)
		}

		if retrieved.Name != "Player One" {
			t.Errorf("Expected player name 'Player One', got %s", retrieved.Name)
		}
	})

	t.Run("get_nonexistent_player", func(t *testing.T) {
		_, err := core.GetPlayer("nonexistent")
		if err == nil {
			t.Error("Expected error for nonexistent player, got nil")
		}
	})
}

func TestBaseCore_ForceStart(t *testing.T) {
	t.Run("force_start_with_enough_players", func(t *testing.T) {
		mapManager := createTestMapManager()
		core := NewBaseCore("test-game", TestMode, mapManager)

		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		err1 := core.ForceStart("player1", true)
		if err1 != nil {
			t.Fatalf("First force start vote failed: %v", err1)
		}

		err2 := core.ForceStart("player2", true)
		if err2 != nil {
			t.Fatalf("Second force start vote failed: %v", err2)
		}

		if core.Status() != StatusInProgress {
			t.Errorf("Expected game status %s, got %s", StatusInProgress, core.Status())
		}

		core.Stop()
	})

	t.Run("force_start_insufficient_players", func(t *testing.T) {
		mapManager := createTestMapManager()
		core := NewBaseCore("test-game", TestMode, mapManager)

		player1 := Player{Id: "player1", Name: "Player One"}
		core.Join(player1)

		err := core.ForceStart("player1", true)
		if err != nil {
			t.Fatalf("Force start vote failed: %v", err)
		}

		if core.Status() != StatusWaiting {
			t.Errorf("Expected game status %s, got %s", StatusWaiting, core.Status())
		}
	})

	t.Run("force_start_cancel_vote", func(t *testing.T) {
		mapManager := createTestMapManager()
		core := NewBaseCore("test-game", TestMode, mapManager)

		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		core.ForceStart("player1", true)
		core.ForceStart("player1", false)
		core.ForceStart("player2", true)

		if core.Status() != StatusWaiting {
			t.Errorf("Expected game status %s after vote cancel, got %s", StatusWaiting, core.Status())
		}
	})
}

func TestBaseCore_Status(t *testing.T) {
	mapManager := createTestMapManager()
	core := NewBaseCore("test-game", TestMode, mapManager)

	t.Run("initial_status", func(t *testing.T) {
		if core.Status() != StatusWaiting {
			t.Errorf("Expected initial status %s, got %s", StatusWaiting, core.Status())
		}
	})

	t.Run("status_after_start", func(t *testing.T) {
		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		core.Start()

		if core.Status() != StatusInProgress {
			t.Errorf("Expected status %s after start, got %s", StatusInProgress, core.Status())
		}

		core.Stop()
	})

	t.Run("status_after_stop", func(t *testing.T) {
		if core.Status() != StatusFinished {
			t.Errorf("Expected status %s after stop, got %s", StatusFinished, core.Status())
		}
	})
}

func TestBaseCore_GetActivePlayerCount(t *testing.T) {
	mapManager := createTestMapManager()
	core := NewBaseCore("test-game", TestMode, mapManager)

	t.Run("no_players", func(t *testing.T) {
		count := core.GetActivePlayerCount()
		if count != 0 {
			t.Errorf("Expected 0 active players, got %d", count)
		}
	})

	t.Run("waiting_players", func(t *testing.T) {
		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		count := core.GetActivePlayerCount()
		if count != 0 {
			t.Errorf("Expected 0 active players in waiting state, got %d", count)
		}
	})

	t.Run("in_game_players", func(t *testing.T) {
		core.Start()

		count := core.GetActivePlayerCount()
		if count != 2 {
			t.Errorf("Expected 2 active players in game, got %d", count)
		}

		core.Stop()
	})
}

func TestBaseCore_IsGameReady(t *testing.T) {
	mapManager := createTestMapManager()
	core := NewBaseCore("test-game", TestMode, mapManager)

	t.Run("not_ready_no_players", func(t *testing.T) {
		if core.IsGameReady() {
			t.Error("Expected game not ready with no players")
		}
	})

	t.Run("not_ready_insufficient_players", func(t *testing.T) {
		player1 := Player{Id: "player1", Name: "Player One"}
		core.Join(player1)

		if core.IsGameReady() {
			t.Error("Expected game not ready with insufficient players")
		}
	})

	t.Run("ready_with_enough_players", func(t *testing.T) {
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player2)

		if !core.IsGameReady() {
			t.Error("Expected game ready with enough players")
		}
	})

	t.Run("not_ready_after_start", func(t *testing.T) {
		core.Start()

		if core.IsGameReady() {
			t.Error("Expected game not ready after start")
		}

		core.Stop()
	})
}

func TestBaseCore_TurnTimer(t *testing.T) {
	t.Run("turn_timer_functionality", func(t *testing.T) {
		testMode := GameMode{
			Name:        "test_mode",
			MaxPlayers:  2,
			MinPlayers:  2,
			TeamSize:    1,
			TurnTime:    time.Millisecond * 100,
			Description: "Test mode",
		}

		mapManager := createTestMapManager()
		core := NewBaseCore("test-game", testMode, mapManager)

		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		var receivedEvents []queue.Event
		core.SetEventHandlers(
			func(event queue.Event) {
				receivedEvents = append(receivedEvents, event)
			},
			func(event queue.Event) {
				receivedEvents = append(receivedEvents, event)
			},
		)

		core.Start()
		time.Sleep(time.Millisecond * 200)
		core.Stop()

		if len(receivedEvents) == 0 {
			t.Error("Expected to receive timer events, got none")
		}
	})
}

func TestBaseCore_EventHandlers(t *testing.T) {
	mapManager := createTestMapManager()
	core := NewBaseCore("test-game", TestMode, mapManager)

	var broadcastEvents []queue.Event
	var controlEvents []queue.Event

	core.SetEventHandlers(
		func(event queue.Event) {
			broadcastEvents = append(broadcastEvents, event)
		},
		func(event queue.Event) {
			controlEvents = append(controlEvents, event)
		},
	)

	player1 := Player{Id: "player1", Name: "Player One"}
	player2 := Player{Id: "player2", Name: "Player Two"}
	core.Join(player1)
	core.Join(player2)

	core.Start()

	if len(broadcastEvents) == 0 {
		t.Error("Expected broadcast events, got none")
	}

	core.Stop()
}

func BenchmarkCore_PlayerOperations(b *testing.B) {
	mapManager := createTestMapManager()
	core := NewBaseCore("bench-test", TestMode, mapManager)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playerId := fmt.Sprintf("player-%d", i%100)
		player := Player{Id: playerId, Name: fmt.Sprintf("Player %d", i)}
		core.Join(player)
		core.Leave(playerId)
	}
}

func BenchmarkCore_EventGeneration(b *testing.B) {
	mapManager := createTestMapManager()
	core := NewBaseCore("bench-events", TestMode, mapManager)

	player1 := Player{Id: "player1", Name: "Player One"}
	player2 := Player{Id: "player2", Name: "Player Two"}
	core.Join(player1)
	core.Join(player2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		core.ForceStart("player1", true)
	}
}
