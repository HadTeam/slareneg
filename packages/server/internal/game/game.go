package game

import (
	"context"
	"fmt"
	"log/slog"
	gamemap "server/internal/game/map"
	"server/internal/queue"
)

// Game 事件转发层 - 负责事件的订阅、解析和转发
type Game struct {
	gameId string
	core   *BaseCore
	queue  queue.Queue

	// 上下文管理
	ctx    context.Context
	cancel context.CancelFunc

	// 消息通道
	commandCh <-chan queue.Event
	controlCh <-chan queue.Event
}

// NewGame 创建新的游戏实例
func NewGame(gameId string, q queue.Queue, mode GameMode, mapManager gamemap.MapManager) *Game {
	core := NewBaseCore(gameId, mode, mapManager)

	game := &Game{
		gameId: gameId,
		core:   core,
		queue:  q,
	}

	// 设置BaseCore的事件回调
	core.SetEventHandlers(
		game.forwardBroadcastEvent,
		game.forwardControlEvent,
	)

	return game
}

// Start 启动游戏事件处理循环
func (g *Game) Start() error {
	// 订阅消息通道
	g.commandCh = g.queue.Subscribe(fmt.Sprintf("%s/commands", g.gameId))
	g.controlCh = g.queue.Subscribe(fmt.Sprintf("%s/control", g.gameId))

	// 启动游戏上下文
	g.ctx, g.cancel = context.WithCancel(context.Background())

	// 启动事件处理循环
	go g.eventLoop()

	slog.Info("game event handler started", "gameId", g.gameId)
	return nil
}

// Stop 停止游戏事件处理
func (g *Game) Stop() error {
	if g.cancel != nil {
		g.cancel()
	}

	// 停止游戏核心
	if err := g.core.Stop(); err != nil {
		slog.Error("failed to stop game core", "error", err, "gameId", g.gameId)
	}

	slog.Info("game event handler stopped", "gameId", g.gameId)
	return nil
}

// Core 获取游戏核心（用于直接访问游戏状态）
func (g *Game) Core() Core {
	return g.core
}

// eventLoop 事件处理主循环
func (g *Game) eventLoop() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("game event loop panic", "error", r, "gameId", g.gameId)
		}
		if g.cancel != nil {
			g.cancel()
		}
	}()

	for {
		select {
		case <-g.ctx.Done():
			slog.Info("game event loop stopped", "gameId", g.gameId)
			return

		case event := <-g.controlCh:
			g.handleControlEvent(event)

		case event := <-g.commandCh:
			g.handleCommandEvent(event)

		}
	}
}

// handleCommandEvent 处理玩家指令事件并调用BaseCore相应方法
func (g *Game) handleCommandEvent(event queue.Event) {
	var err error

	switch cmd := event.(type) {
	case JoinCommand:
		err = g.handleJoinCommand(cmd)
	case LeaveCommand:
		err = g.handleLeaveCommand(cmd)
	case MoveCommand:
		err = g.handleMoveCommand(cmd)
	case ForceStartCommand:
		err = g.handleForceStartCommand(cmd)
	case SurrenderCommand:
		err = g.handleSurrenderCommand(cmd)
	default:
		slog.Warn("unknown command event", "type", fmt.Sprintf("%T", event), "gameId", g.gameId)
		return
	}

	// 处理错误
	if err != nil {
		playerId := g.getPlayerIdFromCommand(event)
		g.publishPlayerError(playerId, err)
	}
}

// handleControlEvent 处理控制事件
func (g *Game) handleControlEvent(event queue.Event) {
	switch e := event.(type) {
	case StartGameControl:
		if err := g.core.Start(); err != nil {
			slog.Error("failed to start game core", "error", err, "gameId", g.gameId)
		}
	case StopGameControl:
		if err := g.core.Stop(); err != nil {
			slog.Error("failed to stop game core", "error", err, "gameId", g.gameId)
		}
	case TurnAdvanceControl:
		if err := g.core.NextTurn(e.TurnNumber); err != nil {
			slog.Error("failed to advance turn", "error", err, "gameId", g.gameId)
		}
	default:
		slog.Warn("unknown control event", "type", fmt.Sprintf("%T", event), "gameId", g.gameId)
	}
}

// =============================================================================
// 指令处理方法
// =============================================================================

