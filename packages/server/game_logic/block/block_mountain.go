package block

import (
	"server/game_logic/block_manager"
	"server/game_logic/game_def"
)

type Mountain struct {
	BaseBlock
}

var MountainMeta = _type.BlockMeta{
	BlockId:           4,
	Name:              "mountain",
	Description:       "",
	VisitFallBackType: 4,
}

func init() {
	block_manager.Register(MountainMeta, toBlockMountain)
}

func toBlockMountain(_type.Block) _type.Block {
	return _type.Block(&Mountain{})
}

func (*Mountain) Meta() _type.BlockMeta {
	return MountainMeta
}

func (*Mountain) GetMoveStatus() _type.MoveStatus { // same as `BaseBlock`'s, in order to remain the function
	return _type.MoveStatus{false, false}
}
