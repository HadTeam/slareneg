package block

import (
	"server/game_logic/game_def"
)

var _ _type.Block = (*BaseBuilding)(nil)

type BaseBuilding struct {
	BaseBlock
}

func (block *BaseBuilding) Number() uint16 {
	return block.number
}

func (block *BaseBuilding) RoundStart(_ uint16) {
	if block.OwnerId() != 0 {
		block.number += 1
	}
}

func (*BaseBuilding) GetMoveStatus() _type.MoveStatus {
	return _type.MoveStatus{AllowMoveFrom: true, AllowMoveTo: true}
}

func (block *BaseBuilding) MoveFrom(number uint16) uint16 {
	var ret uint16
	if block.number <= number {
		ret = block.number - 1
		block.number = 1
	} else {
		ret = number
		block.number -= number
	}
	return ret
}

func (block *BaseBuilding) MoveTo(info _type.BlockVal) _type.Block {
	if block.ownerId != info.OwnerId {
		if block.number < info.Number {
			block.ownerId = info.OwnerId
			block.number = info.Number - block.number
		} else {
			block.number -= info.Number
		}
	} else {
		block.number += info.Number
	}
	return nil
}
