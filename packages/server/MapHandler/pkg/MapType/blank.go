package MapType

var _ Block = (*BlockBlank)(nil)

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

func (BlockBlank) MoveTo(ownerId uint8, _ uint8) Block {
	return toBlockSoldier(0, ownerId)
}
