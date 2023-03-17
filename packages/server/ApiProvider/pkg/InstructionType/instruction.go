package InstructionType

import "server/MapHandler/pkg/MapType"

type Instruction interface{}

type MoveTowardsType uint8

const (
	MoveTowardsLeft  MoveTowardsType = 1
	MoveTowardsRight MoveTowardsType = 2
	MoveTowardsUp    MoveTowardsType = 3
	MoveTowardsDown  MoveTowardsType = 4
)

type MoveInstruction struct {
	UserId   uint8
	Position MapType.BlockPosition
	Towards  MoveTowardsType
	Number   uint8
}
