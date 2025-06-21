package block

var _ Block = (*Blank)(nil)

type Blank struct {
	BaseBlock
}

var BlankMeta Meta
var BlankName Name

func init() {
	BlankName = Register("blank", "", toBlockBlank)
	BlankMeta = GetMetaByName[BlankName]
}

func (*Blank) Meta() Meta {
	return BlankMeta
}

func (b *Blank) AllowMove() AllowMove {
	return AllowMove{
		From:   false,
		To:     true,
		Reason: "Blank block can only be moved to",
	}
}

func (b *Blank) Fog(isOwner bool, isSight bool) Block {
	// Blank blocks are always visible as blank
	return b
}

func toBlockBlank(Block) Block {
	return &Blank{}
}

func (b *Blank) MoveTo(num Num, owner Owner) Block {
	return toBlockSoldier(&BaseBlock{num: num, owner: owner})
}
