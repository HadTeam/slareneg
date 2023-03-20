package DataOperator

import (
	"server/ApiProvider/pkg/InstructionType"
	"server/GameJudge/pkg/GameType"
	"server/MapHandler/pkg/MapType"
)

type DataSource interface {
	GetOriginalMap(mapId uint32) *MapType.Map
	GetCurrentGame(id GameType.GameId) *GameType.Game
	CreateGame(game *GameType.Game) GameType.GameId
	UpdateGame(game *GameType.Game) bool
	PutInstruction(instruction InstructionType.Instruction) bool
	AnnounceGameStart(gameId GameType.GameId) bool
}
