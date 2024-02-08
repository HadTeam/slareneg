package judge

import (
	"github.com/sirupsen/logrus"
	"server/game/block"
	_map "server/game/map"
	"server/game/mode"
	"server/game/user"
)

type Status uint8
type Result uint8
type Winner uint8

const (
	StatusWaiting Status = iota + 1
	StatusWorking
)

const (
	ResultContinue Result = iota
	ResultEnd
)

const (
	WinnerNone Winner = iota
)

type Id uint16

type Judge struct {
	gameId Id
	status Status
	c      chan Status
}

func NewGameJudge(id Id) *Judge {
	j := &Judge{
		gameId: id,
		status: StatusWaiting,
		c:      make(chan Status),
	}
	go judgeWorking(j)
	return j
}

func (j *Judge) StartGame() {
	j.c <- StatusWorking
}

func judgeWorking(j *Judge) {
	for {
		j.status = <-j.c
		if j.status == StatusWorking {
			judgeLogger := logrus.WithFields(logrus.Fields{
				"gameId": j.gameId,
			})

			judgeLogger.Infof("Working")

		}
	}
}

func CheckUserNumber(userList []user.User) Result {
	// Check online player number
	onlinePlayerNum := uint8(0)
	for _, u := range userList {
		if u.Status == user.Connected {
			onlinePlayerNum++
		}
	}

	if onlinePlayerNum <= 0 {
		return ResultEnd
	}
	if onlinePlayerNum == 1 {
		return ResultEnd
	}
	return ResultContinue
}

// Execute TODO: Add unit test
func Execute(m *_map.Map, userList []user.User, kingPos []block.Position, md mode.Mode) (Result, Winner) {
	wt := WinnerNone
	if (CheckUserNumber(userList)) == ResultEnd {
		for _, u := range userList {
			if u.Status == user.Connected {
				wt = Winner(u.TeamId)
				break
			}
		}
		return ResultEnd, wt
	}

	// Check king status
	if md == mode.Mode1v1 {
		return judgeGameMode1v1(m, userList, kingPos)
	}

	return ResultContinue, wt
}
