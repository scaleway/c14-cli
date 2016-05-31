package main

import (
	"os"

	"github.com/QuentinPerez/c14-cli/pkg/commands"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

func main() {
	root := commands.NewRoot()

	log.SetHandler(text.New(os.Stderr))
	if err := root.Parse(); err != nil {
		log.Fatalf("%v", err)
	}
}
