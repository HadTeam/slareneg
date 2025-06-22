package game

import (
	"fmt"
	"server/internal/queue"
	"testing"
	"time"
)

// TestGameCore_BasicOperations 测试 GameCore 的基本操作
func TestGameCore_BasicOperations(t *testing.T) {
	t.Run("create_core", func(t *testing.T) {
		core := NewBaseCore("test-game-1", TestMode)
		if core == nil {
			t.Fatal("Failed to create BaseCore")
		}

		if core.Status() != StatusWaiting {
			t.Errorf("Expected status waiting, got %v", core.Status())
		}
	})

	t.Run("player_management", func(t *testing.T) {
		core := NewBaseCore("test-game-2", TestMode)

		// 测试添加玩家
		player := Player{Id: "player1", Name: "Player One"}
		err := core.Join(player)
		if err != nil {
			t.Errorf("Failed to join player: %v", err)
		}

		players := core.Players()
		if len(players) != 1 {
			t.Errorf("Expected 1 player, got %d", len(players))
		}

		// 测试移除玩家
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

// TestGameMode_Validation 测试游戏模式验证
func TestGameMode_Validation(t *testing.T) {
	t.Run("classic_1v1_validation", func(t *testing.T) {
		mode := TestMode

		if !mode.ValidatePlayerCount(2) {
			t.Error("TestMode should accept 2 players")
		}

		if mode.ValidatePlayerCount(1) {
			t.Error("TestMode should not accept 1 player")
		}

		if mode.ValidatePlayerCount(3) {
			t.Error("TestMode should not accept 3 players")
		}
	})

	t.Run("custom_mode_creation", func(t *testing.T) {
		customMode := GameMode{
			Name:        "test_custom_mode",
			MaxPlayers:  4,
			MinPlayers:  2,
			TeamSize:    1,
			TurnTime:    time.Second * 15,
			Description: "Custom test mode for validation",
		}

		if !customMode.ValidatePlayerCount(3) {
			t.Error("Custom mode should accept 3 players (within range 2-4)")
		}

		teamCount := customMode.CalculateTeamCount(4)
		if teamCount != 4 {
			t.Errorf("Expected 4 teams for 4 players with TeamSize=1, got %d", teamCount)
		}
	})
}

// TestEvent_Types 测试事件类型
func TestEvent_Types(t *testing.T) {
	t.Run("player_joined_event", func(t *testing.T) {
		event := PlayerJoinedEvent{
			PlayerId:   "test-player",
			PlayerName: "Test Player",
		}

		if event.PlayerId != "test-player" {
			t.Errorf("Expected player ID 'test-player', got '%s'", event.PlayerId)
		}

		if event.PlayerName != "Test Player" {
			t.Errorf("Expected player name 'Test Player', got '%s'", event.PlayerName)
		}
	})

	t.Run("game_started_event", func(t *testing.T) {
		event := GameStartedEvent{
			GameStatus: StatusInProgress,
			TurnNumber: 1,
		}

		if event.GameStatus != StatusInProgress {
			t.Errorf("Expected game status StatusInProgress, got %v", event.GameStatus)
		}

		if event.TurnNumber != 1 {
			t.Errorf("Expected turn number 1, got %d", event.TurnNumber)
		}
	})
}

// TestGame_EdgeCases 测试边界情况
func TestGame_EdgeCases(t *testing.T) {
	t.Run("empty_game_id", func(t *testing.T) {
		core := NewBaseCore("", TestMode)
		if core == nil {
			t.Error("Should handle empty game ID gracefully")
		}
	})

	t.Run("nil_mode", func(t *testing.T) {
		// 这个测试可能会 panic，所以要小心处理
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when creating core with nil mode")
			}
		}()

		var nilMode GameMode
		NewBaseCore("test", nilMode)
	})
}

// BenchmarkCore_PlayerOperations 基准测试：玩家操作
func BenchmarkCore_PlayerOperations(b *testing.B) {
	core := NewBaseCore("bench-test", TestMode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playerId := fmt.Sprintf("player-%d", i%100) // 循环使用100个玩家ID
		player := Player{Id: playerId, Name: fmt.Sprintf("Player %d", i)}
		core.Join(player)
		core.Leave(playerId)
	}
}

// BenchmarkCore_EventGeneration 基准测试：事件生成
func BenchmarkCore_EventGeneration(b *testing.B) {
	core := NewBaseCore("bench-events", TestMode)

	// 预先添加一些玩家
	player1 := Player{Id: "player1", Name: "Player One"}
	player2 := Player{Id: "player2", Name: "Player Two"}
	core.Join(player1)
	core.Join(player2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 生成一些事件
		core.ForceStart("player1", true)
	}
}

// TestGamePackage_FullSuite 完整的游戏包测试套件 - 等同于原来的main.go
func TestGamePackage_FullSuite(t *testing.T) {
	t.Run("complete_game_workflow", func(t *testing.T) {
		fmt.Println("=== Complete Game Workflow Test ===")

		// 1. 创建游戏核心
		core := NewBaseCore("integration-test-game", TestMode)
		if core == nil {
			t.Fatal("Failed to create BaseCore")
		}

		// 2. 添加玩家
		player1 := Player{Id: "player1", Name: "Alice"}
		player2 := Player{Id: "player2", Name: "Bob"}

		err := core.Join(player1)
		if err != nil {
			t.Fatalf("Failed to add player1: %v", err)
		}

		err = core.Join(player2)
		if err != nil {
			t.Fatalf("Failed to add player2: %v", err)
		}

		// 3. 验证玩家数量
		players := core.Players()
		if len(players) != 2 {
			t.Fatalf("Expected 2 players, got %d", len(players))
		}

		// 4. 开始游戏
		err = core.Start()
		if err != nil {
			t.Fatalf("Failed to start game: %v", err)
		}

		// 5. 验证游戏状态
		if core.Status() != StatusInProgress {
			t.Errorf("Expected game status to be InProgress, got %v", core.Status())
		}

		// 6. 测试地图
		gameMap := core.Map()
		if gameMap == nil {
			t.Error("Game map should not be nil")
		}

		// 7. 测试回合系统
		initialTurn := core.TurnNumber()
		if initialTurn != 0 {
			t.Errorf("Expected initial turn to be 0, got %d", initialTurn)
		}

		// 8. 停止游戏
		err = core.Stop()
		if err != nil {
			t.Errorf("Failed to stop game: %v", err)
		}

		fmt.Println("=== Complete Game Workflow Test PASSED ===")
	})

	t.Run("error_handling", func(t *testing.T) {
		fmt.Println("=== Error Handling Test ===")

		core := NewBaseCore("error-test-game", TestMode)

		// 测试在未开始的游戏中进行操作
		err := core.Move("nonexistent", Move{})
		if err == nil {
			t.Error("Expected error when moving in non-started game")
		}

		// 测试添加重复玩家
		player := Player{Id: "duplicate", Name: "Duplicate Player"}
		core.Join(player)
		err = core.Join(player)
		if err == nil {
			t.Error("Expected error when adding duplicate player")
		}

		fmt.Println("=== Error Handling Test PASSED ===")
	})

	t.Run("performance_stress", func(t *testing.T) {
		fmt.Println("=== Performance Stress Test ===")

		// 创建多个游戏实例进行压力测试
		for i := 0; i < 10; i++ {
			gameId := fmt.Sprintf("stress-test-%d", i)
			core := NewBaseCore(gameId, TestMode)

			// 快速添加和移除玩家
			for j := 0; j < 5; j++ {
				player := Player{Id: fmt.Sprintf("player-%d-%d", i, j), Name: fmt.Sprintf("Player %d-%d", i, j)}
				core.Join(player)
				core.Leave(player.Id)
			}
		}

		fmt.Println("=== Performance Stress Test PASSED ===")
	})
}

func TestGame_EventHandling(t *testing.T) {
	// 创建测试环境
	q := queue.NewInMemoryQueue()
	game := NewGame("test-game", q, TestMode)

	// 启动游戏
	if err := game.Start(); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer game.Stop()

	// 订阅广播事件来验证事件转发
	broadcastCh := q.Subscribe("test-game/broadcast")

	t.Run("join_command_handling", func(t *testing.T) {
		// 发送加入命令
		cmd := JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			PlayerName:   "Test Player",
		}
		q.Publish("test-game/commands", cmd)

		// 等待并验证广播事件
		select {
		case event := <-broadcastCh:
			if joinEvent, ok := event.(PlayerJoinedEvent); ok {
				if joinEvent.PlayerId != "player1" {
					t.Errorf("Expected PlayerId 'player1', got '%s'", joinEvent.PlayerId)
				}
				if joinEvent.PlayerName != "Test Player" {
					t.Errorf("Expected PlayerName 'Test Player', got '%s'", joinEvent.PlayerName)
				}
			} else {
				t.Errorf("Expected PlayerJoinedEvent, got %T", event)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for join event")
		}
	})

	t.Run("leave_command_handling", func(t *testing.T) {
		// 发送离开命令
		cmd := LeaveCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
		}
		q.Publish("test-game/commands", cmd)

		// 等待并验证广播事件
		select {
		case event := <-broadcastCh:
			if leaveEvent, ok := event.(PlayerLeftEvent); ok {
				if leaveEvent.PlayerId != "player1" {
					t.Errorf("Expected PlayerId 'player1', got '%s'", leaveEvent.PlayerId)
				}
			} else {
				t.Errorf("Expected PlayerLeftEvent, got %T", event)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for leave event")
		}
	})

	t.Run("force_start_command_handling", func(t *testing.T) {
		// 重新添加两个玩家
		q.Publish("test-game/commands", JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			PlayerName:   "Player 1",
		})
		q.Publish("test-game/commands", JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "player2"},
			PlayerName:   "Player 2",
		})

		// 等待加入事件
		<-broadcastCh
		<-broadcastCh

		// 发送强制开始投票
		cmd := ForceStartCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			IsVote:       true,
		}
		q.Publish("test-game/commands", cmd)

		// 等待投票事件
		select {
		case event := <-broadcastCh:
			if _, ok := event.(ForceStartVoteEvent); !ok {
				t.Errorf("Expected ForceStartVoteEvent, got %T", event)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for force start vote event")
		}
	})
}

