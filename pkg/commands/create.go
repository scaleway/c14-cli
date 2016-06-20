package commands

import "github.com/docker/docker/pkg/namesgenerator"

type create struct {
	Base
	createFlags
}

type createFlags struct {
	flName string
	flDesc string
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
	ret.Flags.StringVar(&ret.flName, []string{"n", "-name"}, "", "Assigns a name")
	ret.Flags.StringVar(&ret.flDesc, []string{"d", "-description"}, "", "Assigns a description")
	return ret
}

func (c *create) GetName() string {
	return "create"
}

func (c *create) Run(args []string) (err error) {
	if err = c.InitAPI(); err != nil {
		return err
	}
	if c.flName == "" {
		c.flName = namesgenerator.GetRandomName(0)
	}
	err = c.OnlineAPI.CreateSafe(c.flName, c.flDesc)
	return
}
