package game

import (
	"server/internal/game/block"
	gamemap "server/internal/game/map"
	"testing"
)

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

func TestCommandEvents(t *testing.T) {
	t.Run("join_command", func(t *testing.T) {
		cmd := JoinCommand{
			CommandEvent: CommandEvent{PlayerId: "test-player"},
			PlayerName:   "Test Player",
		}

		if cmd.PlayerId != "test-player" {
			t.Errorf("Expected PlayerId 'test-player', got '%s'", cmd.PlayerId)
		}
		if cmd.PlayerName != "Test Player" {
			t.Errorf("Expected PlayerName 'Test Player', got '%s'", cmd.PlayerName)
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
			CommandEvent: CommandEvent{PlayerId: "test-player"},
			From:         gamemap.Pos{X: 1, Y: 1},
			Direction:    MoveTowardsRight,
			Troops:       5,
		}

		if cmd.From.X != 1 || cmd.From.Y != 1 {
			t.Errorf("Expected position (1,1), got (%d,%d)", cmd.From.X, cmd.From.Y)
		}
		if cmd.Direction != MoveTowardsRight {
			t.Errorf("Expected direction %s, got %s", MoveTowardsRight, cmd.Direction)
		}
		if cmd.Troops != 5 {
			t.Errorf("Expected 5 troops, got %d", cmd.Troops)
		}
	})

	t.Run("force_start_command", func(t *testing.T) {
		cmd := ForceStartCommand{
			CommandEvent: CommandEvent{PlayerId: "test-player"},
			IsVote:       true,
		}

		if !cmd.IsVote {
			t.Error("Expected IsVote to be true")
		}
	})

	t.Run("surrender_command", func(t *testing.T) {
		cmd := SurrenderCommand{
			CommandEvent: CommandEvent{PlayerId: "test-player"},
		}

		if cmd.PlayerId != "test-player" {
			t.Errorf("Expected PlayerId 'test-player', got '%s'", cmd.PlayerId)
		}
	})
}

func TestControlEvents(t *testing.T) {
	t.Run("start_game_control", func(t *testing.T) {
		event := StartGameControl{}
		// 基本验证控制事件可以创建
		_ = event
	})

	t.Run("stop_game_control", func(t *testing.T) {
		event := StopGameControl{}
		// 基本验证控制事件可以创建
		_ = event
	})

	t.Run("turn_advance_control", func(t *testing.T) {
		event := TurnAdvanceControl{
			TurnNumber: 5,
		}

		if event.TurnNumber != 5 {
			t.Errorf("Expected TurnNumber 5, got %d", event.TurnNumber)
		}
	})
}

func TestBroadcastEvents(t *testing.T) {
	players := []Player{
		{Id: "player1", Name: "Player One", Status: PlayerStatusInGame},
		{Id: "player2", Name: "Player Two", Status: PlayerStatusInGame},
	}

	t.Run("player_joined_event", func(t *testing.T) {
		event := PlayerJoinedEvent{
			PlayerId:   "player1",
			PlayerName: "Player One",
			GameStatus: StatusWaiting,
			Players:    players,
		}

		if event.PlayerId != "player1" {
			t.Errorf("Expected PlayerId 'player1', got '%s'", event.PlayerId)
		}
		if event.GameStatus != StatusWaiting {
			t.Errorf("Expected GameStatus %s, got %s", StatusWaiting, event.GameStatus)
		}
		if len(event.Players) != 2 {
			t.Errorf("Expected 2 players, got %d", len(event.Players))
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
		event := GameStartedEvent{
			GameStatus: StatusInProgress,
			Players:    players,
			TurnNumber: 1,
		}

		if event.GameStatus != StatusInProgress {
			t.Errorf("Expected GameStatus %s, got %s", StatusInProgress, event.GameStatus)
		}
		if event.TurnNumber != 1 {
			t.Errorf("Expected TurnNumber 1, got %d", event.TurnNumber)
		}
	})

	t.Run("game_ended_event", func(t *testing.T) {
		event := GameEndedEvent{
			Winner:     "player1",
			GameStatus: StatusFinished,
			Players:    players,
		}

		if event.Winner != "player1" {
			t.Errorf("Expected winner 'player1', got '%s'", event.Winner)
		}
		if event.GameStatus != StatusFinished {
			t.Errorf("Expected GameStatus %s, got %s", StatusFinished, event.GameStatus)
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

func TestPlayerEvents(t *testing.T) {
	t.Run("player_error_event", func(t *testing.T) {
		event := PlayerErrorEvent{
			PlayerId: "test-player",
			Error:    "Test error message",
		}

		if event.PlayerId != "test-player" {
			t.Errorf("Expected PlayerId 'test-player', got '%s'", event.PlayerId)
		}
		if event.Error != "Test error message" {
			t.Errorf("Expected error 'Test error message', got '%s'", event.Error)
		}
	})
}

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

func TestMoveOffset(t *testing.T) {
	tests := []struct {
		direction MoveTowards
		expected  moveOffset
	}{
		{MoveTowardsLeft, moveOffset{-1, 0}},
		{MoveTowardsRight, moveOffset{1, 0}},
		{MoveTowardsUp, moveOffset{0, -1}},
		{MoveTowardsDown, moveOffset{0, 1}},
		{"invalid", moveOffset{0, 0}},
	}

	for _, test := range tests {
		t.Run(string(test.direction), func(t *testing.T) {
			result := getMoveOffset(test.direction)
			if result != test.expected {
				t.Errorf("Expected %+v, got %+v", test.expected, result)
			}
		})
	}
}