// handleJoinCommand 处理加入游戏指令
func (g *Game) handleJoinCommand(cmd JoinCommand) error {
	player := Player{
		Id:   cmd.PlayerId,
		Name: cmd.PlayerName,
	}

	if err := g.core.Join(player); err != nil {
		return err
	}

	// 发布玩家加入事件
	g.forwardBroadcastEvent(PlayerJoinedEvent{
		BroadcastEvent: BroadcastEvent{},
		PlayerId:       cmd.PlayerId,
		PlayerName:     cmd.PlayerName,
		GameStatus:     g.core.Status(),
		Players:        g.core.Players(),
	})

	return nil
}

// handleLeaveCommand 处理离开游戏指令
func (g *Game) handleLeaveCommand(cmd LeaveCommand) error {
	if err := g.core.Leave(cmd.PlayerId); err != nil {
		return err
	}

	// 发布玩家离开事件
	g.forwardBroadcastEvent(PlayerLeftEvent{
		BroadcastEvent: BroadcastEvent{},
		PlayerId:       cmd.PlayerId,
		GameStatus:     g.core.Status(),
		Players:        g.core.Players(),
	})

	return nil
}

// handleMoveCommand 处理移动指令
func (g *Game) handleMoveCommand(cmd MoveCommand) error {
	move := Move{
		Pos:     cmd.From,
		Towards: cmd.Direction,
		Num:     cmd.Troops,
	}

	if err := g.core.Move(cmd.PlayerId, move); err != nil {
		return err
	}

	// 发布地图更新事件
	g.forwardBroadcastEvent(MapUpdateEvent{
		BroadcastEvent: BroadcastEvent{},
		Map:            g.core.Map(),
		TurnNumber:     g.core.TurnNumber(),
	})

	return nil
}

// handleForceStartCommand 处理强制开始指令
func (g *Game) handleForceStartCommand(cmd ForceStartCommand) error {
	if err := g.core.ForceStart(cmd.PlayerId, cmd.IsVote); err != nil {
		return err
	}

	// 发布强制开始投票事件
	g.forwardBroadcastEvent(ForceStartVoteEvent{
		BroadcastEvent: BroadcastEvent{},
		PlayerId:       cmd.PlayerId,
		IsVote:         cmd.IsVote,
		GameStatus:     g.core.Status(),
		Players:        g.core.Players(),
	})

	return nil
}

// handleSurrenderCommand 处理投降指令
func (g *Game) handleSurrenderCommand(cmd SurrenderCommand) error {
	if err := g.core.Surrender(cmd.PlayerId); err != nil {
		return err
	}

	// 发布玩家投降事件
	g.forwardBroadcastEvent(PlayerSurrenderedEvent{
		BroadcastEvent: BroadcastEvent{},
		PlayerId:       cmd.PlayerId,
		GameStatus:     g.core.Status(),
		Players:        g.core.Players(),
	})

	return nil
}

// =============================================================================
// 事件转发方法
// =============================================================================

// forwardBroadcastEvent 转发广播事件
func (g *Game) forwardBroadcastEvent(event queue.Event) {
	g.queue.Publish(fmt.Sprintf("%s/broadcast", g.gameId), event)
}

// forwardControlEvent 转发控制事件
func (g *Game) forwardControlEvent(event queue.Event) {
	g.queue.Publish(fmt.Sprintf("%s/control", g.gameId), event)
}

// publishPlayerError 发布玩家错误消息
func (g *Game) publishPlayerError(playerId string, err error) {
	if playerId == "" {
		slog.Error("cannot publish error for empty player id", "error", err.Error(), "gameId", g.gameId)
		return
	}

	errorEvent := PlayerErrorEvent{
		PlayerEvent: PlayerEvent{},
		PlayerId:    playerId,
		Error:       err.Error(),
	}
	g.queue.Publish(fmt.Sprintf("%s/player/%s", g.gameId, playerId), errorEvent)
}

// getPlayerIdFromCommand 从指令事件中提取玩家ID
func (g *Game) getPlayerIdFromCommand(event queue.Event) string {
	switch e := event.(type) {
	case JoinCommand:
		return e.PlayerId
	case LeaveCommand:
		return e.PlayerId
	case MoveCommand:
		return e.PlayerId
	case ForceStartCommand:
		return e.PlayerId
	case SurrenderCommand:
		return e.PlayerId
	default:
		return ""
	}
}
