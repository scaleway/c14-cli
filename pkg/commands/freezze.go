package commands

import (
	"os"
	"time"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/apex/log"
	"github.com/scaleway/c14-cli/pkg/api"
	"github.com/scaleway/c14-cli/pkg/utils/pgbar"
)

type freeze struct {
	Base
	freezeFlags
}

type freezeFlags struct {
	flQuiet  bool
	flNoWait bool
}

// Freeze returns a new command "freeze"
func Freeze() Command {
	ret := &freeze{}
	ret.Init(Config{
		UsageLine:   "freeze [OPTIONS] [ARCHIVE]+",
		Description: "Lock an archive",
		Help:        "Lock an archive, your archive will be stored in highly secure Online data centers and will stay available On Demand (0.01€/GB).",
		Examples: `
        $ c14 freeze 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	ret.Flags.BoolVar(&ret.flQuiet, []string{"q", "-quiet"}, false, "")
	ret.Flags.BoolVar(&ret.flNoWait, []string{"-nowait"}, false, "")
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

	var (
		safe                    api.OnlineGetSafe
		archiveWait             api.OnlineGetArchive
		uuidArchive, newArchive string
	)

	for _, archive := range args {
		if safe, uuidArchive, err = f.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
			if safe, uuidArchive, err = f.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
				return
			}
		}
		if newArchive, err = f.OnlineAPI.PostArchive(safe.UUIDRef, uuidArchive); err != nil {
			log.Warnf("%s: %s", archive, err)
			err = nil
			continue
		}
		if newArchive != "" {
			uuidArchive = newArchive
		}
		if !f.flNoWait {
			var bar *pb.ProgressBar

			if !f.flQuiet {
				bar = pgbar.NewProgressBar(uuidArchive)
				bar.Start()
			}
			lastLength := 6
			for {
				if archiveWait, err = f.OnlineAPI.GetArchive(safe.UUIDRef, uuidArchive, false); err != nil {
					log.Warnf("%s: %s", args, err)
					err = nil
					continue
				}
				if lastLength != len(archiveWait.Jobs) {
					lastLength = len(archiveWait.Jobs)
					if !f.flQuiet {
						bar.Add(20)
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
