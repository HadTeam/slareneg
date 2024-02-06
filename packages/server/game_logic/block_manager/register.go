package block_manager

import (
	"github.com/sirupsen/logrus"
	"server/game_logic/game_def"
)

type tranFunc func(game_def.Block) game_def.Block

var transBlockTypeFunc map[uint8]tranFunc
var GetBlockIdByName map[string]uint8
var GetMetaById map[uint8]game_def.BlockMeta

func Register(meta game_def.BlockMeta, transFunc tranFunc) {
	if transBlockTypeFunc == nil {
		transBlockTypeFunc = make(map[uint8]tranFunc)
	}
	if GetBlockIdByName == nil {
		GetBlockIdByName = make(map[string]uint8)
	}
	if GetMetaById == nil {
		GetMetaById = make(map[uint8]game_def.BlockMeta)
	}

	GetBlockIdByName[meta.Name] = meta.BlockId
	GetMetaById[meta.BlockId] = meta
	transBlockTypeFunc[meta.BlockId] = transFunc
	logrus.Println("Registered a block game_def", "id:", meta.BlockId, " name:", meta.Name, " description:", meta.Description)
}

func ToBlockByTypeId(typeId uint8, block game_def.Block) game_def.Block {
	transFunc, err := transBlockTypeFunc[typeId]
	if !err {
		logrus.Warningln("Get an unknown blockTypeId", typeId)
		transFunc = transBlockTypeFunc[0] // Note: Must ensure blocks.BlockBlankMeta.BlockId=0
	}
	return transFunc(block)
}

func NewBlock(typeId uint8, number uint16, ownerId uint16) game_def.Block {
	return ToBlockByTypeId(typeId, &BaseBlock{number: number, ownerId: ownerId})
}
