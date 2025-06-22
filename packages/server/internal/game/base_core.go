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
		if p.Status == PlayerStatusInGame {
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

	// 设置新加入玩家的状态
	player.Status = PlayerStatusWaiting
	gc.players = append(gc.players, player)

	slog.Info("player joined", "player", player.Id, "gameId", gc.gameId)

	// 检查是否可以自动开始游戏
	gc.checkGameTransition()

	return nil
}

func (gc *BaseCore) Leave(playerId string) error {
	i, p := gc.findPlayerIndex(playerId)
	if p == nil {
		return errors.New("player not found: " + playerId)
	}

	// 根据游戏状态设置不同的离开状态
	if gc.status == StatusInProgress {
		gc.players[i].Status = PlayerStatusDisconnected
		slog.Info("player left during game", "player", playerId, "gameId", gc.gameId)
		gc.checkGameTransition()
	} else {
		// 如果游戏还没开始，直接移除玩家
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
	// 返回副本而不是指针，避免外部修改
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

	// 创建新的上下文
	gc.ctx, gc.cancel = context.WithCancel(context.Background())

	// 更新游戏状态
	gc.status = StatusInProgress
	for i := range gc.players {
		gc.players[i].Status = PlayerStatusInGame
	}

	// 启动定时器
	gc.startTurnTimer()

	// 发布游戏开始事件
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
	// 停止定时器
	gc.stopTurnTimer()

	// 取消上下文
	if gc.cancel != nil {
		gc.cancel()
		gc.cancel = nil
	}

	// 更新状态
	if gc.status == StatusInProgress {
		gc.status = StatusFinished
	}

	slog.Info("game stopped", "gameId", gc.gameId)
	return nil
}

func (gc *BaseCore) NextTurn(turnNumber uint16) error {
	if gc.status != StatusInProgress {
		return fmt.Errorf("cannot advance turn in status: %s", gc.status)
	}
	if turnNumber != gc.turnNumber+1 {
		return fmt.Errorf("invalid turn number: %d, expected: %d", turnNumber, gc.turnNumber+1)
	}

	// 处理当前回合结束
	if gc._map != nil {
		gc._map.RoundEnd(gc.turnNumber)
	}

	gc.turnNumber = turnNumber

	if gc._map != nil {
		gc._map.RoundStart(turnNumber)
	}

	slog.Info("next turn", "turn", gc.turnNumber, "gameId", gc.gameId)
	return nil
}

func (gc *BaseCore) Move(playerID string, move Move) error {
	// 添加游戏状态检查
	if gc.status != StatusInProgress {
		return fmt.Errorf("cannot move in status: %s", gc.status)
	}

	// 检查玩家状态
	playerIndex, player := gc.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}
	if player.Status != PlayerStatusInGame {
		return fmt.Errorf("player not in game status: %s", player.Status)
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

	// 处理特殊的移动数量值
	if move.Num == 0 { // 0 表示选择所有可移动的兵
		move.Num = fromBlock.Num() - 1
	} else if move.Num == 1 { // 1 表示选择一半兵力
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

	if len(gc.players) < 2 {
		return errors.New("not enough players to start the game")
	}

	slog.Info("force start vote", "player", playerID, "isVote", isVote, "gameId", gc.gameId)

	// 检查是否可以开始游戏
	gc.checkGameTransition()

	return nil
}

func (gc *BaseCore) Surrender(playerID string) error {
	i, player := gc.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}

	if player.Status == PlayerStatusInGame {
		gc.players[i].Status = PlayerStatusSurrendered
		slog.Info("player surrendered", "player", playerID, "gameId", gc.gameId)
		gc.checkGameTransition()
		return nil
	} else {
		return fmt.Errorf("cannot surrender in status: %s", player.Status)
	}
}

// =============================================================================
// 私有方法
// =============================================================================

// findPlayerIndex 查找玩家索引，返回索引和玩家指针
func (gc *BaseCore) findPlayerIndex(playerID string) (int, *Player) {
	for i, p := range gc.players {
		if p.Id == playerID {
			return i, &gc.players[i]
		}
	}
	return -1, nil
}

// checkGameTransition 检查是否需要自动进行游戏状态转换
func (gc *BaseCore) checkGameTransition() {
	oldStatus := gc.status

	switch gc.status {
	case StatusWaiting:
		// 检查是否可以自动开始游戏
		if gc.canStartGame() {
			gc.autoStartGame()
		}

	case StatusInProgress:
		if gc.isGameOver() {
			gc.autoEndGame()
		}
	}

	// 通知状态变化
	if gc.status != oldStatus && gc.onBroadcastEvent != nil {
		gc.onBroadcastEvent(GameStatusUpdateEvent{
			BroadcastEvent: BroadcastEvent{},
			Status:         gc.status,
			Players:        gc.players,
			TurnNumber:     gc.turnNumber,
		})
	}
}

// canStartGame 检查是否应该自动开始游戏
func (gc *BaseCore) canStartGame() bool {
	if !gc.IsGameReady() {
		return false
	}

	// 检查强制开始投票
	votes := 0
	totalPlayers := len(gc.players)

	for _, p := range gc.players {
		if p.Status == PlayerStatusRequestForceStart {
			votes++
		}
	}

	// 如果达到最大玩家数，自动开始
	if totalPlayers >= int(gc.mode.MaxPlayers) {
		return true
	}

	// 否则需要超过半数投票且达到最小玩家数
	return votes >= totalPlayers/2+1 && totalPlayers >= int(gc.mode.MinPlayers)
}

// autoStartGame 自动开始游戏
func (gc *BaseCore) autoStartGame() {
	// 更新游戏状态
	gc.status = StatusInProgress
	for i := range gc.players {
		gc.players[i].Status = PlayerStatusInGame
	}

	// 启动定时器
	gc.startTurnTimer()

	// 通知游戏开始
	if gc.onControlEvent != nil {
		gc.onControlEvent(StartGameControl{})
	}

	slog.Info("game auto-started", "players", len(gc.players), "gameId", gc.gameId)
}

// isGameOver 检查游戏是否应该结束
func (gc *BaseCore) isGameOver() bool {
	return gc.GetActivePlayerCount() <= 1
}

// autoEndGame 自动结束游戏
func (gc *BaseCore) autoEndGame() {
	gc.status = StatusFinished
	var winner *Player

	// 找到获胜者
	for i := range gc.players {
		if gc.players[i].Status == PlayerStatusInGame {
			gc.players[i].Status = PlayerStatusFinished
			winner = &gc.players[i]
		}
	}

	// 停止定时器
	gc.stopTurnTimer()

	// 通知游戏结束
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

// startTurnTimer 启动回合定时器
func (gc *BaseCore) startTurnTimer() {
	gc.timerMutex.Lock()
	defer gc.timerMutex.Unlock()

	// 停止旧的定时器
	if gc.timer != nil {
		gc.timer.Stop()
	}

	interval := gc.mode.TurnTime // 使用游戏模式的回合时间
	gc.timer = time.AfterFunc(interval, gc.handleTurnTimeout)
}

// stopTurnTimer 停止回合定时器
func (gc *BaseCore) stopTurnTimer() {
	gc.timerMutex.Lock()
	defer gc.timerMutex.Unlock()

	if gc.timer != nil {
		gc.timer.Stop()
		gc.timer = nil
	}
}

// handleTurnTimeout 处理回合超时
func (gc *BaseCore) handleTurnTimeout() {
	// 检查游戏是否还在进行中
	if gc.status != StatusInProgress {
		return
	}

	// 自动进入下一回合
	if err := gc.NextTurn(gc.turnNumber + 1); err != nil {
		slog.Error("failed to advance turn", "error", err, "gameId", gc.gameId)
		return
	}

	if gc.isGameOver() {
		gc.autoEndGame()
		return
	}

	// 重新启动定时器
	gc.startTurnTimer()
}
