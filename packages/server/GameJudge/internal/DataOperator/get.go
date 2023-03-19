package DataOperator

import (
	"encoding/json"
	"server/ApiProvider/pkg/InstructionType"
	"server/GameJudge/internal/DataOperator/local"
	"server/GameJudge/internal/GameType"
	"server/MapHandler/pkg/MapType"
	_ "server/MapHandler/pkg/MapType"
)

func GetOriginalGameMap(mapId uint32) MapType.Map {
	originMapStr := local.GetOriginalGameMapStr(mapId)
	var result [][]uint8
	if json.Unmarshal([]byte(originMapStr), &result) == nil {
		// TODO: Error return
	}
	ret := make([][]MapType.Block, len(result))
	for rowNum, row := range result {
		ret[rowNum] = make([]MapType.Block, len(row))
		for colNum, typeId := range row {
			ret[rowNum][colNum] = MapType.ToBlockByTypeId(typeId, &MapType.BaseBlock{})
		}
	}
	return MapType.Map{Blocks: ret, Size: MapType.MapSize{X: uint8(len(ret[0])), Y: uint8(len(ret))}, MapId: mapId}
}

var MapTemp map[GameType.GameId]*MapType.Map

func init() {
	MapTemp = make(map[GameType.GameId]*MapType.Map)
}

func GetCurrentMap(id GameType.GameId) *MapType.Map {
	return MapTemp[id]
}

func GetInstruction(id GameType.GameId) []InstructionType.Instruction {
	// TODO: Stop receiving new instruction
	return local.ExampleInstruction
}
