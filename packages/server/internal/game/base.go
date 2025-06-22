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

type BaseCore struct {
	status     Status
	players    []Player
	_map       gamemap.Map
	turnNumber uint16
	mode       GameMode // 游戏模式

	// 定时器相关
	ctx        context.Context
	cancel     context.CancelFunc
	timer      *time.Timer
	timerMutex sync.Mutex

	// 消息队列和事件
	queue  queue.Queue
	gameId string
}

// NewBaseCore 创建新的BaseCore实例
func NewBaseCore(gameId string, q queue.Queue, mode GameMode) *BaseCore {
	return &BaseCore{
		status:     StatusWaiting,
		players:    make([]Player, 0),
		turnNumber: 0,
		queue:      q,
		gameId:     gameId,
		mode:       mode, // 使用传入的游戏模式
	}
}

func (c *BaseCore) Status() Status {
	return c.status
}

func (c *BaseCore) Players() []Player {
	return c.players
}

func (c *BaseCore) Join(player Player) error {
	if c.status != StatusWaiting {
		return fmt.Errorf("cannot join game in status: %s", c.status)
	}

	_, existingPlayer := c.findPlayerIndex(player.Id)
	if existingPlayer != nil {
		return errors.New("player already exists: " + player.Id)
	}

	// 设置新加入玩家的状态
	player.Status = PlayerStatusWaiting
	c.players = append(c.players, player)

	c.queue.Publish("broadcast", Join{
		PlayerEvent{Player: player.Id},
	})

	// 检查是否可以自动开始游戏
	c.checkGameTransition()

	return nil
}

func (c *BaseCore) Leave(playerId string) error {
	i, p := c.findPlayerIndex(playerId)
	if p == nil {
		return errors.New("player not found: " + playerId)
	}

	// 根据游戏状态设置不同的离开状态
	if c.status == StatusInProgress {
		c.players[i].Status = PlayerStatusDisconnected
		slog.Info("player left during game", "player", playerId)

		c.checkGameTransition()
	} else {
		// 如果游戏还没开始，直接移除玩家
		c.players = append(c.players[:i], c.players[i+1:]...)
		slog.Info("player left before game start", "player", playerId)
	}

	c.queue.Publish("broadcast", Disconnect{
		PlayerEvent: PlayerEvent{Player: playerId},
	})
	return nil
}

func (c *BaseCore) Map() gamemap.Map {
	if c._map == nil {
		slog.Error("map is not initialized")
		return nil
	}
	return c._map
}

func (c *BaseCore) NextTurn(turnNumber uint16) error {
	if c.status != StatusInProgress {
		return fmt.Errorf("cannot advance turn in status: %s", c.status)
	}
	if turnNumber != c.turnNumber+1 {
		return fmt.Errorf("invalid turn number: %d, expected: %d", turnNumber, c.turnNumber+1)
	}

	// 处理当前回合结束
	c._map.RoundEnd(c.turnNumber)

	c.turnNumber = turnNumber
	c._map.RoundStart(turnNumber)

	c.queue.Publish("broadcast", queue.NewEvent(AdvanceTurn{}))

	slog.Info("next turn", "turn", c.turnNumber)
	return nil
}

