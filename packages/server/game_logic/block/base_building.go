package block

import (
	"server/game_logic/game_def"
)

var _ game_def.Block = (*BaseBuilding)(nil)

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

func (*BaseBuilding) GetMoveStatus() game_def.MoveStatus {
	return game_def.MoveStatus{AllowMoveFrom: true, AllowMoveTo: true}
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

func (block *BaseBuilding) MoveTo(info game_def.BlockVal) game_def.Block {
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
