package block

import (
	"github.com/sirupsen/logrus"
	"server/game_logic/game_def"
)

var _ game_def.Block = (*BaseBlock)(nil)

type BaseBlock struct {
	ownerId uint16
	number  uint16
}

func (*BaseBlock) Meta() game_def.BlockMeta {
	logrus.Panic("no block meta can be provided")
	return game_def.BlockMeta{}
}

func (block *BaseBlock) Number() uint16 {
	return block.number
}

func (block *BaseBlock) OwnerId() uint16 {
	return block.ownerId
}

func (*BaseBlock) RoundStart(_ uint16) {}

func (*BaseBlock) RoundEnd(_ uint16) {}

func (*BaseBlock) GetMoveStatus() game_def.MoveStatus {
	return game_def.MoveStatus{}
}

func (*BaseBlock) MoveFrom(_ uint16) uint16 {
	return 0
}

func (*BaseBlock) MoveTo(game_def.BlockVal) game_def.Block {
	return nil
}
