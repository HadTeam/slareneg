package Blocks

import "server/Untils/pkg/MapType"

var _ MapType.Block = (*BlockBlank)(nil)

type BlockBlank struct {
	BaseBlock
}

var blockBlankMeta = MapType.BlockMeta{
	BlockId:     0,
	Name:        "blank",
	Description: "",
}

func init() {
	MapType.RegisterBlockType(blockBlankMeta, toBlockBlank)
}

func toBlockBlank(_ uint8, _ uint8) MapType.Block {
	return MapType.Block(&BlockBlank{})
}

func (*BlockBlank) MoveTo(ownerId uint8, _ uint8) MapType.Block {
	return toBlockSoldier(0, ownerId)
}
