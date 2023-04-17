package Blocks

import "server/Utils/pkg/MapType"

var _ MapType.Block = (*BlockKing)(nil)

type BlockKing struct {
	BaseBuilding
}

var blockKingMeta = MapType.BlockMeta{
	BlockId:     2,
	Name:        "king",
	Description: "",
}

func init() {
	MapType.RegisterBlockType(blockKingMeta, toBlockKing)
}

func toBlockKing(number uint16, ownerId uint16) MapType.Block {
	var ret BlockKing
	ret.number = number
	ret.ownerId = ownerId
	return MapType.Block(&ret)
}

func (block *BlockKing) RoundEnd(_ uint16) (bool, MapType.GameOverSign) {
	if block.number <= 0 {
		return true, true
	}
	return false, false
}

func (*BlockKing) GetMeta() MapType.BlockMeta {
	return blockKingMeta
}
