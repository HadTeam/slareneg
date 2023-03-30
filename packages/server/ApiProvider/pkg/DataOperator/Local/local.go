package Local

import (
	"math/rand"
	"server/ApiProvider/pkg/DataOperator"
	"server/ApiProvider/pkg/InstructionType"
	"server/GameJudge/pkg/GameType"
	"server/MapHandler/pkg/MapOperator"
	"server/MapHandler/pkg/MapType"
	"sync"
)

var _ DataOperator.DataSource = (*Local)(nil)

type Local struct {
	m                   sync.Mutex
	GamePool            map[uint16]*GameType.Game
	OriginalMapStrPool  map[uint32]string
	InstructionTempPool map[uint16][]InstructionType.Instruction
	InstructionLog      map[uint16]map[uint8][]InstructionType.Instruction
}

var Pool Local

func init() {
	Pool = Local{
		GamePool:            make(map[uint16]*GameType.Game),
		OriginalMapStrPool:  make(map[uint32]string),
		InstructionTempPool: make(map[uint16][]InstructionType.Instruction),
		InstructionLog:      make(map[uint16]map[uint8][]InstructionType.Instruction),
	}

	Pool.OriginalMapStrPool[0] = "[\n[0,0,0,0,2],\n[0,2,0,0,0],\n[0,0,0,0,0],\n[0,3,3,0,3],\n[0,3,0,2,0]\n]"
}

func (l *Local) GetOriginalMap(mapId uint32) *MapType.Map {
	if l.m.TryLock() {
		defer l.m.Unlock()
		return MapOperator.Str2GameMap(mapId, l.OriginalMapStrPool[mapId])
	} else {
		return nil
	}
}
func (l *Local) GetCurrentGame(id GameType.GameId) *GameType.Game {
	if l.m.TryLock() {
		defer l.m.Unlock()
		ret, ok := l.GamePool[uint16(id)]
		if ok {
			return ret
		} else {
			return nil
		}
	} else {
		return nil
	}
}
func (l *Local) CreateGame(game *GameType.Game) GameType.GameId {
	id := uint16(rand.Uint32())
	game.Id = GameType.GameId(id)
	if l.m.TryLock() {
		defer l.m.Unlock()
		l.GamePool[id] = game
		l.InstructionLog[id] = make(map[uint8][]InstructionType.Instruction)
		return GameType.GameId(id)
	} else {
		return 0
	}
}

func (l *Local) PutInstructions(id GameType.GameId, instructions []InstructionType.Instruction) bool {
	if l.m.TryLock() {
		defer l.m.Unlock()
		l.InstructionTempPool[uint16(id)] = instructions
		return true
	} else {
		return false
	}
}

func (l *Local) AnnounceGameStart(gameId GameType.GameId) bool {
	//TODO implement me
	panic("implement me")
}

var ExampleInstruction = []InstructionType.Instruction{
	InstructionType.MoveInstruction{UserId: 1, Position: MapType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
	InstructionType.MoveInstruction{UserId: 2, Position: MapType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
}

func (l *Local) GetInstructionsFromTemp(id GameType.GameId, roundNum uint8) []InstructionType.Instruction {
	//if l.m.TryLock() {
	//	defer l.m.Unlock()
	//	return l.InstructionLog[uint16(id)][roundNum]
	//} else {
	//	return nil
	//}
	return ExampleInstruction
}

func (l *Local) AchieveInstructionTemp(id GameType.GameId, roundNum uint8) bool {
	if l.m.TryLock() {
		defer l.m.Unlock()
		l.InstructionLog[uint16(id)][roundNum] = l.InstructionTempPool[uint16(id)]
		l.InstructionTempPool[uint16(id)] = []InstructionType.Instruction{}
		return true
	} else {
		return false
	}
}

func (l *Local) PutMap(id GameType.GameId, m *MapType.Map) bool {
	if l.m.TryLock() {
		defer l.m.Unlock()
		l.GamePool[uint16(id)].Map = m
		return true
	} else {
		return false
	}
}
