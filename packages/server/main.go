package main

import (
	"context"
	"server/ApiProvider"
	"server/JudgePool"
	"server/Utils/pkg/DataSource/Local"
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/InstructionType"
	_ "server/Utils/pkg/MapType/BlockType"
	"time"
)

func main() {
	ctx, exit := context.WithCancel(context.Background())
	defer exit()

	data := Local.Local{
		GamePool:            make(map[GameType.GameId]*GameType.Game),
		OriginalMapStrPool:  make(map[uint32]string),
		InstructionTempPool: make(map[GameType.GameId]map[uint16]InstructionType.Instruction),
		InstructionLog:      make(map[GameType.GameId]map[uint16][]InstructionType.Instruction),
	}
	data.OriginalMapStrPool[0] = "[\n[0,0,0,0,2],\n[0,2,0,0,0],\n[0,0,0,0,0],\n[0,3,3,0,3],\n[0,3,0,2,0]\n]"

	JudgePool.ApplyDataSource(&data)
	p := JudgePool.CreatePool([]GameType.GameMode{GameType.GameMode1v1})

	time.Sleep(200 * time.Millisecond)

	ApiProvider.ApplyDataSource(&data)
	ApiProvider.Test(p)

	<-ctx.Done()
}
