package judge

import (
	"server/game/block"
	_map "server/game/map"
	"server/game/user"
)

func judgeGameMode1v1(m *_map.Map, userList []user.User, kingPos []block.Position) (Result, Winner) {
	flag := true
	for _, k := range kingPos {
		if m.GetBlock(k).Meta().BlockId != block.KingMeta.BlockId {
			flag = false
			break
		}
	}
	if !flag {
		var w uint16
		for _, k := range kingPos {
			if m.GetBlock(k).Meta().BlockId == block.KingMeta.BlockId {
				w = m.GetBlock(k).OwnerId()
			}
		}
		var wt Winner
		for _, u := range userList {
			if u.UserId == w {
				wt = Winner(u.TeamId)
				break
			}
		}
		return ResultEnd, wt
	}
	return ResultContinue, WinnerNone
}
