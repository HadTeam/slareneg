package pkg

import (
	"encoding/json"
	"server/Utils/pkg/MapType"
	"server/Utils/pkg/MapType/Blocks"
)

// Str2GameMap TODO: Add unit test
func Str2GameMap(mapId uint32, originalMapStr string) *MapType.Map {
	var result [][]uint8
	if err := json.Unmarshal([]byte(originalMapStr), &result); err != nil {
		panic(err)
		return nil
	}
	size := MapType.MapSize{X: uint8(len(result[0])), Y: uint8(len(result))}
	ret := make([][]MapType.Block, size.Y)
	for rowNum, row := range result {
		ret[rowNum] = make([]MapType.Block, size.X)
		for colNum, typeId := range row {
			ret[rowNum][colNum] = MapType.ToBlockByTypeId(typeId, &Blocks.BaseBlock{})
		}
	}
	return &MapType.Map{
		Blocks: ret,
		Size:   size,
		MapId:  mapId,
	}
}
