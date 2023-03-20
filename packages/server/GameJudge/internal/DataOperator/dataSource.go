package DataOperator

import (
	"server/ApiProvider/pkg/InstructionType"
	"server/GameJudge/internal/GameType"
	"server/MapHandler/pkg/MapType"
)

type DataSource interface {
	GetCurrentGame(id GameType.GameId) *GameType.Game
	GetCurrentInstruction(id GameType.GameId) []InstructionType.Instruction
	NewTurn(id GameType.GameId, m *MapType.Map, turnNum uint8) bool
}
