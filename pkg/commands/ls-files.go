package commands

import (
	"fmt"

	"github.com/QuentinPerez/c14-cli/pkg/api"
)

type lsFiles struct {
	Base
	lsFilesFlags
}

type lsFilesFlags struct {
}

// LsFiles returns a new command "lsFiles"
func LsFiles() Command {
	ret := &lsFiles{}
	ret.Init(Config{
		UsageLine:   "ls-files [ARCHIVE]+",
		Description: "",
		Help:        "",
		Examples: `
        $ c14 ls-files 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	return ret
}

func (l *lsFiles) GetName() string {
	return "ls-files"
}

func (l *lsFiles) Run(args []string) (err error) {
	if len(args) == 0 {
		l.PrintUsage()
		return
	}
	if err = l.InitAPI(); err != nil {
		return
	}
	l.FetchRessources(true, true)

	var (
		safe        api.OnlineGetSafe
		bucket      api.OnlineGetBucket
		uuidArchive = args[0]
	)
	if safe, err = l.OnlineAPI.FindSafeUUIDFromArchive(uuidArchive, true); err != nil {
		return
	}
	if bucket, err = l.OnlineAPI.GetBucket(safe.UUIDRef, uuidArchive); err != nil {
		return
	}
	fmt.Println(bucket)
	return
}
