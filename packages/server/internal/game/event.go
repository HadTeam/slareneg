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
// 基础事件类型定义
// =============================================================================

// CommandEvent 指令事件基类 - 来自WebSocket的玩家指令
type CommandEvent struct {
	PlayerId string
}

// ControlEvent 控制事件基类 - 用于内部状态管理
type ControlEvent struct{}

// BroadcastEvent 广播事件基类 - 发送给所有玩家
type BroadcastEvent struct{}

// PlayerEvent 玩家事件基类 - 发送给特定玩家
type PlayerEvent struct{}

// =============================================================================
// 指令事件类型 (${gameId}/commands)
// 这些事件来自 WebSocket，由玩家触发，经Game转发层处理后调用GameCore
// =============================================================================

// JoinCommand 加入游戏指令
type JoinCommand struct {
	CommandEvent
	PlayerName string
}

// LeaveCommand 离开游戏指令
type LeaveCommand struct {
	CommandEvent
}

// MoveCommand 移动指令
type MoveCommand struct {
	CommandEvent
	From      gamemap.Pos
	Direction MoveTowards
	Troops    block.Num
}

// ForceStartCommand 强制开始投票指令
type ForceStartCommand struct {
	CommandEvent
	IsVote bool
}

// SurrenderCommand 投降指令
type SurrenderCommand struct {
	CommandEvent
}

// =============================================================================
// 控制事件类型 (${gameId}/control)
// 这些事件用于内部状态管理和游戏生命周期控制
// =============================================================================

// StartGameControl 启动游戏控制事件
type StartGameControl struct {
	ControlEvent
}

// StopGameControl 停止游戏控制事件
type StopGameControl struct {
	ControlEvent
}

// TurnAdvanceControl 回合推进控制事件
type TurnAdvanceControl struct {
	ControlEvent
	TurnNumber uint16
}

// =============================================================================
// 广播事件类型 (${gameId}/broadcast)
// 这些事件发送给房间内的所有玩家
// =============================================================================

// PlayerJoinedEvent 玩家加入事件
type PlayerJoinedEvent struct {
	BroadcastEvent
	PlayerId   string
	PlayerName string
	GameStatus Status
	Players    []Player
}

// PlayerLeftEvent 玩家离开事件
type PlayerLeftEvent struct {
	BroadcastEvent
	PlayerId   string
	GameStatus Status
	Players    []Player
}

// MapUpdateEvent 地图更新事件
type MapUpdateEvent struct {
	BroadcastEvent
	Map        gamemap.Map
	TurnNumber uint16
}

// GameStatusUpdateEvent 游戏状态更新事件
type GameStatusUpdateEvent struct {
	BroadcastEvent
	Status     Status
	Players    []Player
	TurnNumber uint16
}

// ForceStartVoteEvent 强制开始投票事件
type ForceStartVoteEvent struct {
	BroadcastEvent
	PlayerId   string
	IsVote     bool
	GameStatus Status
	Players    []Player
}

// PlayerSurrenderedEvent 玩家投降事件
type PlayerSurrenderedEvent struct {
	BroadcastEvent
	PlayerId   string
	GameStatus Status
	Players    []Player
}

// GameStartedEvent 游戏开始事件
type GameStartedEvent struct {
	BroadcastEvent
	GameStatus Status
	Players    []Player
	TurnNumber uint16
}

// GameEndedEvent 游戏结束事件
type GameEndedEvent struct {
	BroadcastEvent
	Winner     string
	GameStatus Status
	Players    []Player
}

// =============================================================================
// 玩家特定事件类型 (${gameId}/player/${playerId})
// 这些事件发送给特定玩家
// =============================================================================

// PlayerErrorEvent 玩家错误事件
type PlayerErrorEvent struct {
	PlayerEvent
	PlayerId string
	Error    string
}

// =============================================================================
// GameCore内部使用的数据结构（保持兼容性）
// =============================================================================

// Move GameCore内部使用的移动数据结构
type Move struct {
	Pos     gamemap.Pos
	Towards MoveTowards
	Num     block.Num
}
