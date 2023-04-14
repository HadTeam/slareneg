package MapType

type BlockMeta struct {
	Name        string
	Description string
	BlockId     uint8
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
	RoundStart(roundNumber uint8) bool
	RoundEnd(roundNumber uint8) (bool, GameOverSign)

	GetMoveStatus() MoveStatus
	MoveFrom(number uint8)
	// MoveTo Ret: a new block to replace this place
	MoveTo(ownerId uint8, number uint8) Block
}

type BlockPosition struct{ X, Y uint8 }
