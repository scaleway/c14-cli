package commands

import "fmt"

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

func (r *Root) Parse() error {
	for _, cmd := range r.Commands() {
		fmt.Println("->", cmd)
	}
	return nil
}

// Commands returns a string array with the commands name
func (r *Root) Commands() (commands []string) {
	commands = make([]string, len(r.commands))
	for i, cmd := range r.commands {
		commands[i] = cmd.GetName()
	}
	return
}
