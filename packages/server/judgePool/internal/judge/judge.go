package judge

import (
	"log"
	"server/utils/pkg/dataSource"
	"server/utils/pkg/game"
	"server/utils/pkg/instruction"
	"server/utils/pkg/map"
	"server/utils/pkg/map/block"
	"time"
)

var RoundTime = time.Millisecond * 100

var data dataSource.TempDataSource
var pData dataSource.PersistentDataSource

type Status uint8

const (
	StatusWaiting Status = iota + 1
	StatusWorking
)

type GameJudge struct {
	gameId game.GameId
	status Status
	c      chan Status
}

func ApplyDataSource(source interface{}) {
	data = source.(dataSource.TempDataSource)
	pData = source.(dataSource.PersistentDataSource)
}

func NewGameJudge(id game.GameId) *GameJudge {
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
			log.Printf("[judge] Working for GameId %d\n", j.gameId)

			g := data.GetGameInfo(j.gameId)
			if !g.Map.HasBlocks() {
				g.Map = pData.GetOriginalMap(g.Map.Id())
			}
			data.SetGameStatus(j.gameId, game.GameStatusRunning)
			g.UserList = data.GetCurrentUserList(j.gameId)

			kingPos := getKingPos(g)

			allocateKing(g, kingPos)
			allocateTeam(g)
			data.SetGameMap(j.gameId, g.Map)
			data.NewInstructionTemp(j.gameId, 0)
			t := time.NewTicker(RoundTime)
			for range t.C {
				//Round End
				if g.RoundNum != 0 {
					log.Printf("[Round] Round %d end\n", g.RoundNum)
					data.NewInstructionTemp(j.gameId, g.RoundNum)
					instructionList := data.GetInstructions(j.gameId, g.RoundNum)

					ok := true
					for _, ins := range instructionList {
						if !executeInstruction(j.gameId, ins) {
							ok = false
						}
					}
					if !ok {
						log.Printf("[Warn] Instructions execution failed\n")
					}
					g.UserList = data.GetCurrentUserList(g.Id)
					g.Map.RoundEnd(g.RoundNum)
					if judgeGame(g, kingPos) != game.GameStatusRunning {
						// Game Over
						data.SetGameStatus(g.Id, game.GameStatusEnd)
						j.status = StatusWaiting

						var winnerTeam []string
						for _, n := range g.UserList {
							if n.TeamId == g.Winner {
								winnerTeam = append(winnerTeam, n.Name)
							}
						}
						log.Printf("[judge] Done for GameId %d, winner %#v\n", j.gameId, winnerTeam)
						return
					}
				}
				g.RoundNum++
				// Round Start
				log.Printf("[Round] Round %d start\n", g.RoundNum)
				g.Map.RoundStart(g.RoundNum)
				data.SetGameMap(j.gameId, g.Map)

				_map.DebugOutput(g.Map, func(block block.Block) uint16 {
					return uint16(block.GetMeta().BlockId)
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
			if b.GetMeta().BlockId == block.BlockKingMeta.BlockId {
				kingPos = append(kingPos, block.Position{X: x, Y: y})
			}
		}
	}
	return kingPos
}

// judgeGame TODO: Add unit test
func judgeGame(g *game.Game, kingPos []block.Position) game.GameStatus {
	// Check online player number
	onlinePlayerNum := uint8(0)
	for _, u := range g.UserList {
		if u.Status == game.UserStatusConnected {
			onlinePlayerNum++
		}
	}
	if onlinePlayerNum <= 0 {
		return game.GameStatusEnd
	}
	if onlinePlayerNum == 1 {
		// TODO: Announce game-over
		for _, u := range g.UserList {
			if u.Status == game.UserStatusConnected {
				g.Winner = u.TeamId
				break
			}
		}
		return game.GameStatusEnd
	}

	// Check king status
	if g.Mode == game.GameMode1v1 {
		flag := true
		for _, k := range kingPos {
			if g.Map.GetBlock(k).GetMeta().BlockId != block.BlockKingMeta.BlockId {
				flag = false
				break
			}
		}
		if !flag {
			var w uint16
			for _, k := range kingPos {
				if g.Map.GetBlock(k).GetMeta().BlockId == block.BlockKingMeta.BlockId {
					w = g.Map.GetBlock(k).GetOwnerId()
				}
			}
			var wt uint8
			for _, u := range g.UserList {
				if u.UserId == w {
					wt = u.TeamId
				}
			}
			g.Winner = wt
			return game.GameStatusEnd
		}
	}

	return game.GameStatusRunning
}

func allocateKing(g *game.Game, kingPos []block.Position) {
	allocatableKingNum := 0
	for _, k := range kingPos {
		if g.Map.GetBlock(k).GetOwnerId() == 0 {
			allocatableKingNum++
		}
	}

	for i, u := range g.UserList { // allocate king blocks by order, ignoring the part out of user number
		if allocatableKingNum <= 0 { // check for debug creating behaviour

			break
		}
		g.Map.SetBlock(kingPos[i],
			block.NewBlock(block.BlockKingMeta.BlockId, g.Map.GetBlock(kingPos[i]).GetNumber(), u.UserId))
		allocatableKingNum--
	}
}

func allocateTeam(g *game.Game) {
	if g.Mode == game.GameMode1v1 {
		for i, _ := range g.UserList {
			g.UserList[i].TeamId = uint8(i) + 1
		}
	} else {
		panic("unexpected game mod")
	}
}

func executeInstruction(id game.GameId, ins instruction.Instruction) bool {
	var ret bool
	var m *_map.Map
	switch ins.(type) {
	case instruction.Move:
		{
			m = data.GetCurrentMap(id)
			ret = m.Move(ins.(instruction.Move))

		}
	}
	if !ret {
		log.Printf("[Warn] Execute instruction failed: %#v \n", ins)
	}
	return ret
}
