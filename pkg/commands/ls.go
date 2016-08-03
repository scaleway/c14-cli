package commands

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/QuentinPerez/c14-cli/pkg/api"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

type ls struct {
	Base
	lsFlags
}

type lsFlags struct {
	flQuiet    bool
	flPlatform bool
	flAll      bool
}

// Ls returns a new command "ls"
func Ls() Command {
	ret := &ls{}
	ret.Init(Config{
		UsageLine:   "ls [OPTIONS] [ARCHIVE]*",
		Description: "list the archives",
		Help:        "list the archives.",
		Examples: `
        $ c14 ls
        $ c14 ls 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	ret.Flags.BoolVar(&ret.flQuiet, []string{"q", "-quiet"}, false, "Only display UUIDs")
	ret.Flags.BoolVar(&ret.flPlatform, []string{"p", "-platform"}, false, "Show the platforms")
	ret.Flags.BoolVar(&ret.flAll, []string{"a", "-all"}, false, "Show all information on archives")
	return ret
}

func (l *ls) GetName() string {
	return "ls"
}

func (l *ls) Run(args []string) (err error) {
	if err = l.InitAPI(); err != nil {
		return
	}
	if l.flPlatform {
		var (
			val []api.OnlineGetPlatform
		)
		if len(args) == 0 {
			if val, err = l.OnlineAPI.GetPlatforms(); err != nil {
				return
			}
		} else {
			val = make([]api.OnlineGetPlatform, len(args))

			for i, len := 0, len(args); i < len; i++ {
				if val[i], err = l.OnlineAPI.GetPlatform(args[i]); err != nil {
					return
				}
			}
		}
		l.displayPlatforms(val)
	} else {
		l.OnlineAPI.CleanUpCache()
		l.OnlineAPI.FetchRessources()

		var (
			safes []api.OnlineGetSafe
		)

		if len(args) == 0 {
			if safes, err = l.OnlineAPI.GetSafes(true); err != nil {
				return
			}
			l.displayArchives(safes)
		} else {
			err = errors.Errorf("Not implemented yet")
			// safes = make([]api.OnlineGetSafe, len(args))
			//
			// for i, len := 0, len(args); i < len; i++ {
			// 	if safes[i], err = l.OnlineAPI.GetSafe(args[i]); err != nil {
			// 		return
			// 	}
			// }
		}
	}
	return
}

func (l *ls) displayArchives(val []api.OnlineGetSafe) {
	var (
		archives []api.OnlineGetArchive
		archive  api.OnlineGetArchive
		err      error
		w        *tabwriter.Writer
	)

	w = tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !l.flQuiet {
		if l.flAll {
			fmt.Fprintf(w, "NAME\tSTATUS\tUUID\tPARITY\tUUID SAFE\tCREATION DATE\tSIZE\tDESCRIPTION\n")
			wait := sync.WaitGroup{}

			for i := range val {
				archives, err = l.OnlineAPI.GetArchives(val[i].UUIDRef, true)
				if err == nil {
					for j := range archives {
						wait.Add(1)
						go func(uuidSafe, uuidArchive string, w *sync.WaitGroup) {
							l.OnlineAPI.GetArchive(uuidSafe, uuidArchive, false)
							w.Done()
						}(val[i].UUIDRef, archives[j].UUIDRef, &wait)
					}
				}
			}
			wait.Wait()
		} else {
			fmt.Fprintf(w, "NAME\tSTATUS\tUUID\n")
		}
	}
	for i := range val {
		archives, err = l.OnlineAPI.GetArchives(val[i].UUIDRef, true)
		if err == nil {
			sort.Sort(api.OnlineGetArchives(archives))
			for j := range archives {
				if l.flQuiet {
					if l.flAll {
						fmt.Fprintf(w, "%s %s\n", archives[j].UUIDRef, val[i].UUIDRef)
					} else {
						fmt.Fprintf(w, "%s\n", archives[j].UUIDRef)
					}
				} else {
					if l.flAll {
						if archive, err = l.OnlineAPI.GetArchive(val[i].UUIDRef, archives[j].UUIDRef, true); err != nil {
							return
						}
						t, _ := time.Parse(time.RFC3339, archive.CreationDate)
						humanSize := "Unavailable"
						if archive.Size != "" {
							size, _ := strconv.Atoi(archive.Size)
							humanSize = humanize.Bytes(uint64(size))
						}
						fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", archive.Name, archive.Status, archive.UUIDRef, archive.Parity, val[i].UUIDRef, t.Format(time.Stamp), humanSize, archive.Description)
					} else {
						fmt.Fprintf(w, "%s\t%s\t%s\n", archives[j].Name, archives[j].Status, archives[j].UUIDRef)
					}
				}
			}
		}
	}
}

func (l *ls) displayPlatforms(val []api.OnlineGetPlatform) {
	for i := range val {
		if l.flQuiet {
			fmt.Println(val[i].ID)
		} else {
			fmt.Println(val[i])
		}
	}
}
