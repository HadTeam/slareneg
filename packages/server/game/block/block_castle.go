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
	Register(CastleMeta, toBlockCastle)
}

func toBlockCastle(b Block) Block {
	var ret Castle
	if b.Number() == 0 {
		ret.number = uint16(30) + uint16(rand.Intn(30))
	} else {
		ret.number = b.Number()
	}
	ret.ownerId = b.OwnerId()
	return Block(&ret)
}

func (*Castle) Meta() Meta {
	return CastleMeta
}
