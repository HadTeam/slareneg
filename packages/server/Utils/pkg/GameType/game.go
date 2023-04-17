package GameType

import (
	"server/Utils/pkg/MapType"
)

type GameStatus uint8
type GameId uint16

const (
	GameStatusWaiting GameStatus = iota + 1
	GameStatusRunning
	GameStatusEnd
)

type GameScore struct {
	Num   uint8
	Place uint8
}

type Game struct {
	Map        *MapType.Map
	Mode       GameMode
	Id         GameId
	UserList   []User
	CreateTime int64
	Status     GameStatus
	RoundNum   uint8
	Winner     uint8
}
