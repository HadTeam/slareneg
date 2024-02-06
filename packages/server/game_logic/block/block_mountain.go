package block

import (
	"server/game_logic/block_manager"
	"server/game_logic/game_def"
)

type Mountain struct {
	BaseBlock
}

var MountainMeta = game_def.BlockMeta{
	BlockId:           4,
	Name:              "mountain",
	Description:       "",
	VisitFallBackType: 4,
}

func init() {
	block_manager.Register(MountainMeta, toBlockMountain)
}

func toBlockMountain(game_def.Block) game_def.Block {
	return game_def.Block(&Mountain{})
}

func (*Mountain) Meta() game_def.BlockMeta {
	return MountainMeta
}

func (*Mountain) GetMoveStatus() game_def.MoveStatus { // same as `BaseBlock`'s, in order to remain the function
	return game_def.MoveStatus{false, false}
}
