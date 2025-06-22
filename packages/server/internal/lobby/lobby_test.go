package lobby

import (
	"fmt"
	"server/internal/game"
	gamemap "server/internal/game/map"
	"server/internal/queue"
	"sync"
	"testing"
	"time"
)

func createTestLobby() *Lobby {
	q := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()
	return NewLobby(q, mapManager)
}

func TestLobby_NewLobby(t *testing.T) {
	q := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()

	lobby := NewLobby(q, mapManager)

	if lobby == nil {
		t.Fatal("Expected non-nil lobby")
	}

	if lobby.games == nil {
		t.Error("Expected games map to be initialized")
	}

	if len(lobby.games) != 0 {
		t.Errorf("Expected empty games map, got %d games", len(lobby.games))
	}
}

func TestLobby_GetOrCreateGame(t *testing.T) {
	lobby := createTestLobby()

	gameId := "test-game-1"
	gameMode := game.Classic1v1

	t.Run("create_new_game", func(t *testing.T) {
		game1 := lobby.getOrCreateGame(gameId, gameMode)

		if game1 == nil {
			t.Fatal("Expected non-nil game")
		}

		gameList := lobby.GetGameList()
		if len(gameList) != 1 {
			t.Errorf("Expected 1 game in list, got %d", len(gameList))
		}

		if _, exists := gameList[gameId]; !exists {
			t.Error("Expected game to exist in game list")
		}
	})

	t.Run("get_existing_game", func(t *testing.T) {
		game2 := lobby.getOrCreateGame(gameId, gameMode)

		gameList := lobby.GetGameList()
		if len(gameList) != 1 {
			t.Errorf("Expected 1 game in list after second call, got %d", len(gameList))
		}

		if gameList[gameId] != game2 {
			t.Error("Expected same game instance to be returned")
		}
	})
}

func TestLobby_GetGameList(t *testing.T) {
	lobby := createTestLobby()

	t.Run("empty_game_list", func(t *testing.T) {
		gameList := lobby.GetGameList()
		if len(gameList) != 0 {
			t.Errorf("Expected empty game list, got %d games", len(gameList))
		}
	})

	t.Run("populated_game_list", func(t *testing.T) {
		lobby.getOrCreateGame("game1", game.Classic1v1)
		lobby.getOrCreateGame("game2", game.Classic1v1)

		gameList := lobby.GetGameList()
		if len(gameList) != 2 {
			t.Errorf("Expected 2 games in list, got %d", len(gameList))
		}

		if _, exists := gameList["game1"]; !exists {
			t.Error("Expected game1 to exist")
		}

		if _, exists := gameList["game2"]; !exists {
			t.Error("Expected game2 to exist")
		}
	})
}

func TestLobby_RemoveGame(t *testing.T) {
	lobby := createTestLobby()

	gameId := "test-game-remove"
	lobby.getOrCreateGame(gameId, game.Classic1v1)

	if len(lobby.GetGameList()) != 1 {
		t.Error("Expected 1 game before removal")
	}

	lobby.RemoveGame(gameId)

	if len(lobby.GetGameList()) != 0 {
		t.Error("Expected 0 games after removal")
	}
}

func TestLobby_HandleCommand(t *testing.T) {
	lobby := createTestLobby()

	t.Run("create_game_command", func(t *testing.T) {
		cmd := LobbyCommand{
			Type:   "createGame",
			GameId: "cmd-test-game",
			Payload: CreateGamePayload{
				GameMode: game.Classic1v1,
			},
		}

		err := lobby.handleCommand(cmd)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		gameList := lobby.GetGameList()
		if len(gameList) != 1 {
			t.Errorf("Expected 1 game after create command, got %d", len(gameList))
		}

		if _, exists := gameList["cmd-test-game"]; !exists {
			t.Error("Expected created game to exist")
		}
	})

	t.Run("get_game_info_command", func(t *testing.T) {
		cmd := LobbyCommand{
			Type:   "getGameInfo",
			GameId: "cmd-test-game",
		}

		err := lobby.handleCommand(cmd)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("unknown_command", func(t *testing.T) {
		cmd := LobbyCommand{
			Type:   "unknownCommand",
			GameId: "test-game",
		}

		err := lobby.handleCommand(cmd)
		if err == nil {
			t.Error("Expected error for unknown command")
		}
	})
}

func TestLobby_ConcurrentAccess(t *testing.T) {
	lobby := createTestLobby()

	const numGoroutines = 3
	const gamesPerGoroutine = 2

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineId int) {
			defer wg.Done()

			for j := 0; j < gamesPerGoroutine; j++ {
				gameId := fmt.Sprintf("concurrent-game-%d-%d", goroutineId, j)
				game := lobby.getOrCreateGame(gameId, game.Classic1v1)
				if game != nil {
					time.Sleep(1 * time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()

	gameList := lobby.GetGameList()
	expectedCount := numGoroutines * gamesPerGoroutine
	if len(gameList) != expectedCount {
		t.Errorf("Expected %d games after concurrent creation, got %d", expectedCount, len(gameList))
	}
}

func TestLobby_StartStop(t *testing.T) {
	lobby := createTestLobby()

	t.Run("start_lobby", func(t *testing.T) {
		err := lobby.Start()
		if err != nil {
			t.Fatalf("Expected no error starting lobby, got %v", err)
		}

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("stop_lobby", func(t *testing.T) {
		lobby.getOrCreateGame("test-game-stop", game.Classic1v1)

		err := lobby.Stop()
		if err != nil {
			t.Fatalf("Expected no error stopping lobby, got %v", err)
		}
	})
}

func TestLobby_EventPublishing(t *testing.T) {
	q := queue.NewInMemoryQueue()
	mapManager := gamemap.NewMapManager()
	lobby := NewLobby(q, mapManager)

	eventChan := q.Subscribe("lobby/events")

	cmd := LobbyCommand{
		Type:   "createGame",
		GameId: "event-test-game",
		Payload: CreateGamePayload{
			GameMode: game.Classic1v1,
		},
	}

	err := lobby.handleCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	select {
	case event := <-eventChan:
		eventMap, ok := event.(map[string]interface{})
		if !ok {
			t.Fatal("Expected event to be map[string]interface{}")
		}

		if eventType, exists := eventMap["type"]; !exists || eventType != "gameCreated" {
			t.Error("Expected gameCreated event type")
		}

		if gameId, exists := eventMap["gameId"]; !exists || gameId != "event-test-game" {
			t.Error("Expected correct gameId in event")
		}

	case <-time.After(100 * time.Millisecond):
		t.Error("Expected to receive event within timeout")
	}
}

func BenchmarkLobby_CreateGame(b *testing.B) {
	lobby := createTestLobby()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gameId := fmt.Sprintf("bench-game-%d", i)
		lobby.getOrCreateGame(gameId, game.Classic1v1)
	}
}

func BenchmarkLobby_GetGameList(b *testing.B) {
	lobby := createTestLobby()

	for i := 0; i < 100; i++ {
		gameId := fmt.Sprintf("bench-game-%d", i)
		lobby.getOrCreateGame(gameId, game.Classic1v1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lobby.GetGameList()
	}
}