func TestGame_ErrorHandling(t *testing.T) {
	q := queue.NewInMemoryQueue()
	game := NewGame("test-game", q, TestMode)

	if err := game.Start(); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer game.Stop()

	// 订阅玩家错误事件
	playerErrorCh := q.Subscribe("test-game/player/invalid-player")

	t.Run("invalid_player_command", func(t *testing.T) {
		// 发送来自不存在玩家的命令
		cmd := LeaveCommand{
			CommandEvent: CommandEvent{PlayerId: "invalid-player"},
		}
		q.Publish("test-game/commands", cmd)

		// 等待错误事件
		select {
		case event := <-playerErrorCh:
			if errorEvent, ok := event.(PlayerErrorEvent); ok {
				if errorEvent.PlayerId != "invalid-player" {
					t.Errorf("Expected error for 'invalid-player', got '%s'", errorEvent.PlayerId)
				}
			} else {
				t.Errorf("Expected PlayerErrorEvent, got %T", event)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Timeout waiting for error event")
		}
	})
}

func TestGame_ControlEvents(t *testing.T) {
	q := queue.NewInMemoryQueue()
	game := NewGame("test-game", q, TestMode)

	if err := game.Start(); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer game.Stop()

	t.Run("start_game_control", func(t *testing.T) {
		// 添加足够的玩家
		q.Publish("test-game/commands", JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			PlayerName:   "Player 1",
		})
		q.Publish("test-game/commands", JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "player2"},
			PlayerName:   "Player 2",
		})

		time.Sleep(50 * time.Millisecond) // 等待玩家加入

		// 发送开始游戏控制事件
		controlEvent := StartGameControl{}
		q.Publish("test-game/control", controlEvent)

		time.Sleep(50 * time.Millisecond) // 等待处理

		// 验证游戏状态
		if game.Core().Status() != StatusInProgress {
			t.Errorf("Expected game status to be %s, got %s", StatusInProgress, game.Core().Status())
		}
	})

	t.Run("stop_game_control", func(t *testing.T) {
		// 发送停止游戏控制事件
		controlEvent := StopGameControl{}
		q.Publish("test-game/control", controlEvent)

		time.Sleep(50 * time.Millisecond) // 等待处理

		// 验证游戏状态
		if game.Core().Status() != StatusFinished {
			t.Errorf("Expected game status to be %s, got %s", StatusFinished, game.Core().Status())
		}
	})
}

