package MapType

var _ Block = (*BlockKing)(nil)

type BlockKing struct {
	BaseBuilding
}

var blockKingMeta = BlockMeta{
	blockId:     2,
	name:        "king",
	description: "",
}

func init() {
	RegisterBlockType(blockKingMeta, toBlockKing)
}

func toBlockKing(number uint8, ownerId uint8) Block {
	var ret BlockKing
	ret.number = number
	ret.ownerId = ownerId
	return Block(&ret)
}

func (block *BlockKing) roundEnd(_ uint8) (bool, GameOverSign) {
	if block.number <= 0 {
		return true, true
	}
	return false, false
}
