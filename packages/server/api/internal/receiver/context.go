package receiver

import (
	"context"
	"server/game"
	"server/game/user"
)

type Context struct {
	context.Context
	Game    *game.Game
	User    user.User
	Command chan string
	Message chan string
}

// ?
