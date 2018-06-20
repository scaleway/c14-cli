package commands

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/scaleway/c14-cli/pkg/api"
)

type verify struct {
	Base
	verifyFlags
}

type verifyFlags struct {
	// flQuiet  bool
	// flNoWait bool
}

// Verify returns a new command "verify"
func Verify() Command {
	ret := &verify{}
	ret.Init(Config{
		UsageLine:   "verify [ARCHIVE]+",
		Description: "Schedules a verification of the files on an archive's location",
		Help:        "Schedules a verification of the files on an archive's location.",
		Examples: `
        $ c14 verify 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	// ret.Flags.BoolVar(&ret.flQuiet, []string{"q", "-quiet"}, false, "")
	// ret.Flags.BoolVar(&ret.flNoWait, []string{"-nowait"}, false, "")
	return ret
}

func (f *verify) GetName() string {
	return "verify"
}

func (f *verify) CheckFlags(args []string) (err error) {
	if len(args) == 0 {
		f.PrintUsage()
		os.Exit(1)
	}
	return
}

func (f *verify) Run(args []string) (err error) {
	if err = f.InitAPI(); err != nil {
		return
	}

	var (
		safe             api.OnlineGetSafe
		archiveWait      api.OnlineGetArchive
		archiveLocations []api.OnlineGetLocation
		uuidArchive      string
	)

	for _, archive := range args {
		if safe, uuidArchive, err = f.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
			if safe, uuidArchive, err = f.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
				return
			}
		}
		if archiveLocations, err = f.OnlineAPI.GetLocations(safe.UUIDRef, uuidArchive); err != nil {
			log.Warnf("%s: %s", archive, err)
			err = nil
			continue
		}
		if len(archiveLocations) == 0 {
			log.Warnf("%s: no location", archive)
			continue
		}
		for _, location := range archiveLocations {
			if err = f.OnlineAPI.PostVerify(safe.UUIDRef, uuidArchive, location.UUIDRef); err != nil {
				log.Warnf("%s: %s", archive, err)
				err = nil
				continue
			}
			for {
				if archiveWait, err = f.OnlineAPI.GetArchive(safe.UUIDRef, uuidArchive, false); err != nil {
					err = nil
					continue
				}
				if len(archiveWait.Jobs) == 0 {
					fmt.Println(uuidArchive)
					break
				}
				if archiveWait.Jobs[0].Status != "ready" && archiveWait.Jobs[0].Status != "doing" {
					log.Warnf("%s: Wrong job status %s", archive, archiveWait.Jobs[0].Status)
				}
			}
		}
	}
	return
}
