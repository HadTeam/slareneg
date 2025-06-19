package game

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"server/game/block"
	"server/game/instruction"
	"server/game/judge"
	"server/game/map"
	"server/game/mode"
	"server/game/user"
	"server/utils/pkg"
	"time"
)

type Status uint8
type Id uint16

const (
	StatusWaiting Status = iota + 1
	StatusRunning
	StatusEnd
)

type Game struct {
	Map        *_map.Map
	Mode       mode.Mode
	Id         Id
	UserList   []user.User
	CreateTime int64
	Status     Status
	RoundNum   uint16
	Winner     uint8 // TeamId
	KingPos    []block.Position
}

var Data TempDataSource
var pData PersistentDataSource

func ApplyDataSource(source interface{}) {
	Data = source.(TempDataSource)
	pData = source.(PersistentDataSource)
}

func Create(mode mode.Mode) *Game {
	return &Game{
		Map:        nil,
		Mode:       mode,
		Id:         Id(rand.Intn(10000)),
		CreateTime: time.Now().UnixMicro(),
		Status:     StatusWaiting,
		RoundNum:   0,
		UserList:   []user.User{},
	}
}

func (g *Game) Start() {
	g.Status = StatusRunning
	if m := Data.GetCurrentMap(g.Id); !m.HasBlocks() {
		g.Map = pData.GetOriginalMap(g.Map.Id())
	} else {
		g.Map = m
	}
	g.UserList = Data.GetCurrentUserList(g.Id)

	g.allocateKing()
	g.KingPos = g.getKingPos()
	g.allocateTeam()

	pkg.TempPoolCreate(g.Id) // create temp pool for game
	Data.SetGameStatus(g.Id, StatusRunning)
	Data.SetGameMap(g.Id, g.Map)
	Data.NewInstructionTemp(g.Id, 1)
}

var RoundTime = time.Millisecond * 100

func (g *Game) MainLoop() {
	t := time.NewTicker(RoundTime)
	for range t.C {
		if g.Status != StatusRunning {
			break
		}
		g.Tick()
	}
	t.Stop()
}

func (g *Game) Tick() {
	roundLogger := logrus.WithFields(logrus.Fields{
		"gameId": g.Id,
		"round":  g.RoundNum,
	})
	roundStart := func() {
		g.RoundNum++ // NOTE: ONLY increase the LOCAL value
		roundLogger.Infof("Round start")
		g.Map.RoundStart(g.RoundNum)
		Data.SetGameMap(g.Id, g.Map)

		_map.DebugOutput(g.Map, func(block block.Block) uint16 {
			return uint16(block.Meta().BlockId)
		}) // TODO
	}
	roundEnd := func() {
		roundLogger.Infof("Round end")

		Data.NewInstructionTemp(g.Id, g.RoundNum)
		instructionList := Data.GetInstructions(g.Id, g.RoundNum)

		for userId, ins := range instructionList {
			if !instruction.Execute(userId, g.Map, ins) {
				roundLogger.Infof("Instruction %#v failed to execute", ins)
			}
		}
		g.UserList = Data.GetCurrentUserList(g.Id)
		g.Map.RoundEnd(g.RoundNum)
	}
	roundJudge := func() {
		if res, wt := judge.Execute(g.Map, g.UserList, g.KingPos, g.Mode); res != judge.ResultContinue {
			// Game Over
			g.Winner = uint8(wt)
			g.End()

			roundLogger.Infof("Game over, winner %d", g.Winner)
			return
		}
	}

	if g.RoundNum != 0 {
		roundEnd()
		roundJudge()
	}
	roundStart()
}

func (g *Game) End() {
	Data.SetGameStatus(g.Id, StatusEnd)
	Data.SetWinner(g.Id, g.Winner)
	g.Status = StatusEnd
	pkg.TempPoolDelete(g.Id)

	var winnerTeam []string
	for _, n := range g.UserList {
		if n.TeamId == g.Winner {
			winnerTeam = append(winnerTeam, n.Name)
		}
	}
}
