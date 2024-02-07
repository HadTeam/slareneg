package api

import (
	"server/api/internal/command"
	"server/api/internal/receiver"
	judge_pool "server/game/judge_pool"
	"time"
)

func ApplyDataSource(source any) {
	receiver.ApplyDataSource(source)
	command.ApplyDataSource(source)
}

func DebugStartFileReceiver(pool *judge_pool.Pool) {
	time.Sleep(200 * time.Millisecond)
	receiver.NewFileReceiver(pool)
}

func Start() {
}
