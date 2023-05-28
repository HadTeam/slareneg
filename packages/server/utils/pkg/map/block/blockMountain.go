package block

type BlockMountain struct {
	BaseBlock
}

var BlockMountainMeta = BlockMeta{
	BlockId:           4,
	Name:              "mountain",
	Description:       "",
	VisitFallBackType: 4,
}

func init() {
	RegisterBlockType(BlockMountainMeta, toBlockMountain)
}

func toBlockMountain(number uint16, ownerId uint16) Block {
	return Block(&BlockMountain{})
}

func (*BlockMountain) GetMeta() BlockMeta {
	return BlockMountainMeta
}

func (*BlockMountain) GetMoveStatus() MoveStatus { // same as `BaseBlock`'s, in order to remain the function
	return MoveStatus{false, false}
}
