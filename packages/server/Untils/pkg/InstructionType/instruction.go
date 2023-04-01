package InstructionType

type Instruction interface{}

type MoveTowardsType uint8

const (
	MoveTowardsLeft MoveTowardsType = iota + 1
	MoveTowardsRight
	MoveTowardsUp
	MoveTowardsDown
)

type BlockPosition struct {
	X, Y uint8
}

type MoveInstruction struct {
	UserId   uint8
	Position BlockPosition
	Towards  MoveTowardsType
	Number   uint8
}
