package game

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"server/internal/game/block"
	gamemap "server/internal/game/map"
	"server/internal/queue"
	"sync"
	"time"
)

// BaseCore 纯粹的游戏逻辑层，无外部依赖，提供标准化的事件返回接口
type BaseCore struct {
	gameId     string
	status     Status
	players    []Player
	_map       gamemap.Map
	turnNumber uint16
	mode       GameMode

	// 定时器相关
	ctx        context.Context
	cancel     context.CancelFunc
	timer      *time.Timer
	timerMutex sync.Mutex

	// 事件回调
	onBroadcastEvent func(queue.Event)
	onControlEvent   func(queue.Event)
}

// NewBaseCore 创建新的BaseCore实例
func NewBaseCore(gameId string, mode GameMode) *BaseCore {
	if mode.Name == "" {
		panic("game mode cannot be nil or empty")
	}
	return &BaseCore{
		gameId:     gameId,
		status:     StatusWaiting,
		players:    make([]Player, 0),
		turnNumber: 0,
		mode:       mode,
	}
}

// SetEventHandlers 设置事件处理回调函数
func (gc *BaseCore) SetEventHandlers(
	onBroadcast func(queue.Event),
	onControl func(queue.Event),
) {
	gc.onBroadcastEvent = onBroadcast
	gc.onControlEvent = onControl
}

// =============================================================================
// 实现Core接口
// =============================================================================

func (gc *BaseCore) Status() Status {
	return gc.status
}

func (gc *BaseCore) Players() []Player {
	return gc.players
}

func (gc *BaseCore) TurnNumber() uint16 {
	return gc.turnNumber
}

func (gc *BaseCore) IsGameReady() bool {
	if gc.status != StatusWaiting {
		return false
	}
	return gc.mode.ValidatePlayerCount(len(gc.players))
}

func (gc *BaseCore) GetActivePlayerCount() int {
	count := 0
	for _, p := range gc.players {
		if p.IsActive() {
			count++
		}
	}
	return count
}

func (gc *BaseCore) Join(player Player) error {
	if gc.status != StatusWaiting {
		return fmt.Errorf("cannot join game in status: %s", gc.status)
	}

	_, existingPlayer := gc.findPlayerIndex(player.Id)
	if existingPlayer != nil {
		return errors.New("player already exists: " + player.Id)
	}

	player.Status = PlayerStatusWaiting
	gc.players = append(gc.players, player)

	slog.Info("player joined", "player", player.Id, "gameId", gc.gameId)

	gc.checkGameTransition()

	return nil
}

func (gc *BaseCore) Leave(playerId string) error {
	i, p := gc.findPlayerIndex(playerId)
	if p == nil {
		return errors.New("player not found: " + playerId)
	}

	if gc.status == StatusInProgress {
		gc.players[i].Status = PlayerStatusDisconnected
		slog.Info("player left during game", "player", playerId, "gameId", gc.gameId)
		gc.checkGameTransition()
	} else {
		gc.players = append(gc.players[:i], gc.players[i+1:]...)
		slog.Info("player left before game start", "player", playerId, "gameId", gc.gameId)
	}

	return nil
}

func (gc *BaseCore) GetPlayer(playerId string) (*Player, error) {
	_, player := gc.findPlayerIndex(playerId)
	if player == nil {
		return nil, fmt.Errorf("player not found: %s", playerId)
	}
	playerCopy := *player
	return &playerCopy, nil
}

func (gc *BaseCore) Map() gamemap.Map {
	if gc._map == nil {
		slog.Error("map is not initialized", "gameId", gc.gameId)
		return nil
	}
	return gc._map
}

func (gc *BaseCore) Start() error {
	if gc.status != StatusWaiting {
		return fmt.Errorf("cannot start game in status: %s", gc.status)
	}

	gc.ctx, gc.cancel = context.WithCancel(context.Background())

	if err := gc.initializeMap(); err != nil {
		return fmt.Errorf("failed to initialize map: %w", err)
	}

	gc.status = StatusInProgress
	for i := range gc.players {
		gc.players[i].Status = PlayerStatusInGame
		gc.players[i].Moves = gc.mode.MovesPerTurn
	}

	gc.startTurnTimer()

	if gc.onBroadcastEvent != nil {
		gc.onBroadcastEvent(GameStartedEvent{
			BroadcastEvent: BroadcastEvent{},
			GameStatus:     gc.status,
			Players:        gc.players,
			TurnNumber:     gc.turnNumber,
		})
	}

	slog.Info("game started", "players", len(gc.players), "gameId", gc.gameId)
	return nil
}

