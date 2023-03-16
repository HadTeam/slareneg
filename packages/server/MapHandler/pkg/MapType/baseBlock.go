package MapType

type BlockMeta struct {
	name        string
	description string
	blockId     uint8
}

type Block interface {
	GetNumber() uint8
	GetOwnerId() uint8

	// RoundEvent Ret: whether updated(true) or not(false)
	roundStart(roundNumber uint8) bool
	roundEnd(roundNumber uint8) bool

	// MoveRequest Ret: whether allow to move here
	moveRequest(ownerId uint8, number uint8) (bool, Block)
}

type BlockPosition struct{ X, Y uint8 }

type BaseBlock struct {
	ownerId  uint8
	typeId   uint8
	blockId  uint8
	position BlockPosition
}

func (BaseBlock) GetNumber() uint8 {
	return 0
}

func (block BaseBlock) GetOwnerId() uint8 {
	return block.ownerId
}

func (BaseBlock) roundStart(_ uint8) bool {
	return false
}

func (BaseBlock) roundEnd(_ uint8) bool {
	return false
}
