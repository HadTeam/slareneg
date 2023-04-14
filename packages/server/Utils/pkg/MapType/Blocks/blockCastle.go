package Blocks

import (
	"math/rand"
	"server/Utils/pkg/MapType"
)

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
	if number == 0 {
		ret.number = uint8(30) + uint8(rand.Intn(30))
	} else {
		ret.number = number
	}
	ret.ownerId = ownerId
	return MapType.Block(&ret)
}
