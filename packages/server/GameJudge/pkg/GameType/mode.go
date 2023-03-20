package GameType

type GameMode struct {
	MaxUserNum uint8
	MinUserNum uint8
}

var GameMode1v1 = GameMode{MaxUserNum: 2, MinUserNum: 2}
