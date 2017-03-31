package commands

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/spf13/pflag"
)

// Env containts the global options
type Env struct {
	Debug   bool
	Verbose bool
}

type root struct {
	commands []Command
	Streams
	Env
}

// Root handles the commands
var Root *root

func init() {
	Root = &root{
		commands: []Command{
			Create(),
			Files(),
			Freeze(),
			Help(),
			Login(),
			Ls(),
			Rename(),
			Remove(),
			Unfreeze(),
			Upload(),
			Verify(),
		},
	}
}

func (r *root) Parse() (err error) {
	var (
		flDebug   = pflag.BoolP("debug", "D", false, "Enable debug mode")
		flVerbose = pflag.BoolP("verbose", "V", false, "Enable verbose mode")
	)

	args := os.Args[1:]
	if err = pflag.CommandLine.Parse(args); err != nil {
		return
	}
	r.Verbose = *flVerbose || os.Getenv("C14_VERBOSE") == "1"
	r.Debug = *flDebug || os.Getenv("C14_DEBUG") == "1"
	if r.Debug {
		log.SetLevel(log.DebugLevel)
	}

	args = pflag.Args()
	if len(args) < 1 {
		r.printUsage(args)
		return
	}
	for _, cmd := range r.commands {
		if cmd.GetName() == args[0] {
			if args, err = cmd.Parse(args[1:]); err != nil {
				return
			}
			if err = cmd.CheckFlags(args); err != nil {
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
			_ = cmd.Run(args)
			os.Exit(1)
		}
	}
	log.Fatalf("No help method")
}

// Commands returns a string array with the commands name
func (r *root) Commands() (commands []string) {
	commands = make([]string, len(r.commands))
	for i, cmd := range r.commands {
		commands[i] = cmd.GetName()
	}
	return
}
