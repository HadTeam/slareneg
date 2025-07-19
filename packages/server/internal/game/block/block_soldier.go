package block

var _ Block = (*Soldier)(nil)

type Soldier struct {
	BaseBlock
}

var SoldierMeta Meta
var SoldierName Name

func init() {
	SoldierName = Register("soldier", "Basic military unit", toBlockSoldier)
	SoldierMeta = GetMetaByName[SoldierName]
}

func toBlockSoldier(b Block) Block {
	var ret Soldier
	ret.num = b.Num()
	ret.owner = b.Owner()
	return &ret
}

func (*Soldier) Meta() Meta {
	return SoldierMeta
}

func (block *Soldier) Num() Num {
	return block.num
}

func (block *Soldier) RoundStart(roundNum uint16) {
	if (roundNum%25)-1 == 0 && roundNum != 1 {
		block.num += 1
	}
}

func (block *Soldier) AllowMove() AllowMove {
	return AllowMove{
		From:   true,
		To:     true,
		Reason: "Soldier can always move",
	}
}

func (block *Soldier) MoveFrom(num Num) Num {
	var ret Num
	if block.num <= num {
		ret = block.num - 1
		block.num = 1
	} else {
		ret = num
		block.num -= num
	}
	return ret
}

func (block *Soldier) MoveTo(num Num, owner Owner) Block {
	if block.owner != owner {
		if block.num < num {
			block.owner = owner
			block.num = num - block.num
		} else {
			block.num -= num
		}
	} else {
		block.num += num
	}
	return nil
}

func (block *Soldier) Fog(isOwner bool, isSight bool) Block {
	if isOwner || isSight {
		// Owner can always see the real soldier
		return block
	}
	// Out of sight, show as blank
	return &Blank{}
}
