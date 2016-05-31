package commands

import "fmt"

type create struct {
	Base
	createFlags
}

type createFlags struct {
}

func Create() Command {
	ret := &create{}
	ret.Init(Config{
		UsageLine:   "",
		Description: "test",
		Help:        "",
		Examples:    "",
	})
	return ret
}

func (c *create) GetName() string {
	return "create"
}

func (c *create) Run(args []string) (err error) {
	fmt.Printf("args %v\n", args)
	return
}
