package Local

import (
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/InstructionType"
	"server/Utils/pkg/MapType"
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

	// This function must be correct
	l.DebugCreateGame(&GameType.Game{
		Map:        MapType.FullStr2GameMap(1, "[[[0,0,0]]]"), // Wait
		Mode:       GameType.GameMode1v1,
		Id:         1,
		UserList:   []GameType.User{getUser()},
		CreateTime: 0,
		Status:     GameType.GameStatusWaiting,
		RoundNum:   0,
		Winner:     0,
	})

	t.Run("set game status", func(t *testing.T) {
		l.SetGameStatus(1, GameType.GameStatusRunning)
		if l.GamePool[1].Status != GameType.GameStatusRunning {
			t.Fatalf("the status has unchanged")
		}
	})
	t.Run("get game status", func(t *testing.T) {
		g := l.GetGameInfo(1)
		if g == l.GamePool[1] {
			t.Fatalf("game struct copied")
		}
	})
	t.Run("get game list", func(t *testing.T) {
		gl := l.GetGameList(GameType.GameMode1v1)
		if len(gl) != 1 {
			t.Fatalf("game count is incorrect")
		}
		if gl[0].Id != 1 {
			t.Fatalf("game info is incorrect")
		}
	})

	t.Run("instruction temp", func(t *testing.T) {
		l.NewInstructionTemp(1, 1)
		if l.InstructionLog[1][1] == nil {
			t.Fatalf("instruction temp has not been created")
		}
	})
	t.Run("update instruction", func(t *testing.T) {
		u := l.GamePool[1].UserList[0]
		ins := InstructionType.Move{
			Position: InstructionType.BlockPosition{1, 1},
			Towards:  "down",
			Number:   1,
		}
		l.UpdateInstruction(1, u, ins)
		if l.InstructionLog[1][1][u.UserId] != ins {
			t.Fatalf("instruction temp havn't created")
		}
	})

	t.Run("cancel game", func(t *testing.T) {
		l.CancelGame(1)
		if l.GamePool[1].Status != GameType.GameStatusEnd {
			t.Fatalf("the status has unchanged")
		}
		if l.GamePool[1].UserList != nil {
			t.Fatalf("the userlist has not been cleared")
		}
	})
}

func TestLocal_User(t *testing.T) {
	l := create()
	l.DebugCreateGame(&GameType.Game{
		Map:        MapType.FullStr2GameMap(1, "[[[0,0,0]]]"),
		Mode:       GameType.GameMode1v1,
		Id:         2,
		UserList:   []GameType.User{},
		CreateTime: 0,
		Status:     GameType.GameStatusWaiting,
		RoundNum:   0,
		Winner:     0,
	})
	u := getUser()
	t.Run("set user status", func(t *testing.T) {
		l.SetUserStatus(2, u)
		if l.GamePool[2].UserList != nil && len(l.GamePool[2].UserList) != 0 {
			t.Fatalf("user has been added unexpectedly")
		}

		u.Status = GameType.UserStatusConnected
		l.SetUserStatus(2, u)

		if len(l.GamePool[2].UserList) != 1 || l.GamePool[2].UserList[0].UserId != u.UserId {
			t.Fatalf("user has not been added")
		}

		u.Status = GameType.UserStatusDisconnected

		t.Log(l.GamePool[2].UserList)
		l.SetUserStatus(2, u)

		t.Log(l.GamePool[2].UserList)
		//l.GamePool[2].UserList = append(l.GamePool[2].UserList[:0], l.GamePool[2].UserList[0+1:]...) // TODO ?
		//
		//t.Log(l.GamePool[2].UserList)
		if len(l.GamePool[2].UserList) != 0 {
			t.Fatalf("user has not been removed")
		}
	})
}
