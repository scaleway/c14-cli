package commands

type Command interface {
	GetName() string
	Parse() error
	Run() error
}

type Base struct {
	Debug bool
}
