package game_def

type Instruction interface{}

type MoveTowardsType string

const (
	MoveTowardsLeft  MoveTowardsType = "left"
	MoveTowardsRight MoveTowardsType = "right"
	MoveTowardsUp    MoveTowardsType = "up"
	MoveTowardsDown  MoveTowardsType = "down"
)

type Move struct {
	Position Position
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
