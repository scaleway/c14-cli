package commands

import (
	"os"

	"github.com/QuentinPerez/c14-cli/pkg/api"
	"github.com/apex/log"
	"github.com/pkg/errors"
)

type freeze struct {
	Base
	freezeFlags
}

type freezeFlags struct {
	flUnFreeze bool
}

// Freeze returns a new command "freeze"
func Freeze() Command {
	ret := &freeze{}
	ret.Init(Config{
		UsageLine:   "freeze [OPTIONS] [ARCHIVE]+",
		Description: "",
		Help:        "",
		Examples: `
        $ c14 freeze 83b93179-32e0-11e6-be10-10604b9b0ad9,
        $ c14 freeze --unfreeze 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	ret.Flags.BoolVar(&ret.flUnFreeze, []string{"u", "-unfreeze"}, false, "Unfreeze an archive")
	return ret
}

func (f *freeze) GetName() string {
	return "freeze"
}

func (f *freeze) CheckFlags(args []string) (err error) {
	if len(args) == 0 {
		f.PrintUsage()
		os.Exit(0)
	}
	return
}

func (f *freeze) Run(args []string) (err error) {
	if err = f.InitAPI(); err != nil {
		return
	}
	f.FetchRessources(true, false)

	var (
		safe        api.OnlineGetSafe
		keys        []api.OnlineGetSSHKey
		uuidArchive string
	)

	if f.flUnFreeze {
		if keys, err = f.OnlineAPI.GetSSHKeys(); err != nil {
			return
		}
		if len(keys) == 0 {
			err = errors.Errorf("Please add an SSH Key here: https://console.online.net/en/account/ssh-keys")
			return
		}
	}

	for _, archive := range args {
		if safe, uuidArchive, err = f.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
			return
		}
		if f.flUnFreeze {
			var (
				loc []api.OnlineGetLocation
			)

			if loc, err = f.OnlineAPI.GetLocations(safe.UUIDRef, uuidArchive); err != nil {
				log.Warnf("%s", err)
				err = nil
				continue
			}
			if err = f.OnlineAPI.PostUnArchive(safe.UUIDRef, uuidArchive, api.OnlinePostUnArchive{
				LocationID: loc[0].UUIDRef,
				Protocols:  []string{"SSH"},
				SSHKeys:    []string{keys[0].UUIDRef},
			}); err != nil {
				log.Warnf("%s", err)
				err = nil
			}
		} else {
			if err = f.OnlineAPI.PostArchive(safe.UUIDRef, uuidArchive); err != nil {
				log.Warnf("%s", err)
				err = nil
			}
		}
	}
	return
}
