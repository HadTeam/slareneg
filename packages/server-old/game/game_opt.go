package game

import (
	"math"
	"math/rand"
	"server/game/block"
	"server/game/mode"
	"server/game/user"
)

// allocateKing allocates a King to the users who don't own one yet.
func (g *Game) allocateKing() {
	m := g.Map
	var kingPos []block.Position               // To store the positions of Kings that haven't been allocated
	userKingMap := map[uint16]block.Position{} // To store the user who already own a King

	for _, k := range g.getKingPos() {
		id := m.GetBlock(k).OwnerId()
		if id != 0 {
			userKingMap[id] = k
		} else {
			kingPos = append(kingPos, k)
		}
	}

	rand.Shuffle(len(kingPos), func(i, j int) {
		kingPos[i], kingPos[j] = kingPos[j], kingPos[i]
	})

	var userNeedKing []user.User
	for _, user := range g.UserList {
		if _, exists := userKingMap[user.UserId]; !exists {
			userNeedKing = append(userNeedKing, user)
		}
	}

	allocNum := int(math.Min(float64(len(userNeedKing)), float64(len(kingPos))))
	userNeedKing = userNeedKing[:allocNum]
	kingNeedTrans := kingPos[allocNum:]
	kingPos = kingPos[:allocNum]

	for i, user := range userNeedKing {
		pos := kingPos[i]
		b := m.GetBlock(pos)
		m.SetBlock(pos, block.NewBlock(block.KingMeta.BlockId, b.Number(), user.UserId))
		// add new King position to the map
		userKingMap[user.UserId] = pos
	}

	for _, pos := range kingNeedTrans {
		m.SetBlock(pos, block.NewBlock(block.CastleMeta.BlockId, 0, 0))
	}
}

func (g *Game) allocateTeam() {
	if g.Mode == mode.Mode1v1 {
		for i := range g.UserList {
			g.UserList[i].TeamId = uint8(i) + 1
		}
	} else {
		panic("unexpected game mod")
	}
}

func (g *Game) getKingPos() []block.Position {
	var kingPos []block.Position
	for y := uint8(1); y <= g.Map.Size().H; y++ {
		for x := uint8(1); x <= g.Map.Size().W; x++ {
			b := g.Map.GetBlock(block.Position{X: x, Y: y})
			if b.Meta().BlockId == block.KingMeta.BlockId {
				kingPos = append(kingPos, block.Position{X: x, Y: y})
			}
		}
	}
	return kingPos
}
