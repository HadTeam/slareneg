package main

import (
	"context"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"server/api"
	judge_pool "server/judgepool"
	"server/utils/pkg/datasource/local"
	"server/utils/pkg/game"
	"server/utils/pkg/instruction"
	_ "server/utils/pkg/map/block"
	"time"
)

func main() {
	ctx, exit := context.WithCancel(context.Background())
	defer exit()

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

	<-ctx.Done()
}
