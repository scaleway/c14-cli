package commands

import "fmt"

type create struct {
	Base
	createFlags
}

type createFlags struct {
}

// Create returns a new command "create"
func Create() Command {
	ret := &create{}
	ret.Init(Config{
		UsageLine:   "",
		Description: "",
		Help:        "",
		Examples:    "",
	})
	return ret
}

func (c *create) GetName() string {
	return "create"
}

func (c *create) Run(args []string) (err error) {
	fmt.Println(Root.Debug)
	fmt.Printf("args %v\n", args)
	return
}
