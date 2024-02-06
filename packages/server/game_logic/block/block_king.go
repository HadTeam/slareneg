package block

import (
	"server/game_logic/block_manager"
	"server/game_logic/game_def"
)

var _ _type.Block = (*King)(nil)

type King struct {
	BaseBuilding
	originalOwnerId uint16
}

var KingMeta = _type.BlockMeta{
	BlockId:           2,
	Name:              "king",
	Description:       "",
	VisitFallBackType: CastleMeta.BlockId,
}

func init() {
	block_manager.Register(KingMeta, toBlockKing)
}

func toBlockKing(b _type.Block) _type.Block {
	var ret King
	ret.number = b.Number()
	ret.ownerId = b.OwnerId()
	ret.originalOwnerId = b.OwnerId()
	return _type.Block(&ret)
}

func (block *King) IsDied() bool {
	return block.originalOwnerId != block.ownerId
}

func (*King) Meta() _type.BlockMeta {
	return KingMeta
}

func (block *King) MoveTo(info _type.BlockVal) _type.Block {
	if !block.IsDied() {
		block.BaseBuilding.MoveTo(info)
	}
	if block.IsDied() {
		return toBlockCastle(block)
	}
	return nil
}
