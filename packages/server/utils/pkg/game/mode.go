package game

type GameMode struct {
	MaxUserNum uint8
	MinUserNum uint8
	NameStr    string
}

var GameMode1v1 = GameMode{MaxUserNum: 2, MinUserNum: 2, NameStr: "1v1"}
