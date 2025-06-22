package game

import (
	"testing"
	"time"
)

// TestGameMode_Classic1v1 测试经典1v1模式
func TestGameMode_Classic1v1(t *testing.T) {
	mode := Classic1v1

	if mode.Name != "classic_1v1" {
		t.Errorf("Expected name 'classic_1v1', got '%s'", mode.Name)
	}

	if mode.MaxPlayers != 2 {
		t.Errorf("Expected max players 2, got %d", mode.MaxPlayers)
	}

	if mode.MinPlayers != 2 {
		t.Errorf("Expected min players 2, got %d", mode.MinPlayers)
	}

	if mode.TeamSize != 1 {
		t.Errorf("Expected team size 1, got %d", mode.TeamSize)
	}

	if mode.TurnTime != time.Second {
		t.Errorf("Expected turn time 1s, got %v", mode.TurnTime)
	}

	if mode.Speed != 1.0 {
		t.Errorf("Expected speed 1.0, got %f", mode.Speed)
	}

	if mode.MovesPerTurn != 2 {
		t.Errorf("Expected moves per turn 2, got %d", mode.MovesPerTurn)
	}

	if len(mode.EndConditions) == 0 {
		t.Error("Expected at least one end condition")
	}
}

// TestGameMode_ValidatePlayerCount 测试玩家数量验证
func TestGameMode_ValidatePlayerCount(t *testing.T) {
	mode := Classic1v1

	if !mode.ValidatePlayerCount(2) {
		t.Error("Expected 2 players to be valid")
	}

	if mode.ValidatePlayerCount(1) {
		t.Error("Expected 1 player to be invalid")
	}

	if mode.ValidatePlayerCount(3) {
		t.Error("Expected 3 players to be invalid")
	}

	if mode.ValidatePlayerCount(0) {
		t.Error("Expected 0 players to be invalid")
	}

	if mode.ValidatePlayerCount(-1) {
		t.Error("Expected negative players to be invalid")
	}
}

// TestGameMode_CalculateTeamCount 测试团队计数计算
func TestGameMode_CalculateTeamCount(t *testing.T) {
	mode := Classic1v1

	if teams := mode.CalculateTeamCount(2); teams != 2 {
		t.Errorf("Expected 2 teams for 2 players, got %d", teams)
	}

	if teams := mode.CalculateTeamCount(4); teams != 4 {
		t.Errorf("Expected 4 teams for 4 players, got %d", teams)
	}

	teamMode := GameMode{
		Name:         "team_mode",
		MaxPlayers:   6,
		MinPlayers:   4,
		TeamSize:     2,
		TurnTime:     time.Second * 30,
		Speed:        1.0,
		MovesPerTurn: 3,
		Description:  "Team Mode",
	}

	if teams := teamMode.CalculateTeamCount(6); teams != 3 {
		t.Errorf("Expected 3 teams for 6 players with team size 2, got %d", teams)
	}

	if teams := teamMode.CalculateTeamCount(4); teams != 2 {
		t.Errorf("Expected 2 teams for 4 players with team size 2, got %d", teams)
	}

	if teams := teamMode.CalculateTeamCount(0); teams != 0 {
		t.Errorf("Expected 0 teams for 0 players, got %d", teams)
	}
}

// TestGameMode_Registry 测试游戏模式注册表
func TestGameMode_Registry(t *testing.T) {
	allModes := GetAllGameModes()

	if len(allModes) < 2 {
		t.Errorf("Expected at least 2 registered modes, got %d", len(allModes))
	}

	classic, exists := GetGameMode("classic_1v1")
	if !exists {
		t.Error("Expected classic_1v1 mode to be registered")
	}

	if classic.Name != "classic_1v1" {
		t.Errorf("Expected retrieved mode name to be 'classic_1v1', got '%s'", classic.Name)
	}

	testMode, exists := GetGameMode("test_mode")
	if !exists {
		t.Error("Expected test_mode to be registered")
	}

	if testMode.Name != "test_mode" {
		t.Errorf("Expected retrieved mode name to be 'test_mode', got '%s'", testMode.Name)
	}

	_, exists = GetGameMode("nonexistent")
	if exists {
		t.Error("Expected nonexistent mode to not be found")
	}

	customMode := GameMode{
		Name:         "custom",
		MaxPlayers:   4,
		MinPlayers:   2,
		TeamSize:     1,
		TurnTime:     time.Second * 15,
		Speed:        2.0,
		MovesPerTurn: 1,
		Description:  "Custom Mode",
	}

	RegisterGameMode(customMode)

	retrieved, exists := GetGameMode("custom")
	if !exists {
		t.Error("Expected custom mode to be registered")
	}

	if retrieved.Name != "custom" {
		t.Errorf("Expected custom mode name, got '%s'", retrieved.Name)
	}

	if retrieved.Speed != 2.0 {
		t.Errorf("Expected custom mode speed 2.0, got %f", retrieved.Speed)
	}

	allModesAfter := GetAllGameModes()
	if len(allModesAfter) != len(allModes)+1 {
		t.Errorf("Expected %d modes after registration, got %d", len(allModes)+1, len(allModesAfter))
	}

	allModesAfter["modified"] = GameMode{Name: "modified"}
	allModesCheck := GetAllGameModes()
	if len(allModesCheck) != len(allModes)+1 {
		t.Error("External modification affected internal registry")
	}
}