func (c *BaseCore) Move(playerID string, act Move) error {
	// 添加游戏状态检查
	if c.status != StatusInProgress {
		return fmt.Errorf("cannot move in status: %s", c.status)
	}

	// 检查玩家状态
	playerIndex, player := c.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}
	if player.Status != PlayerStatusInGame {
		return fmt.Errorf("player not in game status: %s", player.Status)
	}
	if player.Moves == 0 {
		return fmt.Errorf("player has no moves left: %s", playerID)
	}
	offset := getMoveOffset(act.Towards)
	newPos := gamemap.Pos{
		X: act.Pos.X + uint16(offset.X),
		Y: act.Pos.Y + uint16(offset.Y),
	}
	if !c._map.Size().IsPosValid(act.Pos) {
		return errors.New("invalid position: " + act.Pos.String())
	}
	if !c._map.Size().IsPosValid(newPos) {
		return errors.New("invalid position: " + newPos.String())
	}

	fromblock, err := c._map.Block(act.Pos)
	if err != nil {
		return err
	}
	targetBlock, err := c._map.Block(newPos)
	if err != nil {
		return err
	}

	owner := uint16(playerIndex)
	if fromblock.Owner() != block.Owner(owner) {
		return errors.New("not the owner of the block at position: " + act.Pos.String())
	}

	// 处理特殊的移动数量值
	if act.Num == 0 { // 0 表示选择所有可移动的兵
		act.Num = fromblock.Num() - 1
	} else if act.Num == 1 { // 1 表示选择一半兵力
		act.Num = fromblock.Num() / 2
	}

	if act.Num <= 0 {
		return fmt.Errorf("invalid number of blocks to move: %d", act.Num)
	}

	if fromblock.Num() < act.Num {
		return fmt.Errorf("not enough blocks to move: %d, available: %d", act.Num, fromblock.Num())
	}

	if !fromblock.AllowMove().From || !targetBlock.AllowMove().To {
		return errors.New("move not allowed from " + act.Pos.String() + " to " + newPos.String())
	}

	movedNum := fromblock.MoveFrom(act.Num)
	targetBlockNew := targetBlock.MoveTo(movedNum, fromblock.Owner())

	if targetBlockNew == nil {
		return errors.New("move rejected by target block: " + newPos.String())
	}

	c._map.SetBlock(act.Pos, fromblock)
	c._map.SetBlock(newPos, targetBlockNew)

	c.queue.Publish("broadcast", queue.NewEvent(MapUpdate{Map: c._map}))
	return nil
}

func (c *BaseCore) ForceStart(playerID string, isVote bool) error {
	if c.status != StatusWaiting {
		return fmt.Errorf("cannot force start in status: %s", c.status)
	}

	i, player := c.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}

	if isVote && player.Status == PlayerStatusWaiting {
		c.players[i].Status = PlayerStatusRequestForceStart
	} else if !isVote && player.Status == PlayerStatusRequestForceStart {
		c.players[i].Status = PlayerStatusWaiting
	} else {
		return fmt.Errorf("invalid vote state for player %s: %s", playerID, player.Status)
	}

	if len(c.players) < 2 {
		return errors.New("not enough players to start the game")
	}

	c.publishBroadcastEvent(ForceStart{
		PlayerEvent: PlayerEvent{Player: playerID},
		IsVote:     isVote,
	})

	// 检查是否可以开始游戏
	c.checkGameTransition()

	return nil
}

func (c *BaseCore) Surrender(playerID string) error {
	i, player := c.findPlayerIndex(playerID)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerID)
	}

	if player.Status == PlayerStatusInGame {
		c.players[i].Status = PlayerStatusSurrendered

		c.publishBroadcastEvent(Surrender{
			PlayerEvent: PlayerEvent{Player: playerID},
		})
		slog.Info("player surrendered", "player", playerID)

		c.checkGameTransition()
		return nil
	} else {
		return fmt.Errorf("cannot surrender in status: %s", player.Status)
	}
}

func (c *BaseCore) GetActivePlayerCount() int {
	count := 0
	for _, p := range c.players {
		if p.Status == PlayerStatusInGame {
			count++
		}
	}
	return count
}

func (c *BaseCore) IsGameReady() bool {
	if c.status != StatusWaiting {
		return false
	}

	// 使用游戏模式验证玩家数量
	return c.mode.ValidatePlayerCount(len(c.players))
}

// SetGameMode 设置游戏模式
func (c *BaseCore) SetGameMode(mode GameMode) error {
	if c.status != StatusWaiting {
		return fmt.Errorf("cannot change game mode in status: %s", c.status)
	}
	c.mode = mode
	slog.Info("game mode set", "mode", mode.Name, "gameId", c.gameId)
	return nil
}

// GetGameMode 获取当前游戏模式
func (c *BaseCore) GetGameMode() GameMode {
	return c.mode
}

