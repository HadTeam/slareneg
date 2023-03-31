package MapType

import (
	"encoding/json"
	"fmt"
	"server/Untils/pkg/InstructionType"
)

type MapSize struct{ X, Y uint8 }

type Map struct {
	Blocks [][]Block
	Size   MapSize
	MapId  uint32
}

func (p *Map) GetBlock(position BlockPosition) Block {
	return p.Blocks[position.Y][position.X]
}

func (p *Map) SetBlock(position BlockPosition, block Block) {
	p.Blocks[position.Y][position.X] = block
}

func (p *Map) RoundStart(roundNum uint8) {
	for _, col := range p.Blocks {
		for _, block := range col {
			block.roundStart(roundNum)
		}
	}
}

type GameOverSign bool

func (p *Map) RoundEnd(roundNum uint8) GameOverSign {
	var ret GameOverSign
	for _, col := range p.Blocks {
		for _, block := range col {
			if _, s := block.roundEnd(roundNum); s {
				ret = true
			}
		}
	}
	return ret
}

func (p *Map) OutputNumber() { // Only for debugging
	for _, col := range p.Blocks {
		for _, block := range col {
			fmt.Print(block.GetNumber())
		}
		fmt.Print("\n")
	}
}

func isPositionLegal(position BlockPosition, size MapSize) bool {
	return 1 <= position.X && position.X <= size.X && 1 <= position.Y && position.Y <= size.Y
}

func (p *Map) Move(instruction InstructionType.MoveInstruction) bool {
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

	newPosition := BlockPosition{X: uint8(int(instruction.Position.X) + offsetX), Y: uint8(int(instruction.Position.Y) + offsetY)}
	// It won't overflow 'cause the min value is 0

	if !isPositionLegal(instruction.Position, p.Size) && !isPositionLegal(newPosition, p.Size) {
		return false
	}
	thisBlock := p.GetBlock(instruction.Position)
	if thisBlock.GetNumber() < instruction.Number {
		return false
	}

	toBlock := p.GetBlock(newPosition)
	if !thisBlock.GetMoveStatus().AllowMoveFrom || !toBlock.GetMoveStatus().AllowMoveTo {
		return false
	}

	var toBlockNew Block
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
	if json.Unmarshal([]byte(originalMapStr), &result) != nil {
		return nil
	}
	size := MapSize{X: uint8(len(result[0])), Y: uint8(len(result))}
	ret := make([][]Block, size.Y)
	for rowNum, row := range result {
		ret[rowNum] = make([]Block, size.X)
		for colNum, typeId := range row {
			ret[rowNum][colNum] = ToBlockByTypeId(typeId, &BaseBlock{})
		}
	}
	return &Map{
		Blocks: ret,
		Size:   size,
		MapId:  mapId,
	}
}
