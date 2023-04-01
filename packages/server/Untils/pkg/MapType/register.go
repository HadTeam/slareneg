package MapType

import (
	"fmt"
)

var transBlockTypeFunc map[uint8]func(ownerId uint8, number uint8) Block

func RegisterBlockType(meta BlockMeta, transFunc func(uint8, uint8) Block) {
	if transBlockTypeFunc == nil {
		transBlockTypeFunc = make(map[uint8]func(ownerId uint8, number uint8) Block)
	}
	transBlockTypeFunc[meta.BlockId] = transFunc
	fmt.Println("[Info] Registered a block type", "id:", meta.BlockId, " name:", meta.Name, " description:", meta.Description)
}

func ToBlockByTypeId(typeId uint8, block Block) Block {
	transFunc, err := transBlockTypeFunc[typeId]
	if !err {
		fmt.Println("[Warn] Get an unknown blockTypeId", typeId)
		transFunc = transBlockTypeFunc[0] // Use the blank block
	}
	return transFunc(block.GetOwnerId(), block.GetNumber())
}
