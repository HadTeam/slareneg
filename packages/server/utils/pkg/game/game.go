package game

import (
	"server/utils/pkg/map"
)

type GameStatus uint8
type GameId uint16

const (
	GameStatusWaiting GameStatus = iota + 1
	GameStatusRunning
	GameStatusEnd
)

type GameScore struct {
	Num   uint16
	Place uint16
}

type Game struct {
	Map        *_map.Map
	Mode       GameMode
	Id         GameId
	UserList   []User
	CreateTime int64
	Status     GameStatus
	RoundNum   uint16
	Winner     uint8 // TeamId
}
