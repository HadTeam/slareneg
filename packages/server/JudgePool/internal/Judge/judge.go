package Judge

import (
	"log"
	"server/Utils/pkg/DataSource"
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/InstructionType"
	"server/Utils/pkg/MapType"
	"server/Utils/pkg/MapType/BlockType"
	"time"
)

var RoundTime = time.Millisecond * 100

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
			if !game.Map.HasBlocks() {
				game.Map = pData.GetOriginalMap(game.Map.Id())
			}
			data.SetGameStatus(j.gameId, GameType.GameStatusRunning)
			game.UserList = data.GetCurrentUserList(j.gameId)

			kingPos := getKingPos(game)

			allocateKing(game, kingPos)
			allocateTeam(game)
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
						if !executeInstruction(j.gameId, instruction) {
							ok = false
						}
					}
					if !ok {
						log.Printf("[Warn] Instructions execution failed\n")
					}
					game.UserList = data.GetCurrentUserList(game.Id)
					game.Map.RoundEnd(game.RoundNum)
					if judgeGame(game, kingPos) != GameType.GameStatusRunning {
						// Game Over
						data.SetGameStatus(game.Id, GameType.GameStatusEnd)
						j.status = StatusWaiting

						var winnerTeam []string
						for _, n := range game.UserList {
							if n.TeamId == game.Winner {
								winnerTeam = append(winnerTeam, n.Name)
							}
						}
						log.Printf("[Judge] Done for GameId %d, winner %#v\n", j.gameId, winnerTeam)
						return
					}
				}
				game.RoundNum++
				// Round Start
				log.Printf("[Round] Round %d start\n", game.RoundNum)
				game.Map.RoundStart(game.RoundNum)
				data.SetGameMap(j.gameId, game.Map)

				MapType.DebugOutput(game.Map, func(block BlockType.Block) uint16 {
					return uint16(block.GetMeta().BlockId)
				}) // TODO
			}
		}
	}
}

func getKingPos(g *GameType.Game) []BlockType.Position {
	var kingPos []BlockType.Position
	for y := uint8(1); y <= g.Map.Size().H; y++ {
		for x := uint8(1); x <= g.Map.Size().W; x++ {
			b := g.Map.GetBlock(BlockType.Position{X: x, Y: y})
			if b.GetMeta().BlockId == BlockType.BlockKingMeta.BlockId {
				kingPos = append(kingPos, BlockType.Position{X: x, Y: y})
			}
		}
	}
	return kingPos
}

// judgeGame TODO: Add unit test
func judgeGame(g *GameType.Game, kingPos []BlockType.Position) GameType.GameStatus {
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
		for _, u := range g.UserList {
			if u.Status == GameType.UserStatusConnected {
				g.Winner = u.TeamId
				break
			}
		}
		return GameType.GameStatusEnd
	}

	// Check king status
	if g.Mode == GameType.GameMode1v1 {
		flag := true
		for _, k := range kingPos {
			if g.Map.GetBlock(k).GetMeta().BlockId != BlockType.BlockKingMeta.BlockId {
				flag = false
				break
			}
		}
		if !flag {
			var w uint16
			for _, k := range kingPos {
				if g.Map.GetBlock(k).GetMeta().BlockId == BlockType.BlockKingMeta.BlockId {
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
			return GameType.GameStatusEnd
		}
	}

	return GameType.GameStatusRunning
}

func allocateKing(g *GameType.Game, kingPos []BlockType.Position) {
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
			BlockType.NewBlock(BlockType.BlockKingMeta.BlockId, g.Map.GetBlock(kingPos[i]).GetNumber(), u.UserId))
		allocatableKingNum--
	}
}

func allocateTeam(g *GameType.Game) {
	if g.Mode == GameType.GameMode1v1 {
		for i, _ := range g.UserList {
			g.UserList[i].TeamId = uint8(i) + 1
		}
	} else {
		panic("unexpected game mod")
	}
}

func executeInstruction(id GameType.GameId, instruction InstructionType.Instruction) bool {
	var ret bool
	var m *MapType.Map
	switch instruction.(type) {
	case InstructionType.Move:
		{
			m = data.GetCurrentMap(id)
			ret = m.Move(instruction.(InstructionType.Move))

		}
	}
	if !ret {
		log.Printf("[Warn] Execute instruction failed: %#v \n", instruction)
	}
	return ret
}
