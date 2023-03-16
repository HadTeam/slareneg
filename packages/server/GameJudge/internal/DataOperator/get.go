package DataOperator

import (
	"encoding/json"
	"server/ApiProvider/pkg/InstructionType"
	"server/MapHandler/pkg/MapType"
	_ "server/MapHandler/pkg/MapType"
)

const exampleMap string = "[\n[0,0,0,0,2],\n[0,2,0,0,0],\n[0,0,0,0,0],\n[0,3,3,0,3],\n[0,3,0,2,0]\n]"

var exampleInstruction = []InstructionType.Instruction{
	InstructionType.MoveInstruction{UserId: 1, Position: MapType.BlockPosition{X: 1, Y: 1}, Towards: InstructionType.MoveTowardsDown},
}

func GetOriginalGameMap(mapId uint32) MapType.Map {
	var originMapStr string
	if mapId == 0 {
		originMapStr = exampleMap
	} else {
		// TODO: Get from the db
	}
	var result [][]uint8
	if json.Unmarshal([]byte(originMapStr), &result) == nil {
		// TODO: Error return
	}

	ret := make([][]MapType.Block, len(result))
	for rowNum, row := range result {
		ret[rowNum] = make([]MapType.Block, len(row))
		for colNum, typeId := range row {
			ret[rowNum][colNum] = MapType.ToBlockByTypeId(typeId, MapType.BaseBlock{})

		}
	}
	return MapType.Map{Blocks: ret}
}

func GetInstruction() []InstructionType.Instruction {
	// TODO: Stop receiving new instruction
	return exampleInstruction
}
