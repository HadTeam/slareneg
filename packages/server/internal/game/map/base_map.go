package gamemap

import (
	"errors"
	"log/slog"
	"server/internal/game/block"
	"slices"
)

var _ Map = (*BaseMap)(nil)

type BaseMap struct {
	blocks Blocks
	size   Size
	info   Info
}

func (m *BaseMap) IsEmpty() bool {
	return len(m.blocks) == 0 || len(m.blocks[0]) == 0
}
func (m *BaseMap) Block(pos Pos) (block.Block, error) {
	if !m.size.IsPosValid(pos) {
		return nil, errors.New("invalid position: " + pos.String())
	}
	return m.blocks[pos.Y-1][pos.X-1], nil
}
func (m *BaseMap) Blocks() Blocks {
	if m.IsEmpty() {
		return nil
	}
	blocks := make(Blocks, len(m.blocks))
	for i := range m.blocks {
		blocks[i] = make([]block.Block, len(m.blocks[i]))
		copy(blocks[i], m.blocks[i])
	}
	return blocks
}
func (m *BaseMap) SetBlock(pos Pos, b block.Block) error {
	if !m.size.IsPosValid(pos) {
		return errors.New("invalid position: " + pos.String())
	}
	if m.blocks[pos.Y-1][pos.X-1] != nil {
		return errors.New("block already exists at position: " + pos.String())
	}
	m.blocks[pos.Y-1][pos.X-1] = b
	return nil
}
func (m *BaseMap) SetBlocks(blocks Blocks) error {
	if len(blocks) != int(m.size.Height) || len(blocks[0]) != int(m.size.Width) {
		return errors.New("blocks dimensions do not match map size: " + m.size.String())
	}
	for y, row := range blocks {
		for x, b := range row {
			if b != nil {
				if err := m.SetBlock(Pos{X: uint16(x + 1), Y: uint16(y + 1)}, b); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (m *BaseMap) Size() Size {
	return m.size
}
func (m *BaseMap) Info() Info {
	return m.info
}
func (m *BaseMap) RoundStart(roundNum uint16) {
	for _, col := range m.blocks {
		for _, b := range col {
			if b != nil {
				b.RoundStart(roundNum)
			}
		}
	}
}
func (m *BaseMap) RoundEnd(roundNum uint16) {
	for _, col := range m.blocks {
		for _, b := range col {
			if b != nil {
				b.RoundEnd(roundNum)
			}
		}
	}
}
func (m *BaseMap) Fog(owner []block.Owner, sight Sight) error {
	if m.IsEmpty() {
		return errors.New("map is empty")
	}
	if len(owner) == 0 {
		return errors.New("owner list is empty")
	}
	if len(sight) != len(m.blocks) || len(sight[0]) != len(m.blocks[0]) {
		return errors.New("sight dimensions do not match map dimensions")
	}

	for y, row := range m.blocks {
		for x, b := range row {
			if b != nil {
				isOwner := slices.Contains(owner, b.Owner())
				isSight := sight[y][x]
				m.blocks[y][x] = b.Fog(isOwner, isSight)
			}
		}
	}
	return nil
}

func NewBaseMap(blocks Blocks, size Size, info Info) *BaseMap {
	if len(blocks) == 0 || len(blocks[0]) == 0 {
		return &BaseMap{
			blocks: nil,
			size:   size,
			info:   info,
		}
	}
	if len(blocks) != int(size.Height) || len(blocks[0]) != int(size.Width) {
		slog.Error("blocks size does not match specified size",
			"expected", size.String(),
			"actual", Size{Width: uint16(len(blocks[0])), Height: uint16(len(blocks))}.String(),
			"blocks", blocks.String(),
		)
		return nil
	}
	return &BaseMap{
		blocks: blocks,
		size:   size,
		info:   info,
	}
}

func NewEmptyBaseMap(size Size, info Info) *BaseMap {
	if size.Width == 0 || size.Height == 0 {
		return nil
	}
	blocks := make(Blocks, size.Height)
	for i := range blocks {
		blocks[i] = make([]block.Block, size.Width)
	}
	return &BaseMap{
		blocks: blocks,
		size:   size,
		info:   info,
	}
}
