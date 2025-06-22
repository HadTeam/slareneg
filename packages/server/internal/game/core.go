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
	PlayerStatusInGame            PlayerStatus = "in_game"
	PlayerStatusSurrendered       PlayerStatus = "surrendered"
	PlayerStatusFinished          PlayerStatus = "finished"
	PlayerStatusDisconnected      PlayerStatus = "disconnected"
	PlayerStatusError             PlayerStatus = "error"
)

type Player struct {
	Id    string
	Name  string
	Moves uint16 // remaining moves

	Status PlayerStatus
}

// Core 游戏核心接口 - 纯粹的游戏逻辑处理，无外部依赖
type Core interface {
	// 状态查询
	Status() Status
	Players() []Player
	TurnNumber() uint16
	
	IsGameReady() bool

	GetActivePlayerCount() int

	// 玩家管理
	Join(player Player) error
	Leave(playerId string) error
	GetPlayer(playerId string) (*Player, error)

	// 地图和游戏模式
	Map() gamemap.Map

	// 游戏控制
	Start() error
	Stop() error

	// 游戏操作
	NextTurn(turnNumber uint16) error
	Move(playerId string, move Move) error
	ForceStart(playerId string, isVote bool) error
	Surrender(playerId string) error
}
