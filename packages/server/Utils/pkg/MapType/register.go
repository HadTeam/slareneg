package MapType

import (
	"fmt"
)

type tranFunc func(ownerId uint16, number uint16) Block

var transBlockTypeFunc map[uint8]tranFunc

func RegisterBlockType(meta BlockMeta, transFunc tranFunc) {
	if transBlockTypeFunc == nil {
		transBlockTypeFunc = make(map[uint8]tranFunc)
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