func (gc *BaseCore) Stop() error {
	gc.stopTurnTimer()

	if gc.cancel != nil {
		gc.cancel()
		gc.cancel = nil
	}

	gc.status = StatusFinished

	return nil
}

func (gc *BaseCore) NextTurn(turnNumber uint16) error {
	if gc.status != StatusInProgress {
		return fmt.Errorf("cannot advance turn in status: %s", gc.status)
	}
	if turnNumber != gc.turnNumber+1 {
		return fmt.Errorf("invalid turn number: %d, expected: %d", turnNumber, gc.turnNumber+1)
	}

	if gc._map != nil {
		gc._map.RoundEnd(gc.turnNumber)
	}

	gc.turnNumber = turnNumber

	for i := range gc.players {
		if gc.players[i].CanOperate() {
			gc.players[i].Moves = gc.mode.MovesPerTurn
		}
	}

	if gc._map != nil {
		gc._map.RoundStart(turnNumber)
	}

	if gc.isGameOver() {
		gc.autoEndGame()
		return nil
	}

	if gc.onBroadcastEvent != nil {
		gc.onBroadcastEvent(TurnStartedEvent{
			BroadcastEvent: BroadcastEvent{},
			TurnNumber:     gc.turnNumber,
			Players:        gc.players,
		})
	}

	slog.Info("next turn", "turn", gc.turnNumber, "gameId", gc.gameId)
	return nil
}

func (gc *BaseCore) Move(playerID string, move Move) error {
	if gc.status != StatusInProgress {
		return fmt.Errorf("cannot move in status: %s", gc.status)
	}

	playerIndex, player := gc.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}
	if !player.CanOperate() {
		return fmt.Errorf("player cannot operate (status: %s)", player.Status)
	}
	if player.Moves == 0 {
		return fmt.Errorf("player has no moves left: %s", playerID)
	}

	offset := getMoveOffset(move.Towards)
	newPos := gamemap.Pos{
		X: move.Pos.X + uint16(offset.X),
		Y: move.Pos.Y + uint16(offset.Y),
	}

	if !gc._map.Size().IsPosValid(move.Pos) {
		return errors.New("invalid position: " + move.Pos.String())
	}
	if !gc._map.Size().IsPosValid(newPos) {
		return errors.New("invalid position: " + newPos.String())
	}

	fromBlock, err := gc._map.Block(move.Pos)
	if err != nil {
		return err
	}
	targetBlock, err := gc._map.Block(newPos)
	if err != nil {
		return err
	}

	owner := uint16(playerIndex)
	if fromBlock.Owner() != block.Owner(owner) {
		return errors.New("not the owner of the block at position: " + move.Pos.String())
	}

	if move.Num == 0 {
		move.Num = fromBlock.Num() - 1
	} else if move.Num == 1 {
		move.Num = fromBlock.Num() / 2
	}

	if move.Num <= 0 {
		return fmt.Errorf("invalid number of blocks to move: %d", move.Num)
	}

	if fromBlock.Num() < move.Num {
		return fmt.Errorf("not enough blocks to move: %d, available: %d", move.Num, fromBlock.Num())
	}

	if !fromBlock.AllowMove().From || !targetBlock.AllowMove().To {
		return errors.New("move not allowed from " + move.Pos.String() + " to " + newPos.String())
	}

	movedNum := fromBlock.MoveFrom(move.Num)
	targetBlockNew := targetBlock.MoveTo(movedNum, fromBlock.Owner())

	if targetBlockNew == nil {
		return errors.New("move rejected by target block: " + newPos.String())
	}

	gc._map.SetBlock(move.Pos, fromBlock)
	gc._map.SetBlock(newPos, targetBlockNew)

	gc.players[playerIndex].Moves--

	if gc.onBroadcastEvent != nil {
		gc.onBroadcastEvent(PlayerMovedEvent{
			BroadcastEvent: BroadcastEvent{},
			PlayerId:       playerID,
			Move:           move,
			MovesLeft:      gc.players[playerIndex].Moves,
		})
	}

	return nil
}

func (gc *BaseCore) ForceStart(playerID string, isVote bool) error {
	if gc.status != StatusWaiting {
		return fmt.Errorf("cannot force start in status: %s", gc.status)
	}

	i, player := gc.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}

	if isVote && player.Status == PlayerStatusWaiting {
		gc.players[i].Status = PlayerStatusRequestForceStart
	} else if !isVote && player.Status == PlayerStatusRequestForceStart {
		gc.players[i].Status = PlayerStatusWaiting
	} else {
		return fmt.Errorf("invalid vote state for player %s: %s", playerID, player.Status)
	}

	slog.Info("force start vote", "player", playerID, "isVote", isVote, "gameId", gc.gameId)

	gc.checkGameTransition()

	return nil
}

