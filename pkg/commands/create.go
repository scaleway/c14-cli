package commands

import "github.com/docker/docker/pkg/mflag"

type create struct {
	Base
}

func Create() Command {
	return &create{
		Base: Base{
			Config: Config{
				UsageLine:   "",
				Description: "",
				Help:        "",
				Examples:    "",
			},
		},
	}
}

func (c *create) GetName() string {
	return "create"
}

func (c *create) Parse(args []string) (err error) {
	if err = mflag.CommandLine.Parse(args); err != nil {
		return
	}
	return nil
}

func (c *create) Run() error {
	return nil
}