func TestGame_Integration(t *testing.T) {
	// 完整的游戏流程集成测试
	q := queue.NewInMemoryQueue()
	game := NewGame("integration-test", q, TestMode)

	if err := game.Start(); err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	defer game.Stop()

	// 订阅所有相关事件
	broadcastCh := q.Subscribe("integration-test/broadcast")

	// 模拟完整游戏流程
	t.Run("complete_game_flow", func(t *testing.T) {
		// 1. 玩家加入
		q.Publish("integration-test/commands", JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			PlayerName:   "Player One",
		})
		q.Publish("integration-test/commands", JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "player2"},
			PlayerName:   "Player Two",
		})

		// 等待加入事件
		<-broadcastCh // player1 joined
		<-broadcastCh // player2 joined

		// 2. 强制开始投票
		q.Publish("integration-test/commands", ForceStartCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			IsVote:       true,
		})
		q.Publish("integration-test/commands", ForceStartCommand{
			CommandEvent: CommandEvent{PlayerId: "player2"},
			IsVote:       true,
		})

		// 等待投票和游戏开始事件
		<-broadcastCh // vote event 1
		<-broadcastCh // vote event 2

		// 可能收到游戏开始事件
		select {
		case event := <-broadcastCh:
			if _, ok := event.(GameStartedEvent); ok {
				// 游戏已开始
			}
		case <-time.After(100 * time.Millisecond):
			// 可能没有自动开始
		}

		// 验证最终状态
		if len(game.Core().Players()) != 2 {
			t.Errorf("Expected 2 players, got %d", len(game.Core().Players()))
		}
	})
}
