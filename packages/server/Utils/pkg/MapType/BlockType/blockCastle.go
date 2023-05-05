package BlockType

import (
	"math/rand"
)

var _ Block = (*BlockCastle)(nil)

type BlockCastle struct {
	BaseBuilding
}

var BlockCastleMeta = BlockMeta{
	BlockId:           3,
	Name:              "castle",
	Description:       "",
	VisitFallBackType: 3,
}

func init() {
	RegisterBlockType(BlockCastleMeta, toBlockCastle)
}

func toBlockCastle(number uint16, ownerId uint16) Block {
	var ret BlockCastle
	if number == 0 {
		ret.number = uint16(30) + uint16(rand.Intn(30))
	} else {
		ret.number = number
	}
	ret.ownerId = ownerId
	return Block(&ret)
}

func (*BlockCastle) GetMeta() BlockMeta {
	return BlockCastleMeta
}
