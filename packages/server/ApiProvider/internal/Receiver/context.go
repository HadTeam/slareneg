package Receiver

import (
	"context"
	"server/Utils/pkg/GameType"
)

type Context struct {
	context.Context
	Game    *GameType.Game
	User    GameType.User
	Command chan string
	Message chan string
}

// ?
