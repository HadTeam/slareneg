package block

import (
	"server/game_logic/block_manager"
	"server/game_logic/game_def"
)

var _ game_def.Block = (*Soldier)(nil)

type Soldier struct {
	BaseBlock
}

var SoldierMeta = game_def.BlockMeta{
	BlockId:           1,
	Name:              "soldier",
	Description:       "",
	VisitFallBackType: BlankMeta.BlockId,
}

func init() {
	block_manager.Register(SoldierMeta, toBlockSoldier)
}

func toBlockSoldier(b game_def.Block) game_def.Block {
	var ret Soldier
	ret.number = b.Number()
	ret.ownerId = b.OwnerId()
	return game_def.Block(&ret)
}

func (*Soldier) Meta() game_def.BlockMeta {
	return SoldierMeta
}

func (block *Soldier) Number() uint16 {
	return block.number
}

func (block *Soldier) RoundStart(roundNum uint16) {
	if (roundNum%25)-1 == 0 && roundNum != 1 {
		block.number += 1
	}
}

func (*Soldier) GetMoveStatus() game_def.MoveStatus {
	return game_def.MoveStatus{true, true}
}

func (block *Soldier) MoveFrom(number uint16) uint16 {
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

func (block *Soldier) MoveTo(info game_def.BlockVal) game_def.Block {
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
