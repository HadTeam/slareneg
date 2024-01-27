package block

import (
	"math/rand"
	"server/utils/pkg/map/blockManager"
	"server/utils/pkg/map/type"
)

var _ _type.Block = (*Castle)(nil)

type Castle struct {
	BaseBuilding
}

var CastleMeta = _type.Meta{
	BlockId:           3,
	Name:              "castle",
	Description:       "",
	VisitFallBackType: 3,
}

func init() {
	blockManager.Register(CastleMeta, toBlockCastle)
}

func toBlockCastle(b _type.Block) _type.Block {
	var ret Castle
	if b.Number() == 0 {
		ret.number = uint16(30) + uint16(rand.Intn(30))
	} else {
		ret.number = b.Number()
	}
	ret.ownerId = b.OwnerId()
	return _type.Block(&ret)
}

func (*Castle) Meta() _type.Meta {
	return CastleMeta
}
