package commands

import (
	"fmt"
	"html/template"
	"os"

	"github.com/docker/docker/pkg/mflag"
)

type help struct {
	Base
}

func Help() Command {
	ret := &help{}
	ret.Init(Config{
		UsageLine:   "help [COMMAND]",
		Description: "Help of the c14 command line",
		Help:        "",
		Examples:    "",
	})
	return ret
}

func (c *help) GetName() string {
	return "help"
}

func (c *help) Run(args []string) (err error) {
	if len(args) > 0 {
		for _, cmd := range Root.commands {
			if cmd.GetName() == args[0] {
				//
				return
			}
		}
		err = fmt.Errorf("help: unknow command %v", args[0])
		return
	}

	var header, end = `Usage: c14 [OPTIONS] COMMAND [arg...]

Interact with C14 from the command line.

Options:`, `
Commands:
{{range .}}    {{.GetName | printf "%-9s"}} {{.GetBase.Description}}
{{end}}
Run 'c14 COMMAND --help' for more information on a command.
`

	fmt.Printf("%v", header)
	mflag.PrintDefaults()
	t := template.New("helper")
	template.Must(t.Parse(end))
	err = t.Execute(os.Stdout, Root.commands)
	return
}
