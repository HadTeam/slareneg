package block

var _ Block = (*Blank)(nil)

type Blank struct {
	BaseBlock
}

var BlankMeta = Meta{
	BlockId:           0,
	Name:              "blank",
	Description:       "",
	VisitFallBackType: 0,
}

func init() {
	RegisterBlockType(BlankMeta, toBlockBlank)
}

func (*Blank) GetMeta() Meta {
	return BlankMeta
}

func (*Blank) GetMoveStatus() MoveStatus {
	return MoveStatus{false, true}
}

func toBlockBlank(_ uint16, _ uint16) Block {
	return Block(&Blank{})
}

func (*Blank) MoveTo(ownerId uint16, number uint16) Block {
	return toBlockSoldier(number, ownerId)
}
