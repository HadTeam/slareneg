package BlockType

var _ Block = (*BlockKing)(nil)

type BlockKing struct {
	BaseBuilding
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
	return Block(&ret)
}

func (block *BlockKing) RoundEnd(_ uint16) (bool, GameOverSign bool) {
	if block.number <= 0 {
		return true, true
	}
	return false, false
}

func (*BlockKing) GetMeta() BlockMeta {
	return BlockKingMeta
}
