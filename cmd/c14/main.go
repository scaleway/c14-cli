package main

import (
	"fmt"

	"github.com/QuentinPerez/c14-cli/pkg/version"
)

func main() {
	fmt.Println(version.VERSION, ":", version.GITCOMMIT)
}
