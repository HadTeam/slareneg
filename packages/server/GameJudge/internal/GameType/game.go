package GameType

import (
	"server/MapHandler/pkg/MapType"
)

type GameStatus uint8
type GameId uint16

const (
	GameStatusWaiting GameStatus = 1
	GameStatusRunning GameStatus = 2
	GameStatusEnd     GameStatus = 3
)

type Game struct {
	Map        *MapType.Map
	UserList   []User
	CreateTime int64
	Status     GameStatus
	RoundNum   *uint8
}
