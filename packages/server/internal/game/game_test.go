package game

import (
	"fmt"
	gamemap "server/internal/game/map"
	"server/internal/queue"
	"testing"
	"time"
)

func TestGame_FullWorkflow(t *testing.T) {
	t.Run("complete_1v1_game_simulation", func(t *testing.T) {
		fmt.Println("=== 1v1 Game Simulation ===")

		// 1. 创建游戏实例
		gameId := "simulation-1v1"
		queue := queue.NewInMemoryQueue()
		mapManager := createTestMapManager()
		game := NewGame(gameId, queue, TestMode, mapManager)

		// 启动游戏事件处理
		err := game.Start()
		if err != nil {
			t.Fatalf("Failed to start game: %v", err)
		}
		defer game.Stop()

		// 2. 模拟玩家加入
		fmt.Println("Players joining...")

		// 玩家1加入
		joinCmd1 := JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "alice"},
			PlayerName:   "Alice",
		}
		queue.Publish(fmt.Sprintf("%s/commands", gameId), joinCmd1)

		// 玩家2加入
		joinCmd2 := JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "bob"},
			PlayerName:   "Bob",
		}
		queue.Publish(fmt.Sprintf("%s/commands", gameId), joinCmd2)

		// 等待玩家加入处理
		time.Sleep(100 * time.Millisecond)

		// 验证玩家数量
		players := game.Core().Players()
		if len(players) != 2 {
			t.Fatalf("Expected 2 players, got %d", len(players))
		}
		fmt.Printf("Players joined: %d\n", len(players))

		// 3. 模拟强制开始投票
		fmt.Println("Force start voting...")

		forceStartCmd1 := ForceStartCommand{
			CommandEvent: CommandEvent{PlayerId: "alice"},
			IsVote:       true,
		}
		queue.Publish(fmt.Sprintf("%s/commands", gameId), forceStartCmd1)

		forceStartCmd2 := ForceStartCommand{
			CommandEvent: CommandEvent{PlayerId: "bob"},
			IsVote:       true,
		}
		queue.Publish(fmt.Sprintf("%s/commands", gameId), forceStartCmd2)

		// 等待投票处理和游戏开始
		time.Sleep(200 * time.Millisecond)

		// 验证游戏状态
		status := game.Core().Status()
		fmt.Printf("Game status after voting: %s\n", status)

		// 4. 验证地图初始化
		gameMap := game.Core().Map()
		if gameMap == nil {
			t.Error("Game map should be initialized")
		} else {
			fmt.Printf("Map initialized: %dx%d\n", gameMap.Size().Width, gameMap.Size().Height)
		}

		// 5. 模拟游戏回合
		if status == StatusInProgress {
			fmt.Println("Simulating game turns...")

			// 模拟几个回合的移动
			for turn := 1; turn <= 3; turn++ {
				fmt.Printf("Turn %d\n", turn)

				// 模拟玩家移动（这里只是示例，实际需要有效的地图坐标）
				moveCmd := MoveCommand{
					CommandEvent: CommandEvent{PlayerId: "alice"},
					From:         gamemap.Pos{X: 1, Y: 1},
					Direction:    MoveTowardsRight,
					Troops:       1,
				}
				queue.Publish(fmt.Sprintf("%s/commands", gameId), moveCmd)

				time.Sleep(50 * time.Millisecond)
			}
		}

		fmt.Println("=== 1v1 Game Simulation Completed ===")
	})
}

func TestGame_MultiPlayerStress(t *testing.T) {
	t.Run("concurrent_games_stress", func(t *testing.T) {
		fmt.Println("=== Concurrent Games Stress Test ===")

		queue := queue.NewInMemoryQueue()
		mapManager := createTestMapManager()

		const numGames = 5
		const playersPerGame = 2

		var games []*Game

		// 创建多个并发游戏
		for i := 0; i < numGames; i++ {
			gameId := fmt.Sprintf("stress-game-%d", i)
			game := NewGame(gameId, queue, TestMode, mapManager)
			games = append(games, game)

			err := game.Start()
			if err != nil {
				t.Fatalf("Failed to start game %d: %v", i, err)
			}

			// 为每个游戏添加玩家
			for j := 0; j < playersPerGame; j++ {
				playerId := fmt.Sprintf("game%d-player%d", i, j)
				joinCmd := JoinCommand{
					CommandEvent: CommandEvent{PlayerId: playerId},
					PlayerName:   fmt.Sprintf("Player %d-%d", i, j),
				}
				queue.Publish(fmt.Sprintf("%s/commands", gameId), joinCmd)
			}
		}

		// 等待所有玩家加入
		time.Sleep(200 * time.Millisecond)

		// 验证所有游戏状态
		for i, game := range games {
			players := game.Core().Players()
			if len(players) != playersPerGame {
				t.Errorf("Game %d: expected %d players, got %d", i, playersPerGame, len(players))
			}
		}

		// 清理
		for _, game := range games {
			game.Stop()
		}

		fmt.Printf("Successfully tested %d concurrent games\n", numGames)
		fmt.Println("=== Concurrent Games Stress Test Completed ===")
	})
}

