package GameType

import (
	"server/MapHandler/pkg/MapType"
)

type GameStatus uint8
type GameId uint16

const (
	GameStatusWaiting GameStatus = iota + 1
	GameStatusRunning
	GameStatusEnd
)

type Game struct {
	Map        *MapType.Map
	UserList   []User
	CreateTime int64
	Status     GameStatus
	RoundNum   *uint8
}
