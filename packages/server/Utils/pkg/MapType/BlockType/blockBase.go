package BlockType

var _ Block = (*BaseBlock)(nil)

type BaseBlock struct {
	ownerId uint16
	number  uint16
}

func (block *BaseBlock) GetMeta() BlockMeta {
	panic("no block meta can be provided")
}

func (block *BaseBlock) GetNumber() uint16 {
	return block.number
}

func (block *BaseBlock) GetOwnerId() uint16 {
	return block.ownerId
}

func (*BaseBlock) RoundStart(_ uint16) {
}

func (*BaseBlock) RoundEnd(_ uint16) {
}

func (*BaseBlock) GetMoveStatus() MoveStatus {
	return MoveStatus{false, false}
}

func (*BaseBlock) MoveFrom(number uint16) {}

func (*BaseBlock) MoveTo(ownerId uint16, number uint16) Block {
	return nil
}
