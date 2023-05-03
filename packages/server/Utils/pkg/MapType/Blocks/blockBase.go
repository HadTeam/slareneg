package Blocks

import (
	"server/Utils/pkg/MapType"
)

var _ MapType.Block = (*BaseBlock)(nil)

type BaseBlock struct {
	ownerId uint16
	number  uint16
}

func NewBaseBlock(number uint16, ownerId uint16) *BaseBlock {
	return &BaseBlock{
		ownerId: ownerId,
		number:  number,
	}
}

func (block *BaseBlock) GetMeta() MapType.BlockMeta {
	panic("no block meta can be provided")
}

func (block *BaseBlock) GetNumber() uint16 {
	return block.number
}

func (block *BaseBlock) GetOwnerId() uint16 {
	return block.ownerId
}

func (*BaseBlock) RoundStart(_ uint16) bool {
	return false
}

func (*BaseBlock) RoundEnd(_ uint16) (bool, MapType.GameOverSign) {
	return false, false
}

func (*BaseBlock) GetMoveStatus() MapType.MoveStatus {
	return MapType.MoveStatus{false, false}
}

func (*BaseBlock) MoveFrom(number uint16) {}

func (*BaseBlock) MoveTo(ownerId uint16, number uint16) MapType.Block {
	return nil
}
