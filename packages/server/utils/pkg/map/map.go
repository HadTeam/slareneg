package _map

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"server/utils/pkg/instruction"
	"server/utils/pkg/map/blockManager"
	"server/utils/pkg/map/type"
	"strconv"
	"strings"
)

type MapSize struct{ W, H uint8 }

type mapInfo struct {
	size MapSize
	id   uint32
}

type Map struct {
	blocks [][]_type.Block
	mapInfo
}

func (p *Map) Size() MapSize {
	return p.size
}

func (p *Map) Id() uint32 {
	return p.id
}

func (p *Map) GetBlock(position _type.Position) _type.Block {
	return p.blocks[position.Y-1][position.X-1]
}

func (p *Map) SetBlock(position _type.Position, block _type.Block) {
	p.blocks[position.Y-1][position.X-1] = block
}

func (p *Map) HasBlocks() bool {
	if p.blocks == nil {
		return false
	} else {
		return true
	}
}

func (p *Map) RoundStart(roundNum uint16) {
	for _, col := range p.blocks {
		for _, b := range col {
			b.RoundStart(roundNum)
		}
	}
}

func (p *Map) RoundEnd(roundNum uint16) {
	for _, col := range p.blocks {
		for _, b := range col {
			b.RoundEnd(roundNum)
		}
	}
}

func CreateMapWithInfo(mapId uint32, size MapSize) *Map {
	return &Map{
		blocks: nil,
		mapInfo: mapInfo{
			size: size,
			id:   mapId,
		},
	}
}

func DebugOutput(p *Map, f func(_type.Block) uint16) { // Only for debugging
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
	for rowNum, row := range p.blocks {
		tmp += ex(uint16(rowNum)) + ": "
		for _, b := range row {
			tmp += ex(f(b)) + " "
		}
		tmp += "\n"
	}
	logrus.Tracef("\n%s\n", tmp)
}

func isPositionLegal(position _type.Position, size MapSize) bool {
	return 1 <= position.X && position.X <= size.W && 1 <= position.Y && position.Y <= size.H
}

func (p *Map) Move(ins instruction.Move) bool {
	var offsetX, offsetY int
	switch ins.Towards {
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

	instructionPosition := _type.Position{X: ins.Position.X, Y: ins.Position.Y}
	if !isPositionLegal(instructionPosition, p.size) {
		return false
	}

	newPosition := _type.Position{X: uint8(int(ins.Position.X) + offsetX), Y: uint8(int(ins.Position.Y) + offsetY)}
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
	if ins.Number == 0 {
		ins.Number = thisBlock.Number()
	} else if ins.Number == 65535 {
		ins.Number = thisBlock.Number() / 2
	}

	if thisBlock.Number() < ins.Number {
		return false
	}

	toBlock := p.GetBlock(newPosition)
	if !thisBlock.GetMoveStatus().AllowMoveFrom || !toBlock.GetMoveStatus().AllowMoveTo {
		return false
	}

	var toBlockNew _type.Block
	thisBlock.MoveFrom(ins.Number)
	toBlockNew = toBlock.MoveTo(_type.BlockVal{Number: thisBlock.OwnerId(), OwnerId: ins.Number})
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
	ret := make([][]_type.Block, size.H)
	for rowNum, row := range result {
		ret[rowNum] = make([]_type.Block, size.W)
		for colNum, typeId := range row {
			ret[rowNum][colNum] = blockManager.NewBlock(typeId, 0, 0)
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
	ret := make([][]_type.Block, size.H)
	for rowNum, row := range result {
		ret[rowNum] = make([]_type.Block, size.W)
		for colNum, blockInfo := range row {
			blockId := blockInfo[0]
			ownerId := blockInfo[1]
			number := blockInfo[2]

			newBlock := blockManager.NewBlock(uint8(blockId), number, ownerId)

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
	Type   []uint8Array `json:"type"`
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
		blockMapping[uint8(i)] = blockManager.GetBlockIdByName[v]
	}

	if res.Mappings.Owner == nil || res.Owner == nil {
		res.Mappings.Owner = nil
		res.Owner = nil
	} else {
		// add blank at 0 owner id
		res.Mappings.Owner = append([]uint16{0}, res.Mappings.Owner...)
	}

	var blocks [][]_type.Block
	if (res.Number != nil && len(res.Type) != len(res.Number)) || (res.Owner != nil && len(res.Type) != len(res.Owner)) {
		logrus.Panic("original block type, number, owner id must have the same size")
	}
	blocks = make([][]_type.Block, len(res.Type))

	for i, v := range res.Type {
		if (res.Number != nil && len(v) != len(res.Number[i])) || (res.Owner != nil && len(v) != len(res.Owner[i])) {
			logrus.Panic("original block type, number, owner id must have the same size")
		}
		blocks[i] = make([]_type.Block, len(v))

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
				logrus.Panic("original block type must be less than block type mapping size")
			}

			blocks[i][j] = blockManager.NewBlock(blockMapping[typeId], n, o)
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

	for i, row := range m.blocks {
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
		ret.Mappings.Block = append(ret.Mappings.Block, blockManager.GetMetaById[k].Name)
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
