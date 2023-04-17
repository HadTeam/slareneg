package main

import (
	"context"
	"server/JudgePool"
	"server/Utils/pkg/DataSource/Local"
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/InstructionType"
	_ "server/Utils/pkg/MapType/Blocks"
	"strconv"
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
	_ = JudgePool.CreatePool([]GameType.GameMode{GameType.GameMode1v1})

	time.Sleep(200 * time.Millisecond)

	var id GameType.GameId
	for i, _ := range data.GamePool {
		id = i
	}
	v := data.GetGameInfo(id)

	for i := uint8(1); i <= GameType.GameMode1v1.MaxUserNum; i++ {
		data.SetUserStatus(v.Id, GameType.User{
			Name:   "tester" + strconv.Itoa(int(i)),
			UserId: uint16(i),
			Status: GameType.UserStatusConnected,
		})
	}

	<-ctx.Done()
}
