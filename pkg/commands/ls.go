package commands

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/scaleway/c14-cli/pkg/api"
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
		Description: "Displays the archives",
		Help:        "Displays the archives, by default only the NAME, STATUS, UUID.",
		Examples: `
        $ c14 ls
        $ c14 ls -a`,
	})
	ret.Flags.BoolVar(&ret.flQuiet, []string{"q", "-quiet"}, false, "Only display UUIDs")
	ret.Flags.BoolVar(&ret.flPlatform, []string{"p", "-platform"}, false, "Show the platforms")
	ret.Flags.BoolVar(&ret.flAll, []string{"a", "-all"}, false, "Show all information on archives (size,parity,creationDate,description)")
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
		var (
			archives api.OnlineGetArchives
		)

		if len(args) == 0 {
			if archives, err = l.OnlineAPI.GetAllArchives(); err != nil {
				return
			}
			l.displayArchives(archives)
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

func (l *ls) displayArchives(archives []api.OnlineGetArchive) {
	var (
		w *tabwriter.Writer
	)

	w = tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !l.flQuiet {
		if l.flAll {
			fmt.Fprintf(w, "NAME\tSTATUS\tUUID\tPARITY\tUUID SAFE\tCREATION DATE\tSIZE\tDESCRIPTION\n")
		} else {
			fmt.Fprintf(w, "NAME\tSTATUS\tUUID\n")
		}
	}
	sort.Sort(api.OnlineGetArchives(archives))
	for _, archive := range archives {
		if l.flQuiet {
			if l.flAll {
				fmt.Fprintf(w, "%s %s\n", archive.UUIDRef, archive.Safe.UUIDRef)
			} else {
				fmt.Fprintf(w, "%s\n", archive.UUIDRef)
			}
		} else {
			if l.flAll {
				t, _ := time.Parse(time.RFC3339, archive.CreationDate)
				humanSize := "Unavailable"
				if archive.Size != "" {
					size, _ := strconv.Atoi(archive.Size)
					humanSize = humanize.Bytes(uint64(size))
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", archive.Name, archive.Status, archive.UUIDRef, archive.Parity, archive.Safe.UUIDRef, t.Format(time.Stamp), humanSize, archive.Description)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\n", archive.Name, archive.Status, archive.UUIDRef)
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
