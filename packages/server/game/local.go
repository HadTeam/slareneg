package game

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"server/game/instruction"
	"server/game/map"
	"server/game/mode"
	"server/game/user"
	"sort"
	"sync"
	"time"
)

var _ PersistentDataSource = (*Local)(nil)
var _ TempDataSource = (*Local)(nil)

type Local struct {
	m                  sync.Mutex
	GamePool           map[Id]*Game
	OriginalMapStrPool map[uint32]string
	InstructionLog     map[Id]map[uint16]map[uint16]instruction.Instruction
}

func (l *Local) lock() bool {
	for i := 1; !l.m.TryLock(); i++ {
		time.Sleep(10 * time.Millisecond)
		if i >= 1e3 {
			logrus.Panic("try to lock timeout")
			return false
		}
	}
	//// DEBUG CODE: Use to show the caller function name
	//getCallerFunName := func(skip int) string {
	//	pc, _, _, _ := runtime.Caller(skip)
	//	name := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	//	return name[len(name)-1]
	//}
	//
	//logrus.Tracef("Locked by %s -> %s\n", getCallerFunName(3), getCallerFunName(2))
	return true
}

func (l *Local) unlock() {
	l.m.Unlock()
}

func (l *Local) SetWinner(id Id, teamId uint8) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Winner = teamId
	}
	return false
}

func (l *Local) GetGameList(mode mode.Mode) []Game {
	if l.lock() {
		defer l.unlock()
		var ret []Game
		for _, p := range l.GamePool {
			if p.Status == StatusEnd {
				continue
			}
			g := *p
			g.UserList = nil
			ret = append(ret, g)
		}
		sort.Slice(ret, func(i, j int) bool {
			return ret[i].Id < ret[j].Id // by increasing order
		})
		return ret
	} else {
		return nil
	}
}

// CancelGame 1. Set game status 2. Quit existing users
func (l *Local) CancelGame(id Id) (ok bool) {
	if l.lock() {
		defer l.unlock()
		g := l.GamePool[id]
		if g.Status == StatusWaiting || g.Status == StatusRunning {
			g.Status = StatusEnd
			g.UserList = nil // TODO
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (l *Local) GetGameInfo(id Id) *Game {
	if l.lock() {
		defer l.unlock()
		g := *l.GamePool[id]
		g.UserList = nil
		g.Map = _map.CreateEmptyMapWithInfo(g.Map.Id(), g.Map.Size())
		return &g
	} else {
		return nil
	}
}

func (l *Local) GetCurrentUserList(id Id) []user.User {
	if l.lock() {
		defer l.unlock()
		return l.GamePool[id].UserList
	} else {
		return nil
	}
}

func (l *Local) GetInstructions(id Id, tempId uint16) map[uint16]instruction.Instruction {
	if l.lock() {
		defer l.unlock()
		return l.InstructionLog[id][tempId]
	} else {
		return nil
	}
	//return ExampleInstruction
}

func (l *Local) NewInstructionTemp(id Id, _ uint16) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].RoundNum++
		l.InstructionLog[id][l.GamePool[id].RoundNum] = make(map[uint16]instruction.Instruction)
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameStatus(id Id, status Status) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Status = status
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameMap(id Id, m *_map.Map) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Map = &(*m)
		return true
	} else {
		return false
	}
}

func (l *Local) SetUserStatus(id Id, _user user.User) (ok bool) {
	if l.lock() {
		defer l.unlock()

		g := l.GamePool[id]
		if g.Status == StatusEnd {
			return false
		}
		if g.Status == StatusWaiting {
			// Try to find
			for i, u := range g.UserList {
				if u.UserId == _user.UserId {
					if _user.Status == user.Disconnected {
						// Remove the user from the list if they are disconnected
						g.UserList = append(g.UserList[:i], g.UserList[i+1:]...)
					} else {
						// Update the info
						_user.TeamId = g.UserList[i].TeamId
						g.UserList[i] = _user
					}
					return true
				}
			}

			// Specially check, for unexpected behavior that may exist
			if g.Mode.MaxUserNum == 0 || g.Mode.MinUserNum == 0 || g.Mode.NameStr == "" {
				logrus.Panic("game mode is illegal")
			}
			// If the user is not in the list, try to add him/her if the game is not full
			if _user.Status == user.Connected && uint8(len(g.UserList)) < g.Mode.MaxUserNum {
				_user.TeamId = 0
				g.UserList = append(g.UserList, _user)
				return true
			}
		}
		if g.Status == StatusRunning {
			for i, u := range g.UserList {
				if u.UserId == _user.UserId {
					g.UserList[i].Status = _user.Status
					return true
				}
			}
		}
	}
	return false
}

func (l *Local) UpdateInstruction(id Id, user user.User, instruction instruction.Instruction) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.InstructionLog[id][l.GamePool[id].RoundNum][user.UserId] = instruction
		return true
	} else {
		return false
	}
}

func (l *Local) GetCurrentMap(id Id) *_map.Map {
	if l.lock() {
		defer l.unlock()
		m := *l.GamePool[id].Map
		return &m // Defend the outside modification
	} else {
		return nil
	}
}

func (l *Local) GetOriginalMap(mapId uint32) *_map.Map {
	if l.lock() {
		defer l.unlock()
		return _map.Str2GameMap(mapId, l.OriginalMapStrPool[mapId])
	} else {
		return nil
	}
}

func (l *Local) GetCurrentGame(id Id) *Game {
	if l.lock() {
		defer l.unlock()
		return l.GamePool[id]
	}
	return nil
}

func (l *Local) DebugCreateGame(g *Game) (ok bool) {
	if !g.Map.HasBlocks() {
		g.Map = l.GetOriginalMap(g.Map.Id())
	}
	if l.lock() {
		defer l.unlock()
		if _, ok := l.GamePool[g.Id]; ok {
			logrus.Panic("game id has existed")
			return false
		}

		var gameId Id
		for {
			gameId = Id(rand.Uint32())
			if _, ok := l.GamePool[gameId]; !ok {
				break
			}
		}
		ng := &Game{
			Map:        g.Map,
			Mode:       g.Mode,
			Id:         g.Id,
			UserList:   g.UserList,
			CreateTime: time.Now().UnixMicro(),
			Status:     StatusWaiting,
			RoundNum:   0,
		}
		l.GamePool[g.Id] = ng
		l.InstructionLog[g.Id] = make(map[uint16]map[uint16]instruction.Instruction)
		return true
	} else {
		return false
	}
}
