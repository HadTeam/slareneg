package CommandPauser

type moveCommand struct {
	X, Y    string `validate:"required,min=1,max=8,alphanum"`
	Towards string `validate:"required,oneof='up' 'down' 'right' 'left'"`
	Number  string `validate:"required,min=1,max=8,alphanum"`
}

type forceStartCommand struct {
	Status string `validate:"required,oneof='true' 'false'"`
}
