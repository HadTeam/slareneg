package judge

import (
	"context"
	"github.com/sirupsen/logrus"
	"server/game"
	"server/game/block"
	"server/game/instruction"
	"server/game/map"
	"server/utils/pkg/data_source"
	"server/utils/pkg/game_temp_pool"
	"time"
)

var RoundTime = time.Millisecond * 100

var data data_source.TempDataSource
var pData data_source.PersistentDataSource

type Status uint8

const (
	StatusWaiting Status = iota + 1
	StatusWorking
)

type gameContext struct {
	context.Context
	g       *game.Game
	kingPos []block.Position
}

type Judge struct {
	gameId game.Id
	status Status
	c      chan Status
	ctx    *gameContext
}

func ApplyDataSource(source interface{}) {
	data = source.(data_source.TempDataSource)
	pData = source.(data_source.PersistentDataSource)
}

func NewGameJudge(id game.Id) *Judge {
	j := &Judge{
		gameId: id,
		status: StatusWaiting,
		c:      make(chan Status),
		ctx: &gameContext{
			Context: context.Background(),
			g:       nil,
		},
	}
	go judgeWorking(j)
	return j
}

func (j *Judge) StartGame() {
	g := data.GetGameInfo(j.gameId)
	if m := data.GetCurrentMap(g.Id); !m.HasBlocks() {
		g.Map = pData.GetOriginalMap(g.Map.Id())
	} else {
		g.Map = m
	}
	g.UserList = data.GetCurrentUserList(g.Id)

	j.ctx.kingPos = getKingPos(g)
	j.ctx.g = g

	allocateKing(j.ctx)
	allocateTeam(j.ctx)

	game_temp_pool.Create(g.Id) // create temp pool for game
	data.SetGameStatus(g.Id, game.StatusRunning)
	data.SetGameMap(g.Id, g.Map)
	data.NewInstructionTemp(g.Id, 1)
	j.c <- StatusWorking
}

func judgeWorking(j *Judge) {
	for {
		j.status = <-j.c

		if j.status == StatusWorking {
			judgeLogger := logrus.WithFields(logrus.Fields{
				"gameId": j.gameId,
			})

			g := j.ctx.g
			judgeLogger.Infof("Working")
			t := time.NewTicker(RoundTime)
			for range t.C {
				roundLogger := logrus.WithFields(logrus.Fields{
					"gameId": j.gameId,
					"round":  g.RoundNum,
				})

				// Round Start
				roundLogger.Infof("Round start")
				g.RoundNum++ // NOTE: ONLY increase the LOCAL value
				g.Map.RoundStart(g.RoundNum)
				data.SetGameMap(j.gameId, g.Map)

				_map.DebugOutput(g.Map, func(block block.Block) uint16 {
					return uint16(block.Meta().BlockId)
				}) // TODO

				//Round End
				roundLogger.Infof("Round end")
				{
					data.NewInstructionTemp(j.gameId, g.RoundNum)
					instructionList := data.GetInstructions(j.gameId, g.RoundNum)

					for userId, ins := range instructionList {
						if !executeInstruction(j.gameId, userId, ins) {
							roundLogger.Infof("Instruction %#v failed to execute", ins)
						}
					}
					g.UserList = data.GetCurrentUserList(g.Id)
					g.Map.RoundEnd(g.RoundNum)
				}

				if judgeGame(g, j.ctx.kingPos) != game.StatusRunning {
					// Game Over
					data.SetGameStatus(g.Id, game.StatusEnd)
					data.SetWinner(g.Id, g.Winner)
					j.status = StatusWaiting
					game_temp_pool.Delete(g.Id)

					var winnerTeam []string
					for _, n := range g.UserList {
						if n.TeamId == g.Winner {
							winnerTeam = append(winnerTeam, n.Name)
						}
					}
					judgeLogger.Infof("Game end, winner %#v", winnerTeam)
					return
				}
			}
		}
	}
}

// judgeGame TODO: Add unit test
func judgeGame(g *game.Game, kingPos []block.Position) game.Status {
	// Check online player number
	onlinePlayerNum := uint8(0)
	for _, u := range g.UserList {
		if u.Status == game.UserStatusConnected {
			onlinePlayerNum++
		}
	}

	if onlinePlayerNum <= 0 {
		return game.StatusEnd
	}
	if onlinePlayerNum == 1 {
		// TODO: Announce game-over
		for _, u := range g.UserList {
			if u.Status == game.UserStatusConnected {
				g.Winner = u.TeamId
				break
			}
		}
		return game.StatusEnd
	}

	// Check king status
	if g.Mode == game.Mode1v1 {
		return judgeGameMode1v1(g, kingPos)
	}

	return game.StatusRunning
}

func executeInstruction(id game.Id, userId uint16, ins instruction.Instruction) bool {
	var ret bool
	var m *_map.Map
	switch ins.(type) {
	case instruction.Move:
		{
			i := ins.(instruction.Move)
			m = data.GetCurrentMap(id)
			if m.GetBlock(i.Position).OwnerId() != userId {
				return false
			}
			ret = m.Move(i)
		}
	}
	return ret
}
