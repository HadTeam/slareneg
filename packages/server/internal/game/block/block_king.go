package block

var _ Block = (*King)(nil)

type King struct {
	BaseBuilding
	originalOwner Owner
}

var KingMeta Meta
var KingName Name

func init() {
	KingName = Register("king", "", toBlockKing)
	KingMeta = GetMetaByName[KingName]
}

func toBlockKing(b Block) Block {
	var ret King
	ret.num = b.Num()
	ret.owner = b.Owner()
	ret.originalOwner = b.Owner()
	return &ret
}

func (*King) Meta() Meta {
	return KingMeta
}

func (block *King) MoveTo(num Num, owner Owner) Block {
	block.BaseBuilding.MoveTo(num, owner)
	if block.originalOwner != block.owner {
		return toBlockCastle(block)
	}
	return nil
}

func (block *King) Fog(isOwner bool, isSight bool) Block {
	if isOwner || isSight {
		return block
	}
	return &Blank{}
}
