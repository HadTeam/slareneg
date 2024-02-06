package data_source

import (
	"server/game_logic"
	"server/game_logic/game_def"
	"server/game_logic/map"
)

type TempDataSource interface {
	CreateGame(mode game_def.Mode) game_logic.Id
	GetGameList(mode game_def.Mode) []game_logic.Game // Returns `Game` structs with only basic info of not end game
	CancelGame(id game_logic.Id) (ok bool)

	GetCurrentUserList(id game_logic.Id) []game_def.User
	GetInstructions(id game_logic.Id, tempId uint16) map[uint16]game_def.Instruction
	GetGameInfo(id game_logic.Id) *game_logic.Game // Returns a `Game` struct with only basic info
	NewInstructionTemp(id game_logic.Id, tempId uint16) (ok bool)
	SetGameStatus(id game_logic.Id, status game_logic.Status) (ok bool)
	SetGameMap(id game_logic.Id, m *_map.Map) (ok bool)

	SetUserStatus(id game_logic.Id, user game_def.User) (ok bool)
	SetWinner(id game_logic.Id, teamId uint8) (ok bool)
	UpdateInstruction(id game_logic.Id, user game_def.User, instruction game_def.Instruction) (ok bool)
	GetCurrentMap(id game_logic.Id) *_map.Map

	// Functions only for debug
	DebugCreateGame(game *game_logic.Game) (ok bool)
}