// TestGameMode_CustomModes 测试自定义模式
func TestGameMode_CustomModes(t *testing.T) {
	freeForAllMode := GameMode{
		Name:         "free_for_all",
		MaxPlayers:   8,
		MinPlayers:   3,
		TeamSize:     1,
		TurnTime:     time.Second * 20,
		Speed:        1.5,
		MovesPerTurn: 3,
		Description:  "Free for All",
	}

	if !freeForAllMode.ValidatePlayerCount(5) {
		t.Error("Expected 5 players to be valid for free for all")
	}

	if freeForAllMode.ValidatePlayerCount(2) {
		t.Error("Expected 2 players to be invalid for free for all")
	}

	if freeForAllMode.ValidatePlayerCount(9) {
		t.Error("Expected 9 players to be invalid for free for all")
	}

	if teams := freeForAllMode.CalculateTeamCount(6); teams != 6 {
		t.Errorf("Expected 6 teams for 6 players, got %d", teams)
	}

	teamMode := GameMode{
		Name:         "team_battle",
		MaxPlayers:   9,
		MinPlayers:   6,
		TeamSize:     3,
		TurnTime:     time.Second * 45,
		Speed:        0.8,
		MovesPerTurn: 2,
		Description:  "Team Battle",
	}

	if !teamMode.ValidatePlayerCount(6) {
		t.Error("Expected 6 players to be valid for team battle")
	}

	if !teamMode.ValidatePlayerCount(9) {
		t.Error("Expected 9 players to be valid for team battle")
	}

	if teams := teamMode.CalculateTeamCount(9); teams != 3 {
		t.Errorf("Expected 3 teams for 9 players, got %d", teams)
	}

	flexMode := GameMode{
		Name:         "flex",
		MaxPlayers:   10,
		MinPlayers:   2,
		TeamSize:     1,
		TurnTime:     time.Second * 10,
		Speed:        3.0,
		MovesPerTurn: 1,
		Description:  "Flexible Mode",
	}

	for i := 2; i <= 10; i++ {
		if !flexMode.ValidatePlayerCount(i) {
			t.Errorf("Expected %d players to be valid for flex mode", i)
		}
	}

	if flexMode.ValidatePlayerCount(1) {
		t.Error("Expected 1 player to be invalid for flex mode")
	}

	if flexMode.ValidatePlayerCount(11) {
		t.Error("Expected 11 players to be invalid for flex mode")
	}
}

