package command

type moveCommand struct {
	X, Y    string `validate:"required,min=1,max=8,alphanum,gt=0"`
	Towards string `validate:"required,oneof='up' 'down' 'right' 'left'"`
	Number  string `validate:"required,min=1,max=8,alphanum"`
}

type forceStartCommand struct {
	Status string `validate:"required,oneof='true' 'false'"`
}
