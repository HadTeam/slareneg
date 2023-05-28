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
	RegisterBlockType(MountainMeta, toBlockMountain)
}

func toBlockMountain(number uint16, ownerId uint16) Block {
	return Block(&Mountain{})
}

func (*Mountain) GetMeta() Meta {
	return MountainMeta
}

func (*Mountain) GetMoveStatus() MoveStatus { // same as `BaseBlock`'s, in order to remain the function
	return MoveStatus{false, false}
}
