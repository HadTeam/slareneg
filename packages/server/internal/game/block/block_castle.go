package block

import (
	"math/rand"
)

var _ Block = (*Castle)(nil)

type Castle struct {
	BaseBuilding
}

var CastleMeta Meta
var CastleName Name

func init() {
	CastleName = Register("castle", "", toBlockCastle)
	CastleMeta = GetMetaByName[CastleName]
}

func toBlockCastle(b Block) Block {
	var ret Castle
	if b.Num() == 0 {
		ret.num = Num(30) + Num(rand.Intn(30))
	} else {
		ret.num = b.Num()
	}
	ret.owner = b.Owner()
	return &ret
}

func (*Castle) Meta() Meta {
	return CastleMeta
}

func (block *Castle) Fog(isOwner bool, isSight bool) Block {
	if isOwner || isSight {
		// Owner can always see the real castle
		return block
	}
	return &Mountain{}
}
