package Blocks

import "server/Utils/pkg/MapType"

var _ MapType.Block = (*BlockBlank)(nil)

type BlockBlank struct {
	BaseBlock
}

var BlockBlankMeta = MapType.BlockMeta{
	BlockId:           0,
	Name:              "blank",
	Description:       "",
	VisitFallBackType: 0,
}

func init() {
	MapType.RegisterBlockType(BlockBlankMeta, toBlockBlank)
}

func (*BlockBlank) GetMeta() MapType.BlockMeta {
	return BlockBlankMeta
}

func toBlockBlank(_ uint16, _ uint16) MapType.Block {
	return MapType.Block(&BlockBlank{})
}

func (*BlockBlank) MoveTo(ownerId uint16, _ uint16) MapType.Block {
	return toBlockSoldier(0, ownerId)
}
