package commands

type Config struct {
	UsageLine   string
	Description string
	Help        string
	Examples    string
}

type Env struct {
	Debug bool
}

type Command interface {
	GetBase() *Base
	GetName() string
	Parse(args []string) error
	Run() error
}

type Base struct {
	Config
	Env
}

func (b *Base) GetBase() *Base {
	return b
}
