package BlockType

type BlockMeta struct {
	Name              string
	Description       string
	BlockId           uint8
	VisitFallBackType uint8
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
	GetNumber() uint16
	GetOwnerId() uint16

	// Round Events
	RoundStart(roundNumber uint16) bool
	RoundEnd(roundNumber uint16) (bool, GameOverSign bool)

	GetMoveStatus() MoveStatus
	MoveFrom(number uint16)
	// MoveTo Ret: a new block to replace this place
	MoveTo(ownerId uint16, number uint16) Block

	GetMeta() BlockMeta
}

type Position struct{ X, Y uint8 }
