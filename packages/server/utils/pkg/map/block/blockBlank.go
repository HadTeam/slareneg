package block

var _ Block = (*BlockBlank)(nil)

type BlockBlank struct {
	BaseBlock
}

var BlockBlankMeta = BlockMeta{
	BlockId:           0,
	Name:              "blank",
	Description:       "",
	VisitFallBackType: 0,
}

func init() {
	RegisterBlockType(BlockBlankMeta, toBlockBlank)
}

func (*BlockBlank) GetMeta() BlockMeta {
	return BlockBlankMeta
}

func (*BlockBlank) GetMoveStatus() MoveStatus {
	return MoveStatus{false, true}
}

func toBlockBlank(_ uint16, _ uint16) Block {
	return Block(&BlockBlank{})
}

func (*BlockBlank) MoveTo(ownerId uint16, number uint16) Block {
	return toBlockSoldier(number, ownerId)
}