func (gc *BaseCore) Surrender(playerID string) error {
	i, player := gc.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}

	if player.CanOperate() {
		gc.players[i].Status = PlayerStatusSurrendered
		gc.players[i].FinishReason = FinishReasonSurrendered
		slog.Info("player surrendered", "player", playerID, "gameId", gc.gameId)
		gc.checkGameTransition()
		return nil
	} else {
		return fmt.Errorf("cannot surrender in status: %s", player.Status)
	}
}

// =============================================================================
// 连接管理方法
// =============================================================================

func (gc *BaseCore) PlayerConnect(playerID string) error {
	i, player := gc.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}

	gc.players[i].Connection.IsConnected = true
	gc.players[i].Connection.DisconnectedAt = 0

	slog.Info("player connected", "player", playerID, "gameId", gc.gameId)
	return nil
}

func (gc *BaseCore) PlayerDisconnect(playerID string) error {
	i, player := gc.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}

	gc.players[i].Connection.IsConnected = false
	gc.players[i].Connection.DisconnectedAt = time.Now().UnixMilli()
	gc.players[i].Connection.ReconnectTimeout = time.Now().UnixMilli() + 30000

	if player.Status == PlayerStatusInGame {
		gc.players[i].Status = PlayerStatusDisconnected
	}

	slog.Info("player disconnected", "player", playerID, "gameId", gc.gameId)
	return nil
}

func (gc *BaseCore) PlayerReconnect(playerID string) error {
	i, player := gc.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}

	gc.players[i].Connection.IsConnected = true
	gc.players[i].Connection.DisconnectedAt = 0

	switch player.Status {
	case PlayerStatusDisconnected:
		gc.players[i].Status = PlayerStatusInGame
		slog.Info("player reconnected to game", "player", playerID, "gameId", gc.gameId)

	case PlayerStatusSurrendered, PlayerStatusLost, PlayerStatusSpectator:
		gc.players[i].Status = PlayerStatusSpectator
		slog.Info("player reconnected as spectator", "player", playerID, "gameId", gc.gameId)

	default:
		slog.Info("player reconnected", "player", playerID, "status", player.Status, "gameId", gc.gameId)
	}

	return nil
}

func (gc *BaseCore) CheckDisconnectedPlayers(currentTimeMs int64) error {
	hasChangedPlayers := false

	for i, player := range gc.players {
		if player.Status == PlayerStatusDisconnected &&
			!player.Connection.IsConnected &&
			player.Connection.ReconnectTimeout > 0 &&
			currentTimeMs > player.Connection.ReconnectTimeout {

			gc.players[i].Status = PlayerStatusSpectator
			gc.players[i].FinishReason = FinishReasonDisconnected
			hasChangedPlayers = true

			slog.Info("player disconnected timeout, now spectator",
				"player", player.Id, "gameId", gc.gameId)
		}
	}

	if hasChangedPlayers {
		gc.checkGameTransition()
	}

	return nil
}

// =============================================================================
// 私有方法
// =============================================================================

func (gc *BaseCore) findPlayerIndex(playerID string) (int, *Player) {
	for i, p := range gc.players {
		if p.Id == playerID {
			return i, &gc.players[i]
		}
	}
	return -1, nil
}

func (gc *BaseCore) checkGameTransition() {
	oldStatus := gc.status

	switch gc.status {
	case StatusWaiting:
		if gc.canStartGame() {
			gc.autoStartGame()
		}

	case StatusInProgress:
		if gc.isGameOver() {
			gc.autoEndGame()
		}
	}

	if gc.status != oldStatus && gc.onBroadcastEvent != nil {
		gc.onBroadcastEvent(GameStatusUpdateEvent{
			BroadcastEvent: BroadcastEvent{},
			Status:         gc.status,
			Players:        gc.players,
			TurnNumber:     gc.turnNumber,
		})
	}
}

func (gc *BaseCore) canStartGame() bool {
	if !gc.IsGameReady() {
		return false
	}

	votes := 0
	totalPlayers := len(gc.players)

	for _, p := range gc.players {
		if p.Status == PlayerStatusRequestForceStart {
			votes++
		}
	}

	if votes >= totalPlayers/2+1 && totalPlayers >= int(gc.mode.MinPlayers) {
		return true
	}

	if totalPlayers >= int(gc.mode.MaxPlayers) {
		if gc.mode.Name == "test_mode" {
			return false
		}
		return true
	}

	return false
}

