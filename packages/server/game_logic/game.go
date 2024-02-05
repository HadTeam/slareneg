package game_logic

import (
	"server/game_logic/game_def"
	"server/game_logic/map"
)

type Status uint8
type Id uint16

const (
	StatusWaiting Status = iota + 1
	StatusRunning
	StatusEnd
)

type Game struct {
	Map        *_map.Map
	Mode       _type.Mode
	Id         Id
	UserList   []_type.User
	CreateTime int64
	Status     Status
	RoundNum   uint16
	Winner     uint8 // TeamId
}
