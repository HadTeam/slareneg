package game

import (
	"server/internal/queue"
	"testing"
	"time"
)

// TestBaseCore_NewBaseCore 测试 BaseCore 创建
func TestBaseCore_NewBaseCore(t *testing.T) {
	gameId := "test-game-1"
	mode := TestMode

	core := NewBaseCore(gameId, mode)

	if core.gameId != gameId {
		t.Errorf("Expected gameId %s, got %s", gameId, core.gameId)
	}

	if core.status != StatusWaiting {
		t.Errorf("Expected initial status %s, got %s", StatusWaiting, core.status)
	}

	if len(core.players) != 0 {
		t.Errorf("Expected empty players slice, got %d players", len(core.players))
	}

	if core.turnNumber != 0 {
		t.Errorf("Expected initial turn number 0, got %d", core.turnNumber)
	}
}

// TestBaseCore_PlayerJoin 测试玩家加入功能
func TestBaseCore_PlayerJoin(t *testing.T) {
	core := NewBaseCore("test-game", TestMode)

	t.Run("successful_join", func(t *testing.T) {
		player := Player{
			Id:   "player1",
			Name: "Player One",
		}

		err := core.Join(player)
		if err != nil {
			t.Fatalf("Expected successful join, got error: %v", err)
		}

		if len(core.players) != 1 {
			t.Errorf("Expected 1 player, got %d", len(core.players))
		}

		if core.players[0].Status != PlayerStatusWaiting {
			t.Errorf("Expected player status %s, got %s", PlayerStatusWaiting, core.players[0].Status)
		}
	})

	t.Run("duplicate_player_join", func(t *testing.T) {
		player := Player{
			Id:   "player1", // 相同ID
			Name: "Player One Again",
		}

		err := core.Join(player)
		if err == nil {
			t.Error("Expected error for duplicate player join, got nil")
		}

		if len(core.players) != 1 {
			t.Errorf("Expected 1 player after duplicate join, got %d", len(core.players))
		}
	})

	t.Run("join_after_game_started", func(t *testing.T) {
		// 添加第二个玩家并开始游戏
		player2 := Player{
			Id:   "player2",
			Name: "Player Two",
		}
		core.Join(player2)
		core.Start()

		// 尝试在游戏开始后加入
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

// TestBaseCore_PlayerLeave 测试玩家离开功能
func TestBaseCore_PlayerLeave(t *testing.T) {
	t.Run("leave_before_game_start", func(t *testing.T) {
		core := NewBaseCore("test-game", TestMode)

		// 添加玩家
		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		// 玩家离开
		err := core.Leave("player1")
		if err != nil {
			t.Fatalf("Expected successful leave, got error: %v", err)
		}

		if len(core.players) != 1 {
			t.Errorf("Expected 1 player after leave, got %d", len(core.players))
		}

		if core.players[0].Id != "player2" {
			t.Errorf("Expected remaining player to be player2, got %s", core.players[0].Id)
		}
	})

	t.Run("leave_during_game", func(t *testing.T) {
		core := NewBaseCore("test-game", TestMode)

		// 添加玩家并开始游戏
		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)
		core.Start()

		// 玩家在游戏中离开
		err := core.Leave("player1")
		if err != nil {
			t.Fatalf("Expected successful leave, got error: %v", err)
		}

		// 检查玩家状态变为断线而不是移除
		if len(core.players) != 2 {
			t.Errorf("Expected 2 players (disconnected), got %d", len(core.players))
		}

		// 找到离开的玩家
		var leftPlayer *Player
		for i := range core.players {
			if core.players[i].Id == "player1" {
				leftPlayer = &core.players[i]
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
		core := NewBaseCore("test-game", TestMode)

		err := core.Leave("nonexistent")
		if err == nil {
			t.Error("Expected error for leaving nonexistent player, got nil")
		}
	})
}

// TestBaseCore_GetPlayer 测试获取玩家功能
func TestBaseCore_GetPlayer(t *testing.T) {
	core := NewBaseCore("test-game", TestMode)

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

// TestBaseCore_ForceStart 测试强制开始功能
func TestBaseCore_ForceStart(t *testing.T) {
	t.Run("force_start_with_enough_players", func(t *testing.T) {
		core := NewBaseCore("test-game", TestMode)

		// 添加足够的玩家
		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		// 两个玩家都投票强制开始
		err1 := core.ForceStart("player1", true)
		if err1 != nil {
			t.Fatalf("First force start vote failed: %v", err1)
		}

		err2 := core.ForceStart("player2", true)
		if err2 != nil {
			t.Fatalf("Second force start vote failed: %v", err2)
		}

		// 游戏应该已经开始
		if core.status != StatusInProgress {
			t.Errorf("Expected game status %s, got %s", StatusInProgress, core.status)
		}

		core.Stop()
	})

	t.Run("force_start_insufficient_players", func(t *testing.T) {
		core := NewBaseCore("test-game", TestMode)

		// 只添加一个玩家
		player1 := Player{Id: "player1", Name: "Player One"}
		core.Join(player1)

		err := core.ForceStart("player1", true)
		if err != nil {
			t.Fatalf("Force start vote failed: %v", err)
		}

		// 游戏不应该开始（玩家不足）
		if core.status != StatusWaiting {
			t.Errorf("Expected game status %s, got %s", StatusWaiting, core.status)
		}
	})

	t.Run("force_start_cancel_vote", func(t *testing.T) {
		core := NewBaseCore("test-game", TestMode)

		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		// 玩家1投票，然后取消
		core.ForceStart("player1", true)
		core.ForceStart("player1", false)

		// 玩家2投票
		core.ForceStart("player2", true)

		// 游戏不应该开始（玩家1取消了投票）
		if core.status != StatusWaiting {
			t.Errorf("Expected game status %s after vote cancel, got %s", StatusWaiting, core.status)
		}
	})
}

// TestBaseCore_Status 测试状态相关功能
func TestBaseCore_Status(t *testing.T) {
	core := NewBaseCore("test-game", TestMode)

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

// TestBaseCore_GetActivePlayerCount 测试活跃玩家计数
func TestBaseCore_GetActivePlayerCount(t *testing.T) {
	core := NewBaseCore("test-game", TestMode)

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

// TestBaseCore_IsGameReady 测试游戏就绪检查
func TestBaseCore_IsGameReady(t *testing.T) {
	core := NewBaseCore("test-game", TestMode)

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

// TestBaseCore_TurnTimer 测试回合定时器
func TestBaseCore_TurnTimer(t *testing.T) {
	t.Run("turn_timer_functionality", func(t *testing.T) {
		// 创建短回合时间的模式用于测试
		testMode := GameMode{
			Name:        "test_mode",
			MaxPlayers:  2,
			MinPlayers:  2,
			TeamSize:    1,
			TurnTime:    time.Millisecond * 100, // 短时间用于测试
			Description: "Test mode",
		}

		core := NewBaseCore("test-game", testMode)

		player1 := Player{Id: "player1", Name: "Player One"}
		player2 := Player{Id: "player2", Name: "Player Two"}
		core.Join(player1)
		core.Join(player2)

		// 设置事件回调来捕获定时器事件
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

		// 等待定时器触发
		time.Sleep(time.Millisecond * 200)

		core.Stop()

		// 验证是否收到了定时器相关的事件
		if len(receivedEvents) == 0 {
			t.Error("Expected to receive timer events, got none")
		}
	})
}

// TestBaseCore_EventHandlers 测试事件处理器
func TestBaseCore_EventHandlers(t *testing.T) {
	core := NewBaseCore("test-game", TestMode)

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

	// 添加玩家应该触发广播事件
	player1 := Player{Id: "player1", Name: "Player One"}
	player2 := Player{Id: "player2", Name: "Player Two"}
	core.Join(player1)
	core.Join(player2)

	// 开始游戏应该触发事件
	core.Start()

	// 验证事件被正确调用
	if len(broadcastEvents) == 0 {
		t.Error("Expected broadcast events, got none")
	}

	core.Stop()
}
