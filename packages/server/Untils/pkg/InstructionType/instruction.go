package InstructionType

import "server/Untils/pkg/MapType"

type Instruction interface{}

type MoveTowardsType uint8

const (
	MoveTowardsLeft MoveTowardsType = iota + 1
	MoveTowardsRight
	MoveTowardsUp
	MoveTowardsDown
)

type MoveInstruction struct {
	UserId   uint8
	Position MapType.BlockPosition
	Towards  MoveTowardsType
	Number   uint8
}
