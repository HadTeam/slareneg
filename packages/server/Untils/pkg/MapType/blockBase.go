package MapType

var _ Block = (*BaseBlock)(nil)

type BaseBlock struct {
	ownerId  uint8
	typeId   uint8
	blockId  uint8
	position BlockPosition
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

func (*BaseBlock) RoundEnd(_ uint8) (bool, GameOverSign) {
	return false, false
}

func (*BaseBlock) GetMoveStatus() MoveStatus {
	return MoveStatus{false, false}
}

func (*BaseBlock) MoveFrom(_ uint8) {}

func (*BaseBlock) MoveTo(_ uint8, _ uint8) Block {
	return nil
}
