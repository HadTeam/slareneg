package blockManager

import (
	"github.com/sirupsen/logrus"
	"server/utils/pkg/map/type"
)

type tranFunc func(_type.Block) _type.Block

var transBlockTypeFunc map[uint8]tranFunc
var GetBlockIdByName map[string]uint8
var GetMetaById map[uint8]_type.Meta

func Register(meta _type.Meta, transFunc tranFunc) {
	if transBlockTypeFunc == nil {
		transBlockTypeFunc = make(map[uint8]tranFunc)
	}
	if GetBlockIdByName == nil {
		GetBlockIdByName = make(map[string]uint8)
	}
	if GetMetaById == nil {
		GetMetaById = make(map[uint8]_type.Meta)
	}

	GetBlockIdByName[meta.Name] = meta.BlockId
	GetMetaById[meta.BlockId] = meta
	transBlockTypeFunc[meta.BlockId] = transFunc
	logrus.Println("Registered a block type", "id:", meta.BlockId, " name:", meta.Name, " description:", meta.Description)
}

func ToBlockByTypeId(typeId uint8, block _type.Block) _type.Block {
	transFunc, err := transBlockTypeFunc[typeId]
	if !err {
		logrus.Warningln("Get an unknown blockTypeId", typeId)
		transFunc = transBlockTypeFunc[0] // Note: Must ensure blocks.BlockBlankMeta.BlockId=0
	}
	return transFunc(block)
}

func NewBlock(typeId uint8, number uint16, ownerId uint16) _type.Block {
	return ToBlockByTypeId(typeId, &BaseBlock{number: number, ownerId: ownerId})
}
