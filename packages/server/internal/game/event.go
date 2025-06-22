package game

import (
	"server/internal/game/block"
	gamemap "server/internal/game/map"
)

type MoveTowards string

type moveOffset struct {
	X int16
	Y int16
}

const (
	MoveTowardsLeft  MoveTowards = "left"
	MoveTowardsRight MoveTowards = "right"
	MoveTowardsUp    MoveTowards = "up"
	MoveTowardsDown  MoveTowards = "down"
)

func getMoveOffset(t MoveTowards) moveOffset {
	switch t {
	case MoveTowardsLeft:
		return moveOffset{-1, 0}
	case MoveTowardsRight:
		return moveOffset{1, 0}
	case MoveTowardsUp:
		return moveOffset{0, -1}
	case MoveTowardsDown:
		return moveOffset{0, 1}
	default:
		return moveOffset{0, 0} // No movement
	}
}

// =============================================================================
// 流类型定义
// =============================================================================

// 玩家消息流事件 - 入口来自 ws，出口在 core
type PlayerEvent struct {
	Player string
}

// 控制流事件 - 入口在外部和 core，出口在 core
type ControlEvent struct{}

// 广播流事件 - 入口在 game 和 core，出口在 ws
type BroadcastEvent struct{}

// =============================================================================
// 玩家消息流事件类型 (${gameId}/player)
// 这些事件来自 WebSocket 和测试，由玩家触发。可能同时出现在多个流中
// =============================================================================

type Move struct {
	PlayerEvent
	Pos     gamemap.Pos
	Towards MoveTowards
	Num     block.Num
}

type ForceStart struct {
	PlayerEvent
	IsVote bool
}
type Surrender struct{ PlayerEvent }
type Join struct{ PlayerEvent }

type Disconnect struct{ PlayerEvent }
type Reconnect struct{ PlayerEvent }

// =============================================================================
// 控制流事件类型 (${gameId}/control)
// 这些事件来自游戏底层模块，通过上下文调用，用于内部状态管理
// =============================================================================

type AdvanceTurn struct{ ControlEvent }
type StartGame struct{ ControlEvent }
type EndGame struct {
	ControlEvent
	Winner string // 玩家 ID / 团队 ID
}

// =============================================================================
// 广播流事件类型 (${gameId}/broadcast)
// 这些事件发送给玩家，在交付前可能还需要处理
// =============================================================================

type GameStatusUpdate struct {
	BroadcastEvent
	Status     Status
	Players    []Player
	TurnNumber uint16
}

type MapUpdate struct {
	BroadcastEvent
	Map gamemap.Map // 完整地图信息，包括块和元数据，增量在
}
