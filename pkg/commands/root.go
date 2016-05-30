package commands

import (
	"errors"
	"os"

	"github.com/apex/log"
	"github.com/docker/docker/pkg/mflag"
)

// Root handles the commands
type Root struct {
	commands []Command
}

// NewRoot returns a structure to handle the commands
func NewRoot() *Root {
	return &Root{
		commands: []Command{
			Create(),
		},
	}
}

func (r *Root) Parse() (err error) {
	var (
		flDebug = mflag.Bool([]string{"D", "-debug"}, false, "Enable debug mode")
	)

	args := os.Args[1:]
	if err = mflag.CommandLine.Parse(args); err != nil {
		return
	}
	env := Env{
		Debug: *flDebug,
	}
	if env.Debug {
		log.SetLevel(log.DebugLevel)
	}

	args = mflag.Args()
	if len(args) < 1 {
		err = errors.New("TODO: help message with the commands available")
		return
	}
	for _, cmd := range r.commands {
		if cmd.GetName() == args[0] {
			cmd.GetBase().Env = env
			err = cmd.Parse(args[1:])
			return
		}
	}
	return
}

// Commands returns a string array with the commands name
func (r *Root) Commands() (commands []string) {
	commands = make([]string, len(r.commands))
	for i, cmd := range r.commands {
		commands[i] = cmd.GetName()
	}
	return
}
