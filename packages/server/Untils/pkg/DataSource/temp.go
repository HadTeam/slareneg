package DataSource

import (
	"server/Untils/pkg/GameType"
	"server/Untils/pkg/InstructionType"
	"server/Untils/pkg/MapType"
)

type TempDataSource interface {
	CreateGame(mode GameType.GameMode) GameType.GameId

	GetCurrentUserList(id GameType.GameId) []GameType.User
	GetInstructions(id GameType.GameId, tempId uint8) []InstructionType.Instruction
	GetGameInfo(id GameType.GameId) *GameType.Game // Returns a `Game` struct with only basic info
	NewInstructionTemp(id GameType.GameId, tempId uint8) (ok bool)
	SetGameStatus(id GameType.GameId, status GameType.GameStatus) (ok bool)
	SetGameMap(id GameType.GameId, m *MapType.Map) (ok bool)

	SetUserStatus(id GameType.GameId, user GameType.User) (ok bool)
	UpdateInstruction(id GameType.GameId, user GameType.User, instruction InstructionType.Instruction) (ok bool)
	GetCurrentMap(id GameType.GameId) *MapType.Map
}
