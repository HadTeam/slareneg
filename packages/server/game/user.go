package game

type UserStatus uint8

const (
	UserStatusConnected    UserStatus = 1
	UserStatusDisconnected UserStatus = 2
)

type User struct {
	Name             string
	UserId           uint16
	Status           UserStatus
	TeamId           uint8
	ForceStartStatus bool
}
