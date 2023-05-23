package _map

import (
	"encoding/json"
	"log"
	"server/utils/pkg/instruction"
	"server/utils/pkg/map/block"
	"strconv"
)

type MapSize struct{ W, H uint8 }

type mapInfo struct {
	size MapSize
	id   uint32
}

type Map struct {
	blocks [][]block.Block
	mapInfo
}

func (p *Map) Size() MapSize {
	return p.size
}

func (p *Map) Id() uint32 {
	return p.id
}

func (p *Map) GetBlock(position block.Position) block.Block {
	return p.blocks[position.Y-1][position.X-1]
}

func (p *Map) SetBlock(position block.Position, block block.Block) {
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
	for rowNum, row := range p.blocks {
		tmp += ex(uint16(rowNum)) + ": "
		for _, b := range row {
			tmp += ex(f(b)) + " "
		}
		tmp += "\n"
	}
	log.Printf("\n%s\n", tmp)
}

func isPositionLegal(position block.Position, size MapSize) bool {
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

	instructionPosition := block.Position{X: ins.Position.X, Y: ins.Position.Y}
	if !isPositionLegal(instructionPosition, p.size) {
		//log.Printf("[Move] instruction position illegal")
		return false
	}

	newPosition := block.Position{X: uint8(int(ins.Position.X) + offsetX), Y: uint8(int(ins.Position.Y) + offsetY)}
	// It won't overflow 'cause the min value is 0
	if !isPositionLegal(newPosition, p.size) {
		//log.Printf("[Move] new position illegal")
		return false
	}

	thisBlock := p.GetBlock(instructionPosition)

	/*
	 * Special case
	 * 0 => select all
	 * 1 => select half
	 */
	if ins.Number == 0 {
		ins.Number = thisBlock.GetNumber()
	} else if ins.Number == 65535 {
		ins.Number = thisBlock.GetNumber() / 2
	}

	if thisBlock.GetNumber() < ins.Number {
		//log.Printf("[Move] number isn't enough: %#v %#v", thisBlock.GetNumber(), instruction.Number)
		return false
	}

	toBlock := p.GetBlock(newPosition)
	if !thisBlock.GetMoveStatus().AllowMoveFrom || !toBlock.GetMoveStatus().AllowMoveTo {
		//log.Printf("[Move] not allow to move: %s(%#v) -> %s(%#v)", thisBlock.GetMeta().Name,
		//	thisBlock.GetMoveStatus().AllowMoveFrom, toBlock.GetMeta().Name, toBlock.GetMoveStatus().AllowMoveTo)
		return false
	}

	var toBlockNew block.Block
	thisBlock.MoveFrom(ins.Number)
	toBlockNew = toBlock.MoveTo(thisBlock.GetOwnerId(), ins.Number)
	if toBlockNew != nil {
		p.SetBlock(newPosition, toBlockNew)
	}
	return true
}

// Str2GameMap TODO: Add unit test
func Str2GameMap(mapId uint32, originalMapStr string) *Map {
	var result [][]uint8
	if err := json.Unmarshal([]byte(originalMapStr), &result); err != nil {
		panic(err)
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
		panic(err)
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

			block := block.NewBlock(uint8(blockId), number, ownerId)

			ret[rowNum][colNum] = block
		}
	}
	return &Map{
		ret,
		mapInfo{size, mapId},
	}
}
