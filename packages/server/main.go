// 游戏模式测试程序
// 作为独立程序运行以测试游戏模式功能
package main

import (
	"fmt"
	"log/slog"
	"os"
	"server/internal/game"
	gamemap "server/internal/game/map"
	"server/internal/queue"
	"time"
)

func main() {
	// 设置日志
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	fmt.Println("=== Starting Game Mode Test ===")

	// 创建消息队列
	q := queue.NewInMemoryQueue()

	// 创建游戏实例 (使用经典1v1模式)
	gameInstance := game.NewGame("test-game-1", q, game.Classic1v1)

	// 启动游戏事件处理
	if err := gameInstance.Start(); err != nil {
		fmt.Printf("Failed to start game: %v\n", err)
		return
	}

	// 测试游戏模式功能
	testGameMode(gameInstance, q)

	// 等待一段时间观察效果
	time.Sleep(2 * time.Second)

	// 停止游戏
	if err := gameInstance.Stop(); err != nil {
		fmt.Printf("Failed to stop game: %v\n", err)
	}

	fmt.Println("Game mode test completed!")
}

func testGameMode(gameInstance *game.Game, q queue.Queue) {
	fmt.Println("=== Testing Game Mode Features ===")

	// 测试游戏模式信息
	fmt.Printf("Testing game mode: Classic1v1\n")

	// 模拟两个玩家加入
	fmt.Println("\n--- Testing Player Join ---")
	publishJoinCommand(q, "test-game-1", "player1", "Player One")
	publishJoinCommand(q, "test-game-1", "player2", "Player Two")

	time.Sleep(500 * time.Millisecond) // 等待事件处理

	// 模拟强制开始投票
	fmt.Println("\n--- Testing Force Start Vote ---")
	publishForceStartCommand(q, "test-game-1", "player1", true)
	publishForceStartCommand(q, "test-game-1", "player2", true)

	time.Sleep(1 * time.Second) // 等待游戏开始

	// 模拟移动操作
	fmt.Println("\n--- Testing Game Actions ---")
	publishMoveCommand(q, "test-game-1", "player1")

	time.Sleep(500 * time.Millisecond)
}

// publishJoinCommand 发布加入游戏指令
func publishJoinCommand(q queue.Queue, gameId, playerId, playerName string) {
	cmd := game.JoinCommand{
		CommandEvent: game.CommandEvent{PlayerId: playerId},
		PlayerName:   playerName,
	}
	q.Publish(fmt.Sprintf("%s/commands", gameId), cmd)
	fmt.Printf("Player %s (%s) join command sent\n", playerId, playerName)
}

// publishForceStartCommand 发布强制开始指令
func publishForceStartCommand(q queue.Queue, gameId, playerId string, isVote bool) {
	cmd := game.ForceStartCommand{
		CommandEvent: game.CommandEvent{PlayerId: playerId},
		IsVote:      isVote,
	}
	q.Publish(fmt.Sprintf("%s/commands", gameId), cmd)
	fmt.Printf("Player %s force start vote: %v\n", playerId, isVote)
}

// publishMoveCommand 发布移动指令
func publishMoveCommand(q queue.Queue, gameId, playerId string) {
	// 发送移动指令，从位置(1,1)向右移动1个兵
	cmd := game.MoveCommand{
		CommandEvent: game.CommandEvent{PlayerId: playerId},
		From:         gamemap.Pos{X: 1, Y: 1},
		Direction:    game.MoveTowardsRight,
		Troops:       1,
	}
	q.Publish(fmt.Sprintf("%s/commands", gameId), cmd)
	fmt.Printf("Move command sent for player %s\n", playerId)
}
