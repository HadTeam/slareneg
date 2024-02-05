package data_source

import (
	"server/game_logic"
	"server/game_logic/game_def"
	"server/game_logic/map"
)

type TempDataSource interface {
	CreateGame(mode _type.Mode) game_logic.Id
	GetGameList(mode _type.Mode) []game_logic.Game // Returns `Game` structs with only basic info of not end game
	CancelGame(id game_logic.Id) (ok bool)

	GetCurrentUserList(id game_logic.Id) []_type.User
	GetInstructions(id game_logic.Id, tempId uint16) map[uint16]_type.Instruction
	GetGameInfo(id game_logic.Id) *game_logic.Game // Returns a `Game` struct with only basic info
	NewInstructionTemp(id game_logic.Id, tempId uint16) (ok bool)
	SetGameStatus(id game_logic.Id, status game_logic.Status) (ok bool)
	SetGameMap(id game_logic.Id, m *_map.Map) (ok bool)

	SetUserStatus(id game_logic.Id, user _type.User) (ok bool)
	SetWinner(id game_logic.Id, teamId uint8) (ok bool)
	UpdateInstruction(id game_logic.Id, user _type.User, instruction _type.Instruction) (ok bool)
	GetCurrentMap(id game_logic.Id) *_map.Map

	// Functions only for debug
	DebugCreateGame(game *game_logic.Game) (ok bool)
}
