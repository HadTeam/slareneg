package Local

import (
	"fmt"
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/InstructionType"
	"server/Utils/pkg/MapType"
	"server/Utils/pkg/MapType/BlockType"
	"strconv"
	"testing"
)

func init() {

}

func create() *Local {
	return &Local{
		GamePool:           make(map[GameType.GameId]*GameType.Game),
		OriginalMapStrPool: make(map[uint32]string),
		InstructionLog:     make(map[GameType.GameId]map[uint16]map[uint16]InstructionType.Instruction),
	}
}

var userCount = uint16(1)

func getUser() GameType.User {
	return GameType.User{
		Name:             strconv.Itoa(int(userCount)),
		UserId:           userCount,
		Status:           GameType.UserStatusDisconnected,
		TeamId:           uint8(userCount),
		ForceStartStatus: false,
	}
}

func TestLocal_lock(t *testing.T) {
	l := create()
	t.Run("lock", func(t *testing.T) {
		l.lock()
		if l.m.TryLock() {
			t.Fatalf("the lock has not locked")
		}
	})
	t.Run("unlock", func(t *testing.T) {
		l.unlock()
		if !l.m.TryLock() {
			t.Fatalf("the lock has not unlocked")
		}
	})
}

func TestLocal_Game(t *testing.T) {
	l := create()
	id := GameType.GameId(1)

	// This function must be correct
	l.DebugCreateGame(&GameType.Game{
		Map:        MapType.FullStr2GameMap(1, "[[[0,0,0]]]"), // Wait
		Mode:       GameType.GameMode1v1,
		Id:         id,
		UserList:   []GameType.User{getUser()},
		CreateTime: 0,
		Status:     GameType.GameStatusWaiting,
		RoundNum:   0,
		Winner:     0,
	})

	t.Run("set game status", func(t *testing.T) {
		l.SetGameStatus(id, GameType.GameStatusRunning)
		if l.GamePool[id].Status != GameType.GameStatusRunning {
			t.Fatalf("the status has unchanged")
		}
	})
	t.Run("get game status", func(t *testing.T) {
		g := l.GetGameInfo(id)
		if g == l.GamePool[id] {
			t.Fatalf("game struct copied")
		}
	})
	t.Run("get game list", func(t *testing.T) {
		gl := l.GetGameList(GameType.GameMode1v1)
		if len(gl) != 1 {
			t.Fatalf("game count is incorrect")
		}
		if gl[0].Id != id {
			t.Fatalf("game info is incorrect")
		}
	})

	t.Run("instruction temp", func(t *testing.T) {
		l.NewInstructionTemp(id, 1)
		if l.InstructionLog[id][1] == nil {
			t.Fatalf("instruction temp has not been created")
		}
	})
	t.Run("update instruction", func(t *testing.T) {
		u := l.GamePool[id].UserList[0]
		ins := InstructionType.Move{
			Position: InstructionType.BlockPosition{1, 1},
			Towards:  "down",
			Number:   1,
		}
		l.UpdateInstruction(id, u, ins)
		if l.InstructionLog[id][1][u.UserId] != ins {
			t.Fatalf("instruction temp havn't created")
		}
	})

	t.Run("cancel game", func(t *testing.T) {
		l.CancelGame(id)
		if l.GamePool[id].Status != GameType.GameStatusEnd {
			t.Fatalf("the status has unchanged")
		}
		if l.GamePool[id].UserList != nil {
			t.Fatalf("the userlist has not been cleared")
		}
	})
}

func TestLocal_User(t *testing.T) {
	l := create()
	id := GameType.GameId(2)
	l.DebugCreateGame(&GameType.Game{
		Map:        MapType.FullStr2GameMap(1, "[[[0,0,0]]]"),
		Mode:       GameType.GameMode1v1,
		Id:         id,
		UserList:   []GameType.User{},
		CreateTime: 0,
		Status:     GameType.GameStatusWaiting,
		RoundNum:   0,
		Winner:     0,
	})
	u := getUser()
	t.Run("set user status", func(t *testing.T) {
		l.SetUserStatus(id, u)
		if l.GamePool[id].UserList != nil && len(l.GamePool[id].UserList) != 0 {
			t.Fatalf("user has been added unexpectedly")
		}

		u.Status = GameType.UserStatusConnected
		l.SetUserStatus(id, u)

		if len(l.GamePool[id].UserList) != 1 || l.GamePool[id].UserList[0].UserId != u.UserId {
			t.Fatalf("user has not been added")
		}

		u.Status = GameType.UserStatusDisconnected

		l.SetUserStatus(id, u)

		if len(l.GamePool[2].UserList) != 0 {
			t.Fatalf("user has not been removed")
		}
	})
}

func TestLocal_Map(t *testing.T) {
	l := create()
	id := GameType.GameId(3)
	l.DebugCreateGame(&GameType.Game{
		Map:        MapType.FullStr2GameMap(1, "[[[0,0,0]]]"),
		Mode:       GameType.GameMode1v1,
		Id:         id,
		UserList:   []GameType.User{},
		CreateTime: 0,
		Status:     GameType.GameStatusWaiting,
		RoundNum:   0,
		Winner:     0,
	})
	t.Run("test initial map", func(t *testing.T) {
		m := l.GamePool[id].Map
		if !m.HasBlocks() || m.Id() != 1 {
			t.Fatalf("initial map isn't correct")
		}
	})
	t.Run("get current map", func(t *testing.T) {
		m := l.GetCurrentMap(id)
		if !m.HasBlocks() || m.Id() != 1 {
			t.Fatalf("current map isn't correct")
		}
	})
	t.Run("set map", func(t *testing.T) {
		compareMap := func(a *MapType.Map, b *MapType.Map) bool { // if a == b, returns true
			ret := true
			if a.Size() == b.Size() && a.Id() == b.Id() {
				for y := uint8(1); y <= a.Size().H; y++ {
					for x := uint8(1); x <= a.Size().W; x++ {
						blockA := a.GetBlock(BlockType.Position{x, y})
						blockB := b.GetBlock(BlockType.Position{x, y})
						if blockA != blockB {
							t.Log(fmt.Sprintf("%#v", blockA), fmt.Sprintf("%#v", blockB), x, y)
							ret = false
							break
						}
					}
				}
			} else {
				t.Log("id")
				ret = false
			}

			return ret
		}

		m := l.GetCurrentMap(id)
		newM := MapType.FullStr2GameMap(1, "[[[1,0,0]]]")
		l.SetGameMap(id, newM)
		currM := l.GetCurrentMap(id)
		if compareMap(m, currM) != false || compareMap(newM, currM) != true {
			t.Fatalf("the map in the pool is unexpected")
		}
		if currM == newM {
			t.Fatalf("map copied")
		}
	})
}
