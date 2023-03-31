package Judge

import (
	"fmt"
	"math/rand"
	"server/ApiProvider/pkg/DataOperator"
	"server/JudgePool/internal/InstructionExecutor"
	GameType2 "server/Untils/pkg/GameType"
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
	gameId GameType2.GameId
	status Status
	id     uint8
	P      chan GameType2.GameId
}

func ApplyDataSource(source DataOperator.DataSource) {
	data = source
	InstructionExecutor.ApplyDataSource(source)
}

func NewGameJudge(pool chan GameType2.GameId) *GameJudge {
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
				if gameOverSign || judgeGame(game) != GameType2.GameStatusRunning {
					// Game Over
					// TODO: Announce game-over
					game.Status = GameType2.GameStatusEnd
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
func judgeGame(g *GameType2.Game) GameType2.GameStatus {
	// Check online player number
	onlinePlayerNum := uint8(0)
	for _, u := range g.UserList {
		if u.Status == GameType2.UserStatusConnected {
			onlinePlayerNum++
		}
	}
	if onlinePlayerNum <= 0 {
		return GameType2.GameStatusEnd
	}
	if onlinePlayerNum == 1 {
		// TODO: Announce game-over
		return GameType2.GameStatusEnd
	}

	return GameType2.GameStatusRunning
}
