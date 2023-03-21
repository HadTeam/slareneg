package Local

import (
	"math/rand"
	"server/ApiProvider/pkg/DataOperator"
	"server/ApiProvider/pkg/InstructionType"
	"server/GameJudge/pkg/GameType"
	"server/MapHandler/pkg/MapOperator"
	"server/MapHandler/pkg/MapType"
)

var _ DataOperator.DataSource = (*Local)(nil)

type Local struct {
	GamePool           map[uint16]*GameType.Game
	OriginalMapStrPool map[uint32]string
	InstructionPool    map[uint16][]InstructionType.Instruction
}

var Pool Local

func init() {
	Pool.GamePool = make(map[uint16]*GameType.Game)
	Pool.OriginalMapStrPool = make(map[uint32]string)
	Pool.OriginalMapStrPool[0] = "[\n[0,0,0,0,2],\n[0,2,0,0,0],\n[0,0,0,0,0],\n[0,3,3,0,3],\n[0,3,0,2,0]\n]"

}

func (l *Local) GetOriginalMap(mapId uint32) *MapType.Map {
	return MapOperator.Str2GameMap(mapId, l.OriginalMapStrPool[mapId])
}
func (l *Local) GetCurrentGame(id GameType.GameId) *GameType.Game {
	ret, ok := l.GamePool[uint16(id)]
	if ok {
		return ret
	} else {
		return nil
	}
}
func (l *Local) CreateGame(game *GameType.Game) GameType.GameId {
	id := uint16(rand.Uint32())
	game.Id = GameType.GameId(id)
	l.GamePool[id] = game
	return GameType.GameId(id)
}

func (l *Local) PutInstructions(id GameType.GameId, instructions []InstructionType.Instruction) bool {
	l.InstructionPool[uint16(id)] = instructions
	return true
}

func (l *Local) AnnounceGameStart(gameId GameType.GameId) bool {
	// TODO
	return true
}

//var ExampleInstruction = []InstructionType.Instruction{
//	InstructionType.MoveInstruction{UserId: 1, Position: MapType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
//	InstructionType.MoveInstruction{UserId: 2, Position: MapType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
//}

func (l *Local) GetCurrentInstruction(id GameType.GameId) []InstructionType.Instruction {
	return l.InstructionPool[uint16(id)]
}

func (l *Local) NewRound(id GameType.GameId, m *MapType.Map, roundNum uint8) bool {
	g := l.GamePool[uint16(id)]
	g.Map = m
	g.RoundNum = roundNum
	return true
}