func TestGame_ErrorRecovery(t *testing.T) {
	t.Run("handles_player_disconnection", func(t *testing.T) {
		fmt.Println("=== Player Disconnection Recovery Test ===")

		gameId := "disconnection-test"
		queue := queue.NewInMemoryQueue()
		mapManager := createTestMapManager()
		game := NewGame(gameId, queue, TestMode, mapManager)

		err := game.Start()
		if err != nil {
			t.Fatalf("Failed to start game: %v", err)
		}
		defer game.Stop()

		// 添加玩家
		joinCmd := JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "test-player"},
			PlayerName:   "Test Player",
		}
		queue.Publish(fmt.Sprintf("%s/commands", gameId), joinCmd)
		time.Sleep(50 * time.Millisecond)

		// 模拟玩家离开
		leaveCmd := LeaveCommand{
			CommandEvent: CommandEvent{PlayerId: "test-player"},
		}
		queue.Publish(fmt.Sprintf("%s/commands", gameId), leaveCmd)
		time.Sleep(50 * time.Millisecond)

		// 验证玩家已离开
		players := game.Core().Players()
		if len(players) != 0 {
			t.Errorf("Expected 0 players after leave, got %d", len(players))
		}

		fmt.Println("Player disconnection handled successfully")
		fmt.Println("=== Player Disconnection Recovery Test Completed ===")
	})

	t.Run("handles_invalid_commands_gracefully", func(t *testing.T) {
		fmt.Println("=== Invalid Commands Handling Test ===")

		gameId := "invalid-commands-test"
		queue := queue.NewInMemoryQueue()
		mapManager := createTestMapManager()
		game := NewGame(gameId, queue, TestMode, mapManager)

		err := game.Start()
		if err != nil {
			t.Fatalf("Failed to start game: %v", err)
		}
		defer game.Stop()

		// 发送无效命令 - 不存在的玩家离开
		leaveCmd := LeaveCommand{
			CommandEvent: CommandEvent{PlayerId: "nonexistent-player"},
		}
		queue.Publish(fmt.Sprintf("%s/commands", gameId), leaveCmd)
		time.Sleep(50 * time.Millisecond)

		// 游戏应该仍然正常运行
		if game.Core().Status() != StatusWaiting {
			t.Errorf("Game should still be waiting after invalid command")
		}

		fmt.Println("Invalid commands handled gracefully")
		fmt.Println("=== Invalid Commands Handling Test Completed ===")
	})
}

func TestGame_SystemIntegration(t *testing.T) {
	t.Run("full_system_integration", func(t *testing.T) {
		fmt.Println("=== Full System Integration Test ===")

		// 这个测试模拟完整的系统流程：
		// 1. 游戏创建 -> 2. 玩家加入 -> 3. 游戏开始 -> 4. 游戏进行 -> 5. 游戏结束

		gameId := "integration-test"
		queue := queue.NewInMemoryQueue()
		mapManager := createTestMapManager()
		game := NewGame(gameId, queue, TestMode, mapManager)

		// 启动游戏
		err := game.Start()
		if err != nil {
			t.Fatalf("Failed to start game: %v", err)
		}
		defer game.Stop()

		// 阶段1：玩家加入
		fmt.Println("Phase 1: Players joining")
		players := []string{"alice", "bob"}
		for _, playerId := range players {
			joinCmd := JoinCommand{
				CommandEvent: CommandEvent{PlayerId: playerId},
				PlayerName:   playerId,
			}
			queue.Publish(fmt.Sprintf("%s/commands", gameId), joinCmd)
		}
		time.Sleep(100 * time.Millisecond)

		if len(game.Core().Players()) != 2 {
			t.Fatalf("Expected 2 players, got %d", len(game.Core().Players()))
		}

		// 阶段2：强制开始
		fmt.Println("Phase 2: Force starting game")
		for _, playerId := range players {
			forceStartCmd := ForceStartCommand{
				CommandEvent: CommandEvent{PlayerId: playerId},
				IsVote:       true,
			}
			queue.Publish(fmt.Sprintf("%s/commands", gameId), forceStartCmd)
		}
		time.Sleep(200 * time.Millisecond)

		// 阶段3：验证游戏状态
		fmt.Println("Phase 3: Validating game state")
		status := game.Core().Status()
		fmt.Printf("Final game status: %s\n", status)

		// 验证地图存在
		if game.Core().Map() == nil {
			t.Error("Game map should exist after game start")
		}

		// 阶段4：清理
		fmt.Println("Phase 4: Cleanup")
		stopControl := StopGameControl{}
		queue.Publish(fmt.Sprintf("%s/control", gameId), stopControl)
		time.Sleep(50 * time.Millisecond)

		if game.Core().Status() != StatusFinished {
			t.Errorf("Expected game to be finished, got %s", game.Core().Status())
		}

		fmt.Println("=== Full System Integration Test Completed ===")
	})
}

func BenchmarkGame_FullWorkflow(b *testing.B) {
	mapManager := createTestMapManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gameId := fmt.Sprintf("bench-game-%d", i)
		queue := queue.NewInMemoryQueue()
		game := NewGame(gameId, queue, TestMode, mapManager)

		game.Start()

		// 快速添加玩家和开始游戏
		for j := 0; j < 2; j++ {
			joinCmd := JoinCommand{
				CommandEvent: CommandEvent{PlayerId: fmt.Sprintf("player-%d-%d", i, j)},
				PlayerName:   fmt.Sprintf("Player %d", j),
			}
			queue.Publish(fmt.Sprintf("%s/commands", gameId), joinCmd)
		}

		game.Stop()
	}
}
