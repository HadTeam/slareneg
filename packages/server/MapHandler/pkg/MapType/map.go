package MapType

type MapSize struct{ X, Y uint8 }

type Map struct {
	Blocks [][]Block
	Size   MapSize
}
