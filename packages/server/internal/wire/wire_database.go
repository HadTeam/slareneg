//go:build wireinject && database
// +build wireinject,database

// Use go build -tags database

package wire

import (
	"server/internal/auth"
	"server/internal/config"
	"server/internal/lobby"
	"server/internal/queue"
	"server/internal/websocket"

	"github.com/google/wire"
)

func InitializeApplicationWithDatabase(cfg *config.Config) (*Application, error) {
	wire.Build(
		queue.NewInMemoryQueue,
		wire.Bind(new(queue.Queue), new(*queue.InMemoryQueue)),

		auth.NewDatabaseUserRepository,
		wire.Bind(new(auth.UserRepository), new(*auth.DatabaseUserRepository)),

		provideJWTTokenService,
		wire.Bind(new(auth.TokenService), new(*auth.JWTTokenService)),

		auth.NewArgon2PasswordService,
		wire.Bind(new(auth.PasswordService), new(*auth.Argon2PasswordService)),

		auth.NewAuthService,

		provideCacheService,

		lobby.NewLobby,

		websocket.NewWebSocketServer,

		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
