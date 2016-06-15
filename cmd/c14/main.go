package main

import (
	"os"

	"github.com/QuentinPerez/c14-cli/pkg/commands"
	"github.com/QuentinPerez/c14-cli/pkg/version"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/juju/errors"
)

func main() {
	log.SetHandler(text.New(os.Stderr))
	if err := commands.Root.Parse(); err != nil {
		log.Fatalf("↴\n%s\n\nVersion : %v\nCommit : %v", errors.ErrorStack(err), version.VERSION, version.GITCOMMIT)
	}
}
