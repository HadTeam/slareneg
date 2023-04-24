package ApiProvider

import (
	"server/ApiProvider/internal/CommandPauser"
	"server/ApiProvider/internal/Receiver"
	"time"
)

func ApplyDataSource(source any) {
	Receiver.ApplyDataSource(source)
	CommandPauser.ApplyDataSource(source)
}

func Test() {
	time.Sleep(200 * time.Millisecond)
	Receiver.NewFileReceiver()
}
