package block

import (
	"server/utils/pkg/map/blockManager"
	"server/utils/pkg/map/type"
)

type Mountain struct {
	BaseBlock
}

var MountainMeta = _type.Meta{
	BlockId:           4,
	Name:              "mountain",
	Description:       "",
	VisitFallBackType: 4,
}

func init() {
	blockManager.Register(MountainMeta, toBlockMountain)
}

func toBlockMountain(_type.Block) _type.Block {
	return _type.Block(&Mountain{})
}

func (*Mountain) Meta() _type.Meta {
	return MountainMeta
}

func (*Mountain) GetMoveStatus() _type.MoveStatus { // same as `BaseBlock`'s, in order to remain the function
	return _type.MoveStatus{false, false}
}
