package Blocks

import "server/Utils/pkg/MapType"

type BlockSoldier struct {
	BaseBlock
	number uint16
}

var BlockSoldierMeta = MapType.BlockMeta{
	BlockId:           1,
	Name:              "soldier",
	Description:       "",
	VisitFallBackType: BlockBlankMeta.BlockId,
}

func init() {
	MapType.RegisterBlockType(BlockSoldierMeta, toBlockSoldier)
}

func toBlockSoldier(number uint16, ownerId uint16) MapType.Block {
	var ret BlockSoldier
	ret.number = number
	ret.ownerId = ownerId
	return MapType.Block(&ret)
}

func (*BlockSoldier) GetMeta() MapType.BlockMeta {
	return BlockSoldierMeta
}

func (block *BlockSoldier) GetNumber() uint16 {
	return block.number
}

func (block *BlockSoldier) RoundStart(roundNum uint16) bool {
	if (roundNum%25)-1 == 0 && roundNum != 1 {
		block.number += 1
		return true
	}
	return false
}

func (block *BlockSoldier) MoveFrom(number uint16) {
	block.number -= number
}

func (block *BlockSoldier) MoveTo(ownerId uint16, number uint16) MapType.Block {

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
