package Blocks

import "server/Utils/pkg/MapType"

var _ MapType.Block = (*BlockKing)(nil)

type BlockKing struct {
	BaseBuilding
}

var BlockKingMeta = MapType.BlockMeta{
	BlockId:           2,
	Name:              "king",
	Description:       "",
	VisitFallBackType: BlockCastleMeta.BlockId,
}

func init() {
	MapType.RegisterBlockType(BlockKingMeta, toBlockKing)
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
	return BlockKingMeta
}
