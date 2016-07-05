package commands

import (
	"os"

	"github.com/QuentinPerez/c14-cli/pkg/api"
	"github.com/apex/log"
)

type freeze struct {
	Base
}

type freezeFlags struct {
}

// Freeze returns a new command "freeze"
func Freeze() Command {
	ret := &freeze{}
	ret.Init(Config{
		UsageLine:   "freeze [OPTIONS] [ARCHIVE]+",
		Description: "",
		Help:        "",
		Examples: `
        $ c14 freeze 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	return ret
}

func (f *freeze) GetName() string {
	return "freeze"
}

func (f *freeze) CheckFlags(args []string) (err error) {
	if len(args) == 0 {
		f.PrintUsage()
		os.Exit(1)
	}
	return
}

func (f *freeze) Run(args []string) (err error) {
	if err = f.InitAPI(); err != nil {
		return
	}
	// f.FetchRessources(true, false)

	var (
		safe        api.OnlineGetSafe
		uuidArchive string
	)

	for _, archive := range args {
		if safe, uuidArchive, err = f.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
			return
		}
		if err = f.OnlineAPI.PostArchive(safe.UUIDRef, uuidArchive); err != nil {
			log.Warnf("%s", err)
			err = nil
		}
	}
	return
}
