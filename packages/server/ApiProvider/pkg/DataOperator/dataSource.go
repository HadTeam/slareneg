package DataOperator

import (
	"server/ApiProvider/pkg/InstructionType"
	"server/GameJudge/pkg/GameType"
	"server/MapHandler/pkg/MapType"
)

type DataSource interface {
	GetOriginalMap(mapId uint32) *MapType.Map
	GetCurrentGame(id GameType.GameId) *GameType.Game
	GetCurrentInstruction(id GameType.GameId) []InstructionType.Instruction
	CreateGame(game *GameType.Game) GameType.GameId
	PutInstructions(id GameType.GameId, instructions []InstructionType.Instruction) bool
	AnnounceGameStart(gameId GameType.GameId) bool
	NewRound(id GameType.GameId, m *MapType.Map, turnNum uint8) bool
}
