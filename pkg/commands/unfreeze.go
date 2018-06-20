package commands

import (
	"os"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/apex/log"
	"github.com/scaleway/c14-cli/pkg/api"
	"github.com/scaleway/c14-cli/pkg/utils/pgbar"
	"github.com/pkg/errors"
)

type unfreeze struct {
	Base
	unfreezeFlags
}

type unfreezeFlags struct {
	flQuiet  bool
	flNoWait bool
}

// Unfreeze returns a new command "unfreeze"
func Unfreeze() Command {
	ret := &unfreeze{}
	ret.Init(Config{
		UsageLine:   "unfreeze [OPTIONS] [ARCHIVE]+",
		Description: "Unlock an archive",
		Help:        "Unlock an archive, extraction of the archive's data (0.01€/GB).",
		Examples: `
        $ c14 unfreeze 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	ret.Flags.BoolVar(&ret.flQuiet, []string{"q", "-quiet"}, false, "")
	ret.Flags.BoolVar(&ret.flNoWait, []string{"-nowait"}, false, "")
	return ret
}

func (f *unfreeze) GetName() string {
	return "unfreeze"
}

func (f *unfreeze) CheckFlags(args []string) (err error) {
	if len(args) == 0 {
		f.PrintUsage()
		os.Exit(1)
	}
	return
}

func (f *unfreeze) Run(args []string) (err error) {
	if err = f.InitAPI(); err != nil {
		return
	}

	var (
		safe        api.OnlineGetSafe
		keys        []api.OnlineGetSSHKey
		archiveWait api.OnlineGetArchive
		uuidArchive string
	)

	if keys, err = f.OnlineAPI.GetSSHKeys(); err != nil {
		return
	}
	if len(keys) == 0 {
		err = errors.Errorf("Please add an SSH Key here: https://console.online.net/en/account/ssh-keys")
		return
	}

	for _, archive := range args {
		if safe, uuidArchive, err = f.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
			if safe, uuidArchive, err = f.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
				return
			}
		}

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
			continue
		}
		if !f.flNoWait {
			var bar *pb.ProgressBar

			if !f.flQuiet {
				bar = pgbar.NewProgressBar(uuidArchive)
				bar.Start()
			}
			lastLength := 5
			for {
				if archiveWait, err = f.OnlineAPI.GetArchive(safe.UUIDRef, uuidArchive, false); err != nil {
					log.Warnf("%s: %s", args, err)
					err = nil
					break
				}
				if lastLength != len(archiveWait.Jobs) {
					lastLength = len(archiveWait.Jobs)
					if !f.flQuiet {
						bar.Add(25)
					}
					if len(archiveWait.Jobs) == 0 {
						break
					}
				}
				time.Sleep(1 * time.Second)
			}
			if !f.flQuiet {
				bar.Finish()
			}
		}
	}
	return
}
