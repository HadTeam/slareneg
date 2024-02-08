package api

import (
	"server/api/internal/command"
	"server/api/internal/receiver"
	"server/pool"
	"time"
)

func ApplyDataSource(source any) {
	receiver.ApplyDataSource(source)
	command.ApplyDataSource(source)
}

func DebugStartFileReceiver(pool *pool.Pool) {
	time.Sleep(200 * time.Millisecond)
	receiver.NewFileReceiver(pool)
}

func Start() {
}
