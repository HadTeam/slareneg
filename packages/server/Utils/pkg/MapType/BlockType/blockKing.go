package BlockType

var _ Block = (*BlockKing)(nil)

type BlockKing struct {
	BaseBuilding
	originalOwnerId uint16
}

var BlockKingMeta = BlockMeta{
	BlockId:           2,
	Name:              "king",
	Description:       "",
	VisitFallBackType: BlockCastleMeta.BlockId,
}

func init() {
	RegisterBlockType(BlockKingMeta, toBlockKing)
}

func toBlockKing(number uint16, ownerId uint16) Block {
	var ret BlockKing
	ret.number = number
	ret.ownerId = ownerId
	ret.originalOwnerId = ownerId
	return Block(&ret)
}

func (block *BlockKing) IsDied() bool {
	return block.originalOwnerId != block.ownerId
}

func (*BlockKing) GetMeta() BlockMeta {
	return BlockKingMeta
}

func (block *BlockKing) MoveTo(ownerId uint16, number uint16) Block {
	if !block.IsDied() {
		block.BaseBuilding.MoveTo(ownerId, number)
	}
	if block.IsDied() {
		return toBlockCastle(block.number, ownerId)
	}
	return nil
}
