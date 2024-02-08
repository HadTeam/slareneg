package pool

import (
	"github.com/sirupsen/logrus"
	"server/game"
	"server/game/instruction"
	"server/game/judge"
	"server/game/mode"
	"sync"
	"time"
)

type Pool struct {
	games         sync.Map
	AllowGameMode []mode.Mode
}

func (p *Pool) NewGame(mode mode.Mode) *game.Game {
	g := game.Create(mode)
	l := game.Data.(*game.Local)
	l.GamePool[g.Id] = g
	l.InstructionLog[g.Id] = make(map[uint16]map[uint16]instruction.Instruction)
	if g.Id == 0 {
		logrus.Panic("cannot create game")
	}
	p.games.Store(g.Id, judge.NewGameJudge(judge.Id(g.Id)))
	return g
}

func (p *Pool) DebugNewGame(g *game.Game) {
	if ok := game.Data.DebugCreateGame(g); !ok {
		logrus.Panic("cannot create game in debug mode")
	}
	if g.Id == 0 {
		logrus.Panic("Cannot create game")
	}
	p.games.Store(g.Id, judge.NewGameJudge(judge.Id(g.Id)))
}

func CreatePool(allowGameMode []mode.Mode) *Pool {
	p := &Pool{AllowGameMode: allowGameMode}
	go poolWorking(p)
	return p
}

func poolWorking(p *Pool) {
	t := time.NewTicker(100 * time.Millisecond)
	for _, mode := range p.AllowGameMode {
		p.NewGame(mode)
	}
	for range t.C {
		// Ensure there is a game always in waiting status
		tryStartGame := func(g game.Game) {
			if uint8(len(game.Data.GetCurrentUserList(g.Id))) == g.Mode.MaxUserNum {
				//jAny, _ := p.games.Load(game.Id)
				//j := jAny.(*judge.Judge)
				g.Start()
				go g.MainLoop()
			}
		}

		for _, mode := range p.AllowGameMode {
			list := game.Data.GetGameList(mode)
			for _, g := range list {
				if g.Status == game.StatusWaiting {
					tryStartGame(g)
					break
				}
			}
		}
	}
}
