package block

import (
	"server/utils/pkg/map/blockManager"
	"server/utils/pkg/map/type"
)

var _ _type.Block = (*Blank)(nil)

type Blank struct {
	BaseBlock
}

var BlankMeta = _type.Meta{
	BlockId:           0,
	Name:              "blank",
	Description:       "",
	VisitFallBackType: 0,
}

func init() {
	blockManager.Register(BlankMeta, toBlockBlank)
}

func (*Blank) Meta() _type.Meta {
	return BlankMeta
}

func (*Blank) GetMoveStatus() _type.MoveStatus {
	return _type.MoveStatus{false, true}
}

func toBlockBlank(_type.Block) _type.Block {
	return _type.Block(&Blank{})
}

func (b *Blank) MoveTo(_type.BlockVal) _type.Block {
	return toBlockSoldier(b)
}
