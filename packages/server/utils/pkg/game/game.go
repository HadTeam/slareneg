package game

import (
	"server/utils/pkg/map"
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
	Mode       Mode
	Id         Id
	UserList   []User
	CreateTime int64
	Status     Status
	RoundNum   uint16
	Winner     uint8 // TeamId
}
