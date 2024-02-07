package judge

import (
	"server/game"
	"server/game/block"
)

func judgeGameMode1v1(g *game.Game, kingPos []block.Position) game.Status {
	flag := true
	for _, k := range kingPos {
		if g.Map.GetBlock(k).Meta().BlockId != block.KingMeta.BlockId {
			flag = false
			break
		}
	}
	if !flag {
		var w uint16
		for _, k := range kingPos {
			if g.Map.GetBlock(k).Meta().BlockId == block.KingMeta.BlockId {
				w = g.Map.GetBlock(k).OwnerId()
			}
		}
		var wt uint8
		for _, u := range g.UserList {
			if u.UserId == w {
				wt = u.TeamId
				break
			}
		}
		g.Winner = wt
		return game.StatusEnd
	}
	return game.StatusRunning
}
