package main

import (
	"github.com/QuentinPerez/c14-cli/pkg/commands"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func main() {
	log.SetHandler(cli.Default)
	if err := commands.Root.Parse(); err != nil {
		log.Fatalf("%v", err)
	}
}
