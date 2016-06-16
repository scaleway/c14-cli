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
	flQuiet bool
}

// List returns a new command "list"
func List() Command {
	ret := &list{}
	ret.Init(Config{
		UsageLine:   "list [OPTIONS]",
		Description: "",
		Help:        "",
		Examples: `
        $ c14 list
        $ c14 list 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	ret.Flags.BoolVar(&ret.flQuiet, []string{"q", "-quiet"}, false, "Only display UUIDs")
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
		if l.flQuiet {
			fmt.Println(val[i].UUIDRef)
		} else {
			fmt.Println(val[i])
		}
	}
	return
}
