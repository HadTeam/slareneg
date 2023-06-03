package judge

import (
	"github.com/sirupsen/logrus"
	data_source "server/utils/pkg/datasource"
	"server/utils/pkg/game"
	game_temp_pool "server/utils/pkg/gametemppool"
	"server/utils/pkg/instruction"
	"server/utils/pkg/map"
	"server/utils/pkg/map/block"
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

type GameJudge struct {
	gameId game.Id
	status Status
	c      chan Status
}

func ApplyDataSource(source interface{}) {
	data = source.(data_source.TempDataSource)
	pData = source.(data_source.PersistentDataSource)
}

func NewGameJudge(id game.Id) *GameJudge {
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
			judgeLogger := logrus.WithFields(logrus.Fields{
				"gameId": j.gameId,
			})

			judgeLogger.Infof("Working")

			g := data.GetGameInfo(j.gameId)
			game_temp_pool.Create(g.Id)
			g.Map = data.GetCurrentMap(g.Id)
			if !g.Map.HasBlocks() {
				g.Map = pData.GetOriginalMap(g.Map.Id())
			}
			data.SetGameStatus(j.gameId, game.StatusRunning)
			g.UserList = data.GetCurrentUserList(j.gameId)

			kingPos := getKingPos(g)

			allocateKing(g, kingPos)
			allocateTeam(g)
			data.SetGameMap(j.gameId, g.Map)
			data.NewInstructionTemp(j.gameId, 0)
			t := time.NewTicker(RoundTime)
			for range t.C {
				roundLogger := logrus.WithFields(logrus.Fields{
					"gameId": j.gameId,
					"round":  g.RoundNum,
				})
				//Round End
				if g.RoundNum != 0 {
					roundLogger.Infof("Round end")
					data.NewInstructionTemp(j.gameId, g.RoundNum)
					instructionList := data.GetInstructions(j.gameId, g.RoundNum)

					for _, ins := range instructionList {
						if !executeInstruction(j.gameId, ins) {
							roundLogger.Infof("Instruction %#v failed to execute", ins)
						}
					}
					g.UserList = data.GetCurrentUserList(g.Id)
					g.Map.RoundEnd(g.RoundNum)
					if judgeGame(g, kingPos) != game.StatusRunning {
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
				g.RoundNum++ // NOTE: ONLY increase the LOCAL value
				// Round Start
				roundLogger.Infof("Round start")
				g.Map.RoundStart(g.RoundNum)
				data.SetGameMap(j.gameId, g.Map)

				_map.DebugOutput(g.Map, func(block block.Block) uint16 {
					return uint16(block.Meta().BlockId)
				}) // TODO
			}
		}
	}
}

func getKingPos(g *game.Game) []block.Position {
	var kingPos []block.Position
	for y := uint8(1); y <= g.Map.Size().H; y++ {
		for x := uint8(1); x <= g.Map.Size().W; x++ {
			b := g.Map.GetBlock(block.Position{X: x, Y: y})
			if b.Meta().BlockId == block.KingMeta.BlockId {
				kingPos = append(kingPos, block.Position{X: x, Y: y})
			}
		}
	}
	return kingPos
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
		flag := true
		for _, k := range kingPos {
			if g.Map.GetBlock(k).Meta().BlockId != block.KingMeta.BlockId {
				flag = false
				break
			}
		}
		if !flag {
			var w uint16
			for _, k := range kingPos {
				if g.Map.GetBlock(k).Meta().BlockId == block.KingMeta.BlockId {
					w = g.Map.GetBlock(k).OwnerId()
				}
			}
			var wt uint8
			for _, u := range g.UserList {
				if u.UserId == w {
					wt = u.TeamId
					break
				}
			}
			g.Winner = wt
			return game.StatusEnd
		}
	}

	return game.StatusRunning
}

func allocateKing(g *game.Game, kingPos []block.Position) {
	allocatableKingNum := 0
	for _, k := range kingPos {
		if g.Map.GetBlock(k).OwnerId() == 0 {
			allocatableKingNum++
		}
	}

	for i, u := range g.UserList { // allocate king blocks by order, ignoring the part out of user number
		if allocatableKingNum <= 0 { // check for debug creating behaviour

			break
		}
		g.Map.SetBlock(kingPos[i],
			block.NewBlock(block.KingMeta.BlockId, g.Map.GetBlock(kingPos[i]).Number(), u.UserId))
		allocatableKingNum--
	}
}

func allocateTeam(g *game.Game) {
	if g.Mode == game.Mode1v1 {
		for i := range g.UserList {
			g.UserList[i].TeamId = uint8(i) + 1
		}
	} else {
		panic("unexpected game mod")
	}
}

func executeInstruction(id game.Id, ins instruction.Instruction) bool {
	var ret bool
	var m *_map.Map
	switch ins.(type) {
	case instruction.Move:
		{
			m = data.GetCurrentMap(id)
			ret = m.Move(ins.(instruction.Move))

		}
	}
	return ret
}
