package main

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"server/api"
	"server/game"
	_ "server/game/block"
	"server/game/instruction"
	judge_pool "server/game/judge_pool"
	"server/utils/pkg/data_source/local"
	"testing"
	"time"
)

func TestServer_main(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: time.RFC3339,
	})

	data := local.Local{
		GamePool:           make(map[game.Id]*game.Game),
		OriginalMapStrPool: make(map[uint32]string),
		InstructionLog:     make(map[game.Id]map[uint16]map[uint16]instruction.Instruction),
	}
	data.OriginalMapStrPool[0] = "[\n[0,0,0,0,2],\n[0,2,0,0,0],\n[0,0,0,0,0],\n[0,3,3,0,3],\n[0,3,0,2,0]\n]"

	judge_pool.ApplyDataSource(&data)
	p := judge_pool.CreatePool([]game.Mode{game.Mode1v1})

	time.Sleep(200 * time.Millisecond)

	api.ApplyDataSource(&data)
	api.DebugStartFileReceiver(p)

	id := game.Id(1000)
	ticker := time.NewTicker(1 * time.Second)
	for i := 1; true; i++ {
		<-ticker.C
		if i >= 20 {
			t.Fatalf("game running timeout(20s)")
		}
		if data.GetGameInfo(id).Status == game.StatusEnd {
			ticker.Stop()
			break
		}
	}

	var userList []string
	for _, u := range data.GamePool[id].UserList {
		if u.TeamId == data.GamePool[id].Winner {
			userList = append(userList, u.Name)
		}
	}
	if userList[0] != "test2" {
		t.Fatalf("game result is unexpected: expected [\"2\"], got %v (team %d)", userList, data.GamePool[id].Winner)
	}
}
