package game

import (
	"server/internal/game/block"
	gamemap "server/internal/game/map"
	"testing"
)

// TestMoveTowards 测试移动方向
func TestMoveTowards(t *testing.T) {
	t.Run("move_offsets", func(t *testing.T) {
		testCases := []struct {
			direction MoveTowards
			expectedX int16
			expectedY int16
		}{
			{MoveTowardsLeft, -1, 0},
			{MoveTowardsRight, 1, 0},
			{MoveTowardsUp, 0, -1},
			{MoveTowardsDown, 0, 1},
		}

		for _, tc := range testCases {
			offset := getMoveOffset(tc.direction)
			if offset.X != tc.expectedX {
				t.Errorf("Direction %s: expected X offset %d, got %d", tc.direction, tc.expectedX, offset.X)
			}
			if offset.Y != tc.expectedY {
				t.Errorf("Direction %s: expected Y offset %d, got %d", tc.direction, tc.expectedY, offset.Y)
			}
		}
	})

	t.Run("invalid_direction", func(t *testing.T) {
		invalidDirection := MoveTowards("invalid")
		offset := getMoveOffset(invalidDirection)

		// 无效方向应该返回零偏移
		if offset.X != 0 || offset.Y != 0 {
			t.Errorf("Invalid direction should return zero offset, got (%d, %d)", offset.X, offset.Y)
		}
	})
}

// TestCommandEvents 测试指令事件
func TestCommandEvents(t *testing.T) {
	t.Run("join_command", func(t *testing.T) {
		cmd := JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			PlayerName:   "Test Player",
		}

		if cmd.PlayerId != "player1" {
			t.Errorf("Expected player ID 'player1', got %s", cmd.PlayerId)
		}
		if cmd.PlayerName != "Test Player" {
			t.Errorf("Expected player name 'Test Player', got %s", cmd.PlayerName)
		}
	})

	t.Run("leave_command", func(t *testing.T) {
		cmd := LeaveCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
		}

		if cmd.PlayerId != "player1" {
			t.Errorf("Expected player ID 'player1', got %s", cmd.PlayerId)
		}
	})

	t.Run("move_command", func(t *testing.T) {
		cmd := MoveCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			From:         gamemap.Pos{X: 5, Y: 10},
			Direction:    MoveTowardsRight,
			Troops:       block.Num(3),
		}

		if cmd.PlayerId != "player1" {
			t.Errorf("Expected player ID 'player1', got %s", cmd.PlayerId)
		}
		if cmd.From.X != 5 || cmd.From.Y != 10 {
			t.Errorf("Expected position (5, 10), got (%d, %d)", cmd.From.X, cmd.From.Y)
		}
		if cmd.Direction != MoveTowardsRight {
			t.Errorf("Expected direction right, got %s", cmd.Direction)
		}
		if cmd.Troops != 3 {
			t.Errorf("Expected 3 troops, got %d", cmd.Troops)
		}
	})

	t.Run("force_start_command", func(t *testing.T) {
		// 测试投票开始
		cmd1 := ForceStartCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			IsVote:       true,
		}

		if !cmd1.IsVote {
			t.Error("Expected IsVote to be true")
		}

		// 测试取消投票
		cmd2 := ForceStartCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
			IsVote:       false,
		}

		if cmd2.IsVote {
			t.Error("Expected IsVote to be false")
		}
	})

	t.Run("surrender_command", func(t *testing.T) {
		cmd := SurrenderCommand{
			CommandEvent: CommandEvent{PlayerId: "player1"},
		}

		if cmd.PlayerId != "player1" {
			t.Errorf("Expected player ID 'player1', got %s", cmd.PlayerId)
		}
	})
}

// TestControlEvents 测试控制事件
func TestControlEvents(t *testing.T) {
	t.Run("start_game_control", func(t *testing.T) {
		ctrl := StartGameControl{
			ControlEvent: ControlEvent{},
		}

		// 验证类型
		_ = ctrl.ControlEvent
	})

	t.Run("stop_game_control", func(t *testing.T) {
		ctrl := StopGameControl{
			ControlEvent: ControlEvent{},
		}

		// 验证类型
		_ = ctrl.ControlEvent
	})

	t.Run("turn_advance_control", func(t *testing.T) {
		ctrl := TurnAdvanceControl{
			ControlEvent: ControlEvent{},
			TurnNumber:   5,
		}

		if ctrl.TurnNumber != 5 {
			t.Errorf("Expected turn number 5, got %d", ctrl.TurnNumber)
		}
	})
}

