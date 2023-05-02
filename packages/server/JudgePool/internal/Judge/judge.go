package Judge

import (
	"log"
	"server/JudgePool/internal/InstructionExecutor"
	"server/Utils/pkg/DataSource"
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/MapType"
	"server/Utils/pkg/MapType/Blocks"
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
	gameId GameType.GameId
	status Status
	c      chan Status
}

func ApplyDataSource(source interface{}) {
	data = source.(DataSource.TempDataSource)
	pData = source.(DataSource.PersistentDataSource)
}

func NewGameJudge(id GameType.GameId) *GameJudge {
	j := &GameJudge{
		gameId: id,
		status: StatusWaiting,
		c:      make(chan Status),
	}
	go judgeWorking(j)
	return j
}

func (j *GameJudge) StartGame() {
	j.c <- StatusWorking
}

func judgeWorking(j *GameJudge) {
	for {
		j.status = <-j.c

		if j.status == StatusWorking {
			log.Printf("[Judge] Working for GameId %d\n", j.gameId)

			game := data.GetGameInfo(j.gameId)
			if game.Map.Blocks == nil {
				game.Map = pData.GetOriginalMap(game.Map.MapId)
			}
			data.SetGameStatus(j.gameId, GameType.GameStatusRunning)
			AllocateKing(game)
			data.SetGameMap(j.gameId, game.Map)
			data.NewInstructionTemp(j.gameId, 0)
			t := time.NewTicker(RoundTime)
			for range t.C {
				//Round End
				if game.RoundNum != 0 {
					log.Printf("[Round] Round %d end\n", game.RoundNum)
					data.NewInstructionTemp(j.gameId, game.RoundNum)
					instructionList := data.GetInstructions(j.gameId, game.RoundNum)

					ok := true
					for _, instruction := range instructionList {
						if !InstructionExecutor.ExecuteInstruction(j.gameId, instruction) {
							ok = false
						}
					}
					if !ok {
						log.Printf("[Warn] Instructions execution failed\n")
					}
					game.UserList = data.GetCurrentUserList(game.Id)
					gameOverSign := game.Map.RoundEnd(game.RoundNum) // TODO: Refactor the way to spread the game-over sign
					if gameOverSign || judgeGame(game) != GameType.GameStatusRunning {
						// Game Over
						// TODO: Announce game-over
						game.Status = GameType.GameStatusEnd
						j.status = StatusWaiting
						log.Printf("[Judge] Done for GameId %d\n", j.gameId)
						return
					}
				}
				game.RoundNum++
				// Round Start
				log.Printf("[Round] Round %d start\n", game.RoundNum)
				game.Map.RoundStart(game.RoundNum)
				data.SetGameMap(j.gameId, game.Map)
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
