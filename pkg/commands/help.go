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

// Help returns a new command "help"
func Help() Command {
	ret := &help{}
	ret.Init(Config{
		UsageLine:   "help [COMMAND]",
		Description: "Help of the c14 command line",
		Help: `Help prints help information about c14 and its commands.
By default, help lists available commands.
When invoked with a command name, it prints the usage and the help of
the command.`,
		Examples: `
    $ c14 help
    $ c14 help create
`,
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
				cmd.PrintUsage()
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
{{range .}}    {{.GetName | printf "%-9s"}} {{.Base.Description}}
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
