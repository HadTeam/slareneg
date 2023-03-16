package MapType

type BaseBuilding struct {
	BaseBlock
	number uint8
}

func (block BaseBuilding) getNumber() uint8 {
	return block.number
}

func (block BaseBuilding) roundStart(_ uint8) bool {
	block.number += 1
	return true
}

func (block BaseBuilding) MoveRequest(ownerId uint8, number uint8) (bool, Block) {
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
