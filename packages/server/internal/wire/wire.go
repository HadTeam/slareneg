//go:generate wire
//go:build wireinject
// +build wireinject

// Basic wire injector, used for in-memory database

package wire

import (
	"server/internal/auth"
	"server/internal/cache"
	"server/internal/config"
	gamemap "server/internal/game/map"
	"server/internal/lobby"
	"server/internal/queue"
	"server/internal/websocket"

	"github.com/google/wire"
)

type Application struct {
	Config      *config.Config
	AuthService *auth.AuthService
	Lobby       *lobby.Lobby
	WSServer    *websocket.WebSocketServer
	Cache       *cache.CacheService
	MapManager  gamemap.MapManager
}

func InitializeApplication(cfg *config.Config) (*Application, error) {
	wire.Build(
		queue.NewInMemoryQueue,
		wire.Bind(new(queue.Queue), new(*queue.InMemoryQueue)),

		auth.NewInMemoryUserRepository,
		wire.Bind(new(auth.UserRepository), new(*auth.InMemoryUserRepository)),

		provideJWTTokenService,
		wire.Bind(new(auth.TokenService), new(*auth.JWTTokenService)),

		auth.NewArgon2PasswordService,
		wire.Bind(new(auth.PasswordService), new(*auth.Argon2PasswordService)),

		auth.NewAuthService,

		provideCacheService,

		provideMapManager,
		wire.Bind(new(gamemap.MapManager), new(*gamemap.DefaultMapManager)),

		lobby.NewLobby,

		websocket.NewWebSocketServer,

		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}

func provideJWTTokenService(cfg *config.Config) *auth.JWTTokenService {
	return auth.NewJWTTokenService(cfg.Auth.JWTSecret)
}

func provideCacheService(cfg *config.Config) *cache.CacheService {
	inMemoryCache := cache.NewInMemoryCache(cfg.Cache.CleanupInterval)
	return cache.NewCacheService(inMemoryCache)
}

func provideMapManager() *gamemap.DefaultMapManager {
	return gamemap.NewMapManager()
}
