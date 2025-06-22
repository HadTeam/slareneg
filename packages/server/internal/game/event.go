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

type CommandEvent struct {
	PlayerId string
}

type ControlEvent struct{}

type BroadcastEvent struct{}

type PlayerEvent struct{}

type JoinCommand struct {
	CommandEvent
	PlayerName string
}

type LeaveCommand struct {
	CommandEvent
}

type MoveCommand struct {
	CommandEvent
	From      gamemap.Pos
	Direction MoveTowards
	Troops    block.Num
}

type ForceStartCommand struct {
	CommandEvent
	IsVote bool
}

type SurrenderCommand struct {
	CommandEvent
}

type StartGameControl struct {
	ControlEvent
}

type StopGameControl struct {
	ControlEvent
}

type TurnAdvanceControl struct {
	ControlEvent
	TurnNumber uint16
}

type PlayerJoinedEvent struct {
	BroadcastEvent
	PlayerId   string
	PlayerName string
	GameStatus Status
	Players    []Player
}

type PlayerLeftEvent struct {
	BroadcastEvent
	PlayerId   string
	GameStatus Status
	Players    []Player
}

type MapUpdateEvent struct {
	BroadcastEvent
	Map        gamemap.Map
	TurnNumber uint16
}

type GameStatusUpdateEvent struct {
	BroadcastEvent
	Status     Status
	Players    []Player
	TurnNumber uint16
}

type ForceStartVoteEvent struct {
	BroadcastEvent
	PlayerId   string
	IsVote     bool
	GameStatus Status
	Players    []Player
}

type PlayerSurrenderedEvent struct {
	BroadcastEvent
	PlayerId   string
	GameStatus Status
	Players    []Player
}

type GameStartedEvent struct {
	BroadcastEvent
	GameStatus Status
	Players    []Player
	TurnNumber uint16
}

type GameEndedEvent struct {
	BroadcastEvent
	Winner     string
	GameStatus Status
	Players    []Player
}

type TurnStartedEvent struct {
	BroadcastEvent
	TurnNumber uint16
	Players    []Player
}

type PlayerMovedEvent struct {
	BroadcastEvent
	PlayerId  string
	Move      Move
	MovesLeft uint16
}

type PlayerErrorEvent struct {
	PlayerEvent
	PlayerId string
	Error    string
}

type Move struct {
	Pos     gamemap.Pos
	Towards MoveTowards
	Num     block.Num
}
