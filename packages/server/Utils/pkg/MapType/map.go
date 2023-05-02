package MapType

import (
	"log"
	"server/Utils/pkg/InstructionType"
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

func (p *Map) RoundStart(roundNum uint16) {
	for _, col := range p.Blocks {
		for _, block := range col {
			block.RoundStart(roundNum)
		}
	}
}

type GameOverSign bool

func (p *Map) RoundEnd(roundNum uint16) GameOverSign {
	var ret GameOverSign
	for _, col := range p.Blocks {
		for _, block := range col {
			if _, s := block.RoundEnd(roundNum); s {
				ret = true
			}
		}
	}
	return ret
}

func OutputNumber(p *Map) { // Only for debugging
	tmp := ""
	ex := func(i uint16) string {
		ex := ""
		if i < 10 {
			ex = " "
		}
		return ex + strconv.Itoa(int(i))
	}

	tmp += " *  "
	for i := uint16(1); i <= uint16(len(p.Blocks[0])); i++ {
		tmp += ex(i) + " "
	}
	tmp += "\n"
	for colNum, col := range p.Blocks {
		tmp += ex(uint16(colNum)) + ": "
		for _, block := range col {
			tmp += ex(block.GetNumber()) + " "
		}
		tmp += "\n"
	}
	log.Printf("\n%s\n", tmp)
}

func isPositionLegal(position BlockPosition, size MapSize) bool {
	return 1 <= position.X && position.X <= size.X && 1 <= position.Y && position.Y <= size.Y
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

	newPosition := BlockPosition{X: uint8(int(instruction.Position.X) + offsetX), Y: uint8(int(instruction.Position.Y) + offsetY)}
	// It won't overflow 'cause the min value is 0

	instructionPosition := BlockPosition{instruction.Position.X, instruction.Position.Y}
	if !isPositionLegal(instructionPosition, p.Size) && !isPositionLegal(newPosition, p.Size) {
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

	var toBlockNew Block
	thisBlock.MoveFrom(instruction.Number)
	toBlockNew = toBlock.MoveTo(thisBlock.GetOwnerId(), instruction.Number)
	if toBlockNew != nil {
		p.SetBlock(newPosition, toBlockNew)
	}
	return true
}
