package DataOperator

import (
	"encoding/json"
	"server/MapHandler/pkg/MapType"
	_ "server/MapHandler/pkg/MapType"
)

type position struct{ x, y uint8 }
type moveArg struct{ towards string }

type Instruction struct {
	action     string
	instructor string
	position   position
	moveArg    moveArg
}

const exampleMap string = "[\n[0,0,0,0,2],\n[0,2,0,0,0],\n[0,0,0,0,0],\n[0,3,3,0,3],\n[0,3,0,2,0]\n]"

var exampleInstruction = []Instruction{
	{action: "move", instructor: "CornWorld", position: position{x: 1, y: 1}, moveArg: moveArg{towards: "up"}},
}

func GetOriginalGameMap(mapId uint32) MapType.Map {
	// TODO: Get from the db
	var result [][]uint8
	if json.Unmarshal([]byte(exampleMap), &result) == nil {
		// TODO: Error return
	}

	ret := make([][]MapType.Block, len(result))
	for rowNum, row := range result {
		ret[rowNum] = make([]MapType.Block, len(row))
		for colNum, typeId := range row {
			ret[rowNum][colNum] = MapType.ToBlockByTypeId(uint8(typeId), MapType.BaseBlock{})

		}
	}
	return MapType.Map{Blocks: ret}
}

func GetInstruction() []Instruction {
	// TODO: Stop receiving new instruction
	return exampleInstruction
}
