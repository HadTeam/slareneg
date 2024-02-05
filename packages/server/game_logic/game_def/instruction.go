package _type

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
	Position BlockPosition
	Towards  MoveTowardsType
	Number   uint16
}

type ForceStart struct {
	UserId uint16
	Status bool
}

type Surrender struct {
	UserId uint16
}
