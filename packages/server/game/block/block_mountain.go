package block

type Mountain struct {
	BaseBlock
}

var MountainMeta = Meta{
	BlockId:           4,
	Name:              "mountain",
	Description:       "",
	VisitFallBackType: 4,
}

func init() {
	Register(MountainMeta, toBlockMountain)
}

func toBlockMountain(Block) Block {
	return Block(&Mountain{})
}

func (*Mountain) Meta() Meta {
	return MountainMeta
}

func (*Mountain) GetMoveStatus() MoveStatus { // same as `BaseBlock`'s, in order to remain the function
	return MoveStatus{false, false}
}
