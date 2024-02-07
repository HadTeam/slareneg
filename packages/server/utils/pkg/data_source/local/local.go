package local

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"server/game"
	"server/game/instruction"
	"server/game/map"
	"server/utils/pkg/data_source"
	"sort"
	"sync"
	"time"
)

var _ data_source.PersistentDataSource = (*Local)(nil)
var _ data_source.TempDataSource = (*Local)(nil)

type Local struct {
	m                  sync.Mutex
	GamePool           map[game.Id]*game.Game
	OriginalMapStrPool map[uint32]string
	InstructionLog     map[game.Id]map[uint16]map[uint16]instruction.Instruction
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

func (l *Local) SetWinner(id game.Id, teamId uint8) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Winner = teamId
	}
	return false
}

func (l *Local) GetGameList(mode game.Mode) []game.Game {
	if l.lock() {
		defer l.unlock()
		var ret []game.Game
		for _, p := range l.GamePool {
			if p.Status == game.StatusEnd {
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
func (l *Local) CancelGame(id game.Id) (ok bool) {
	if l.lock() {
		defer l.unlock()
		g := l.GamePool[id]
		if g.Status == game.StatusWaiting || g.Status == game.StatusRunning {
			g.Status = game.StatusEnd
			g.UserList = nil // TODO
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (l *Local) GetGameInfo(id game.Id) *game.Game {
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

func (l *Local) GetCurrentUserList(id game.Id) []game.User {
	if l.lock() {
		defer l.unlock()
		return l.GamePool[id].UserList
	} else {
		return nil
	}
}

func (l *Local) GetInstructions(id game.Id, tempId uint16) map[uint16]instruction.Instruction {
	if l.lock() {
		defer l.unlock()
		return l.InstructionLog[id][tempId]
	} else {
		return nil
	}
	//return ExampleInstruction
}

func (l *Local) NewInstructionTemp(id game.Id, _ uint16) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].RoundNum++
		l.InstructionLog[id][l.GamePool[id].RoundNum] = make(map[uint16]instruction.Instruction)
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameStatus(id game.Id, status game.Status) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Status = status
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameMap(id game.Id, m *_map.Map) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Map = &(*m)
		return true
	} else {
		return false
	}
}

func (l *Local) SetUserStatus(id game.Id, user game.User) (ok bool) {
	if l.lock() {
		defer l.unlock()

		g := l.GamePool[id]
		if g.Status == game.StatusEnd {
			return false
		}
		if g.Status == game.StatusWaiting {
			// Try to find
			for i, u := range g.UserList {
				if u.UserId == user.UserId {
					if user.Status == game.UserStatusDisconnected {
						// Remove the user from the list if they are disconnected
						g.UserList = append(g.UserList[:i], g.UserList[i+1:]...)
					} else {
						// Update the info
						user.TeamId = g.UserList[i].TeamId
						g.UserList[i] = user
					}
					return true
				}
			}

			// Specially check, for unexpected behavior that may exist
			if g.Mode.MaxUserNum == 0 || g.Mode.MinUserNum == 0 || g.Mode.NameStr == "" {
				logrus.Panic("game mode is illegal")
			}
			// If the user is not in the list, try to add him/her if the game is not full
			if user.Status == game.UserStatusConnected && uint8(len(g.UserList)) < g.Mode.MaxUserNum {
				user.TeamId = 0
				g.UserList = append(g.UserList, user)
				return true
			}
		}
		if g.Status == game.StatusRunning {
			for i, u := range g.UserList {
				if u.UserId == user.UserId {
					g.UserList[i].Status = user.Status
					return true
				}
			}
		}
	}
	return false
}

func (l *Local) UpdateInstruction(id game.Id, user game.User, instruction instruction.Instruction) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.InstructionLog[id][l.GamePool[id].RoundNum][user.UserId] = instruction
		return true
	} else {
		return false
	}
}

func (l *Local) GetCurrentMap(id game.Id) *_map.Map {
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

func (l *Local) GetCurrentGame(id game.Id) *game.Game {
	if l.lock() {
		defer l.unlock()
		return l.GamePool[id]
	}
	return nil
}

func (l *Local) CreateGame(mode game.Mode) game.Id {
	//m := l.GetOriginalMap(rand.Uint32())
	m := l.GetOriginalMap(0) // TODO DEBUG ONLY
	if l.lock() {
		defer l.unlock()
		var gameId game.Id
		for {
			gameId = game.Id(rand.Uint32())
			if _, ok := l.GamePool[gameId]; !ok && gameId >= 100 { // gameId 1-99 is for debugging usage
				break
			}
		}
		g := &game.Game{
			Map:        m,
			Mode:       mode,
			Id:         gameId,
			CreateTime: time.Now().UnixMicro(),
			Status:     game.StatusWaiting,
			RoundNum:   0,
			UserList:   []game.User{},
		}

		l.GamePool[g.Id] = g
		l.InstructionLog[g.Id] = make(map[uint16]map[uint16]instruction.Instruction)
		return g.Id
	} else {
		return 0
	}
}

func (l *Local) DebugCreateGame(g *game.Game) (ok bool) {
	if !g.Map.HasBlocks() {
		g.Map = l.GetOriginalMap(g.Map.Id())
	}
	if l.lock() {
		defer l.unlock()
		if _, ok := l.GamePool[g.Id]; ok {
			logrus.Panic("game id has existed")
			return false
		}

		var gameId game.Id
		for {
			gameId = game.Id(rand.Uint32())
			if _, ok := l.GamePool[gameId]; !ok {
				break
			}
		}
		ng := &game.Game{
			Map:        g.Map,
			Mode:       g.Mode,
			Id:         g.Id,
			UserList:   g.UserList,
			CreateTime: time.Now().UnixMicro(),
			Status:     game.StatusWaiting,
			RoundNum:   0,
		}
		l.GamePool[g.Id] = ng
		l.InstructionLog[g.Id] = make(map[uint16]map[uint16]instruction.Instruction)
		return true
	} else {
		return false
	}
}