// TestBroadcastEvents 测试广播事件
func TestBroadcastEvents(t *testing.T) {
	t.Run("player_joined_event", func(t *testing.T) {
		players := []Player{
			{Id: "player1", Name: "Player One", Status: PlayerStatusWaiting},
		}

		event := PlayerJoinedEvent{
			BroadcastEvent: BroadcastEvent{},
			PlayerId:       "player1",
			PlayerName:     "Player One",
			GameStatus:     StatusWaiting,
			Players:        players,
		}

		if event.PlayerId != "player1" {
			t.Errorf("Expected player ID 'player1', got %s", event.PlayerId)
		}
		if event.PlayerName != "Player One" {
			t.Errorf("Expected player name 'Player One', got %s", event.PlayerName)
		}
		if event.GameStatus != StatusWaiting {
			t.Errorf("Expected game status waiting, got %s", event.GameStatus)
		}
		if len(event.Players) != 1 {
			t.Errorf("Expected 1 player, got %d", len(event.Players))
		}
	})

	t.Run("player_left_event", func(t *testing.T) {
		players := []Player{
			{Id: "player2", Name: "Player Two", Status: PlayerStatusWaiting},
		}

		event := PlayerLeftEvent{
			BroadcastEvent: BroadcastEvent{},
			PlayerId:       "player1",
			GameStatus:     StatusWaiting,
			Players:        players,
		}

		if event.PlayerId != "player1" {
			t.Errorf("Expected player ID 'player1', got %s", event.PlayerId)
		}
		if len(event.Players) != 1 {
			t.Errorf("Expected 1 remaining player, got %d", len(event.Players))
		}
	})

	t.Run("game_started_event", func(t *testing.T) {
		players := []Player{
			{Id: "player1", Name: "Player One", Status: PlayerStatusInGame},
			{Id: "player2", Name: "Player Two", Status: PlayerStatusInGame},
		}

		event := GameStartedEvent{
			BroadcastEvent: BroadcastEvent{},
			GameStatus:     StatusInProgress,
			Players:        players,
			TurnNumber:     1,
		}

		if event.GameStatus != StatusInProgress {
			t.Errorf("Expected game status in_progress, got %s", event.GameStatus)
		}
		if event.TurnNumber != 1 {
			t.Errorf("Expected turn number 1, got %d", event.TurnNumber)
		}
		if len(event.Players) != 2 {
			t.Errorf("Expected 2 players, got %d", len(event.Players))
		}
	})

	t.Run("game_ended_event", func(t *testing.T) {
		players := []Player{
			{Id: "player1", Name: "Player One", Status: PlayerStatusFinished},
			{Id: "player2", Name: "Player Two", Status: PlayerStatusFinished},
		}

		event := GameEndedEvent{
			BroadcastEvent: BroadcastEvent{},
			Winner:         "player1",
			GameStatus:     StatusFinished,
			Players:        players,
		}

		if event.Winner != "player1" {
			t.Errorf("Expected winner 'player1', got %s", event.Winner)
		}
		if event.GameStatus != StatusFinished {
			t.Errorf("Expected game status finished, got %s", event.GameStatus)
		}
	})

	t.Run("game_status_update_event", func(t *testing.T) {
		event := GameStatusUpdateEvent{
			BroadcastEvent: BroadcastEvent{},
			Status:         StatusInProgress,
			Players:        []Player{},
			TurnNumber:     3,
		}

		if event.Status != StatusInProgress {
			t.Errorf("Expected status in_progress, got %s", event.Status)
		}
		if event.TurnNumber != 3 {
			t.Errorf("Expected turn number 3, got %d", event.TurnNumber)
		}
	})
}

// TestPlayerEvents 测试玩家事件
func TestPlayerEvents(t *testing.T) {
	t.Run("player_error_event", func(t *testing.T) {
		event := PlayerErrorEvent{
			PlayerEvent: PlayerEvent{},
			PlayerId:    "player1",
			Error:       "Invalid move",
		}

		if event.PlayerId != "player1" {
			t.Errorf("Expected player ID 'player1', got %s", event.PlayerId)
		}
		if event.Error != "Invalid move" {
			t.Errorf("Expected error 'Invalid move', got %s", event.Error)
		}
	})
}

