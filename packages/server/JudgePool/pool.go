package JudgePool

import (
	"server/JudgePool/internal/Judge"
	"server/Utils/pkg/DataSource"
	"server/Utils/pkg/GameType"
	"sync"
	"time"
)

type Pool struct {
	judges        sync.Map
	AllowGameMode []GameType.GameMode
}

var data DataSource.TempDataSource

func ApplyDataSource(source interface{}) {
	data = source.(DataSource.TempDataSource)
	Judge.ApplyDataSource(source)

}
func (p *Pool) NewGame(mode GameType.GameMode) {
	id := data.CreateGame(mode)
	if id == 0 {
		panic("Cannot create game")
	}
	p.judges.Store(id, Judge.NewGameJudge(id))
}

func (p *Pool) DebugNewGame(g *GameType.Game) {
	if ok := data.DebugCreateGame(g); !ok {
		panic("cannot create game in debug mode")
	}
	if g.Id == 0 {
		panic("Cannot create game")
	}
	p.judges.Store(g.Id, Judge.NewGameJudge(g.Id))
}

func CreatePool(allowGameMode []GameType.GameMode) *Pool {
	p := &Pool{AllowGameMode: allowGameMode}
	go poolWorking(p)
	return p
}

func poolWorking(p *Pool) {
	t := time.NewTicker(100 * time.Millisecond)
	for _, mode := range p.AllowGameMode {
		p.NewGame(mode)
	}
	for _ = range t.C {
		// Ensure there is a game always in waiting status
		tryStartGame := func(game GameType.Game) {

			if uint8(len(data.GetCurrentUserList(game.Id))) == game.Mode.MaxUserNum {
				jAny, _ := p.judges.Load(game.Id)
				j := jAny.(*Judge.GameJudge)
				j.StartGame()
				p.NewGame(game.Mode)
			}
		}

		for _, mode := range p.AllowGameMode {
			list := data.GetGameList(mode)
			for _, g := range list {
				if g.Status == GameType.GameStatusWaiting {
					tryStartGame(g)
					break
				}
			}
		}
	}
}
