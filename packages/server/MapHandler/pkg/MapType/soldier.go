package MapType

type BlockSoldier struct {
	BaseBlock
	number uint8
}

var blockSoldierMeta = BlockMeta{
	blockId:     1,
	name:        "king",
	description: "",
}

func init() {
	RegisterBlockType(blockSoldierMeta, toBlockSoldier)
}

func toBlockSoldier(number uint8, ownerId uint8) Block {
	var ret BlockSoldier
	ret.number = number
	ret.ownerId = ownerId
	return Block(ret)
}

func (block BlockSoldier) getNumber() uint8 {
	return block.number
}

func (block BlockSoldier) roundStart(roundNum uint8) bool {
	if (roundNum%25)-1 == 0 && roundNum != 1 {
		block.number += 1
		return true
	}
	return false
}

func (block BlockSoldier) moveRequest(ownerId uint8, number uint8) (bool, Block) {
	if block.ownerId != ownerId {
		if block.number < number {
			block.ownerId = ownerId
			block.number = number - block.number
		} else {
			block.number -= number
		}
	} else {
		block.number += number
	}
	return true, nil
}
