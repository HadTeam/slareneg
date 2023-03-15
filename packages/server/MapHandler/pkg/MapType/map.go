package MapType

type MapSize struct{ x, y uint8 }

type Map struct {
	Blocks [][]Block
	size   MapSize
}
