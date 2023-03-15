package MapType

import "fmt"

var transBlockTypeFunc map[uint8]func(ownerId uint8, number uint8) Block

func RegisterBlockType(typeId uint8, transFunc func(uint8, uint8) Block) {
	if transBlockTypeFunc == nil {
		transBlockTypeFunc = make(map[uint8]func(ownerId uint8, number uint8) Block)
	}
	transBlockTypeFunc[typeId] = transFunc
	fmt.Println("[Info] Registered a block type ", typeId)
}

func ToBlockByTypeId(typeId uint8, block Block) Block {
	transFunc, err := transBlockTypeFunc[typeId]
	if !err {
		fmt.Println("[Warn] Get an unknown blockTypeId", typeId)
		transFunc = transBlockTypeFunc[0] // Use the blank block
	}
	return transFunc(block.GetOwnerId(), block.GetNumber())
}
