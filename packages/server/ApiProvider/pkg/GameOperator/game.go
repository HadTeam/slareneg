package GameOperator

import (
	"server/Untils/pkg/GameType"
	"time"
)

// NewGame TODO: Add unit test
func NewGame(mapId uint32, mode GameType.GameMode) GameType.GameId {
	m := data.GetOriginalMap(mapId)
	g := &GameType.Game{
		Map:        m,
		UserList:   []GameType.User{},
		CreateTime: time.Now().UnixMicro(),
		Status:     GameType.GameStatusWaiting,
		RoundNum:   0,
		Mode:       mode,
		// The Id field will be filled in `data.CreateGame`
	}
	return data.CreateGame(g)
}

func StartGame(id GameType.GameId) {
	g := data.GetCurrentGame(id)
	g.Status = GameType.GameStatusRunning
	// TODO: Announce
}

func TryForceStart(id GameType.GameId) {
	// TODO
}
