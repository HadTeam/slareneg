package game

import (
	"context"
	"fmt"
	"log/slog"
	"server/internal/queue"
	"time"
)

type Game struct {
	BaseCore

	Id    string
	Queue queue.Queue

	// 上下文管理
	ctx    context.Context
	cancel context.CancelFunc

	// 消息通道
	playerEventCh  <-chan queue.Event
	controlEventCh <-chan queue.Event
}

// NewGame 创建新的游戏实例
func NewGame(gameId string, q queue.Queue) *Game {
	return &Game{
		BaseCore: *NewBaseCore(gameId, q, Classic1v1),
		Id:       gameId,
		Queue:    q,
	}
}

// InitGame 初始化游戏
func (g *Game) InitGame() error {
	// 订阅消息通道
	g.playerEventCh = g.Queue.Subscribe("player")
	g.controlEventCh = g.Queue.Subscribe("control")

	return nil
}

// Listen 启动游戏主循环，处理事件并转发
func (g *Game) Listen() {
	// 启动游戏上下文
	g.ctx, g.cancel = context.WithCancel(context.Background())

	go func() {
		defer func() {
			if g.cancel != nil {
				g.cancel()
			}
		}()

		// 消息处理主循环
		for {
			select {
			case <-g.ctx.Done():
				slog.Info("Game context cancelled", "gameId", g.Id)
				return

			// 处理控制事件
			case e := <-g.controlEventCh:
				g.handleControlEvent(e)

			// 处理玩家事件（主要处理逻辑）
			case e := <-g.playerEventCh:
				g.handlePlayerEvent(e)

			// 定期维护
			case <-time.After(100 * time.Millisecond):
				g.performPeriodicMaintenance()
			}
		}
	}()
}

// handlePlayerEvent 处理玩家事件并转发相应事件
func (g *Game) handlePlayerEvent(e queue.Event) {
	var err error

	switch event := e.(type) {
	case Join:
		err = g.handleJoin(event)
	case Disconnect:
		err = g.handleLeave(event)
	case Move:
		err = g.handleMove(event)
	case ForceStart:
		err = g.handleForceStart(event)
	case Surrender:
		err = g.handleSurrender(event)
	default:
		slog.Warn("unknown player event", "type", fmt.Sprintf("%T", event))
		return
	}

	// 处理错误
	if err != nil {
		// TODO
	}
}

// handleControlEvent 处理控制事件
func (g *Game) handleControlEvent(e queue.Event) {
	switch event := e.(type) {
	case StartGame:
		g.handleStartGame(event)
	case EndGame:
		g.handleEndGame(event)
	case string:
		switch event {
		case "force_stop":
			g.handleForceStop()
		default:
			slog.Warn("unknown control command", "command", event)
		}
	default:
		slog.Warn("unknown control event", "type", fmt.Sprintf("%T", event))
	}
}

// === 玩家事件处理方法 ===

func (g *Game) handleJoin(event Join) error {
	player := Player{
		Id:     event.Player,
		Name:   event.Player, // TODO: 从用户系统获取真实用户名
		Status: PlayerStatusWaiting,
		Moves:  0,
	}

	err := g.Join(player)
	if err != nil {
		return err
	}

	// 转发事件：发布玩家加入广播事件
	g.publishPlayerUpdateEvent("joined", "Player joined the game")
	g.publishRoomInfoUpdate()

	slog.Info("player joined", "gameId", g.Id, "player", event.Player)
	return nil
}

func (g *Game) handleLeave(event Disconnect) error {
	err := g.Leave(event.Player)
	if err != nil {
		return err
	}

	// 转发事件：发布玩家离开广播事件
	g.publishPlayerUpdateEvent("left", "Player left the game")
	g.publishRoomInfoUpdate()

	slog.Info("player left", "gameId", g.Id, "player", event.Player)
	return nil
}

func (g *Game) handleMove(event Move) error {
	err := g.Move(event.Player, event)
	if err != nil {
		return err
	}

	// 转发事件：发布地图更新广播事件
	g.publishMapUpdateEvent()

	slog.Info("player moved", "gameId", g.Id, "player", event.Player, "from", event.Pos, "towards", event.Towards)
	return nil
}

