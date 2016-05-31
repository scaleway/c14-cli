package commands

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/docker/docker/pkg/mflag"
)

// Root handles the commands
type root struct {
	commands []Command
}

var Root *root

func init() {
	Root = &root{
		commands: []Command{
			Help(),
			Create(),
		},
	}
}

func (r *root) Parse() (err error) {
	var (
		flDebug = mflag.Bool([]string{"D", "-debug"}, false, "Enable debug mode")
	)

	args := os.Args[1:]
	if err = mflag.CommandLine.Parse(args); err != nil {
		return
	}

	env := Env{
		Debug: *flDebug || os.Getenv("C14_DEBUG") == "1",
	}
	if env.Debug {
		log.SetLevel(log.DebugLevel)
	}

	args = mflag.Args()
	if len(args) < 1 {
		r.printUsage(args)
		return
	}
	for _, cmd := range r.commands {
		if cmd.GetName() == args[0] {
			cmd.GetBase().Env = env
			if args, err = cmd.Parse(args[1:]); err != nil {
				return
			}
			err = cmd.Run(args)
			return
		}
	}
	err = fmt.Errorf(`c14: unknow command %v
Run 'c14 help' for usage`, args[0])
	return
}

func (r *root) printUsage(args []string) {
	for _, cmd := range r.commands {
		if cmd.GetName() == "help" {
			cmd.Run(args)
			os.Exit(1)
		}
	}
	panic("No help method")
}

// Commands returns a string array with the commands name
func (r *root) Commands() (commands []string) {
	commands = make([]string, len(r.commands))
	for i, cmd := range r.commands {
		commands[i] = cmd.GetName()
	}
	return
}
