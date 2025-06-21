package gamemap

import (
	"server/internal/game/block"
	"strconv"
)

type Size struct {
	Width  uint16
	Height uint16
}

type Pos struct {
	X uint16
	Y uint16
}

type Info struct {
	Id   string
	Name string
	Desc string
}

type Sight [][]bool // 2D slice for visibility, true if visible

type Blocks [][]block.Block // 2D slice for blocks

type Map interface {
	IsEmpty() bool
	Block(pos Pos) (block.Block, error)
	Blocks() Blocks
	SetBlock(pos Pos, b block.Block) error
	SetBlocks(blocks Blocks) error // Set all blocks at once
	Size() Size
	Info() Info

	RoundStart(roundNum uint16)
	RoundEnd(roundNum uint16)

	Fog(owner []block.Owner, sight Sight) error
}

func (b Blocks) String() string {
	// return a debug string representation of the blocks
	// TODO
	return "Blocks(...)"
}

func (s Size) String() string {
	return "Size(" + strconv.Itoa(int(s.Width)) + "x" + strconv.Itoa(int(s.Height)) + ")"
}

func (p Pos) String() string {
	return "Pos(" + strconv.Itoa(int(p.X)) + "," + strconv.Itoa(int(p.Y)) + ")"
}

func (s Size) IsPosValid(p Pos) bool {
	return p.X > 0 && p.X <= s.Width && p.Y > 0 && p.Y <= s.Height
}

func (i Info) String() string {
	return "Info(#" + i.Id + ", " + i.Name + ", Desc: " + i.Desc + ")"
}
