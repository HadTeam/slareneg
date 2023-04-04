package Blocks

import "server/Untils/pkg/MapType"

var _ MapType.Block = (*BaseBuilding)(nil)

type BaseBuilding struct {
	BaseBlock
	number uint8
}

func (block *BaseBuilding) GetNumber() uint8 {
	return block.number
}

func (block *BaseBuilding) roundStart(_ uint8) bool {
	block.number += 1
	return true
}

func (block *BaseBuilding) MoveFrom(number uint8) {
	block.number -= number
}

func (block *BaseBuilding) MoveTo(ownerId uint8, number uint8) MapType.Block {
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
