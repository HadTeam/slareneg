package InstructionType

type Instruction interface{}

type MoveTowardsType string

const (
	MoveTowardsLeft  MoveTowardsType = "left"
	MoveTowardsRight MoveTowardsType = "right"
	MoveTowardsUp    MoveTowardsType = "up"
	MoveTowardsDown  MoveTowardsType = "down"
)

type BlockPosition struct {
	X, Y uint8
}

type Move struct {
	UserId   uint8
	Position BlockPosition
	Towards  MoveTowardsType
	Number   uint8
}

type ForceStart struct {
	UserId uint8
	Status bool
}

type Surrender struct {
	UserId uint8
}
