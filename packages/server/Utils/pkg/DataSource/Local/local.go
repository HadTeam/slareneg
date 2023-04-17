package Local

import (
	"math/rand"
	"server/Utils/pkg"
	"server/Utils/pkg/DataSource"
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/InstructionType"
	"server/Utils/pkg/MapType"
	"sync"
	"time"
)

var _ DataSource.PersistentDataSource = (*Local)(nil)
var _ DataSource.TempDataSource = (*Local)(nil)

type Local struct {
	m                   sync.Mutex
	GamePool            map[GameType.GameId]*GameType.Game
	OriginalMapStrPool  map[uint32]string
	InstructionTempPool map[GameType.GameId]map[uint16]InstructionType.Instruction
	InstructionLog      map[GameType.GameId]map[uint16][]InstructionType.Instruction
}

func (l *Local) lock() bool {
	for i := 1; !l.m.TryLock(); i++ {
		time.Sleep(10 * time.Millisecond)
		if i >= 1e4 {
			panic("try to lock timeout")
			return false
		}
	}
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
		g.UserList = nil
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

var ExampleInstruction = []InstructionType.Instruction{
	InstructionType.Move{UserId: 1, Position: InstructionType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
	InstructionType.Move{UserId: 2, Position: InstructionType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
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
	//if l.lock("") {
	//	defer l.unlock("")
	//	return l.InstructionLog[id][roundNum]
	//} else {
	//	return nil
	//}
	return ExampleInstruction
}

func (l *Local) NewInstructionTemp(id GameType.GameId, tempId uint16) (ok bool) {
	if l.lock() {
		defer l.unlock()
		var list []InstructionType.Instruction
		for _, v := range l.InstructionTempPool[id] {
			list = append(list, v)
		}
		l.InstructionLog[id][tempId] = list
		l.InstructionTempPool[id] = make(map[uint16]InstructionType.Instruction)
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
		l.GamePool[id].Map = m
		return true
	} else {
		return false
	}
}

func (l *Local) SetUserStatus(id GameType.GameId, user GameType.User) (ok bool) {
	if l.lock() {
		defer l.unlock()
		g := l.GamePool[id]
		// Try to find the user
		for i, u := range g.UserList {
			if u.UserId == user.UserId {
				g.UserList[i].Status = user.Status
				return true
			}
		}
		// Not found, try to join the game
		if g.Status == GameType.GameStatusWaiting && uint8(len(g.UserList)) < g.Mode.MaxUserNum {
			g.UserList = append(g.UserList, user)
			return true
		}
	}
	return false
}

func (l *Local) UpdateInstruction(id GameType.GameId, user GameType.User, instruction InstructionType.Instruction) (ok bool) {
	if l.lock() {
		defer l.unlock()
		l.InstructionTempPool[id][user.UserId] = instruction
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
		return pkg.Str2GameMap(mapId, l.OriginalMapStrPool[mapId])
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
			if _, ok := l.GamePool[gameId]; !ok {
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
		l.InstructionLog[g.Id] = make(map[uint16][]InstructionType.Instruction)
		return g.Id
	} else {
		return 0
	}
}
