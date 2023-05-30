package block

var _ Block = (*BaseBlock)(nil)

type BaseBlock struct {
	ownerId uint16
	number  uint16
}

func (block *BaseBlock) Meta() Meta {
	panic("no block meta can be provided")
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
	return MoveStatus{false, false}
}

func (*BaseBlock) MoveFrom(_ uint16) {}

func (*BaseBlock) MoveTo(_ uint16, _ uint16) Block {
	return nil
}
