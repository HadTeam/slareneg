package block

type Soldier struct {
	BaseBlock
}

var SoldierMeta = Meta{
	BlockId:           1,
	Name:              "soldier",
	Description:       "",
	VisitFallBackType: BlankMeta.BlockId,
}

func init() {
	RegisterBlockType(SoldierMeta, toBlockSoldier)
}

func toBlockSoldier(number uint16, ownerId uint16) Block {
	var ret Soldier
	ret.number = number
	ret.ownerId = ownerId
	return Block(&ret)
}

func (*Soldier) GetMeta() Meta {
	return SoldierMeta
}

func (block *Soldier) GetNumber() uint16 {
	return block.number
}

func (block *Soldier) RoundStart(roundNum uint16) {
	if (roundNum%25)-1 == 0 && roundNum != 1 {
		block.number += 1
	}
}

func (*Soldier) GetMoveStatus() MoveStatus {
	return MoveStatus{true, true}
}

func (block *Soldier) MoveFrom(number uint16) {
	block.number -= number
}

func (block *Soldier) MoveTo(ownerId uint16, number uint16) Block {

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
