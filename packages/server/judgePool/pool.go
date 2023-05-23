package judgePool

import (
	"server/judgePool/internal/judge"
	"server/utils/pkg/dataSource"
	"server/utils/pkg/game"
	"sync"
	"time"
)

type Pool struct {
	judges        sync.Map
	AllowGameMode []game.GameMode
}

var data dataSource.TempDataSource

func ApplyDataSource(source interface{}) {
	data = source.(dataSource.TempDataSource)
	judge.ApplyDataSource(source)

}
func (p *Pool) NewGame(mode game.GameMode) {
	id := data.CreateGame(mode)
	if id == 0 {
		panic("Cannot create game")
	}
	p.judges.Store(id, judge.NewGameJudge(id))
}

func (p *Pool) DebugNewGame(g *game.Game) {
	if ok := data.DebugCreateGame(g); !ok {
		panic("cannot create game in debug mode")
	}
	if g.Id == 0 {
		panic("Cannot create game")
	}
	p.judges.Store(g.Id, judge.NewGameJudge(g.Id))
}

func CreatePool(allowGameMode []game.GameMode) *Pool {
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
		tryStartGame := func(game game.Game) {

			if uint8(len(data.GetCurrentUserList(game.Id))) == game.Mode.MaxUserNum {
				jAny, _ := p.judges.Load(game.Id)
				j := jAny.(*judge.GameJudge)
				j.StartGame()
				p.NewGame(game.Mode)
			}
		}

		for _, mode := range p.AllowGameMode {
			list := data.GetGameList(mode)
			for _, g := range list {
				if g.Status == game.GameStatusWaiting {
					tryStartGame(g)
					break
				}
			}
		}
	}
}
