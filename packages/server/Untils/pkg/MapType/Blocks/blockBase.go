package Blocks

import "server/Untils/pkg/MapType"

var _ MapType.Block = (*BaseBlock)(nil)

type BaseBlock struct {
	ownerId  uint8
	typeId   uint8
	blockId  uint8
	position MapType.BlockPosition
}

func (*BaseBlock) GetNumber() uint8 {
	return 0
}

func (block *BaseBlock) GetOwnerId() uint8 {
	return block.ownerId
}

func (*BaseBlock) RoundStart(_ uint8) bool {
	return false
}

func (*BaseBlock) RoundEnd(_ uint8) (bool, MapType.GameOverSign) {
	return false, false
}

func (*BaseBlock) GetMoveStatus() MapType.MoveStatus {
	return MapType.MoveStatus{false, false}
}

func (*BaseBlock) MoveFrom(_ uint8) {}

func (*BaseBlock) MoveTo(_ uint8, _ uint8) MapType.Block {
	return nil
}
