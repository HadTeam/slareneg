package receiver

import (
	"context"
	"server/game_logic"
	"server/game_logic/game_def"
)

type Context struct {
	context.Context
	Game    *game_logic.Game
	User    _type.User
	Command chan string
	Message chan string
}

// ?
