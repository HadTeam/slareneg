package Judge

import (
	"fmt"
	"math/rand"
	"server/JudgePool/internal/InstructionExecutor"
	"server/Untils/pkg/DataSource"
	"server/Untils/pkg/GameType"
	"server/Untils/pkg/MapType"
	"time"
)

var RoundTime = time.Millisecond * 1000

var data DataSource.TempDataSource
var pData DataSource.PersistentDataSource

type Status uint8

const (
	StatusWaiting Status = iota + 1
	StatusWorking
)

type GameJudge struct {
	gameId       GameType.GameId
	status       Status
	ChangeStatus chan Status
	id           uint8
}

func ApplyDataSource(source interface{}) {
	data = source.(DataSource.TempDataSource)
	pData = source.(DataSource.PersistentDataSource)
}

func NewGameJudge(id GameType.GameId) *GameJudge {
	j := &GameJudge{
		gameId: id,
		status: StatusWaiting,
		id:     uint8(rand.Uint32()),
	}
	go judgeWorking(j)
	return j
}

func judgeWorking(j *GameJudge) {
	for {
		j.status = <-j.ChangeStatus

		if j.status == StatusWorking {
			fmt.Printf("[Judge %d] Working for GameId %d\n", j.id, j.gameId)

			game := data.GetGameInfo(j.gameId)
			game.Map = pData.GetOriginalMap(game.Map.MapId)
			data.SetGameStatus(j.gameId, GameType.GameStatusRunning)
			data.SetGameMap(j.gameId, game.Map)
			fmt.Println("OriginalMap:")
			MapType.OutputNumber(game.Map)
			t := time.NewTicker(RoundTime)
			for range t.C {
				//Round End
				if game.RoundNum != 0 {
					fmt.Printf("[Round] Round %d end\n", game.RoundNum)
					data.NewInstructionTemp(j.gameId, game.RoundNum)
					instructionList := data.GetInstructions(j.gameId, game.RoundNum)

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
					if gameOverSign || judgeGame(&game) != GameType.GameStatusRunning {
						// Game Over
						// TODO: Announce game-over
						game.Status = GameType.GameStatusEnd
						j.status = StatusWaiting
						fmt.Printf("[Judge %d] Done for GameId %d\n", j.id, j.gameId)
						return
					}
				}
				game.RoundNum++
				// Round Start
				fmt.Printf("[Round] Round %d start\n", game.RoundNum)
				game.Map.RoundStart(game.RoundNum)
				data.SetGameMap(j.gameId, game.Map)
				MapType.OutputNumber(game.Map)
			}
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
