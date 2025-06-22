package lobby

import (
	"fmt"
	"log/slog"
	"server/internal/game"
	"server/internal/queue"
	"sync"
)

type Lobby struct {
	games   map[string]*game.Game
	gamesMu sync.RWMutex
	queue   queue.Queue
}

type LobbyCommand struct {
	Type    string      `json:"type"`
	GameId  string      `json:"gameId"`
	Payload interface{} `json:"payload"`
}

type CreateGamePayload struct {
	GameMode game.GameMode `json:"gameMode"`
}

func NewLobby(q queue.Queue) *Lobby {
	return &Lobby{
		games: make(map[string]*game.Game),
		queue: q,
	}
}

func (l *Lobby) Start() error {
	commandChan := l.queue.Subscribe("lobby/commands")

	go func() {
		for cmd := range commandChan {
			if lobbyCmd, ok := cmd.(LobbyCommand); ok {
				if err := l.handleCommand(lobbyCmd); err != nil {
					slog.Error("failed to handle lobby command", "error", err, "type", lobbyCmd.Type)
				}
			}
		}
	}()

	slog.Info("lobby service started")
	return nil
}

func (l *Lobby) Stop() error {
	l.gamesMu.Lock()
	defer l.gamesMu.Unlock()

	for gameId, gameInstance := range l.games {
		if err := gameInstance.Stop(); err != nil {
			slog.Error("failed to stop game", "error", err, "gameId", gameId)
		}
	}

	slog.Info("lobby service stopped")
	return nil
}

func (l *Lobby) handleCommand(cmd LobbyCommand) error {
	switch cmd.Type {
	case "createGame":
		return l.handleCreateGame(cmd)
	case "getGameInfo":
		return l.handleGetGameInfo(cmd)
	default:
		return fmt.Errorf("unknown lobby command type: %s", cmd.Type)
	}
}

func (l *Lobby) handleCreateGame(cmd LobbyCommand) error {
	payload, ok := cmd.Payload.(CreateGamePayload)
	if !ok {
		payload = CreateGamePayload{GameMode: game.Classic1v1} // 默认游戏模式
	}

	gameInstance := l.getOrCreateGame(cmd.GameId, payload.GameMode)

	l.queue.Publish("lobby/events", map[string]interface{}{
		"type":   "gameCreated",
		"gameId": cmd.GameId,
		"game":   gameInstance,
	})

	return nil
}

func (l *Lobby) handleGetGameInfo(cmd LobbyCommand) error {
	l.gamesMu.RLock()
	gameInstance, exists := l.games[cmd.GameId]
	l.gamesMu.RUnlock()

	l.queue.Publish("lobby/events", map[string]interface{}{
		"type":   "gameInfo",
		"gameId": cmd.GameId,
		"exists": exists,
		"game":   gameInstance,
	})

	return nil
}

func (l *Lobby) getOrCreateGame(gameId string, gameMode game.GameMode) *game.Game {
	l.gamesMu.Lock()
	defer l.gamesMu.Unlock()

	if existingGame, exists := l.games[gameId]; exists {
		return existingGame
	}

	newGame := game.NewGame(gameId, l.queue, gameMode)
	l.games[gameId] = newGame

	go func() {
		if err := newGame.Start(); err != nil {
			slog.Error("failed to start game", "error", err, "gameId", gameId)
		}
	}()

	slog.Info("created new game", "gameId", gameId, "gameMode", gameMode)
	return newGame
}

func (l *Lobby) GetGameList() map[string]*game.Game {
	l.gamesMu.RLock()
	defer l.gamesMu.RUnlock()

	result := make(map[string]*game.Game)
	for id, gameInstance := range l.games {
		result[id] = gameInstance
	}
	return result
}

func (l *Lobby) RemoveGame(gameId string) {
	l.gamesMu.Lock()
	defer l.gamesMu.Unlock()

	delete(l.games, gameId)
	slog.Info("removed game", "gameId", gameId)
}
