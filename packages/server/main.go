package main

import (
	"context"
	"server/ApiProvider"
	"server/ApiProvider/pkg/DataOperator"
	"server/ApiProvider/pkg/DataOperator/Local"
	_ "server/ApiProvider/pkg/DataOperator/Local"
	"server/ApiProvider/pkg/GameOperator"
	"server/GameJudge"
	"server/GameJudge/pkg/GameType"
)

var data DataOperator.DataSource

func main() {
	ctx, exit := context.WithCancel(context.Background())
	defer exit()
	data = &Local.Pool


	workPool := make(chan GameType.GameId)

	ApiProvider.ApplyDataSource(data)
	GameJudge.ApplyDataSource(data)
	_ = GameJudge.NewGameJudge(workPool)
	gameId := GameOperator.NewGame(0, GameType.GameMode1v1)
	workPool <- gameId
	<-ctx.Done()
}
