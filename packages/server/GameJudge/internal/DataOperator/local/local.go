package local

import (
	"server/ApiProvider/pkg/InstructionType"
	"server/GameJudge/pkg/GameType"
	"server/MapHandler/pkg/MapType"
)

var originalMap map[uint32]string
var GamePool map[uint16]*GameType.Game
var UserList map[uint16][]GameType.User

func init() {
	originalMap = make(map[uint32]string)
	originalMap[0] = "[\n[0,0,0,0,2],\n[0,2,0,0,0],\n[0,0,0,0,0],\n[0,3,3,0,3],\n[0,3,0,2,0]\n]"

	GamePool = make(map[uint16]*GameType.Game)

	UserList = make(map[uint16][]GameType.User)
	UserList[0] = []GameType.User{
		{"Tester1", 1, GameType.UserStatusConnected},
		{"Tester2", 2, GameType.UserStatusConnected},
	}
}

func GetOriginalGameMapStr(mapId uint32) string {
	return originalMap[mapId]
}

var ExampleInstruction = []InstructionType.Instruction{
	InstructionType.MoveInstruction{UserId: 1, Position: MapType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
	InstructionType.MoveInstruction{UserId: 2, Position: MapType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
}
