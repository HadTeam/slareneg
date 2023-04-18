package Blocks

import (
	"math/rand"
	"server/Utils/pkg/MapType"
)

var _ MapType.Block = (*BlockCastle)(nil)

type BlockCastle struct {
	BaseBuilding
}

var BlockCastleMeta = MapType.BlockMeta{
	BlockId:           3,
	Name:              "castle",
	Description:       "",
	VisitFallBackType: 3,
}

func init() {
	MapType.RegisterBlockType(BlockCastleMeta, toBlockCastle)
}

func toBlockCastle(number uint16, ownerId uint16) MapType.Block {
	var ret BlockKing
	if number == 0 {
		ret.number = uint16(30) + uint16(rand.Intn(30))
	} else {
		ret.number = number
	}
	ret.ownerId = ownerId
	return MapType.Block(&ret)
}

func (*BlockCastle) GetMeta() MapType.BlockMeta {
	return BlockCastleMeta
}
