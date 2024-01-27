package block

import (
	"server/utils/pkg/map/blockManager"
	"server/utils/pkg/map/type"
)

type Soldier struct {
	BaseBlock
}

var SoldierMeta = _type.Meta{
	BlockId:           1,
	Name:              "soldier",
	Description:       "",
	VisitFallBackType: BlankMeta.BlockId,
}

func init() {
	blockManager.Register(SoldierMeta, toBlockSoldier)
}

func toBlockSoldier(b _type.Block) _type.Block {
	var ret Soldier
	ret.number = b.Number()
	ret.ownerId = b.OwnerId()
	return _type.Block(&ret)
}

func (*Soldier) Meta() _type.Meta {
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

func (*Soldier) GetMoveStatus() _type.MoveStatus {
	return _type.MoveStatus{true, true}
}

func (block *Soldier) MoveFrom(number uint16) {
	block.number -= number
}

func (block *Soldier) MoveTo(info _type.BlockVal) _type.Block {

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
