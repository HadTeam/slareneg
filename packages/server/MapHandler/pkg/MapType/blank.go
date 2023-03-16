package MapType

type BlockBlank struct {
	BaseBlock
}

var blockBlankMeta = BlockMeta{
	blockId:     0,
	name:        "blank",
	description: "",
}

func init() {
	RegisterBlockType(blockBlankMeta, toBlockBlank)
}

func toBlockBlank(_ uint8, _ uint8) Block {
	return Block(BlockBlank{})
}

func (BaseBlock) moveRequest(ownerId uint8, _ uint8) (bool, Block) {
	// TODO: Transform this block to BlockSoldier
	return true, toBlockSoldier(ownerId, 0)
}
