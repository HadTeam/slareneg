package block

import (
	"server/game_logic/block_manager"
	"server/game_logic/game_def"
)

var _ _type.Block = (*Blank)(nil)

type Blank struct {
	BaseBlock
}

var BlankMeta = _type.BlockMeta{
	BlockId:           0,
	Name:              "blank",
	Description:       "",
	VisitFallBackType: 0,
}

func init() {
	block_manager.Register(BlankMeta, toBlockBlank)
}

func (*Blank) Meta() _type.BlockMeta {
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
