package user

type Status uint8

const Unknown = uint16(0)

const (
	Connected    Status = 1
	Disconnected Status = 2
)

type User struct {
	Name             string
	UserId           uint16
	Status           Status
	TeamId           uint8
	ForceStartStatus bool
}
