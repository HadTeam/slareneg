package Blocks

import "server/Utils/pkg/MapType"

var _ MapType.Block = (*BaseBuilding)(nil)

type BaseBuilding struct {
	BaseBlock
}

func (block *BaseBuilding) GetNumber() uint16 {
	return block.number
}

func (block *BaseBuilding) RoundStart(_ uint16) bool {
	if block.GetOwnerId() != 0 {
		block.number += 1
	}
	return true
}

func (block *BaseBuilding) MoveFrom(number uint16) {
	block.number -= number
}

func (block *BaseBuilding) MoveTo(ownerId uint16, number uint16) MapType.Block {
	if block.ownerId != ownerId {
		if block.number < number {
			block.ownerId = ownerId
			block.number = number - block.number
		} else {
			block.number -= number
		}
	} else {
		block.number += number
	}
	return nil
}
