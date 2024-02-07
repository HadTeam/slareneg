package block

var _ Block = (*King)(nil)

type King struct {
	BaseBuilding
	originalOwnerId uint16
}

var KingMeta = Meta{
	BlockId:           2,
	Name:              "king",
	Description:       "",
	VisitFallBackType: CastleMeta.BlockId,
}

func init() {
	Register(KingMeta, toBlockKing)
}

func toBlockKing(b Block) Block {
	var ret King
	ret.number = b.Number()
	ret.ownerId = b.OwnerId()
	ret.originalOwnerId = b.OwnerId()
	return Block(&ret)
}

func (*King) Meta() Meta {
	return KingMeta
}

func (block *King) MoveTo(info Val) Block {
	if block.originalOwnerId == block.ownerId {
		block.BaseBuilding.MoveTo(info)
		return nil
	} else {
		return toBlockCastle(block)
	}
}
