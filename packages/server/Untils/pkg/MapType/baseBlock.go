package MapType

type BlockMeta struct {
	name        string
	description string
	blockId     uint8
}

type MoveRequestType uint8

const (
	MoveRequestTypeFrom MoveRequestType = 1
	MoveRequestTypeTo   MoveRequestType = 2
)

type MoveStatus struct {
	AllowMoveFrom bool
	AllowMoveTo   bool
}

type Block interface {
	GetNumber() uint8
	GetOwnerId() uint8

	// Round Events
	roundStart(roundNumber uint8) bool
	roundEnd(roundNumber uint8) (bool, GameOverSign)

	GetMoveStatus() MoveStatus
	MoveFrom(number uint8)
	// MoveTo Ret: a new block to replace this place
	MoveTo(ownerId uint8, number uint8) Block
}

type BlockPosition struct{ X, Y uint8 }

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

func (*BaseBlock) roundStart(_ uint8) bool {
	return false
}

func (*BaseBlock) roundEnd(_ uint8) (bool, GameOverSign) {
	return false, false
}

func (*BaseBlock) GetMoveStatus() MoveStatus {
	return MoveStatus{false, false}
}

func (*BaseBlock) MoveFrom(_ uint8) {}

func (*BaseBlock) MoveTo(_ uint8, _ uint8) Block {
	return nil
}
