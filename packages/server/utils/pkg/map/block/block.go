package block

type Meta struct {
	Name              string
	Description       string
	BlockId           uint8
	VisitFallBackType uint8
}

type MoveRequestType uint8

type MoveStatus struct {
	AllowMoveFrom bool
	AllowMoveTo   bool
}

type Block interface {
	Number() uint16
	OwnerId() uint16

	// Round Events
	RoundStart(roundNumber uint16)
	RoundEnd(roundNumber uint16)

	GetMoveStatus() MoveStatus
	MoveFrom(number uint16)
	// MoveTo Ret: a new block to replace this place
	MoveTo(ownerId uint16, number uint16) Block

	Meta() Meta
}

type Position struct{ X, Y uint8 }
