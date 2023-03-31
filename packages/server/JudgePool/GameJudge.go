package JudgePool

import (
	"fmt"
	"math/rand"
	"server/ApiProvider/pkg/DataOperator"
	"server/JudgePool/internal/InstructionExecutor"
	"server/JudgePool/pkg/GameType"
	"time"
)

var RoundTime = time.Millisecond * 1000

var data DataOperator.DataSource

type Status uint8

const (
	StatusWaiting Status = iota + 1
	StatusWorking
)

type GameJudge struct {
	gameId GameType.GameId
	status Status
	id     uint8
	P      chan GameType.GameId
}

func ApplyDataSource(source DataOperator.DataSource) {
	data = source
	InstructionExecutor.ApplyDataSource(source)
}

func NewGameJudge(pool chan GameType.GameId) *GameJudge {
	j := &GameJudge{
		gameId: 0,
		status: StatusWaiting,
		id:     uint8(rand.Uint32()),
		P:      pool,
	}
	go judgeWorking(j)
	return j
}

func judgeWorking(j *GameJudge) {
	for {
		j.gameId = <-j.P
		j.status = StatusWorking
		fmt.Printf("[Judge %d] Working for GameId %d\n", j.id, j.gameId)
		game := data.GetCurrentGame(j.gameId)
		game.RoundNum = 0
		game.Map = data.GetOriginalMap(game.Map.MapId)
		data.PutMap(j.gameId, game.Map)
		fmt.Println("OriginalMap:")
		game.Map.OutputNumber()
		t := time.NewTicker(RoundTime)
		for range t.C {
			//Round End
			if game.RoundNum != 0 {
				fmt.Printf("[Round] Round %d end\n", game.RoundNum)
				data.AchieveInstructionTemp(j.gameId, game.RoundNum)
				instructionList := data.GetInstructionsFromTemp(j.gameId, game.RoundNum)

				ok := true
				for _, instruction := range instructionList {
					if !InstructionExecutor.ExecuteInstruction(j.gameId, instruction) {
						ok = false
					}
				}
				if !ok {
					fmt.Printf("[Warn] Instructions execution failed\n")
				}
				gameOverSign := game.Map.RoundEnd(game.RoundNum) // TODO: Refactor the way to spread the game-over sign
				if gameOverSign || judgeGame(game) != GameType.GameStatusRunning {
					// Game Over
					// TODO: Announce game-over
					game.Status = GameType.GameStatusEnd
					j.status = StatusWaiting
					fmt.Printf("[Judge %d] Done for GameId %d\n", j.id, j.gameId)
					break
				}
			}
			game.RoundNum++
			// Round Start
			fmt.Printf("[Round] Round %d start\n", game.RoundNum)
			game.Map.RoundStart(game.RoundNum)
			data.PutMap(j.gameId, game.Map)
			game.Map.OutputNumber()
		}
	}
}

// judgeGame TODO: Add unit test
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
	if onlinePlayerNum == 1 {
		// TODO: Announce game-over
		return GameType.GameStatusEnd
	}

	return GameType.GameStatusRunning
}
