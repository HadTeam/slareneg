package receiver

import (
	"context"
	"server/utils/pkg/game"
)

type Context struct {
	context.Context
	Game    *game.Game
	User    game.User
	Command chan string
	Message chan string
}

// ?
