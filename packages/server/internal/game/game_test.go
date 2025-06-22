package game

import (
	gamemap "server/internal/game/map"
	"server/internal/queue"
	"testing"
	"time"
)

func TestGame_NewGame(t *testing.T) {
	gameId := "test-game"
	queue := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()

	game := NewGame(gameId, queue, TestMode, mapManager)

	if game == nil {
		t.Fatal("Expected non-nil game")
	}

	if game.Core() == nil {
		t.Error("Expected non-nil core")
	}

	if game.Core().Status() != StatusWaiting {
		t.Errorf("Expected initial status %s, got %s", StatusWaiting, game.Core().Status())
	}
}

func TestGame_StartStop(t *testing.T) {
	gameId := "test-game-start-stop"
	queue := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()
	game := NewGame(gameId, queue, TestMode, mapManager)

	t.Run("start_game", func(t *testing.T) {
		err := game.Start()
		if err != nil {
			t.Fatalf("Expected no error starting game, got %v", err)
		}

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("stop_game", func(t *testing.T) {
		err := game.Stop()
		if err != nil {
			t.Fatalf("Expected no error stopping game, got %v", err)
		}
	})
}

func TestGame_Core(t *testing.T) {
	gameId := "test-game-core"
	queue := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()
	game := NewGame(gameId, queue, TestMode, mapManager)

	core := game.Core()
	if core == nil {
		t.Fatal("Expected non-nil core")
	}

	if core.Status() != StatusWaiting {
		t.Errorf("Expected status %s, got %s", StatusWaiting, core.Status())
	}

	if len(core.Players()) != 0 {
		t.Errorf("Expected 0 players initially, got %d", len(core.Players()))
	}
}

func TestGame_EventHandling(t *testing.T) {
	gameId := "test-game-events"
	queue := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()
	game := NewGame(gameId, queue, TestMode, mapManager)

	err := game.Start()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer game.Stop()

	t.Run("join_command", func(t *testing.T) {
		joinCmd := JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "test-player"},
			PlayerName:   "Test Player",
		}

		queue.Publish(gameId+"/commands", joinCmd)
		time.Sleep(50 * time.Millisecond)

		players := game.Core().Players()
		if len(players) != 1 {
			t.Errorf("Expected 1 player after join, got %d", len(players))
		}
	})

	t.Run("leave_command", func(t *testing.T) {
		leaveCmd := LeaveCommand{
			CommandEvent: CommandEvent{PlayerId: "test-player"},
		}

		queue.Publish(gameId+"/commands", leaveCmd)
		time.Sleep(50 * time.Millisecond)

		players := game.Core().Players()
		if len(players) != 0 {
			t.Errorf("Expected 0 players after leave, got %d", len(players))
		}
	})
}

func TestGame_CommandValidation(t *testing.T) {
	gameId := "test-game-validation"
	queue := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()
	game := NewGame(gameId, queue, TestMode, mapManager)

	err := game.Start()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer game.Stop()

	t.Run("invalid_command_type", func(t *testing.T) {
		invalidCmd := "invalid-command"
		queue.Publish(gameId+"/commands", invalidCmd)
		time.Sleep(50 * time.Millisecond)

		// 游戏应该继续正常运行
		if game.Core().Status() != StatusWaiting {
			t.Error("Game should continue running after invalid command")
		}
	})
}

func BenchmarkGame_Creation(b *testing.B) {
	queue := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gameId := "bench-game"
		NewGame(gameId, queue, TestMode, mapManager)
	}
}

func BenchmarkGame_CommandProcessing(b *testing.B) {
	gameId := "bench-game-commands"
	queue := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()
	game := NewGame(gameId, queue, TestMode, mapManager)

	game.Start()
	defer game.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		joinCmd := JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "bench-player"},
			PlayerName:   "Bench Player",
		}
		queue.Publish(gameId+"/commands", joinCmd)

		leaveCmd := LeaveCommand{
			CommandEvent: CommandEvent{PlayerId: "bench-player"},
		}
		queue.Publish(gameId+"/commands", leaveCmd)
	}
}
