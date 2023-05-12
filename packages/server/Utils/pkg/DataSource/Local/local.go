package Local

import (
	"math/rand"
	"server/Utils/pkg/DataSource"
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/InstructionType"
	"server/Utils/pkg/MapType"
	"sort"
	"sync"
	"time"
)

var _ DataSource.PersistentDataSource = (*Local)(nil)
var _ DataSource.TempDataSource = (*Local)(nil)

type Local struct {
	m                  sync.Mutex
	GamePool           map[GameType.GameId]*GameType.Game
	OriginalMapStrPool map[uint32]string
	InstructionLog     map[GameType.GameId]map[uint16]map[uint16]InstructionType.Instruction
}

func (l *Local) lock() bool {
	for i := 1; !l.m.TryLock(); i++ {
		time.Sleep(10 * time.Millisecond)
		if i >= 1e3 {
			panic("try to lock timeout")
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
	//log.Printf("[Local Lock Debug] Locked by %s -> %s\n", getCallerFunName(3), getCallerFunName(2))
	return true
}

func (l *Local) unlock() {
	l.m.Unlock()
}

func (l *Local) GetGameList(mode GameType.GameMode) []GameType.Game {
	if l.lock() {
		defer l.unlock()
		var ret []GameType.Game
		for _, p := range l.GamePool {
			if p.Status == GameType.GameStatusEnd {
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
func (l *Local) CancelGame(id GameType.GameId) (ok bool) {
	if l.lock() {
		defer l.unlock()
		g := l.GamePool[id]
		g.Status = GameType.GameStatusEnd
		g.UserList = nil // TODO
		return true
	} else {
		return false
	}
}

func (l *Local) GetGameInfo(id GameType.GameId) *GameType.Game {
	if l.lock() {
		defer l.unlock()
		g := *l.GamePool[id]
		g.UserList = nil
		return &g
	} else {
		return nil
	}
}

func (l *Local) GetCurrentUserList(id GameType.GameId) []GameType.User {
	if l.lock() {
		defer l.unlock()
		return l.GamePool[id].UserList
	} else {
		return nil
	}
}

func (l *Local) GetInstructions(id GameType.GameId, tempId uint16) []InstructionType.Instruction {
	if l.lock() {
		defer l.unlock()
		var list []InstructionType.Instruction
		for _, v := range l.InstructionLog[id][tempId] {
			list = append(list, v)
		}
		return list
	} else {
		return nil
	}
	//return ExampleInstruction
}

func (l *Local) NewInstructionTemp(id GameType.GameId, _ uint16) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].RoundNum++
		l.InstructionLog[id][l.GamePool[id].RoundNum] = make(map[uint16]InstructionType.Instruction)
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameStatus(id GameType.GameId, status GameType.GameStatus) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Status = status
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameMap(id GameType.GameId, m *MapType.Map) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Map = &(*m)
		return true
	} else {
		return false
	}
}

func (l *Local) SetUserStatus(id GameType.GameId, user GameType.User) (ok bool) {
	if l.lock() {
		defer l.unlock()

		g := l.GamePool[id]
		if g.Status == GameType.GameStatusEnd {
			return false
		}
		if g.Status == GameType.GameStatusWaiting {
			// Try to find
			for i, u := range g.UserList {
				if u.UserId == user.UserId {
					if user.Status == GameType.UserStatusDisconnected {
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
				panic("game mode is illegal")
			}
			// If the user is not in the list, try to add him/her if the game is not full
			if user.Status == GameType.UserStatusConnected && uint8(len(g.UserList)) < g.Mode.MaxUserNum {
				user.TeamId = 0
				g.UserList = append(g.UserList, user)
				return true
			}
		}
		if g.Status == GameType.GameStatusRunning {
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

func (l *Local) UpdateInstruction(id GameType.GameId, user GameType.User, instruction InstructionType.Instruction) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.InstructionLog[id][l.GamePool[id].RoundNum][user.UserId] = instruction
		return true
	} else {
		return false
	}
}

func (l *Local) GetCurrentMap(id GameType.GameId) *MapType.Map {
	if l.lock() {
		defer l.unlock()
		m := *l.GamePool[id].Map
		return &m // Defend the outside modification
	} else {
		return nil
	}
}

func (l *Local) GetOriginalMap(mapId uint32) *MapType.Map {
	if l.lock() {
		defer l.unlock()
		return MapType.Str2GameMap(mapId, l.OriginalMapStrPool[mapId])
	} else {
		return nil
	}
}

func (l *Local) GetCurrentGame(id GameType.GameId) *GameType.Game {
	if l.lock() {
		defer l.unlock()
		return l.GamePool[id]
	}
	return nil
}

func (l *Local) CreateGame(mode GameType.GameMode) GameType.GameId {
	//m := l.GetOriginalMap(rand.Uint32())
	m := l.GetOriginalMap(0) // TODO DEBUG ONLY
	if l.lock() {
		defer l.unlock()
		var gameId GameType.GameId
		for {
			gameId = GameType.GameId(rand.Uint32())
			if _, ok := l.GamePool[gameId]; !ok && gameId >= 100 { // gameId 1-99 is for debugging usage
				break
			}
		}
		g := &GameType.Game{
			Map:        m,
			Mode:       mode,
			Id:         gameId,
			CreateTime: time.Now().UnixMicro(),
			Status:     GameType.GameStatusWaiting,
			RoundNum:   0,
			UserList:   []GameType.User{},
		}

		l.GamePool[g.Id] = g
		l.InstructionLog[g.Id] = make(map[uint16]map[uint16]InstructionType.Instruction)
		return g.Id
	} else {
		return 0
	}
}

func (l *Local) DebugCreateGame(g *GameType.Game) (ok bool) {
	if !g.Map.HasBlocks() {
		g.Map = l.GetOriginalMap(g.Map.Id())
	}
	if l.lock() {
		defer l.unlock()
		_, ok := l.GamePool[g.Id]
		if ok {
			panic("game id has existed")
			return false
		}

		var gameId GameType.GameId
		for {
			gameId = GameType.GameId(rand.Uint32())
			if _, ok := l.GamePool[gameId]; !ok {
				break
			}
		}
		ng := &GameType.Game{
			Map:        g.Map,
			Mode:       g.Mode,
			Id:         g.Id,
			UserList:   g.UserList,
			CreateTime: time.Now().UnixMicro(),
			Status:     GameType.GameStatusWaiting,
			RoundNum:   0,
		}
		l.GamePool[g.Id] = ng
		l.InstructionLog[g.Id] = make(map[uint16]map[uint16]InstructionType.Instruction)
		return true
	} else {
		return false
	}
}
