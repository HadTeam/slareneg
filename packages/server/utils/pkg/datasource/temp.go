package datasource

import (
	"server/utils/pkg/game"
	"server/utils/pkg/instruction"
	"server/utils/pkg/map"
)

type TempDataSource interface {
	CreateGame(mode game.Mode) game.Id
	GetGameList(mode game.Mode) []game.Game // Returns `Game` structs with only basic info of not end game
	CancelGame(id game.Id) (ok bool)

	GetCurrentUserList(id game.Id) []game.User
	GetInstructions(id game.Id, tempId uint16) []instruction.Instruction
	GetGameInfo(id game.Id) *game.Game // Returns a `Game` struct with only basic info
	NewInstructionTemp(id game.Id, tempId uint16) (ok bool)
	SetGameStatus(id game.Id, status game.Status) (ok bool)
	SetGameMap(id game.Id, m *_map.Map) (ok bool)

	SetUserStatus(id game.Id, user game.User) (ok bool)
	SetWinner(id game.Id, teamId uint8) (ok bool)
	UpdateInstruction(id game.Id, user game.User, instruction instruction.Instruction) (ok bool)
	GetCurrentMap(id game.Id) *_map.Map

	// Functions only for debug
	DebugCreateGame(game *game.Game) (ok bool)
}