func (gc *BaseCore) autoStartGame() {
	if err := gc.Start(); err != nil {
		slog.Error("failed to auto-start game", "error", err, "gameId", gc.gameId)
		return
	}
	slog.Info("game auto-started by vote/player count", "players", len(gc.players), "gameId", gc.gameId)
}

func (gc *BaseCore) isGameOver() bool {
	if gc._map == nil {
		return false
	}

	playerCastles := make(map[block.Owner]int)
	playerHasUnits := make(map[block.Owner]bool)

	size := gc._map.Size()
	for y := uint16(1); y <= size.Height; y++ {
		for x := uint16(1); x <= size.Width; x++ {
			pos := gamemap.Pos{X: x, Y: y}
			b, err := gc._map.Block(pos)
			if err != nil || b == nil {
				continue
			}

			if b.Owner() != block.Owner(0) {
				meta := b.Meta()
				if meta.Name == block.CastleName || meta.Name == block.KingName {
					playerCastles[b.Owner()]++
				}
				if b.Num() > 0 {
					playerHasUnits[b.Owner()] = true
				}
			}
		}
	}

	activePlayers := 0
	inGamePlayers := 0
	for i, player := range gc.players {
		if player.IsActive() {
			owner := block.Owner(i)
			if playerCastles[owner] > 0 || playerHasUnits[owner] {
				activePlayers++
			} else if player.Status == PlayerStatusInGame {
				gc.players[i].Status = PlayerStatusLost
				gc.players[i].FinishReason = FinishReasonDefeated
			}
		}

		if player.Status == PlayerStatusInGame || player.Status == PlayerStatusDisconnected {
			inGamePlayers++
		}
	}

	return activePlayers <= 1 && inGamePlayers <= 1
}

func (gc *BaseCore) autoEndGame() {
	gc.status = StatusFinished
	var winner *Player

	for i := range gc.players {
		if gc.players[i].IsActive() {
			gc.players[i].Status = PlayerStatusWinner
			gc.players[i].FinishReason = FinishReasonVictory
			winner = &gc.players[i]
		}
	}

	gc.stopTurnTimer()

	var winnerId string
	if winner != nil {
		winnerId = winner.Id
	}

	if gc.onBroadcastEvent != nil {
		gc.onBroadcastEvent(GameEndedEvent{
			BroadcastEvent: BroadcastEvent{},
			Winner:         winnerId,
			GameStatus:     gc.status,
			Players:        gc.players,
		})
	}

	if winner != nil {
		slog.Info("game auto-ended", "winner", winner.Id, "gameId", gc.gameId)
	} else {
		slog.Info("game auto-ended, no winner", "gameId", gc.gameId)
	}
}

// =============================================================================
// 定时器相关方法
// =============================================================================

func (gc *BaseCore) startTurnTimer() {
	gc.timerMutex.Lock()
	defer gc.timerMutex.Unlock()

	if gc.timer != nil {
		gc.timer.Stop()
	}

	interval := gc.mode.GetTurnTime()
	gc.timer = time.AfterFunc(interval, gc.handleTurnTimeout)
}

func (gc *BaseCore) stopTurnTimer() {
	gc.timerMutex.Lock()
	defer gc.timerMutex.Unlock()

	if gc.timer != nil {
		gc.timer.Stop()
		gc.timer = nil
	}
}

func (gc *BaseCore) handleTurnTimeout() {
	if gc.status != StatusInProgress {
		return
	}

	if err := gc.NextTurn(gc.turnNumber + 1); err != nil {
		slog.Error("failed to advance turn", "error", err, "gameId", gc.gameId)
		return
	}

	if gc.isGameOver() {
		gc.autoEndGame()
		return
	}

	gc.startTurnTimer()
}

func (gc *BaseCore) initializeMap() error {
	mapSize := gamemap.Size{Width: 20, Height: 20}
	mapInfo := gamemap.Info{
		Id:   gc.gameId + "_map",
		Name: "Game Map",
		Desc: "Generated map for game " + gc.gameId,
	}

	players := make([]gamemap.Player, len(gc.players))
	for i, p := range gc.players {
		players[i] = gamemap.Player{
			Index:    i,
			Owner:    block.Owner(i),
			IsActive: p.IsActive(),
		}
	}

	generatedMap, err := gamemap.GenerateMap("base", mapSize, mapInfo, players)
	if err != nil {
		return fmt.Errorf("failed to generate map: %w", err)
	}

	gc._map = generatedMap
	return nil
}
