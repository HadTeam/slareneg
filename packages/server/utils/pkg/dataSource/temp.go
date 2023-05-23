package dataSource

import (
	"server/utils/pkg/game"
	"server/utils/pkg/instruction"
	"server/utils/pkg/map"
)

type TempDataSource interface {
	CreateGame(mode game.GameMode) game.GameId
	GetGameList(mode game.GameMode) []game.Game // Returns `Game` structs with only basic info of not end game
	CancelGame(id game.GameId) (ok bool)

	GetCurrentUserList(id game.GameId) []game.User
	GetInstructions(id game.GameId, tempId uint16) []instruction.Instruction
	GetGameInfo(id game.GameId) *game.Game // Returns a `Game` struct with only basic info
	NewInstructionTemp(id game.GameId, tempId uint16) (ok bool)
	SetGameStatus(id game.GameId, status game.GameStatus) (ok bool)
	SetGameMap(id game.GameId, m *_map.Map) (ok bool)

	SetUserStatus(id game.GameId, user game.User) (ok bool)
	UpdateInstruction(id game.GameId, user game.User, instruction instruction.Instruction) (ok bool)
	GetCurrentMap(id game.GameId) *_map.Map

	// Functions only for debug
	DebugCreateGame(game *game.Game) (ok bool)
}
