package local

import (
	"math/rand"
	"server/utils/pkg/datasource"
	"server/utils/pkg/game"
	"server/utils/pkg/instruction"
	"server/utils/pkg/map"
	"sort"
	"sync"
	"time"
)

var _ datasource.PersistentDataSource = (*Local)(nil)
var _ datasource.TempDataSource = (*Local)(nil)

type Local struct {
	m                  sync.Mutex
	GamePool           map[game.GameId]*game.Game
	OriginalMapStrPool map[uint32]string
	InstructionLog     map[game.GameId]map[uint16]map[uint16]instruction.Instruction
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
	//log.Printf("[local Lock Debug] Locked by %s -> %s\n", getCallerFunName(3), getCallerFunName(2))
	return true
}

func (l *Local) unlock() {
	l.m.Unlock()
}

func (l *Local) GetGameList(mode game.GameMode) []game.Game {
	if l.lock() {
		defer l.unlock()
		var ret []game.Game
		for _, p := range l.GamePool {
			if p.Status == game.GameStatusEnd {
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
func (l *Local) CancelGame(id game.GameId) (ok bool) {
	if l.lock() {
		defer l.unlock()
		g := l.GamePool[id]
		g.Status = game.GameStatusEnd
		g.UserList = nil // TODO
		return true
	} else {
		return false
	}
}

func (l *Local) GetGameInfo(id game.GameId) *game.Game {
	if l.lock() {
		defer l.unlock()
		g := *l.GamePool[id]
		g.UserList = nil
		return &g
	} else {
		return nil
	}
}

func (l *Local) GetCurrentUserList(id game.GameId) []game.User {
	if l.lock() {
		defer l.unlock()
		return l.GamePool[id].UserList
	} else {
		return nil
	}
}

func (l *Local) GetInstructions(id game.GameId, tempId uint16) []instruction.Instruction {
	if l.lock() {
		defer l.unlock()
		var list []instruction.Instruction
		for _, v := range l.InstructionLog[id][tempId] {
			list = append(list, v)
		}
		return list
	} else {
		return nil
	}
	//return ExampleInstruction
}

func (l *Local) NewInstructionTemp(id game.GameId, _ uint16) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].RoundNum++
		l.InstructionLog[id][l.GamePool[id].RoundNum] = make(map[uint16]instruction.Instruction)
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameStatus(id game.GameId, status game.GameStatus) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Status = status
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameMap(id game.GameId, m *_map.Map) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.GamePool[id].Map = &(*m)
		return true
	} else {
		return false
	}
}

func (l *Local) SetUserStatus(id game.GameId, user game.User) (ok bool) {
	if l.lock() {
		defer l.unlock()

		g := l.GamePool[id]
		if g.Status == game.GameStatusEnd {
			return false
		}
		if g.Status == game.GameStatusWaiting {
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
				panic("game mode is illegal")
			}
			// If the user is not in the list, try to add him/her if the game is not full
			if user.Status == game.UserStatusConnected && uint8(len(g.UserList)) < g.Mode.MaxUserNum {
				user.TeamId = 0
				g.UserList = append(g.UserList, user)
				return true
			}
		}
		if g.Status == game.GameStatusRunning {
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

func (l *Local) UpdateInstruction(id game.GameId, user game.User, instruction instruction.Instruction) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.InstructionLog[id][l.GamePool[id].RoundNum][user.UserId] = instruction
		return true
	} else {
		return false
	}
}

func (l *Local) GetCurrentMap(id game.GameId) *_map.Map {
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

func (l *Local) GetCurrentGame(id game.GameId) *game.Game {
	if l.lock() {
		defer l.unlock()
		return l.GamePool[id]
	}
	return nil
}

func (l *Local) CreateGame(mode game.GameMode) game.GameId {
	//m := l.GetOriginalMap(rand.Uint32())
	m := l.GetOriginalMap(0) // TODO DEBUG ONLY
	if l.lock() {
		defer l.unlock()
		var gameId game.GameId
		for {
			gameId = game.GameId(rand.Uint32())
			if _, ok := l.GamePool[gameId]; !ok && gameId >= 100 { // gameId 1-99 is for debugging usage
				break
			}
		}
		g := &game.Game{
			Map:        m,
			Mode:       mode,
			Id:         gameId,
			CreateTime: time.Now().UnixMicro(),
			Status:     game.GameStatusWaiting,
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
		_, ok := l.GamePool[g.Id]
		if ok {
			panic("game id has existed")
			return false
		}

		var gameId game.GameId
		for {
			gameId = game.GameId(rand.Uint32())
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
			Status:     game.GameStatusWaiting,
			RoundNum:   0,
		}
		l.GamePool[g.Id] = ng
		l.InstructionLog[g.Id] = make(map[uint16]map[uint16]instruction.Instruction)
		return true
	} else {
		return false
	}
}
