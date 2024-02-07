package _map

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"server/game/block"
	"server/game/instruction"
	"strconv"
	"strings"
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

func (p *Map) Size() MapSize {
	return p.size
}

func (p *Map) Id() uint32 {
	return p.id
}

func (p *Map) GetBlock(position block.Position) block.Block {
	return p.Blocks[position.Y-1][position.X-1]
}

func (p *Map) SetBlock(position block.Position, block block.Block) {
	p.Blocks[position.Y-1][position.X-1] = block
}

func (p *Map) HasBlocks() bool {
	if p.Blocks == nil {
		return false
	} else {
		return true
	}
}

func (p *Map) RoundStart(roundNum uint16) {
	for _, col := range p.Blocks {
		for _, b := range col {
			b.RoundStart(roundNum)
		}
	}
}

func (p *Map) RoundEnd(roundNum uint16) {
	for _, col := range p.Blocks {
		for _, b := range col {
			b.RoundEnd(roundNum)
		}
	}
}

func CreateMapWithInfo(mapId uint32, size MapSize) *Map {
	return &Map{
		Blocks: nil,
		mapInfo: mapInfo{
			size: size,
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

func isPositionLegal(position block.Position, size MapSize) bool {
	return 1 <= position.X && position.X <= size.W && 1 <= position.Y && position.Y <= size.H
}

func (p *Map) Move(inst instruction.Move) bool {
	var offsetX, offsetY int
	switch inst.Towards {
	case instruction.MoveTowardsDown:
		{
			offsetX = 0
			offsetY = 1
		}
	case instruction.MoveTowardsUp:
		{
			offsetX = 0
			offsetY = -1
		}
	case instruction.MoveTowardsLeft:
		{
			offsetX = -1
			offsetY = 0
		}
	case instruction.MoveTowardsRight:
		{
			offsetX = 1
			offsetY = 0
		}
	}

	instructionPosition := block.Position{X: inst.Position.X, Y: inst.Position.Y}
	if !isPositionLegal(instructionPosition, p.size) {
		return false
	}

	newPosition := block.Position{X: uint8(int(inst.Position.X) + offsetX), Y: uint8(int(inst.Position.Y) + offsetY)}
	// It won't overflow 'cause the min value is 0
	if !isPositionLegal(newPosition, p.size) {
		return false
	}

	thisBlock := p.GetBlock(instructionPosition)

	/*
	 * Special case
	 * 0 => select all
	 * 1 => select half
	 */
	if inst.Number == 0 {
		inst.Number = thisBlock.Number()
	} else if inst.Number == 65535 {
		inst.Number = thisBlock.Number() / 2
	}

	if thisBlock.Number() < inst.Number {
		return false
	}

	toBlock := p.GetBlock(newPosition)
	if !thisBlock.GetMoveStatus().AllowMoveFrom || !toBlock.GetMoveStatus().AllowMoveTo {
		return false
	}

	var toBlockNew block.Block
	hasMovedNum := thisBlock.MoveFrom(inst.Number)
	toBlockNew = toBlock.MoveTo(block.Val{Number: hasMovedNum, OwnerId: thisBlock.OwnerId()})
	if toBlockNew != nil {
		p.SetBlock(newPosition, toBlockNew)
	}
	return true
}

// Str2GameMap TODO: Add unit test
func Str2GameMap(mapId uint32, originalMapStr string) *Map {
	var result [][]uint8
	if err := json.Unmarshal([]byte(originalMapStr), &result); err != nil {
		logrus.Panic(err)
		return nil
	}
	size := MapSize{W: uint8(len(result[0])), H: uint8(len(result))}
	ret := make([][]block.Block, size.H)
	for rowNum, row := range result {
		ret[rowNum] = make([]block.Block, size.W)
		for colNum, typeId := range row {
			ret[rowNum][colNum] = block.NewBlock(typeId, 0, 0)
		}
	}
	return &Map{
		ret,
		mapInfo{size, mapId},
	}
}

func FullStr2GameMap(mapId uint32, originalMapStr string) *Map {
	var result [][][]uint16
	if err := json.Unmarshal([]byte(originalMapStr), &result); err != nil {
		logrus.Panic(err)
		return nil
	}
	size := MapSize{W: uint8(len(result[0])), H: uint8(len(result))}
	ret := make([][]block.Block, size.H)
	for rowNum, row := range result {
		ret[rowNum] = make([]block.Block, size.W)
		for colNum, blockInfo := range row {
			blockId := blockInfo[0]
			ownerId := blockInfo[1]
			number := blockInfo[2]

			newBlock := block.NewBlock(uint8(blockId), number, ownerId)

			ret[rowNum][colNum] = newBlock
		}
	}
	return &Map{
		ret,
		mapInfo{size, mapId},
	}
}

// To marshal uint8 array to json as number array
type uint8Array []uint8

func (a uint8Array) MarshalJSON() ([]byte, error) {
	var res string
	if a == nil {
		res = "null"
	} else {
		res = strings.Join(strings.Fields(fmt.Sprintf("%d", a)), ",")
	}
	return []byte(res), nil
}

type mapJsonStruct struct {
	Mappings struct {
		Block []string `json:"block"`
		Owner []uint16 `json:"owner,omitempty"`
	} `json:"mappings"`
	Type   []uint8Array `json:"game_def"`
	Owner  []uint8Array `json:"owner,omitempty"`
	Number [][]uint16   `json:"number,omitempty"`
}

func JsonStrToMap(jsonStr string) *Map {
	var res mapJsonStruct
	if err := json.Unmarshal([]byte(jsonStr), &res); err != nil {
		logrus.Panic(err)
		return nil
	}

	// process original mapping
	blockMapping := make(map[uint8]uint8)
	for i, v := range res.Mappings.Block {
		blockMapping[uint8(i)] = block.GetBlockIdByName[v]
	}

	if res.Mappings.Owner == nil || res.Owner == nil {
		res.Mappings.Owner = nil
		res.Owner = nil
	} else {
		// add blank at 0 owner id
		res.Mappings.Owner = append([]uint16{0}, res.Mappings.Owner...)
	}

	var blocks [][]block.Block
	if (res.Number != nil && len(res.Type) != len(res.Number)) || (res.Owner != nil && len(res.Type) != len(res.Owner)) {
		logrus.Panic("original block game_def, number, owner id must have the same size")
	}
	blocks = make([][]block.Block, len(res.Type))

	for i, v := range res.Type {
		if (res.Number != nil && len(v) != len(res.Number[i])) || (res.Owner != nil && len(v) != len(res.Owner[i])) {
			logrus.Panic("original block game_def, number, owner id must have the same size")
		}
		blocks[i] = make([]block.Block, len(v))

		for j, typeId := range v {
			n := uint16(0)
			o := uint16(0)
			if res.Owner != nil {
				if res.Owner[i][j] >= uint8(len(res.Mappings.Owner)) {
					logrus.Panic("original owner id must be less than owner id mapping size")
				}

				o = res.Mappings.Owner[res.Owner[i][j]]
			}
			if res.Number != nil {
				n = res.Number[i][j]
			}

			if typeId >= uint8(len(blockMapping)) {
				logrus.Panic("original block game_def must be less than block game_def mapping size")
			}

			blocks[i][j] = block.NewBlock(blockMapping[typeId], n, o)
		}
	}
	return &Map{
		blocks,
		mapInfo{MapSize{uint8(len(blocks[0])), uint8(len(blocks))}, 0},
	}
}

func MapToJsonStr(m *Map) string {
	var ret = mapJsonStruct{}
	typeMapping := make(map[uint8]uint8)
	ownerMapping := make(map[uint16]uint8)
	ownerMapping[0] = 0 // blank owner id

	ret.Type = make([]uint8Array, m.size.H)
	ret.Owner = make([]uint8Array, m.size.H)
	ret.Number = make([][]uint16, m.size.H)

	for i, row := range m.Blocks {
		ret.Type[i] = make(uint8Array, m.size.W)
		ret.Owner[i] = make(uint8Array, m.size.W)
		ret.Number[i] = make([]uint16, m.size.W)

		for j, b := range row {
			if _, ok := typeMapping[b.Meta().BlockId]; !ok {
				typeMapping[b.Meta().BlockId] = uint8(len(typeMapping))
			}
			if _, ok := ownerMapping[b.OwnerId()]; !ok {
				ownerMapping[b.OwnerId()] = uint8(len(ownerMapping))
			}

			ret.Type[i][j] = typeMapping[b.Meta().BlockId]
			ret.Owner[i][j] = ownerMapping[b.OwnerId()]
			ret.Number[i][j] = b.Number()
		}
	}

	for k := range typeMapping {
		ret.Mappings.Block = append(ret.Mappings.Block, block.GetMetaById[k].Name)
	}
	for k := range ownerMapping {
		ret.Mappings.Owner = append(ret.Mappings.Owner, k)
	}

	ret.Mappings.Owner = ret.Mappings.Owner[1:] // remove blank owner id

	logrus.Info(ret)

	if retJson, err := json.Marshal(ret); err != nil {
		logrus.Panic(err)
		return ""
	} else {
		return string(retJson)
	}
}
