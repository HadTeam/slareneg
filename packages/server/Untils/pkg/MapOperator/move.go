package MapOperator

import (
	"server/Untils/pkg/InstructionType"
	"server/Untils/pkg/MapType"
)

func isPositionLegal(position MapType.BlockPosition, size MapType.MapSize) bool {
	return 1 <= position.X && position.X <= size.X && 1 <= position.Y && position.Y <= size.Y
}

func Move(p *MapType.Map, instruction InstructionType.MoveInstruction) bool {
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

	newPosition := MapType.BlockPosition{X: uint8(int(instruction.Position.X) + offsetX), Y: uint8(int(instruction.Position.Y) + offsetY)}
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

	var toBlockNew MapType.Block
	thisBlock.MoveFrom(instruction.Number)
	toBlockNew = toBlock.MoveTo(thisBlock.GetOwnerId(), instruction.Number)
	if toBlockNew != nil {
		p.SetBlock(newPosition, toBlockNew)
	}
	return true
}