// Start 启动游戏核心
func (c *BaseCore) Start() error {
	if c.status != StatusWaiting {
		return fmt.Errorf("cannot start game in status: %s", c.status)
	}

	// 创建新的上下文
	c.ctx, c.cancel = context.WithCancel(context.Background())

	// 更新游戏状态
	c.status = StatusInProgress
	for i := range c.players {
		c.players[i].Status = PlayerStatusInGame
	}

	// 启动定时器
	c.startTurnTimer()

	// 通知游戏开始
	c.publishControlEvent(StartGame{})
	c.publishBroadcastEvent(GameStatusUpdate{
		BroadcastEvent: BroadcastEvent{},
		Status:         c.status,
		Players:        c.players,
		TurnNumber:     c.turnNumber,
	})

	slog.Info("game started", "players", len(c.players))
	return nil
}

// Stop 停止游戏核心
func (c *BaseCore) Stop() error {
	// 停止定时器
	c.stopTurnTimer()

	// 取消上下文
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}

	// 更新状态
	if c.status == StatusInProgress {
		c.status = StatusFinished
	}

	// 通知游戏状态变化到广播流
	c.publishBroadcastEvent(GameStatusUpdate{
		BroadcastEvent: BroadcastEvent{},
		Status:         c.status,
		Players:        c.players,
		TurnNumber:     c.turnNumber,
	})

	slog.Info("game stopped")
	return nil
}

// startTurnTimer 启动回合定时器
func (c *BaseCore) startTurnTimer() {
	c.timerMutex.Lock()
	defer c.timerMutex.Unlock()

	// 停止旧的定时器
	if c.timer != nil {
		c.timer.Stop()
	}

	interval := c.mode.TurnTime // 使用游戏模式的回合时间
	c.timer = time.AfterFunc(interval, c.handleTurnTimeout)
}

// stopTurnTimer 停止回合定时器
func (c *BaseCore) stopTurnTimer() {
	c.timerMutex.Lock()
	defer c.timerMutex.Unlock()

	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
}

// handleTurnTimeout 处理回合超时
func (c *BaseCore) handleTurnTimeout() {
	// 检查游戏是否还在进行中
	if c.status != StatusInProgress {
		return
	}

	// 自动进入下一回合
	if err := c.NextTurn(c.turnNumber + 1); err != nil {
		slog.Error("failed to advance turn", "error", err)
		return
	}

	if c.isGameOver() {
		c.endGame(nil)
		return
	}

	// 重新启动定时器
	c.startTurnTimer()
}

func (c *BaseCore) TurnNumber() uint16 {
	return c.turnNumber
}

// SetPlayerStatus 设置玩家状态
func (c *BaseCore) SetPlayerStatus(playerId string, status PlayerStatus) error {
	i, player := c.findPlayerIndex(playerId)
	if player == nil {
		return fmt.Errorf("player not found: %s", playerId)
	}

	c.players[i].Status = status

	// 发布游戏状态更新到广播流
	c.queue.Publish("broadcast", GameStatusUpdate{
		Status:     c.status,
		Players:    c.players,
		TurnNumber: c.turnNumber,
	})

	// 检查状态变化是否影响游戏状态
	c.checkGameTransition()

	return nil
}

// GetPlayer 获取玩家信息
func (c *BaseCore) GetPlayer(playerId string) (*Player, error) {
	_, player := c.findPlayerIndex(playerId)
	if player == nil {
		return nil, fmt.Errorf("player not found: %s", playerId)
	}
	// 返回副本而不是指针，避免外部修改
	playerCopy := *player
	return &playerCopy, nil
}

// Event methods implementation for event types
func (e Join) Event() interface{} { return e }
func (e Disconnect) Event() interface{} { return e }
func (e Reconnect) Event() interface{} { return e }
func (e ForceStart) Event() interface{} { return e }
func (e Surrender) Event() interface{} { return e }
func (e StartGame) Event() interface{} { return e }
func (e EndGame) Event() interface{} { return e }
func (e GameStatusUpdate) Event() interface{} { return e }
func (e MapUpdate) Event() interface{} { return e }
func (e AdvanceTurn) Event() interface{} { return e }

