package commands

type list struct {
	Base
	listFlags
}

type listFlags struct {
}

// List returns a new command "list"
func List() Command {
	ret := &list{}
	ret.Init(Config{
		UsageLine:   "",
		Description: "",
		Help:        "",
		Examples:    "",
	})
	return ret
}

func (l *list) GetName() string {
	return "list"
}

func (l *list) Run(args []string) (err error) {
	if err = l.InitAPI(); err != nil {
		return err
	}

	return
}
