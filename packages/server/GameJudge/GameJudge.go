package GameJudge

import (
	"fmt"
	"math/rand"
	"server/ApiProvider/pkg/DataOperator"
	"server/GameJudge/pkg/GameType"
	"time"
)

var RoundTime = time.Millisecond * 1000

var data DataOperator.DataSource

type Status uint8

const (
	StatusWaiting Status = iota + 1
	StatusWorking
)

type JudgeCommand uint8

const (
	JudgeCommandWork JudgeCommand = iota + 1
)

type GameJudge struct {
	gameId GameType.GameId
	status Status
	id     uint8
	m      chan JudgeCommand
}

func ApplyDataSource(source DataOperator.DataSource) {
	data = source
}

func NewGameJudge() *GameJudge {
	j := &GameJudge{
		gameId: 0,
		status: StatusWaiting,
		id:     uint8(rand.Uint32()),
		m:      make(chan JudgeCommand),
	}
	go judgeWorking(j)
	return j
}

func judgeWorking(j *GameJudge) {
	for v := range j.m {
		switch v {
		case JudgeCommandWork:
			{
				if j.gameId == 0 {
					break
				}
				j.status = StatusWorking
				t := time.NewTicker(RoundTime)
				game := data.GetCurrentGame(j.gameId)
				game.RoundNum = 0
				game.Map = data.GetOriginalMap(game.Map.MapId)
				data.PutMap(j.gameId, game.Map)
				for range t.C {
					//Round End
					if game.RoundNum != 0 {
						data.AchieveInstructionTemp(j.gameId, game.RoundNum)
						data.GetInstructionsFromTemp(j.gameId, game.RoundNum)
						game.Map.RoundEnd(game.RoundNum)
					}
					if judgeGame(game) != GameType.GameStatusRunning {
						// Game Over
						// TODO: Announce game-over
						game.Status = GameType.GameStatusEnd
						j.status = StatusWaiting
						break
					}
					game.RoundNum++
					// Round Start
					fmt.Println("Round", game.RoundNum, "start")
					game.Map.RoundStart(game.RoundNum)
					data.PutMap(j.gameId, game.Map)
					game.Map.OutputNumber()
				}
			}
		}
	}
}

func judgeGame(g *GameType.Game) GameType.GameStatus {
	// Check online player number
	onlinePlayerNum := uint8(0)
	for _, u := range g.UserList {
		if u.Status == GameType.UserStatusConnected {
			onlinePlayerNum++
		}
	}
	if onlinePlayerNum <= 0 {
		return GameType.GameStatusEnd
	}

	return GameType.GameStatusRunning
}
