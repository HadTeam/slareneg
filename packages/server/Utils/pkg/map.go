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
			ret[rowNum][colNum] = MapType.ToBlockByTypeId(typeId, Blocks.NewBaseBlock(0, 0))
		}
	}
	return &MapType.Map{
		Blocks: ret,
		Size:   size,
		MapId:  mapId,
	}
}

func FullStr2GameMap(mapId uint32, originalMapStr string) *MapType.Map {
	var result [][][]uint16
	if err := json.Unmarshal([]byte(originalMapStr), &result); err != nil {
		panic(err)
		return nil
	}
	size := MapType.MapSize{X: uint8(len(result[0])), Y: uint8(len(result))}
	ret := make([][]MapType.Block, size.Y)
	for rowNum, row := range result {
		ret[rowNum] = make([]MapType.Block, size.X)
		for colNum, blockInfo := range row {
			blockId := blockInfo[0]
			ownerId := blockInfo[1]
			number := blockInfo[2]

			block := MapType.ToBlockByTypeId(uint8(blockId), Blocks.NewBaseBlock(number, ownerId))

			ret[rowNum][colNum] = block
		}
	}
	return &MapType.Map{
		Blocks: ret,
		Size:   size,
		MapId:  mapId,
	}
}
