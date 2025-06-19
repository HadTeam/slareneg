package game

import (
	"server/game/instruction"
	"server/game/map"
	"server/game/mode"
	"server/game/user"
)

type TempDataSource interface {
	GetGameList(mode mode.Mode) []Game // Returns `Game` structs with only basic info of not end game
	CancelGame(id Id) (ok bool)

	GetCurrentUserList(id Id) []user.User
	GetInstructions(id Id, tempId uint16) map[uint16]instruction.Instruction
	GetGameInfo(id Id) *Game // Returns a `Game` struct with only basic info
	NewInstructionTemp(id Id, tempId uint16) (ok bool)
	SetGameStatus(id Id, status Status) (ok bool)
	SetGameMap(id Id, m *_map.Map) (ok bool)

	SetUserStatus(id Id, user user.User) (ok bool)
	SetWinner(id Id, teamId uint8) (ok bool)
	UpdateInstruction(id Id, user user.User, instruction instruction.Instruction) (ok bool)
	GetCurrentMap(id Id) *_map.Map

	// Functions only for debug
	DebugCreateGame(game *Game) (ok bool)
}
