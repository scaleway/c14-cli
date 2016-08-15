package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/juju/errors"
	"github.com/online-net/c14-cli/pkg/commands"
)

func main() {
	log.SetHandler(text.New(os.Stderr))
	if err := commands.Root.Parse(); err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
}
