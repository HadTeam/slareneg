package game

import (
	"testing"
	"time"
)

// TestGameMode_Classic1v1 测试经典1v1模式
func TestGameMode_Classic1v1(t *testing.T) {
	mode := Classic1v1

	t.Run("mode_properties", func(t *testing.T) {
		if mode.Name != "classic_1v1" {
			t.Errorf("Expected name 'classic_1v1', got %s", mode.Name)
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

		if mode.TurnTime != time.Microsecond*1000 {
			t.Errorf("Expected turn time %v, got %v", time.Microsecond*1000, mode.TurnTime)
		}

		if mode.Description == "" {
			t.Error("Expected non-empty description")
		}
	})
}

// TestGameMode_ValidatePlayerCount 测试玩家数量验证
func TestGameMode_ValidatePlayerCount(t *testing.T) {
	mode := Classic1v1

	t.Run("valid_player_counts", func(t *testing.T) {
		// 2个玩家应该有效
		if !mode.ValidatePlayerCount(2) {
			t.Error("Expected 2 players to be valid for Classic1v1")
		}
	})

	t.Run("invalid_player_counts", func(t *testing.T) {
		// 1个玩家应该无效
		if mode.ValidatePlayerCount(1) {
			t.Error("Expected 1 player to be invalid for Classic1v1")
		}

		// 3个玩家应该无效
		if mode.ValidatePlayerCount(3) {
			t.Error("Expected 3 players to be invalid for Classic1v1")
		}

		// 0个玩家应该无效
		if mode.ValidatePlayerCount(0) {
			t.Error("Expected 0 players to be invalid")
		}

		// 负数玩家应该无效
		if mode.ValidatePlayerCount(-1) {
			t.Error("Expected negative player count to be invalid")
		}
	})
}

// TestGameMode_CalculateTeamCount 测试团队计数计算
func TestGameMode_CalculateTeamCount(t *testing.T) {
	t.Run("classic_1v1_teams", func(t *testing.T) {
		mode := Classic1v1

		// 每个玩家单独一队
		teamCount := mode.CalculateTeamCount(2)
		if teamCount != 2 {
			t.Errorf("Expected 2 teams for 2 players in 1v1, got %d", teamCount)
		}

		teamCount = mode.CalculateTeamCount(1)
		if teamCount != 1 {
			t.Errorf("Expected 1 team for 1 player in 1v1, got %d", teamCount)
		}
	})

	t.Run("team_mode", func(t *testing.T) {
		// 创建一个团队模式进行测试
		teamMode := GameMode{
			Name:       "team_2v2",
			MaxPlayers: 4,
			MinPlayers: 4,
			TeamSize:   2, // 每队2人
			TurnTime:   time.Second,
		}

		teamCount := teamMode.CalculateTeamCount(4)
		if teamCount != 2 {
			t.Errorf("Expected 2 teams for 4 players in 2v2, got %d", teamCount)
		}

		teamCount = teamMode.CalculateTeamCount(2)
		if teamCount != 1 {
			t.Errorf("Expected 1 team for 2 players in 2v2, got %d", teamCount)
		}
	})

	t.Run("edge_cases", func(t *testing.T) {
		mode := Classic1v1

		// 0个玩家
		teamCount := mode.CalculateTeamCount(0)
		if teamCount != 0 {
			t.Errorf("Expected 0 teams for 0 players, got %d", teamCount)
		}
	})
}

// TestGameMode_Registry 测试游戏模式注册表
func TestGameMode_Registry(t *testing.T) {
	t.Run("get_classic_1v1", func(t *testing.T) {
		mode, exists := GetGameMode("classic_1v1")
		if !exists {
			t.Error("Expected classic_1v1 mode to exist")
		}

		if mode.Name != "classic_1v1" {
			t.Errorf("Expected mode name 'classic_1v1', got %s", mode.Name)
		}
	})

	t.Run("get_nonexistent_mode", func(t *testing.T) {
		_, exists := GetGameMode("nonexistent_mode")
		if exists {
			t.Error("Expected nonexistent mode to not exist")
		}
	})

	t.Run("get_all_modes", func(t *testing.T) {
		modes := GetAllGameModes()

		if len(modes) == 0 {
			t.Error("Expected at least one registered mode")
		}

		if _, exists := modes["classic_1v1"]; !exists {
			t.Error("Expected classic_1v1 to be in all modes")
		}
	})

	t.Run("register_new_mode", func(t *testing.T) {
		newMode := GameMode{
			Name:        "test_mode",
			MaxPlayers:  4,
			MinPlayers:  2,
			TeamSize:    2,
			TurnTime:    time.Second * 30,
			Description: "Test mode for unit tests",
		}

		RegisterGameMode(newMode)

		retrieved, exists := GetGameMode("test_mode")
		if !exists {
			t.Error("Expected registered mode to exist")
		}

		if retrieved.Name != "test_mode" {
			t.Errorf("Expected retrieved mode name 'test_mode', got %s", retrieved.Name)
		}

		if retrieved.MaxPlayers != 4 {
			t.Errorf("Expected max players 4, got %d", retrieved.MaxPlayers)
		}
	})

	t.Run("registry_immutability", func(t *testing.T) {
		originalModes := GetAllGameModes()
		originalCount := len(originalModes)

		// 修改返回的map不应该影响内部注册表
		originalModes["malicious_mode"] = GameMode{Name: "malicious"}

		newModes := GetAllGameModes()
		if len(newModes) != originalCount {
			t.Error("Registry should not be affected by external modifications")
		}

		if _, exists := newModes["malicious_mode"]; exists {
			t.Error("Malicious modification should not affect registry")
		}
	})
}

// TestGameMode_CustomModes 测试自定义模式
func TestGameMode_CustomModes(t *testing.T) {
	t.Run("ffa_mode", func(t *testing.T) {
		// 测试自由混战模式
		ffaMode := GameMode{
			Name:        "ffa_4p",
			MaxPlayers:  4,
			MinPlayers:  3,
			TeamSize:    1, // 每人一队
			TurnTime:    time.Second * 15,
			Description: "4-player free-for-all",
		}

		// 验证玩家数量
		if !ffaMode.ValidatePlayerCount(3) {
			t.Error("Expected 3 players to be valid for FFA")
		}
		if !ffaMode.ValidatePlayerCount(4) {
			t.Error("Expected 4 players to be valid for FFA")
		}
		if ffaMode.ValidatePlayerCount(2) {
			t.Error("Expected 2 players to be invalid for FFA")
		}
		if ffaMode.ValidatePlayerCount(5) {
			t.Error("Expected 5 players to be invalid for FFA")
		}

		// 验证团队计数
		if ffaMode.CalculateTeamCount(4) != 4 {
			t.Errorf("Expected 4 teams for 4-player FFA, got %d", ffaMode.CalculateTeamCount(4))
		}
	})

	t.Run("team_mode", func(t *testing.T) {
		// 测试团队模式
		teamMode := GameMode{
			Name:        "team_3v3",
			MaxPlayers:  6,
			MinPlayers:  6,
			TeamSize:    3, // 每队3人
			TurnTime:    time.Second * 20,
			Description: "3v3 team battle",
		}

		// 验证玩家数量
		if !teamMode.ValidatePlayerCount(6) {
			t.Error("Expected 6 players to be valid for 3v3")
		}
		if teamMode.ValidatePlayerCount(5) {
			t.Error("Expected 5 players to be invalid for 3v3")
		}

		// 验证团队计数
		if teamMode.CalculateTeamCount(6) != 2 {
			t.Errorf("Expected 2 teams for 6-player 3v3, got %d", teamMode.CalculateTeamCount(6))
		}
	})

	t.Run("flexible_mode", func(t *testing.T) {
		// 测试灵活人数模式
		flexMode := GameMode{
			Name:        "flexible",
			MaxPlayers:  8,
			MinPlayers:  2,
			TeamSize:    1,
			TurnTime:    time.Second * 10,
			Description: "Flexible player count mode",
		}

		// 验证各种玩家数量
		for i := 2; i <= 8; i++ {
			if !flexMode.ValidatePlayerCount(i) {
				t.Errorf("Expected %d players to be valid for flexible mode", i)
			}
		}

		if flexMode.ValidatePlayerCount(1) {
			t.Error("Expected 1 player to be invalid for flexible mode")
		}
		if flexMode.ValidatePlayerCount(9) {
			t.Error("Expected 9 players to be invalid for flexible mode")
		}
	})
}

// TestGameMode_EdgeCases 测试边界情况
func TestGameMode_EdgeCases(t *testing.T) {
	t.Run("zero_team_size", func(t *testing.T) {
		zeroTeamMode := GameMode{
			Name:        "zero_team",
			MaxPlayers:  4,
			MinPlayers:  2,
			TeamSize:    0, // 特殊情况：0表示每人一队
			TurnTime:    time.Second,
			Description: "Zero team size mode",
		}

		// 应该按每人一队处理
		teamCount := zeroTeamMode.CalculateTeamCount(4)
		if teamCount != 4 {
			t.Errorf("Expected 4 teams for zero team size, got %d", teamCount)
		}
	})

	t.Run("large_team_size", func(t *testing.T) {
		largeTeamMode := GameMode{
			Name:        "large_team",
			MaxPlayers:  10,
			MinPlayers:  5,
			TeamSize:    7, // 团队大小大于玩家数
			TurnTime:    time.Second,
			Description: "Large team size mode",
		}

		// 5个玩家分成7人一队，应该是0队
		teamCount := largeTeamMode.CalculateTeamCount(5)
		if teamCount != 0 {
			t.Errorf("Expected 0 teams when team size > player count, got %d", teamCount)
		}
	})

	t.Run("boundary_player_counts", func(t *testing.T) {
		mode := GameMode{
			Name:        "boundary_test",
			MaxPlayers:  5,
			MinPlayers:  3,
			TeamSize:    1,
			TurnTime:    time.Second,
			Description: "Boundary test mode",
		}

		// 测试边界值
		if !mode.ValidatePlayerCount(3) { // 最小值
			t.Error("Expected min player count to be valid")
		}
		if !mode.ValidatePlayerCount(5) { // 最大值
			t.Error("Expected max player count to be valid")
		}
		if mode.ValidatePlayerCount(2) { // 小于最小值
			t.Error("Expected below-min player count to be invalid")
		}
		if mode.ValidatePlayerCount(6) { // 大于最大值
			t.Error("Expected above-max player count to be invalid")
		}
	})
}
