package instruction

import (
	"server/game/block"
	_map "server/game/map"
)

type Instruction interface{}

type MoveTowardsType string

const (
	MoveTowardsLeft  MoveTowardsType = "left"
	MoveTowardsRight MoveTowardsType = "right"
	MoveTowardsUp    MoveTowardsType = "up"
	MoveTowardsDown  MoveTowardsType = "down"
)

type Move struct {
	Position block.Position
	Towards  MoveTowardsType
	Number   uint16
}

type ForceStart struct {
	UserId uint16
	Status bool
}

type Surrender struct {
	UserId uint16
}

func Execute(userId uint16, m *_map.Map, ins Instruction) bool {
	var ret bool
	defer func() {
		if r := recover(); r != nil {
			ret = false
		}
	}()
	switch ins.(type) {
	case Move:
		{
			i := ins.(Move)
			if m.GetBlock(i.Position).OwnerId() != userId {
				panic("not the owner")
			} else {
				ret = MapMove(m, i)
			}
		}
	}

	return ret
}

func MapMove(p *_map.Map, inst Move) bool {
	var offsetX, offsetY int
	switch inst.Towards {
	case MoveTowardsDown:
		{
			offsetX = 0
			offsetY = 1
		}
	case MoveTowardsUp:
		{
			offsetX = 0
			offsetY = -1
		}
	case MoveTowardsLeft:
		{
			offsetX = -1
			offsetY = 0
		}
	case MoveTowardsRight:
		{
			offsetX = 1
			offsetY = 0
		}
	}

	instructionPosition := block.Position{X: inst.Position.X, Y: inst.Position.Y}
	if !_map.IsPositionLegal(instructionPosition, p.Size()) {
		panic("instruction position is not legal")
		return false
	}

	newPosition := block.Position{X: uint8(int(inst.Position.X) + offsetX), Y: uint8(int(inst.Position.Y) + offsetY)}
	// It won't overflow 'cause the min value is 0
	if !_map.IsPositionLegal(newPosition, p.Size()) {
		panic("new position is not legal")
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
		panic("not enough number")
		return false
	}

	toBlock := p.GetBlock(newPosition)
	if !thisBlock.GetMoveStatus().AllowMoveFrom || !toBlock.GetMoveStatus().AllowMoveTo {
		panic("move status not allowed")
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
