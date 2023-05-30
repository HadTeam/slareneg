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
	RegisterBlockType(KingMeta, toBlockKing)
}

func toBlockKing(number uint16, ownerId uint16) Block {
	var ret King
	ret.number = number
	ret.ownerId = ownerId
	ret.originalOwnerId = ownerId
	return Block(&ret)
}

func (block *King) IsDied() bool {
	return block.originalOwnerId != block.ownerId
}

func (*King) Meta() Meta {
	return KingMeta
}

func (block *King) MoveTo(ownerId uint16, number uint16) Block {
	if !block.IsDied() {
		block.BaseBuilding.MoveTo(ownerId, number)
	}
	if block.IsDied() {
		return toBlockCastle(block.number, ownerId)
	}
	return nil
}
