package block

import (
	"log"
)

type tranFunc func(number uint16, ownerId uint16) Block

var transBlockTypeFunc map[uint8]tranFunc

func RegisterBlockType(meta BlockMeta, transFunc tranFunc) {
	if transBlockTypeFunc == nil {
		transBlockTypeFunc = make(map[uint8]tranFunc)
	}
	transBlockTypeFunc[meta.BlockId] = transFunc
	log.Println("[Info] Registered a block type", "id:", meta.BlockId, " name:", meta.Name, " description:", meta.Description)
}

func ToBlockByTypeId(typeId uint8, block Block) Block {
	transFunc, err := transBlockTypeFunc[typeId]
	if !err {
		log.Println("[Warn] Get an unknown blockTypeId", typeId)
		transFunc = transBlockTypeFunc[0] // Note: Must ensure blocks.BlockBlankMeta.BlockId=0
	}
	return transFunc(block.GetNumber(), block.GetOwnerId())
}

func NewBlock(typeId uint8, number uint16, ownerId uint16) Block {
	return ToBlockByTypeId(typeId, &BaseBlock{ownerId, number})
}
