package commands

import (
	"fmt"

	"github.com/QuentinPerez/c14-cli/pkg/api"
)

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
		return
	}
	var (
		val []api.OnlineGetSafe
	)

	if val, err = l.OnlineAPI.GetSafes(); err != nil {
		return
	}
	for i := range val {
		fmt.Println(val[i])
	}
	return
}
