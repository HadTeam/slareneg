package DataOperator

import (
	"server/ApiProvider/pkg/InstructionType"
	"server/JudgePool/pkg/GameType"
	"server/MapHandler/pkg/MapType"
)

type DataSource interface {
	GetOriginalMap(mapId uint32) *MapType.Map
	GetCurrentGame(id GameType.GameId) *GameType.Game
	CreateGame(game *GameType.Game) GameType.GameId
	PutInstructions(id GameType.GameId, instructions []InstructionType.Instruction) bool
	AnnounceGameStart(gameId GameType.GameId) bool

	GetInstructionsFromTemp(id GameType.GameId, roundNum uint8) []InstructionType.Instruction
	AchieveInstructionTemp(id GameType.GameId, roundNum uint8) bool
	PutMap(id GameType.GameId, m *MapType.Map) bool
}
