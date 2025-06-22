package game

import "time"

// GameMode 游戏模式定义
type GameMode struct {
	Name        string        // 模式名称
	MaxPlayers  uint8         // 最大玩家数
	MinPlayers  uint8         // 最小玩家数
	TeamSize    uint8         // 团队人数 (1表示每个玩家单独一队)
	TurnTime    time.Duration // 回合时间
	Description string        // 模式描述
}

// 预定义的游戏模式
var (
	// Classic1v1 经典1v1模式
	Classic1v1 = GameMode{
		Name:        "classic_1v1",
		MaxPlayers:  2,
		MinPlayers:  2,
		TeamSize:    1,
		TurnTime:    time.Microsecond * 1000,
		Description: "经典1对1对战模式",
	}
)

// 游戏模式注册表
var registeredModes = map[string]GameMode{
	Classic1v1.Name: Classic1v1,
}

// RegisterGameMode 注册新的游戏模式
func RegisterGameMode(mode GameMode) {
	registeredModes[mode.Name] = mode
}

// GetGameMode 根据名称获取游戏模式
func GetGameMode(name string) (GameMode, bool) {
	mode, exists := registeredModes[name]
	return mode, exists
}

// GetAllGameModes 获取所有已注册的游戏模式
func GetAllGameModes() map[string]GameMode {
	// 返回副本，防止外部修改
	result := make(map[string]GameMode)
	for k, v := range registeredModes {
		result[k] = v
	}
	return result
}

// ValidatePlayerCount 验证玩家数量是否符合模式要求
func (gm GameMode) ValidatePlayerCount(count int) bool {
	return count >= int(gm.MinPlayers) && count <= int(gm.MaxPlayers)
}

// CalculateTeamCount 计算团队数量
func (gm GameMode) CalculateTeamCount(playerCount int) int {
	if gm.TeamSize == 0 || gm.TeamSize == 1 {
		// 每个玩家单独一队
		return playerCount
	}
	return playerCount / int(gm.TeamSize)
}
