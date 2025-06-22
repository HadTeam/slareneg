package game

import (
	gamemap "server/internal/game/map"
)

type Status string

const (
	StatusWaiting    Status = "waiting"
	StatusInProgress Status = "in_progress"
	StatusFinished   Status = "finished"
)

type PlayerStatus string

const (
	PlayerStatusWaiting           PlayerStatus = "waiting"
	PlayerStatusRequestForceStart PlayerStatus = "request_force_start"

	PlayerStatusInGame       PlayerStatus = "in_game"
	PlayerStatusDisconnected PlayerStatus = "disconnected"

	PlayerStatusSurrendered PlayerStatus = "surrendered"
	PlayerStatusLost        PlayerStatus = "lost"
	PlayerStatusWinner      PlayerStatus = "winner"

	PlayerStatusSpectator PlayerStatus = "spectator"

	PlayerStatusFinished PlayerStatus = "finished"
)

type FinishReason string

const (
	FinishReasonNone         FinishReason = ""
	FinishReasonVictory      FinishReason = "victory"
	FinishReasonDefeated     FinishReason = "defeated"
	FinishReasonSurrendered  FinishReason = "surrendered"
	FinishReasonDisconnected FinishReason = "disconnected"
	FinishReasonError        FinishReason = "error"
)

type PlayerConnectionInfo struct {
	IsConnected      bool
	DisconnectedAt   int64
	ReconnectTimeout int64
}

type Player struct {
	Id     string
	Name   string
	Moves  uint16
	Status PlayerStatus

	Connection PlayerConnectionInfo

	FinishReason     FinishReason
	IsForceStartVote bool
}

func (p *Player) CanOperate() bool {
	return p.Status == PlayerStatusInGame
}

func (p *Player) CanReceiveUpdates() bool {
	return p.Status == PlayerStatusInGame ||
		p.Status == PlayerStatusSurrendered ||
		p.Status == PlayerStatusSpectator ||
		p.Status == PlayerStatusDisconnected
}

func (p *Player) IsActive() bool {
	return p.Status == PlayerStatusInGame || p.Status == PlayerStatusDisconnected
}

func (p *Player) IsFinished() bool {
	return p.Status == PlayerStatusSurrendered ||
		p.Status == PlayerStatusLost ||
		p.Status == PlayerStatusWinner ||
		p.Status == PlayerStatusSpectator ||
		p.Status == PlayerStatusFinished
}

type Core interface {
	Status() Status
	Players() []Player
	TurnNumber() uint16

	IsGameReady() bool

	GetActivePlayerCount() int

	Join(player Player) error
	Leave(playerId string) error
	GetPlayer(playerId string) (*Player, error)

	Map() gamemap.Map

	Start() error
	Stop() error

	NextTurn(turnNumber uint16) error
	Move(playerId string, move Move) error
	ForceStart(playerId string, isVote bool) error
	Surrender(playerId string) error

	PlayerConnect(playerId string) error
	PlayerDisconnect(playerId string) error
	PlayerReconnect(playerId string) error
	CheckDisconnectedPlayers(currentTimeMs int64) error
}
