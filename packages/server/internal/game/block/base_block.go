package block

import (
	"log/slog"
)

var _ Block = (*BaseBlock)(nil)

type BaseBlock struct {
	num   Num
	owner Owner
}

func (*BaseBlock) Meta() Meta {
	slog.Warn("Meta() called on BaseBlock, which should implement in every Block type")
	return Meta{}
}

func (block *BaseBlock) Num() Num {
	return block.num
}

func (block *BaseBlock) Owner() Owner {
	return block.owner
}

func (*BaseBlock) RoundStart(_ uint16) {}

func (*BaseBlock) RoundEnd(_ uint16) {}

func (*BaseBlock) AllowMove() AllowMove {
	return AllowMove{}
}

func (*BaseBlock) MoveFrom(_ Num) Num {
	return 0
}

func (*BaseBlock) MoveTo(num Num, owner Owner) Block {
	return nil
}

func (block *BaseBlock) Fog(isOwner bool, isSight bool) Block {
	// Base implementation: if owner or in sight, return self; otherwise return blank
	if isOwner || isSight {
		return block
	}
	// Return a blank block with no information
	return &BaseBlock{num: 0, owner: 0}
}
