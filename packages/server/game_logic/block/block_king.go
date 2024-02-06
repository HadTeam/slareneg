package block

import (
	"server/game_logic/block_manager"
	"server/game_logic/game_def"
)

var _ game_def.Block = (*King)(nil)

type King struct {
	BaseBuilding
	originalOwnerId uint16
}

var KingMeta = game_def.BlockMeta{
	BlockId:           2,
	Name:              "king",
	Description:       "",
	VisitFallBackType: CastleMeta.BlockId,
}

func init() {
	block_manager.Register(KingMeta, toBlockKing)
}

func toBlockKing(b game_def.Block) game_def.Block {
	var ret King
	ret.number = b.Number()
	ret.ownerId = b.OwnerId()
	ret.originalOwnerId = b.OwnerId()
	return game_def.Block(&ret)
}

func (block *King) IsDied() bool {
	return block.originalOwnerId != block.ownerId
}

func (*King) Meta() game_def.BlockMeta {
	return KingMeta
}

func (block *King) MoveTo(info game_def.BlockVal) game_def.Block {
	if !block.IsDied() {
		block.BaseBuilding.MoveTo(info)
	}
	if block.IsDied() {
		return toBlockCastle(block)
	}
	return nil
}
