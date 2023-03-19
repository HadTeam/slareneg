package GameType

type UserStatus uint8

const (
	UserStatusConnected    UserStatus = 1
	UserStatusDisconnected UserStatus = 2
)

type User struct {
	Name   string
	UserId uint8
	Status UserStatus
}
