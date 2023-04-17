package Blocks

import "server/Utils/pkg/MapType"

var _ MapType.Block = (*BaseBlock)(nil)

type BaseBlock struct {
	ownerId  uint16
	typeId   uint8
	blockId  uint16
	position MapType.BlockPosition
}

func (block *BaseBlock) GetMeta() MapType.BlockMeta {
	panic("no block meta can be provided")
}

func (*BaseBlock) GetNumber() uint16 {
	return 0
}

func (block *BaseBlock) GetOwnerId() uint16 {
	return block.ownerId
}

func (*BaseBlock) RoundStart(roundNumber uint16) bool {
	return false
}

func (*BaseBlock) RoundEnd(roundNumber uint16) (bool, MapType.GameOverSign) {
	return false, false
}

func (*BaseBlock) GetMoveStatus() MapType.MoveStatus {
	return MapType.MoveStatus{false, false}
}

func (*BaseBlock) MoveFrom(number uint16) {}

func (*BaseBlock) MoveTo(ownerId uint16, number uint16) MapType.Block {
	return nil
}
