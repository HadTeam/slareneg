package judge

import (
	"github.com/sirupsen/logrus"
	"server/game_logic"
	"server/game_logic/block"
	"server/game_logic/block_manager"
	"server/game_logic/game_def"
	"server/game_logic/map"
	data_source "server/utils/pkg/data_source"
	game_temp_pool "server/utils/pkg/game_temp_pool"
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
	gameId game_logic.Id
	status Status
	c      chan Status
}

func ApplyDataSource(source interface{}) {
	data = source.(data_source.TempDataSource)
	pData = source.(data_source.PersistentDataSource)
}

func NewGameJudge(id game_logic.Id) *GameJudge {
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
			data.SetGameStatus(j.gameId, game_logic.StatusRunning)
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
					if g.RoundNum == 1 {
						logrus.Info("debug")
					}
					data.NewInstructionTemp(j.gameId, g.RoundNum)
					instructionList := data.GetInstructions(j.gameId, g.RoundNum)

					for userId, ins := range instructionList {
						if !executeInstruction(j.gameId, userId, ins) {
							roundLogger.Infof("Instruction %#v failed to execute", ins)
						}
					}
					g.UserList = data.GetCurrentUserList(g.Id)
					g.Map.RoundEnd(g.RoundNum)
					if judgeGame(g, kingPos) != game_logic.StatusRunning {
						// Game Over
						data.SetGameStatus(g.Id, game_logic.StatusEnd)
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

				_map.DebugOutput(g.Map, func(block _type.Block) uint16 {
					return uint16(block.Meta().BlockId)
				}) // TODO
			}
		}
	}
}

func getKingPos(g *game_logic.Game) []_type.Position {
	var kingPos []_type.Position
	for y := uint8(1); y <= g.Map.Size().H; y++ {
		for x := uint8(1); x <= g.Map.Size().W; x++ {
			b := g.Map.GetBlock(_type.Position{X: x, Y: y})
			if b.Meta().BlockId == block.KingMeta.BlockId {
				kingPos = append(kingPos, _type.Position{X: x, Y: y})
			}
		}
	}
	return kingPos
}

// judgeGame TODO: Add unit test
func judgeGame(g *game_logic.Game, kingPos []_type.Position) game_logic.Status {
	// Check online player number
	onlinePlayerNum := uint8(0)
	for _, u := range g.UserList {
		if u.Status == _type.UserStatusConnected {
			onlinePlayerNum++
		}
	}
	if onlinePlayerNum <= 0 {
		return game_logic.StatusEnd
	}
	if onlinePlayerNum == 1 {
		// TODO: Announce game-over
		for _, u := range g.UserList {
			if u.Status == _type.UserStatusConnected {
				g.Winner = u.TeamId
				break
			}
		}
		return game_logic.StatusEnd
	}

	// Check king status
	if g.Mode == _type.Mode1v1 {
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
			return game_logic.StatusEnd
		}
	}

	return game_logic.StatusRunning
}

func allocateKing(g *game_logic.Game, kingPos []_type.Position) {
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
			block_manager.NewBlock(block.KingMeta.BlockId, g.Map.GetBlock(kingPos[i]).Number(), u.UserId))
		allocatableKingNum--
	}
}

func allocateTeam(g *game_logic.Game) {
	if g.Mode == _type.Mode1v1 {
		for i := range g.UserList {
			g.UserList[i].TeamId = uint8(i) + 1
		}
	} else {
		panic("unexpected game mod")
	}
}

func executeInstruction(id game_logic.Id, userId uint16, ins _type.Instruction) bool {
	var ret bool
	var m *_map.Map
	switch ins.(type) {
	case _type.Move:
		{
			i := ins.(_type.Move)
			m = data.GetCurrentMap(id)
			if m.GetBlock(_type.Position(i.Position)).OwnerId() != userId {
				return false
			}
			ret = m.Move(i)
		}
	}
	return ret
}
