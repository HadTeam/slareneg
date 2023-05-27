package judgePool

import (
	"github.com/sirupsen/logrus"
	"server/judgepool/internal/judge"
	data_source "server/utils/pkg/datasource"
	"server/utils/pkg/game"
	"sync"
	"time"
)

type Pool struct {
	judges        sync.Map
	AllowGameMode []game.Mode
}

var data data_source.TempDataSource

func ApplyDataSource(source interface{}) {
	data = source.(data_source.TempDataSource)
	judge.ApplyDataSource(source)

}
func (p *Pool) NewGame(mode game.Mode) {
	id := data.CreateGame(mode)
	if id == 0 {
		logrus.Panic("cannot create game")
	}
	p.judges.Store(id, judge.NewGameJudge(id))
}

func (p *Pool) DebugNewGame(g *game.Game) {
	if ok := data.DebugCreateGame(g); !ok {
		logrus.Panic("cannot create game in debug mode")
	}
	if g.Id == 0 {
		logrus.Panic("Cannot create game")
	}
	p.judges.Store(g.Id, judge.NewGameJudge(g.Id))
}

func CreatePool(allowGameMode []game.Mode) *Pool {
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
				if g.Status == game.StatusWaiting {
					tryStartGame(g)
					break
				}
			}
		}
	}
}