// TestMove 测试移动数据结构
func TestMove(t *testing.T) {
	t.Run("move_structure", func(t *testing.T) {
		move := Move{
			Pos:     gamemap.Pos{X: 10, Y: 20},
			Towards: MoveTowardsLeft,
			Num:     block.Num(5),
		}

		if move.Pos.X != 10 || move.Pos.Y != 20 {
			t.Errorf("Expected position (10, 20), got (%d, %d)", move.Pos.X, move.Pos.Y)
		}
		if move.Towards != MoveTowardsLeft {
			t.Errorf("Expected direction left, got %s", move.Towards)
		}
		if move.Num != 5 {
			t.Errorf("Expected num 5, got %d", move.Num)
		}
	})
}

// TestEventInheritance 测试事件继承结构
func TestEventInheritance(t *testing.T) {
	t.Run("command_event_inheritance", func(t *testing.T) {
		// 所有指令事件都应该包含CommandEvent
		joinCmd := JoinCommand{CommandEvent: CommandEvent{PlayerId: "test"}}
		leaveCmd := LeaveCommand{CommandEvent: CommandEvent{PlayerId: "test"}}
		moveCmd := MoveCommand{CommandEvent: CommandEvent{PlayerId: "test"}}
		forceStartCmd := ForceStartCommand{CommandEvent: CommandEvent{PlayerId: "test"}}
		surrenderCmd := SurrenderCommand{CommandEvent: CommandEvent{PlayerId: "test"}}

		// 验证PlayerId字段存在
		if joinCmd.PlayerId != "test" {
			t.Error("JoinCommand should inherit PlayerId from CommandEvent")
		}
		if leaveCmd.PlayerId != "test" {
			t.Error("LeaveCommand should inherit PlayerId from CommandEvent")
		}
		if moveCmd.PlayerId != "test" {
			t.Error("MoveCommand should inherit PlayerId from CommandEvent")
		}
		if forceStartCmd.PlayerId != "test" {
			t.Error("ForceStartCommand should inherit PlayerId from CommandEvent")
		}
		if surrenderCmd.PlayerId != "test" {
			t.Error("SurrenderCommand should inherit PlayerId from CommandEvent")
		}
	})

	t.Run("broadcast_event_inheritance", func(t *testing.T) {
		// 所有广播事件都应该包含BroadcastEvent
		playerJoined := PlayerJoinedEvent{BroadcastEvent: BroadcastEvent{}}
		playerLeft := PlayerLeftEvent{BroadcastEvent: BroadcastEvent{}}
		gameStarted := GameStartedEvent{BroadcastEvent: BroadcastEvent{}}
		gameEnded := GameEndedEvent{BroadcastEvent: BroadcastEvent{}}
		statusUpdate := GameStatusUpdateEvent{BroadcastEvent: BroadcastEvent{}}

		// 验证结构正确性（编译时检查）
		_ = playerJoined.BroadcastEvent
		_ = playerLeft.BroadcastEvent
		_ = gameStarted.BroadcastEvent
		_ = gameEnded.BroadcastEvent
		_ = statusUpdate.BroadcastEvent
	})

	t.Run("control_event_inheritance", func(t *testing.T) {
		// 所有控制事件都应该包含ControlEvent
		startGame := StartGameControl{ControlEvent: ControlEvent{}}
		stopGame := StopGameControl{ControlEvent: ControlEvent{}}
		turnAdvance := TurnAdvanceControl{ControlEvent: ControlEvent{}}

		// 验证结构正确性（编译时检查）
		_ = startGame.ControlEvent
		_ = stopGame.ControlEvent
		_ = turnAdvance.ControlEvent
	})

	t.Run("player_event_inheritance", func(t *testing.T) {
		// 所有玩家事件都应该包含PlayerEvent
		playerError := PlayerErrorEvent{PlayerEvent: PlayerEvent{}}

		// 验证结构正确性（编译时检查）
		_ = playerError.PlayerEvent
	})
}

// TestEventConstants 测试事件相关常量
func TestEventConstants(t *testing.T) {
	t.Run("move_directions", func(t *testing.T) {
		directions := []MoveTowards{
			MoveTowardsLeft,
			MoveTowardsRight,
			MoveTowardsUp,
			MoveTowardsDown,
		}

		// 验证所有方向都有定义
		expectedDirections := []string{"left", "right", "up", "down"}

		for i, direction := range directions {
			if string(direction) != expectedDirections[i] {
				t.Errorf("Expected direction %s, got %s", expectedDirections[i], string(direction))
			}
		}
	})
}
