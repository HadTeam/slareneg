package _map

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"server/game/block"
	"strings"
)

// Str2GameMap TODO: Add unit test
func Str2GameMap(mapId uint32, originalMapStr string) *Map {
	var result [][]uint8
	if err := json.Unmarshal([]byte(originalMapStr), &result); err != nil {
		logrus.Panic(err)
		return nil
	}
	size := MapSize{W: uint8(len(result[0])), H: uint8(len(result))}
	ret := make([][]block.Block, size.H)
	for rowNum, row := range result {
		ret[rowNum] = make([]block.Block, size.W)
		for colNum, typeId := range row {
			ret[rowNum][colNum] = block.NewBlock(typeId, 0, 0)
		}
	}
	return &Map{
		ret,
		mapInfo{size, mapId},
	}
}

func FullStr2GameMap(mapId uint32, originalMapStr string) *Map {
	var result [][][]uint16
	if err := json.Unmarshal([]byte(originalMapStr), &result); err != nil {
		logrus.Panic(err)
		return nil
	}
	size := MapSize{W: uint8(len(result[0])), H: uint8(len(result))}
	ret := make([][]block.Block, size.H)
	for rowNum, row := range result {
		ret[rowNum] = make([]block.Block, size.W)
		for colNum, blockInfo := range row {
			blockId := blockInfo[0]
			ownerId := blockInfo[1]
			number := blockInfo[2]

			newBlock := block.NewBlock(uint8(blockId), number, ownerId)

			ret[rowNum][colNum] = newBlock
		}
	}
	return &Map{
		ret,
		mapInfo{size, mapId},
	}
}

// To marshal uint8 array to json as number array
type uint8Array []uint8

func (a uint8Array) MarshalJSON() ([]byte, error) {
	var res string
	if a == nil {
		res = "null"
	} else {
		res = strings.Join(strings.Fields(fmt.Sprintf("%d", a)), ",")
	}
	return []byte(res), nil
}

type mapJsonStruct struct {
	Mappings struct {
		Block []string `json:"block"`
		Owner []uint16 `json:"owner,omitempty"`
	} `json:"mappings"`
	Type   []uint8Array `json:"game_def"`
	Owner  []uint8Array `json:"owner,omitempty"`
	Number [][]uint16   `json:"number,omitempty"`
}

func JsonStrToMap(jsonStr string) *Map {
	var res mapJsonStruct
	if err := json.Unmarshal([]byte(jsonStr), &res); err != nil {
		logrus.Panic(err)
		return nil
	}

	// process original mapping
	blockMapping := make(map[uint8]uint8)
	for i, v := range res.Mappings.Block {
		blockMapping[uint8(i)] = block.GetBlockIdByName[v]
	}

	if res.Mappings.Owner == nil || res.Owner == nil {
		res.Mappings.Owner = nil
		res.Owner = nil
	} else {
		// add blank at 0 owner id
		res.Mappings.Owner = append([]uint16{0}, res.Mappings.Owner...)
	}

	var blocks [][]block.Block
	if (res.Number != nil && len(res.Type) != len(res.Number)) || (res.Owner != nil && len(res.Type) != len(res.Owner)) {
		logrus.Panic("original block game_def, number, owner id must have the same size")
	}
	blocks = make([][]block.Block, len(res.Type))

	for i, v := range res.Type {
		if (res.Number != nil && len(v) != len(res.Number[i])) || (res.Owner != nil && len(v) != len(res.Owner[i])) {
			logrus.Panic("original block game_def, number, owner id must have the same size")
		}
		blocks[i] = make([]block.Block, len(v))

		for j, typeId := range v {
			n := uint16(0)
			o := uint16(0)
			if res.Owner != nil {
				if res.Owner[i][j] >= uint8(len(res.Mappings.Owner)) {
					logrus.Panic("original owner id must be less than owner id mapping size")
				}

				o = res.Mappings.Owner[res.Owner[i][j]]
			}
			if res.Number != nil {
				n = res.Number[i][j]
			}

			if typeId >= uint8(len(blockMapping)) {
				logrus.Panic("original block game_def must be less than block game_def mapping size")
			}

			blocks[i][j] = block.NewBlock(blockMapping[typeId], n, o)
		}
	}
	return &Map{
		blocks,
		mapInfo{MapSize{uint8(len(blocks[0])), uint8(len(blocks))}, 0},
	}
}

func MapToJsonStr(m *Map) string {
	var ret = mapJsonStruct{}
	typeMapping := make(map[uint8]uint8)
	ownerMapping := make(map[uint16]uint8)
	ownerMapping[0] = 0 // blank owner id

	ret.Type = make([]uint8Array, m.size.H)
	ret.Owner = make([]uint8Array, m.size.H)
	ret.Number = make([][]uint16, m.size.H)

	for i, row := range m.Blocks {
		ret.Type[i] = make(uint8Array, m.size.W)
		ret.Owner[i] = make(uint8Array, m.size.W)
		ret.Number[i] = make([]uint16, m.size.W)

		for j, b := range row {
			if _, ok := typeMapping[b.Meta().BlockId]; !ok {
				typeMapping[b.Meta().BlockId] = uint8(len(typeMapping))
			}
			if _, ok := ownerMapping[b.OwnerId()]; !ok {
				ownerMapping[b.OwnerId()] = uint8(len(ownerMapping))
			}

			ret.Type[i][j] = typeMapping[b.Meta().BlockId]
			ret.Owner[i][j] = ownerMapping[b.OwnerId()]
			ret.Number[i][j] = b.Number()
		}
	}

	for k := range typeMapping {
		ret.Mappings.Block = append(ret.Mappings.Block, block.GetMetaById[k].Name)
	}
	for k := range ownerMapping {
		ret.Mappings.Owner = append(ret.Mappings.Owner, k)
	}

	ret.Mappings.Owner = ret.Mappings.Owner[1:] // remove blank owner id

	logrus.Info(ret)

	if retJson, err := json.Marshal(ret); err != nil {
		logrus.Panic(err)
		return ""
	} else {
		return string(retJson)
	}
}
