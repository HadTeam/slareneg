package block

type BlockSoldier struct {
	BaseBlock
}

var BlockSoldierMeta = BlockMeta{
	BlockId:           1,
	Name:              "soldier",
	Description:       "",
	VisitFallBackType: BlockBlankMeta.BlockId,
}

func init() {
	RegisterBlockType(BlockSoldierMeta, toBlockSoldier)
}

func toBlockSoldier(number uint16, ownerId uint16) Block {
	var ret BlockSoldier
	ret.number = number
	ret.ownerId = ownerId
	return Block(&ret)
}

func (*BlockSoldier) GetMeta() BlockMeta {
	return BlockSoldierMeta
}

func (block *BlockSoldier) GetNumber() uint16 {
	return block.number
}

func (block *BlockSoldier) RoundStart(roundNum uint16) {
	if (roundNum%25)-1 == 0 && roundNum != 1 {
		block.number += 1
	}
}

func (*BlockSoldier) GetMoveStatus() MoveStatus {
	return MoveStatus{true, true}
}

func (block *BlockSoldier) MoveFrom(number uint16) {
	block.number -= number
}

func (block *BlockSoldier) MoveTo(ownerId uint16, number uint16) Block {

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
	return nil

}
