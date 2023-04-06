package main

import (
	"context"
	"fmt"
	"server/JudgePool"
	"server/Untils/pkg/DataSource/Local"
	"server/Untils/pkg/GameType"
	"server/Untils/pkg/InstructionType"
	_ "server/Untils/pkg/MapType/Blocks"
)

func main() {
	ctx, exit := context.WithCancel(context.Background())
	defer exit()
	data = &Local.Pool



	<-ctx.Done()
}
