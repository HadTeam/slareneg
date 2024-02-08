package _map

import (
	"github.com/sirupsen/logrus"
	"server/game/block"
	"server/game/user"
	"server/utils/pkg"
	"strconv"
)

type MapSize struct{ W, H uint8 }

type mapInfo struct {
	size MapSize
	id   uint32
}

type Map struct {
	Blocks [][]block.Block
	mapInfo
}

func (m *Map) Size() MapSize {
	return m.size
}

func (m *Map) Id() uint32 {
	return m.id
}

func (m *Map) GetBlock(position block.Position) block.Block {
	return m.Blocks[position.Y-1][position.X-1]
}

func (m *Map) SetBlock(position block.Position, block block.Block) {
	m.Blocks[position.Y-1][position.X-1] = block
}

func (m *Map) HasBlocks() bool {
	if m.Blocks == nil {
		return false
	} else {
		return true
	}
}

func (m *Map) RoundStart(roundNum uint16) {
	for _, col := range m.Blocks {
		for _, b := range col {
			b.RoundStart(roundNum)
		}
	}
}

func (m *Map) RoundEnd(roundNum uint16) {
	for _, col := range m.Blocks {
		for _, b := range col {
			b.RoundEnd(roundNum)
		}
	}
}

func CreateEmptyMapWithInfo(mapId uint32, size MapSize) *Map {
	return &Map{
		Blocks: nil,
		mapInfo: mapInfo{
			size: size,
			id:   mapId,
		},
	}
}

func CreateMapWithBlocks(mapId uint32, blocks [][]block.Block) *Map {
	return &Map{
		Blocks: blocks,
		mapInfo: mapInfo{
			size: MapSize{uint8(len(blocks[0])), uint8(len(blocks))},
			id:   mapId,
		},
	}
}

func DebugOutput(p *Map, f func(block.Block) uint16) { // Only for debugging
	tmp := ""
	ex := func(i uint16) string {
		ex := ""
		if i < 10 {
			ex = " "
		}
		return ex + strconv.Itoa(int(i))
	}

	tmp += " *  "
	for i := uint16(1); i <= uint16(p.Size().W); i++ {
		tmp += ex(i) + " "
	}
	tmp += "\n"
	for rowNum, row := range p.Blocks {
		tmp += ex(uint16(rowNum+1)) + ": "
		for _, b := range row {
			tmp += ex(f(b)) + " "
		}
		tmp += "\n"
	}
	logrus.Tracef("\n%s\n", tmp)
}

func IsPositionLegal(position block.Position, size MapSize) bool {
	return 1 <= position.X && position.X <= size.W && 1 <= position.Y && position.Y <= size.H
}

type visibilityInFog [][]bool

func (m *Map) getVisibilityInFog(gameId uint16, teamId uint8, userList []user.User) *visibilityInFog {
	teamUserMap := make(map[uint16]bool)
	for _, v := range userList {
		if v.TeamId == teamId {
			teamUserMap[v.UserId] = true
		}
	}

	var ret *visibilityInFog
	if r, ok := pkg.TempPoolGet(gameId, "visibility"); ok {
		ret = r.(*visibilityInFog)
	} else {
		t := make(visibilityInFog, m.Size().H)
		ret = &t
		pkg.TempPoolPut(gameId, "visibility", &t)
	}

	for rowNum := uint8(0); rowNum <= m.Size().H-1; rowNum++ {
		(*ret)[rowNum] = make([]bool, m.Size().W)
	}

	light := func(x int, y int) {
		lightOffset := []struct {
			x, y int
		}{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
		for _, r := range lightOffset {
			ly := y + r.y
			lx := x + r.x
			if 0 <= ly && ly <= int(m.Size().H-1) && 0 <= lx && lx <= int(m.Size().W-1) {
				(*ret)[ly][lx] = true
			}
		}
	}

	for rowNum := uint8(0); rowNum <= m.Size().H-1; rowNum++ {
		for colNum := uint8(0); colNum <= m.Size().W-1; colNum++ {
			b := m.GetBlock(block.Position{X: colNum + 1, Y: rowNum + 1})
			if b.OwnerId() != user.Unknown {
				if _, exists := teamUserMap[b.OwnerId()]; exists {
					light(int(colNum), int(rowNum))
				}
			}
		}
	}

	return ret
}

type mapWithFog [][][]uint16

func (m *Map) GetProcessedMap(gameId uint16, teamId uint8, userList []user.User) *mapWithFog {
	vis := m.getVisibilityInFog(gameId, teamId, userList)

	var ret *mapWithFog
	if r, ok := pkg.TempPoolGet(gameId, "mapWithFog"); ok {
		ret = r.(*mapWithFog)
	} else {
		t := make(mapWithFog, m.Size().H)
		ret = &t
		pkg.TempPoolPut(gameId, "mapWithFog", &t)
	}

	for rowNum := uint8(0); rowNum <= m.Size().H-1; rowNum++ {
		(*ret)[rowNum] = make([][]uint16, m.Size().W)
		for colNum := uint8(0); colNum <= m.Size().W-1; colNum++ {
			b := m.GetBlock(block.Position{X: colNum + 1, Y: rowNum + 1})
			if (*vis)[rowNum][colNum] {
				(*ret)[rowNum][colNum] = []uint16{uint16(b.Meta().BlockId), b.OwnerId(), b.Number()}
			} else {
				(*ret)[rowNum][colNum] = []uint16{uint16(b.Meta().VisitFallBackType), user.Unknown, 0}
			}
		}
	}
	return ret
}
