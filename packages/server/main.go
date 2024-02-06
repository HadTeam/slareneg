package main

import (
	"context"
	"flag"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gookit/ini/v2"
	"github.com/sirupsen/logrus"
	"server/api"
	"server/game_logic"
	_ "server/game_logic/block"
	"server/game_logic/game_def"
	judge_pool "server/judge_pool"
	"server/utils/pkg/data_source/local"
	db "server/utils/pkg/pg"
	"time"
)

var configFile string

const defaultConfigPath = "./slareneg.server.ini"
const defaultConfigOptions = `
	[db]
	host = localhost
	port = 5432
	user = postgres
	password = slareneg
	sslMode = disable
	name = slareneg
	`

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: time.RFC3339,
	})

	flag.StringVar(&configFile, "config", defaultConfigPath, "config file path")
	flag.Parse()

	if err := ini.LoadStrings(defaultConfigOptions); err != nil {
		logrus.Panic(err)
	}

	if err := ini.LoadExists(configFile); err != nil {
		logrus.Panic(err)
	}

	db.Bootstrap()
	defer db.Exit()

	ctx, exit := context.WithCancel(context.Background())
	defer exit()

	data := local.Local{
		GamePool:           make(map[game_logic.Id]*game_logic.Game),
		OriginalMapStrPool: make(map[uint32]string),
		InstructionLog:     make(map[game_logic.Id]map[uint16]map[uint16]game_def.Instruction),
	}
	data.OriginalMapStrPool[0] = "[\n[0,0,0,0,2],\n[0,2,0,0,0],\n[0,0,0,0,0],\n[0,3,3,0,3],\n[0,3,0,2,0]\n]"

	judge_pool.ApplyDataSource(&data)
	p := judge_pool.CreatePool([]game_def.Mode{game_def.Mode1v1})

	time.Sleep(200 * time.Millisecond)

	api.ApplyDataSource(&data)
	api.DebugStartFileReceiver(p)

	<-ctx.Done()
}
