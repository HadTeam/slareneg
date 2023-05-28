package block

import (
	"math/rand"
)

var _ Block = (*Castle)(nil)

type Castle struct {
	BaseBuilding
}

var CastleMeta = Meta{
	BlockId:           3,
	Name:              "castle",
	Description:       "",
	VisitFallBackType: 3,
}

func init() {
	RegisterBlockType(CastleMeta, toBlockCastle)
}

func toBlockCastle(number uint16, ownerId uint16) Block {
	var ret Castle
	if number == 0 {
		ret.number = uint16(30) + uint16(rand.Intn(30))
	} else {
		ret.number = number
	}
	ret.ownerId = ownerId
	return Block(&ret)
}

func (*Castle) GetMeta() Meta {
	return CastleMeta
}