func (g *Game) handleForceStart(event ForceStart) error {
	err := g.ForceStart(event.Player, event.IsVote)
	if err != nil {
		return err
	}

	// 转发事件：发布强制开始投票广播事件
	if event.IsVote {
		g.publishPlayerUpdateEvent("force_start_vote", "Player voted to force start")
	} else {
		g.publishPlayerUpdateEvent("vote_cancelled", "Player cancelled force start vote")
	}
	g.publishRoomInfoUpdate()

	slog.Info("player force start vote", "gameId", g.Id, "player", event.Player, "vote", event.IsVote)
	return nil
}

func (g *Game) handleSurrender(event Surrender) error {
	err := g.Surrender(event.Player)
	if err != nil {
		return err
	}

	// 转发事件：发布玩家投降广播事件
	g.publishPlayerUpdateEvent("surrendered", "Player surrendered")
	g.publishRoomInfoUpdate()

	slog.Info("player surrendered", "gameId", g.Id, "player", event.Player)
	return nil
}

// === 控制事件处理方法 ===

func (g *Game) handleStartGame(event StartGame) {
	// 游戏自动开始，发布游戏开始广播事件
	g.publishGameStartEvent()
	slog.Info("game started", "gameId", g.Id)
}

func (g *Game) handleEndGame(event EndGame) {
	// 游戏结束，发布游戏结束广播事件
	g.publishGameEndEvent(event.Winner)
	slog.Info("game ended", "gameId", g.Id, "winner", event.Winner)
}

func (g *Game) handleForceStop() {
	err := g.Stop()
	if err != nil {
		slog.Error("failed to force stop game", "gameId", g.Id, "error", err)
	} else {
		g.publishGameEndEvent("")
		slog.Info("game force stopped", "gameId", g.Id)
	}
}

// === 事件发布方法 ===

func (g *Game) publishPlayerUpdateEvent(status, message string) {
	// 发布游戏状态更新事件，包含玩家状态变化
	event := GameStatusUpdate{
		Status:     g.Status(),
		Players:    g.Players(),
		TurnNumber: g.TurnNumber(),
	}

	g.Queue.Publish("broadcast", event)
}

func (g *Game) publishRoomInfoUpdate() {
	// 房间信息更新也使用游戏状态更新事件
	event := GameStatusUpdate{
		Status:     g.Status(),
		Players:    g.Players(),
		TurnNumber: g.TurnNumber(),
	}

	g.Queue.Publish("broadcast", event)
}

func (g *Game) publishGameStartEvent() {
	// 游戏开始时发布游戏状态更新
	event := GameStatusUpdate{
		Status:     g.Status(),
		Players:    g.Players(),
		TurnNumber: g.TurnNumber(),
	}

	g.Queue.Publish("broadcast", event)
}

func (g *Game) publishGameEndEvent(winner string) {
	// 游戏结束时发布游戏状态更新
	event := GameStatusUpdate{
		Status:     g.Status(),
		Players:    g.Players(),
		TurnNumber: g.TurnNumber(),
	}

	g.Queue.Publish("broadcast", event)
}

func (g *Game) publishMapUpdateEvent() {
	// 使用已定义的地图更新事件
	event := MapUpdate{
		Map: g.Map(),
	}

	g.Queue.Publish("broadcast", event)
}

func (g *Game) publishErrorEvent(err error) {
	// 错误事件可以通过日志处理，暂时不发布广播事件
	slog.Error("game error for player", "error", err)
}

// === 辅助方法 ===

func (g *Game) performPeriodicMaintenance() {
	// 定期检查和维护
	// 例如：检查玩家连接状态、清理过期数据等
}

// StopGame 停止游戏并清理资源
func (g *Game) StopGame() {
	if g.cancel != nil {
		g.cancel()
		g.cancel = nil
	}

	// 调用BaseCore的Stop方法
	if err := g.Stop(); err != nil {
		slog.Error("failed to stop BaseCore", "gameId", g.Id, "error", err)
	}

	// 取消订阅
	if g.playerEventCh != nil {
		g.Queue.Unsubscribe("player", g.playerEventCh)
	}
	if g.controlEventCh != nil {
		g.Queue.Unsubscribe("control", g.controlEventCh)
	}

	slog.Info("Game stopped and resources cleaned", "gameId", g.Id)
}
