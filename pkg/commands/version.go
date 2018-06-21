package commands

import (
	"fmt"
	"github.com/scaleway/c14-cli/pkg/version"
	"runtime"
)

// struct named cliVersion instead of version du to conflict with version in
// commands package
type cliVersion struct {
	Base
}

func Version() Command {
	ret := &cliVersion{}
	ret.Init(Config{
		UsageLine:   "version ",
		Description: "Show the version information",
		Help:        "Show the version information.",
		Examples:    "$ c14 version",
	})

	return ret
}

func (c *cliVersion) GetName() string {
	return "version"
}

func (c *cliVersion) Run(args []string) (err error) {
	fmt.Println("Client version :", version.VERSION)
	fmt.Println("Go version :", runtime.Version())
	fmt.Println("OS/Arch :", runtime.GOOS, runtime.GOARCH)
	return
}