// TestGameMode_EdgeCases 测试边界情况
func TestGameMode_EdgeCases(t *testing.T) {
	zeroTeamSizeMode := GameMode{
		Name:         "zero_team",
		MaxPlayers:   5,
		MinPlayers:   2,
		TeamSize:     0,
		TurnTime:     time.Second,
		Speed:        1.0,
		MovesPerTurn: 1,
		Description:  "Zero Team Size",
	}

	if teams := zeroTeamSizeMode.CalculateTeamCount(4); teams != 4 {
		t.Errorf("Expected 4 teams for 4 players with team size 0, got %d", teams)
	}

	largeTeamSizeMode := GameMode{
		Name:         "large_team",
		MaxPlayers:   5,
		MinPlayers:   3,
		TeamSize:     7,
		TurnTime:     time.Second,
		Speed:        1.0,
		MovesPerTurn: 1,
		Description:  "Large Team Size",
	}

	if teams := largeTeamSizeMode.CalculateTeamCount(5); teams != 0 {
		t.Errorf("Expected 0 teams for 5 players with team size 7, got %d", teams)
	}

	mode := GameMode{
		Name:         "boundary_test",
		MaxPlayers:   5,
		MinPlayers:   3,
		TeamSize:     1,
		TurnTime:     time.Second,
		Speed:        1.0,
		MovesPerTurn: 1,
		Description:  "Boundary Test",
	}

	if !mode.ValidatePlayerCount(3) {
		t.Error("Expected minimum player count to be valid")
	}
	if !mode.ValidatePlayerCount(5) {
		t.Error("Expected maximum player count to be valid")
	}
	if mode.ValidatePlayerCount(2) {
		t.Error("Expected below minimum to be invalid")
	}
	if mode.ValidatePlayerCount(6) {
		t.Error("Expected above maximum to be invalid")
	}
}

func TestSpeedConfiguration(t *testing.T) {
	t.Run("basic_speed_calculation", func(t *testing.T) {
		mode := Classic1v1
		mode.TurnTime = time.Second
		mode.Speed = 1.0

		if mode.GetTurnTime() != time.Second {
			t.Errorf("Expected 1 second for 1x speed, got %v", mode.GetTurnTime())
		}

		mode.Speed = 2.0
		if mode.GetTurnTime() != 500*time.Millisecond {
			t.Errorf("Expected 500ms for 2x speed, got %v", mode.GetTurnTime())
		}

		mode.Speed = 4.0
		if mode.GetTurnTime() != 250*time.Millisecond {
			t.Errorf("Expected 250ms for 4x speed, got %v", mode.GetTurnTime())
		}
	})

	t.Run("speed_variants", func(t *testing.T) {
		variants := CreateSpeedVariants()

		if len(variants) != 3 {
			t.Errorf("Expected 3 speed variants, got %d", len(variants))
		}

		mode1x, exists := variants["classic_1v1_1x"]
		if !exists {
			t.Error("Expected 1x speed variant to exist")
		}
		if mode1x.Speed != 1.0 {
			t.Errorf("Expected 1x variant to have speed 1.0, got %f", mode1x.Speed)
		}

		mode2x, exists := variants["classic_1v1_2x"]
		if !exists {
			t.Error("Expected 2x speed variant to exist")
		}
		if mode2x.Speed != 2.0 {
			t.Errorf("Expected 2x variant to have speed 2.0, got %f", mode2x.Speed)
		}

		mode4x, exists := variants["classic_1v1_4x"]
		if !exists {
			t.Error("Expected 4x speed variant to exist")
		}
		if mode4x.Speed != 4.0 {
			t.Errorf("Expected 4x variant to have speed 4.0, got %f", mode4x.Speed)
		}
	})

	t.Run("dynamic_speed_changes", func(t *testing.T) {
		mode := Classic1v1
		mode.TurnTime = time.Second
		mode.Speed = 1.0

		mode.SetSpeed(3.0)
		if mode.Speed != 3.0 {
			t.Errorf("Expected speed to be updated to 3.0, got %f", mode.Speed)
		}

		expectedTime := time.Duration(float64(time.Second.Nanoseconds()) / 3.0)
		if mode.GetTurnTime() != expectedTime {
			t.Errorf("Expected turn time %v for 3x speed, got %v", expectedTime, mode.GetTurnTime())
		}

		mode.SetSpeed(0)
		if mode.Speed == 0 {
			t.Error("Expected speed 0 to be rejected")
		}

		mode.SetSpeed(-1.0)
		if mode.Speed < 0 {
			t.Error("Expected negative speed to be rejected")
		}
	})

	t.Run("custom_speed_mode", func(t *testing.T) {
		customMode := NewClassic1v1WithSpeed(1.5)
		if customMode.Speed != 1.5 {
			t.Errorf("Expected custom mode speed 1.5, got %f", customMode.Speed)
		}

		expectedTime := time.Duration(float64(time.Second.Nanoseconds()) / 1.5)
		if customMode.GetTurnTime() != expectedTime {
			t.Errorf("Expected turn time %v for 1.5x speed, got %v", expectedTime, customMode.GetTurnTime())
		}
	})
}
