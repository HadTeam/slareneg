package game

import (
	"server/internal/game/block"
	gamemap "server/internal/game/map"
	"time"
)

type GameEndCondition interface {
	Check(players []Player, gameMap gamemap.Map) (isOver bool, winners []string, reason string)
	Name() string
}

type LastPlayerStandingCondition struct{}

func (c *LastPlayerStandingCondition) Name() string {
	return "last_player_standing"
}

func (c *LastPlayerStandingCondition) Check(players []Player, gameMap gamemap.Map) (bool, []string, string) {
	if gameMap == nil {
		return false, nil, ""
	}

	playerCastles := make(map[block.Owner]int)
	playerHasUnits := make(map[block.Owner]bool)

	size := gameMap.Size()
	for y := uint16(1); y <= size.Height; y++ {
		for x := uint16(1); x <= size.Width; x++ {
			pos := gamemap.Pos{X: x, Y: y}
			b, err := gameMap.Block(pos)
			if err != nil {
				continue
			}

			if b.Owner() != block.Owner(0) { // 不是中立
				if b.Meta().Name == block.CastleName || b.Meta().Name == block.KingName {
					playerCastles[b.Owner()]++
				}
				if b.Num() > 0 {
					playerHasUnits[b.Owner()] = true
				}
			}
		}
	}

	var activePlayers []string
	for i, player := range players {
		if player.Status == PlayerStatusInGame {
			owner := block.Owner(i)
			if playerCastles[owner] > 0 || playerHasUnits[owner] {
				activePlayers = append(activePlayers, player.Id)
			}
		}
	}

	if len(activePlayers) <= 1 {
		reason := "elimination"
		if len(activePlayers) == 1 {
			reason = "last_player_standing"
		} else {
			reason = "all_players_eliminated"
		}
		return true, activePlayers, reason
	}

	return false, nil, ""
}

type GameMode struct {
	Name         string
	MaxPlayers   uint8
	MinPlayers   uint8
	TeamSize     uint8
	TurnTime     time.Duration
	Speed        float64
	MovesPerTurn uint16
	Description  string

	EndConditions []GameEndCondition
}

func (gm GameMode) GetTurnTime() time.Duration {
	if gm.Speed <= 0 {
		return gm.TurnTime
	}
	return time.Duration(float64(gm.TurnTime) / gm.Speed)
}

func (gm *GameMode) SetSpeed(speed float64) {
	if speed > 0 {
		gm.Speed = speed
	}
}

func (gm GameMode) CheckGameEnd(players []Player, gameMap gamemap.Map) (isOver bool, winners []string, reason string) {
	for _, condition := range gm.EndConditions {
		if isOver, winners, reason := condition.Check(players, gameMap); isOver {
			return true, winners, reason
		}
	}
	return false, nil, ""
}

var (
	Classic1v1 = GameMode{
		Name:         "classic_1v1",
		MaxPlayers:   2,
		MinPlayers:   2,
		TeamSize:     1,
		TurnTime:     time.Second,
		Speed:        1.0,
		MovesPerTurn: 2,
		Description:  "经典1对1对战模式",
		EndConditions: []GameEndCondition{
			&LastPlayerStandingCondition{},
		},
	}

	TestMode = GameMode{
		Name:         "test_mode",
		MaxPlayers:   2,
		MinPlayers:   2,
		TeamSize:     1,
		TurnTime:     time.Hour * 24,
		Speed:        1.0,
		MovesPerTurn: 2,
		Description:  "测试专用模式",
		EndConditions: []GameEndCondition{
			&LastPlayerStandingCondition{},
		},
	}
)

var registeredModes = map[string]GameMode{
	Classic1v1.Name: Classic1v1,
	TestMode.Name:   TestMode,
}

func RegisterGameMode(mode GameMode) {
	registeredModes[mode.Name] = mode
}

func GetGameMode(name string) (GameMode, bool) {
	mode, exists := registeredModes[name]
	return mode, exists
}

func GetAllGameModes() map[string]GameMode {
	result := make(map[string]GameMode)
	for k, v := range registeredModes {
		result[k] = v
	}
	return result
}

func (gm GameMode) ValidatePlayerCount(count int) bool {
	return count >= int(gm.MinPlayers) && count <= int(gm.MaxPlayers)
}

func (gm GameMode) CalculateTeamCount(playerCount int) int {
	if gm.TeamSize == 0 || gm.TeamSize == 1 {
		return playerCount
	}
	return playerCount / int(gm.TeamSize)
}

func NewClassic1v1WithSpeed(speed float64) GameMode {
	mode := Classic1v1
	mode.SetSpeed(speed)
	return mode
}

func CreateSpeedVariants() map[string]GameMode {
	return map[string]GameMode{
		"classic_1v1_1x": NewClassic1v1WithSpeed(1.0),
		"classic_1v1_2x": NewClassic1v1WithSpeed(2.0),
		"classic_1v1_4x": NewClassic1v1WithSpeed(4.0),
	}
}
