package block

import (
	"github.com/sirupsen/logrus"
	"server/utils/pkg/map/type"
)

var _ _type.Block = (*BaseBlock)(nil)

type BaseBlock struct {
	ownerId uint16
	number  uint16
}

func (*BaseBlock) Meta() _type.Meta {
	logrus.Panic("no block meta can be provided")
	return _type.Meta{}
}

func (block *BaseBlock) Number() uint16 {
	return block.number
}

func (block *BaseBlock) OwnerId() uint16 {
	return block.ownerId
}

func (*BaseBlock) RoundStart(_ uint16) {}

func (*BaseBlock) RoundEnd(_ uint16) {}

func (*BaseBlock) GetMoveStatus() _type.MoveStatus {
	return _type.MoveStatus{}
}

func (*BaseBlock) MoveFrom(_ uint16) {}

func (*BaseBlock) MoveTo(_type.BlockVal) _type.Block {
	return nil
}
