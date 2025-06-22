package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"server/internal/game"
	"server/internal/game/block"
	gamemap "server/internal/game/map"
	"server/internal/queue"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

type WebSocketServer struct {
	queue queue.Queue
}

type ClientMessage struct {
	Type    string          `json:"type"`
	GameId  string          `json:"gameId"`
	Payload json.RawMessage `json:"payload"`
}

type JoinPayload struct {
	PlayerId   string `json:"playerId"`
	PlayerName string `json:"playerName"`
}

type MovePayload struct {
	PlayerId  string      `json:"playerId"`
	From      gamemap.Pos `json:"from"`
	Direction string      `json:"direction"`
	Troops    uint16      `json:"troops"`
}

type ForceStartPayload struct {
	PlayerId string `json:"playerId"`
	IsVote   bool   `json:"isVote"`
}

func NewWebSocketServer(q queue.Queue) *WebSocketServer {
	return &WebSocketServer{
		queue: q,
	}
}

func (ws *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	slog.Info("websocket client connected", "remote", conn.RemoteAddr())

	for {
		var msg ClientMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket read error", "error", err)
			}
			break
		}

		if err := ws.handleMessage(msg); err != nil {
			slog.Error("message handling failed", "error", err, "type", msg.Type)
			errorMsg := map[string]string{
				"type":  "error",
				"error": err.Error(),
			}
			conn.WriteJSON(errorMsg)
		}
	}

	slog.Info("websocket client disconnected", "remote", conn.RemoteAddr())
}

func (ws *WebSocketServer) handleMessage(msg ClientMessage) error {
	switch msg.Type {
	case "join":
		return ws.handleJoinMessage(msg)
	case "leave":
		return ws.handleLeaveMessage(msg)
	case "move":
		return ws.handleMoveMessage(msg)
	case "forceStart":
		return ws.handleForceStartMessage(msg)
	case "surrender":
		return ws.handleSurrenderMessage(msg)
	case "createGame":
		return ws.handleCreateGameMessage(msg)
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

func (ws *WebSocketServer) handleJoinMessage(msg ClientMessage) error {
	var payload JoinPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("invalid join payload: %w", err)
	}

	joinCmd := game.JoinCommand{
		CommandEvent: game.CommandEvent{PlayerId: payload.PlayerId},
		PlayerName:   payload.PlayerName,
	}

	ws.queue.Publish(fmt.Sprintf("%s/commands", msg.GameId), joinCmd)
	return nil
}

func (ws *WebSocketServer) handleLeaveMessage(msg ClientMessage) error {
	var payload map[string]string
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("invalid leave payload: %w", err)
	}

	playerId := payload["playerId"]
	if playerId == "" {
		return fmt.Errorf("playerId is required")
	}

	leaveCmd := game.LeaveCommand{
		CommandEvent: game.CommandEvent{PlayerId: playerId},
	}

	ws.queue.Publish(fmt.Sprintf("%s/commands", msg.GameId), leaveCmd)
	return nil
}

func (ws *WebSocketServer) handleMoveMessage(msg ClientMessage) error {
	var payload MovePayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("invalid move payload: %w", err)
	}

	var direction game.MoveTowards
	switch payload.Direction {
	case "up":
		direction = game.MoveTowardsUp
	case "down":
		direction = game.MoveTowardsDown
	case "left":
		direction = game.MoveTowardsLeft
	case "right":
		direction = game.MoveTowardsRight
	default:
		return fmt.Errorf("invalid direction: %s", payload.Direction)
	}

	moveCmd := game.MoveCommand{
		CommandEvent: game.CommandEvent{PlayerId: payload.PlayerId},
		From:         payload.From,
		Direction:    direction,
		Troops:       block.Num(payload.Troops),
	}

	ws.queue.Publish(fmt.Sprintf("%s/commands", msg.GameId), moveCmd)
	return nil
}

func (ws *WebSocketServer) handleForceStartMessage(msg ClientMessage) error {
	var payload ForceStartPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("invalid force start payload: %w", err)
	}

	forceStartCmd := game.ForceStartCommand{
		CommandEvent: game.CommandEvent{PlayerId: payload.PlayerId},
		IsVote:       payload.IsVote,
	}

	ws.queue.Publish(fmt.Sprintf("%s/commands", msg.GameId), forceStartCmd)
	return nil
}

func (ws *WebSocketServer) handleSurrenderMessage(msg ClientMessage) error {
	var payload map[string]string
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("invalid surrender payload: %w", err)
	}

	playerId := payload["playerId"]
	if playerId == "" {
		return fmt.Errorf("playerId is required")
	}

	surrenderCmd := game.SurrenderCommand{
		CommandEvent: game.CommandEvent{PlayerId: playerId},
	}

	ws.queue.Publish(fmt.Sprintf("%s/commands", msg.GameId), surrenderCmd)
	return nil
}

func (ws *WebSocketServer) handleCreateGameMessage(msg ClientMessage) error {
	createGameCmd := map[string]interface{}{
		"type":   "createGame",
		"gameId": msg.GameId,
		"payload": map[string]interface{}{
			"gameMode": "classic_1v1", // 默认模式，可以从payload中解析
		},
	}

	ws.queue.Publish("lobby/commands", createGameCmd)
	return nil
}

func (ws *WebSocketServer) StartServer(addr string) error {
	http.HandleFunc("/ws", ws.HandleWebSocket)

	http.Handle("/", http.FileServer(http.Dir("./static/")))

	slog.Info("websocket server starting", "addr", addr)
	return http.ListenAndServe(addr, nil)
}
func (ws *WebSocketServer) StopServer() error {
	slog.Info("websocket server stopped")
	return nil
}
