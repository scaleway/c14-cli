package commands

type login struct {
	Base
	loginFlags
}

type loginFlags struct {
}

// Login returns a new command "login"
func Login() Command {
	ret := &login{}
	ret.Init(Config{
		UsageLine:   "",
		Description: "",
		Help:        "",
		Examples:    "",
	})
	return ret
}

func (l *login) GetName() string {
	return "login"
}

func (l *login) Run(args []string) (err error) {
	return
}
