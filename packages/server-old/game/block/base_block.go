package block

import (
	"github.com/sirupsen/logrus"
)

var _ Block = (*BaseBlock)(nil)

type BaseBlock struct {
	ownerId uint16
	number  uint16
}

func (*BaseBlock) Meta() Meta {
	logrus.Panic("no block meta can be provided")
	return Meta{}
}

func (block *BaseBlock) Number() uint16 {
	return block.number
}

func (block *BaseBlock) OwnerId() uint16 {
	return block.ownerId
}

func (*BaseBlock) RoundStart(_ uint16) {}

func (*BaseBlock) RoundEnd(_ uint16) {}

func (*BaseBlock) GetMoveStatus() MoveStatus {
	return MoveStatus{}
}

func (*BaseBlock) MoveFrom(_ uint16) uint16 {
	return 0
}

func (*BaseBlock) MoveTo(Val) Block {
	return nil
}
