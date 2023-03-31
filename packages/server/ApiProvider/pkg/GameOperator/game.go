package GameOperator

import (
	GameType2 "server/Untils/pkg/GameType"
	"time"
)

// NewGame TODO: Add unit test
func NewGame(mapId uint32, mode GameType2.GameMode) GameType2.GameId {
	m := data.GetOriginalMap(mapId)
	g := &GameType2.Game{
		Map:        m,
		UserList:   []GameType2.User{},
		CreateTime: time.Now().UnixMicro(),
		Status:     GameType2.GameStatusWaiting,
		RoundNum:   0,
		Mode:       mode,
		// The Id field will be filled in `data.CreateGame`
	}
	return data.CreateGame(g)
}

func StartGame(id GameType2.GameId) {
	g := data.GetCurrentGame(id)
	g.Status = GameType2.GameStatusRunning
	// TODO: Announce
}

func TryForceStart(id GameType2.GameId) {
	// TODO
}
