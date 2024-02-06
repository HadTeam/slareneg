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
	Mode       game_def.Mode
	Id         Id
	UserList   []game_def.User
	CreateTime int64
	Status     Status
	RoundNum   uint16
	Winner     uint8 // TeamId
}
