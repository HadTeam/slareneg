package Blocks

import "server/Utils/pkg/MapType"

type BlockSoldier struct {
	BaseBlock
	number uint8
}

var blockSoldierMeta = MapType.BlockMeta{
	BlockId:     1,
	Name:        "soldier",
	Description: "",
}

func init() {
	MapType.RegisterBlockType(blockSoldierMeta, toBlockSoldier)
}

func toBlockSoldier(number uint8, ownerId uint8) MapType.Block {
	var ret BlockSoldier
	ret.number = number
	ret.ownerId = ownerId
	return MapType.Block(&ret)
}

func (*BlockSoldier) GetMeta() MapType.BlockMeta {
	return blockSoldierMeta
}

func (block *BlockSoldier) GetNumber() uint8 {
	return block.number
}

func (block *BlockSoldier) RoundStart(roundNum uint8) bool {
	if (roundNum%25)-1 == 0 && roundNum != 1 {
		block.number += 1
		return true
	}
	return false
}

func (block *BlockSoldier) MoveFrom(number uint8) {
	block.number -= number
}

func (block *BlockSoldier) MoveTo(ownerId uint8, number uint8) MapType.Block {

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