package game

type Mode struct {
	MaxUserNum uint8
	MinUserNum uint8
	NameStr    string
}

var Mode1v1 = Mode{MaxUserNum: 2, MinUserNum: 2, NameStr: "1v1"}
