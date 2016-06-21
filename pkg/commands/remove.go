package commands

import "fmt"

type remove struct {
	Base
	removeFlags
}

type removeFlags struct {
}

// remove returns a new command "remove"
func Remove() Command {
	ret := &remove{}
	ret.Init(Config{
		UsageLine:   "remove [ARGS]+",
		Description: "",
		Help:        "",
		Examples: `
        $ c14 remove 83b93179-32e0-11e6-be10-10604b9b0ad9 2d752399-429f-447f-85cd-c6104dfed5db`,
	})
	return ret
}

func (r *remove) GetName() string {
	return "remove"
}

func (r *remove) Run(args []string) (err error) {
	if err = r.InitAPI(); err != nil {
		return
	}
	for _, uuid := range args {
		fmt.Println(uuid)
	}
	return
}
