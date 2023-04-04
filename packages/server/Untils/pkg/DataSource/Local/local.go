package Local

import (
	"fmt"
	"server/Untils/pkg"
	"server/Untils/pkg/DataSource"
	"server/Untils/pkg/GameType"
	"server/Untils/pkg/InstructionType"
	"server/Untils/pkg/MapType"
	"sync"
	"time"
)

var _ DataSource.PersistentDataSource = (*Local)(nil)
var _ DataSource.TempDataSource = (*Local)(nil)

type Local struct {
	m                   sync.Mutex
	GamePool            map[GameType.GameId]*GameType.Game
	OriginalMapStrPool  map[uint32]string
	InstructionTempPool map[GameType.GameId]map[uint8]InstructionType.Instruction
	InstructionLog      map[GameType.GameId]map[uint8][]InstructionType.Instruction
}

var ExampleInstruction = []InstructionType.Instruction{
	InstructionType.MoveInstruction{UserId: 1, Position: InstructionType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
	InstructionType.MoveInstruction{UserId: 2, Position: InstructionType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
}

func (l *Local) GetCurrentUserList(id GameType.GameId) []GameType.User {
	if l.m.TryLock() {
		defer l.m.Unlock()
		return l.GamePool[id].UserList
	} else {
		return nil
	}
}

func (l *Local) GetInstructions(id GameType.GameId, tempId uint8) []InstructionType.Instruction {
	//if l.m.TryLock() {
	//	defer l.m.Unlock()
	//	return l.InstructionLog[id][roundNum]
	//} else {
	//	return nil
	//}
	return ExampleInstruction
}

func (l *Local) NewInstructionTemp(id GameType.GameId, tempId uint8) (ok bool) {
	if l.m.TryLock() {
		defer l.m.Unlock()
		var list []InstructionType.Instruction
		for _, v := range l.InstructionTempPool[id] {
			list = append(list, v)
		}
		l.InstructionLog[id][tempId] = list
		l.InstructionTempPool[id] = make(map[uint8]InstructionType.Instruction)
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameStatus(id GameType.GameId, status GameType.GameStatus) (ok bool) {
	if l.m.TryLock() {
		defer l.m.Unlock()
		l.GamePool[id].Status = status
		return true
	} else {
		return false
	}
}

func (l *Local) SetGameMap(id GameType.GameId, m *MapType.Map) (ok bool) {
	if l.m.TryLock() {
		defer l.m.Unlock()
		l.GamePool[id].Map = m
		return true
	} else {
		return false
	}
}

func (l *Local) SetUserStatus(id GameType.GameId, user GameType.User) (ok bool) {
	if l.m.TryLock() {
		defer l.m.Unlock()
		for _, u := range l.GamePool[id].UserList {
			if u.UserId == user.UserId {
				u.Status = user.Status
				return true
			}
		}
	}
	return false
}

func (l *Local) UpdateInstruction(id GameType.GameId, user GameType.User, instruction InstructionType.Instruction) (ok bool) {
	if l.m.TryLock() {
		defer l.m.Unlock()
		l.InstructionTempPool[id][user.UserId] = instruction
		return true
	} else {
		return false
	}
}

func (l *Local) GetCurrentMap(id GameType.GameId) *MapType.Map {
	if l.m.TryLock() {
		defer l.m.Unlock()
		m := *l.GamePool[id].Map
		return &m // Defend the outside modification
	} else {
		return nil
	}
}

func (l *Local) GetOriginalMap(mapId uint32) *MapType.Map {
	if l.m.TryLock() {
		defer l.m.Unlock()
		return pkg.Str2GameMap(mapId, l.OriginalMapStrPool[mapId])
	} else {
		return nil
	}
}

func (l *Local) GetCurrentGame(id GameType.GameId) *GameType.Game {
	if l.m.TryLock() {
		defer l.m.Unlock()
		return l.GamePool[id]
	}
	return nil
}

func (l *Local) CreateGame(mode GameType.GameMode) GameType.GameId {
	m := l.GetOriginalMap(rand.Uint32())
	g := &GameType.Game{
		Map:        m,
		Mode:       mode,
		Id:         GameType.GameId(rand.Uint32()),
		UserList:   []GameType.User{},
		CreateTime: time.Now().UnixMicro(),
		Status:     GameType.GameStatusWaiting,
		RoundNum:   0,
	}
	if l.m.TryLock() {
		defer l.m.Unlock()
		l.GamePool[g.Id] = g
		l.InstructionLog[g.Id] = make(map[uint8][]InstructionType.Instruction)
		return g.Id
	} else {
		return 0
	}
}
