package MapType

import (
	"encoding/json"
	"log"
	"server/Utils/pkg/InstructionType"
	"server/Utils/pkg/MapType/BlockType"
	"strconv"
)

type MapSize struct{ W, H uint8 }

type mapInfo struct {
	size MapSize
	id   uint32
}

type Map struct {
	blocks [][]BlockType.Block
	mapInfo
}

func (p *Map) Size() MapSize {
	return p.size
}

func (p *Map) Id() uint32 {
	return p.id
}

func (p *Map) GetBlock(position BlockType.Position) BlockType.Block {
	return p.blocks[position.Y-1][position.X-1]
}

func (p *Map) SetBlock(position BlockType.Position, block BlockType.Block) {
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
		for _, block := range col {
			block.RoundStart(roundNum)
		}
	}
}

func (p *Map) RoundEnd(roundNum uint16) bool {
	var ret bool
	for _, col := range p.blocks {
		for _, block := range col {
			if _, s := block.RoundEnd(roundNum); s {
				ret = true
			}
		}
	}
	return ret
}

func DebugOutput(p *Map, f func(BlockType.Block) uint16) { // Only for debugging
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
		for _, block := range row {
			tmp += ex(f(block)) + " "
		}
		tmp += "\n"
	}
	log.Printf("\n%s\n", tmp)
}

func isPositionLegal(position BlockType.Position, size MapSize) bool {
	return 1 <= position.X && position.X <= size.W && 1 <= position.Y && position.Y <= size.H
}

func (p *Map) Move(instruction InstructionType.Move) bool {
	var offsetX, offsetY int
	switch instruction.Towards {
	case InstructionType.MoveTowardsDown:
		{
			offsetX = 0
			offsetY = 1
		}
	case InstructionType.MoveTowardsUp:
		{
			offsetX = 0
			offsetY = -1
		}
	case InstructionType.MoveTowardsLeft:
		{
			offsetX = -1
			offsetY = 0
		}
	case InstructionType.MoveTowardsRight:
		{
			offsetX = 1
			offsetY = 0
		}
	}

	instructionPosition := BlockType.Position{instruction.Position.X, instruction.Position.Y}
	if !isPositionLegal(instructionPosition, p.size) {
		return false
	}

	newPosition := BlockType.Position{X: uint8(int(instruction.Position.X) + offsetX), Y: uint8(int(instruction.Position.Y) + offsetY)}
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
	if instruction.Number == 0 {
		instruction.Number = thisBlock.GetNumber()
	} else if instruction.Number == 65535 {
		instruction.Number = thisBlock.GetNumber() / 2
	}

	if thisBlock.GetNumber() < instruction.Number {
		return false
	}

	toBlock := p.GetBlock(newPosition)
	if !thisBlock.GetMoveStatus().AllowMoveFrom || !toBlock.GetMoveStatus().AllowMoveTo {
		return false
	}

	var toBlockNew BlockType.Block
	thisBlock.MoveFrom(instruction.Number)
	toBlockNew = toBlock.MoveTo(thisBlock.GetOwnerId(), instruction.Number)
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
	ret := make([][]BlockType.Block, size.H)
	for rowNum, row := range result {
		ret[rowNum] = make([]BlockType.Block, size.W)
		for colNum, typeId := range row {
			ret[rowNum][colNum] = BlockType.NewBlock(typeId, 0, 0)
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
	ret := make([][]BlockType.Block, size.H)
	for rowNum, row := range result {
		ret[rowNum] = make([]BlockType.Block, size.W)
		for colNum, blockInfo := range row {
			blockId := blockInfo[0]
			ownerId := blockInfo[1]
			number := blockInfo[2]

			block := BlockType.NewBlock(uint8(blockId), number, ownerId)

			ret[rowNum][colNum] = block
		}
	}
	return &Map{
		ret,
		mapInfo{size, mapId},
	}
}
