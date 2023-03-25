package MapType

import "fmt"

type MapSize struct{ X, Y uint8 }

type Map struct {
	Blocks [][]Block
	Size   MapSize
	MapId  uint32
}

func (m *Map) GetBlock(position BlockPosition) Block {
	return m.Blocks[position.Y][position.X]
}

func (m *Map) SetBlock(position BlockPosition, block Block) {
	m.Blocks[position.Y][position.X] = block
}

func (m *Map) RoundStart(roundNum uint8) {
	for _, col := range m.Blocks {
		for _, block := range col {
			block.roundStart(roundNum)
		}
	}
}

type GameOverSign bool

func (m *Map) RoundEnd(roundNum uint8) GameOverSign {
	var ret GameOverSign
	for _, col := range m.Blocks {
		for _, block := range col {
			if _, s := block.roundEnd(roundNum); s {
				ret = true
			}
		}
	}
	return ret
}

func (m *Map) OutputNumber() { // Only for debugging
	for _, col := range m.Blocks {
		for _, block := range col {
			fmt.Print(block.GetNumber())
		}
		fmt.Print("\n")
	}
}
