package block

var _ Block = (*Soldier)(nil)

type Soldier struct {
	BaseBlock
}

var SoldierMeta = BlockMeta{
	BlockId:           1,
	Name:              "soldier",
	Description:       "",
	VisitFallBackType: BlankMeta.BlockId,
}

func init() {
	Register(SoldierMeta, toBlockSoldier)
}

func toBlockSoldier(b Block) Block {
	var ret Soldier
	ret.number = b.Number()
	ret.ownerId = b.OwnerId()
	return Block(&ret)
}

func (*Soldier) Meta() BlockMeta {
	return SoldierMeta
}

func (block *Soldier) Number() uint16 {
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

func (block *Soldier) MoveFrom(number uint16) uint16 {
	var ret uint16
	if block.number <= number {
		ret = block.number - 1
		block.number = 1
	} else {
		ret = number
		block.number -= number
	}
	return ret
}

func (block *Soldier) MoveTo(info Val) Block {
	if block.ownerId != info.OwnerId {
		if block.number < info.Number {
			block.ownerId = info.OwnerId
			block.number = info.Number - block.number
		} else {
			block.number -= info.Number
		}
	} else {
		block.number += info.Number
	}
	return nil

}
