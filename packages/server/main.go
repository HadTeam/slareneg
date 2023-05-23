package main

import (
	"context"
	"server/api"
	"server/judgePool"
	"server/utils/pkg/dataSource/local"
	"server/utils/pkg/game"
	"server/utils/pkg/instruction"
	_ "server/utils/pkg/map/block"
	"time"
)

func main() {
	ctx, exit := context.WithCancel(context.Background())
	defer exit()

	data := local.Local{
		GamePool:           make(map[game.GameId]*game.Game),
		OriginalMapStrPool: make(map[uint32]string),
		InstructionLog:     make(map[game.GameId]map[uint16]map[uint16]instruction.Instruction),
	}
	data.OriginalMapStrPool[0] = "[\n[0,0,0,0,2],\n[0,2,0,0,0],\n[0,0,0,0,0],\n[0,3,3,0,3],\n[0,3,0,2,0]\n]"

	judgePool.ApplyDataSource(&data)
	p := judgePool.CreatePool([]game.GameMode{game.GameMode1v1})

	time.Sleep(200 * time.Millisecond)

	api.ApplyDataSource(&data)
	api.DebugStartFileReceiver(p)

	<-ctx.Done()
}
