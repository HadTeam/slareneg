package MapType

type MapSize struct{ X, Y uint8 }

type Map struct {
	Blocks [][]Block
	Size   MapSize
	MapId  uint32
}

func (m Map) GetBlock(position BlockPosition) Block {
	return m.Blocks[position.Y][position.X]
}

func (m Map) SetBlock(position BlockPosition, block Block) {
	m.Blocks[position.Y][position.X] = block
}
