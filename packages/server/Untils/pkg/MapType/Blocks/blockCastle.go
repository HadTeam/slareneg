package Blocks

import "server/Untils/pkg/MapType"

var _ MapType.Block = (*BlockCastle)(nil)

type BlockCastle struct {
	BaseBuilding
}

var blockCastleMeta = MapType.BlockMeta{
	BlockId:     3,
	Name:        "castle",
	Description: "",
}

func init() {
	MapType.RegisterBlockType(blockCastleMeta, toBlockCastle)
}

func toBlockCastle(number uint8, ownerId uint8) MapType.Block {
	var ret BlockKing
	ret.number = number
	ret.ownerId = ownerId
	return MapType.Block(&ret)
}
