package block

import (
	"server/game_logic/block_manager"
	"server/game_logic/game_def"
)

var _ game_def.Block = (*Blank)(nil)

type Blank struct {
	BaseBlock
}

var BlankMeta = game_def.BlockMeta{
	BlockId:           0,
	Name:              "blank",
	Description:       "",
	VisitFallBackType: 0,
}

func init() {
	block_manager.Register(BlankMeta, toBlockBlank)
}

func (*Blank) Meta() game_def.BlockMeta {
	return BlankMeta
}

func (*Blank) GetMoveStatus() game_def.MoveStatus {
	return game_def.MoveStatus{false, true}
}

func toBlockBlank(game_def.Block) game_def.Block {
	return game_def.Block(&Blank{})
}

func (b *Blank) MoveTo(game_def.BlockVal) game_def.Block {
	return toBlockSoldier(b)
}
