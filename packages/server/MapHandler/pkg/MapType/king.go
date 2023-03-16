package MapType

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
	return Block(ret)
}

func (block BlockKing) roundEnd(_ uint8) bool {
	if block.number <= 0 {
		// TODO: Handle game-over
		return true
	}
	return false
}
