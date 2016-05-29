package main

import (
	"fmt"

	"github.com/QuentinPerez/c14-cli/pkg/commands"
	"github.com/QuentinPerez/c14-cli/pkg/version"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func main() {
	root := commands.NewRoot()

	log.SetHandler(cli.Default)
	fmt.Println(version.VERSION, ":", version.GITCOMMIT)
	if err := root.Parse(); err != nil {
		log.Fatalf("%v", err)
	}
}
