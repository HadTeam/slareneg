package block

type Num uint16
type Owner uint16
type Name string

type Meta struct {
	Name        Name
	Description string
}

type AllowMove struct {
	From   bool
	To     bool
	Reason string // debug info
}

type Block interface {
	Num() Num
	Owner() Owner

	// Round Events
	RoundStart(roundNum uint16)
	RoundEnd(roundNum uint16)

	// Move related
	AllowMove() AllowMove
	MoveFrom(num Num) Num
	// MoveTo Ret: a new block to replace this place
	MoveTo(num Num, owner Owner) Block

	Fog(isOwner bool, isSight bool) Block

	Meta() Meta
}

type Position struct{ X, Y uint8 }
