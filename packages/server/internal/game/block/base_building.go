package block

var _ Block = (*BaseBuilding)(nil)

type BaseBuilding struct {
	BaseBlock
}

func (block *BaseBuilding) Num() Num {
	return block.num
}

func (block *BaseBuilding) RoundStart(_ uint16) {
	if block.Owner() != 0 {
		block.num += 1
	}
}

func (block *BaseBuilding) AllowMove() AllowMove {
	return AllowMove{
		From:   true,
		To:     true,
		Reason: "BaseBuilding can always move",
	}
}

func (block *BaseBuilding) MoveFrom(num Num) Num {
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

func (block *BaseBuilding) MoveTo(num Num, owner Owner) Block {
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

func (block *BaseBuilding) Fog(isOwner bool, isSight bool) Block {
	if isOwner {
		// Owner can always see the real building
		return block
	}
	if isSight {
		// In sight, can see the building type but not exact numbers
		// This should be overridden by specific building types
		return block
	}
	// Out of sight, show as blank
	return &Blank{}
}