// 辅助函数用于创建事件
func (c *BaseCore) publishPlayerEvent(player string, event queue.EventWrapper) {
	c.queue.Publish(c.gameId+"/player", event)
}

func (c *BaseCore) publishControlEvent(event queue.EventWrapper) {
	c.queue.Publish(c.gameId+"/control", event)
}

func (c *BaseCore) publishBroadcastEvent(event queue.EventWrapper) {
	c.queue.Publish(c.gameId+"/broadcast", event)
}

// checkGameTransition 检查是否需要自动进行游戏状态转换
func (c *BaseCore) checkGameTransition() {
	oldStatus := c.status

	switch c.status {
	case StatusWaiting:
		// 检查是否可以自动开始游戏
		if c.canStartGame() {
			c.startGame()
		}

	case StatusInProgress:
		if c.isGameOver() {
			c.endGame(nil)
		}
	}

	// 通知状态变化到广播流
	if c.status != oldStatus {
		c.queue.Publish("broadcast", GameStatusUpdate{
			Status:     c.status,
			Players:    c.players,
			TurnNumber: c.turnNumber,
		})
	}
}

// canStartGame 检查是否应该自动开始游戏
func (c *BaseCore) canStartGame() bool {
	if !c.IsGameReady() {
		return false
	}

	// 检查强制开始投票
	votes := 0
	totalPlayers := len(c.players)

	for _, p := range c.players {
		if p.Status == PlayerStatusRequestForceStart {
			votes++
		}
	}

	// 如果达到最大玩家数，自动开始
	if totalPlayers >= int(c.mode.MaxPlayers) {
		return true
	}

	// 否则需要超过半数投票且达到最小玩家数
	return votes >= totalPlayers/2+1 && totalPlayers >= int(c.mode.MinPlayers)
}

// startGame 自动开始游戏
func (c *BaseCore) startGame() {
	// 更新游戏状态
	c.status = StatusInProgress
	for i := range c.players {
		c.players[i].Status = PlayerStatusInGame
	}

	// 发布游戏状态更新到广播流
	c.publishBroadcastEvent(GameStatusUpdate{
		BroadcastEvent: BroadcastEvent{},
		Status:         c.status,
		Players:        c.players,
		TurnNumber:     c.turnNumber,
	})

	// 启动定时器
	c.startTurnTimer()

	// 通知游戏开始到控制流
	c.queue.Publish("control", StartGame{})

	slog.Info("game auto-started", "players", len(c.players))
}

// isGameOver 检查游戏是否应该结束
func (c *BaseCore) isGameOver() bool {
	return c.GetActivePlayerCount() <= 1
}

// endGame 自动结束游戏
func (c *BaseCore) endGame(endResult interface{}) {
	c.status = StatusFinished
	var winner *Player

	// 找到获胜者
	for i := range c.players {
		if c.players[i].Status == PlayerStatusInGame {
			c.players[i].Status = PlayerStatusFinished
			winner = &c.players[i]
		}
	}

	// 停止定时器
	c.stopTurnTimer()

	// 通知游戏结束到控制流
	var winnerId string
	if winner != nil {
		winnerId = winner.Id
	}

	c.publishControlEvent(EndGame{
		ControlEvent: ControlEvent{},
		Winner:       winnerId,
	})

	// 发布游戏状态更新到广播流
	c.queue.Publish("broadcast", GameStatusUpdate{
		Status:     c.status,
		Players:    c.players,
		TurnNumber: c.turnNumber,
	})

	if winner != nil {
		slog.Info("game auto-ended", "winner", winner.Id)
	} else {
		slog.Info("game auto-ended, no winner")
	}
}

// findPlayerIndex 查找玩家索引，返回索引和玩家指针
func (c *BaseCore) findPlayerIndex(playerID string) (int, *Player) {
	for i, p := range c.players {
		if p.Id == playerID {
			return i, &c.players[i]
		}
	}
	return -1, nil
}
