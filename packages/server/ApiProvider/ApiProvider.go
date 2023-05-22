package ApiProvider

import (
	"server/ApiProvider/internal/CommandPauser"
	"server/ApiProvider/internal/Receiver"
	"server/JudgePool"
	"time"
)

func ApplyDataSource(source any) {
	Receiver.ApplyDataSource(source)
	CommandPauser.ApplyDataSource(source)
}

func Test(pool *JudgePool.Pool) {
	time.Sleep(200 * time.Millisecond)
	Receiver.NewFileReceiver(pool)
}
