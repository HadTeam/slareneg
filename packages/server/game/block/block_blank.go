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
	Register(BlankMeta, toBlockBlank)
}

func (*Blank) Meta() Meta {
	return BlankMeta
}

func (*Blank) GetMoveStatus() MoveStatus {
	return MoveStatus{false, true}
}

func toBlockBlank(Block) Block {
	return Block(&Blank{})
}

func (b *Blank) MoveTo(Val) Block {
	return toBlockSoldier(b)
}
