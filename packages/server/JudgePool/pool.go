package JudgePool

import (
	"math/rand"
	"server/JudgePool/internal/InstructionExecutor"
	"server/JudgePool/internal/Judge"
	"server/Untils/pkg/DataSource"
	"server/Untils/pkg/GameType"
	"sync"
)

type Pool struct {
	judges    sync.Map
	judgeChan chan GameType.GameId
}

var data DataSource.TempDataSource

func ApplyDataSource(source interface{}) {
	data = source.(DataSource.TempDataSource)
	Judge.ApplyDataSource(source)
	InstructionExecutor.ApplyDataSource(source)
}

func (p *Pool) NewGame(mode GameType.GameMode) {
	id := data.CreateGame(mode)
	if id == 0 {
		panic("Cannot create game")
	}
	jId := rand.Uint32()

	for {
		if _, ok := p.judges.Load(jId); !ok {
			break
		}
		jId = rand.Uint32()
	}

	p.judges.Store(jId, Judge.NewGameJudge(id))
}

func CreatePool() *Pool {
	return &Pool{}
}
