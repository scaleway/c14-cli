package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/QuentinPerez/c14-cli/pkg/api"
)

type list struct {
	Base
	listFlags
}

type listFlags struct {
	flQuiet    bool
	flPlatform bool
	flAll      bool
}

// List returns a new command "list"
func List() Command {
	ret := &list{}
	ret.Init(Config{
		UsageLine:   "list [OPTIONS]",
		Description: "List archives",
		Help:        "List archives. By default, only archive which are not deleted.",
		Examples: `
        $ c14 list
        $ c14 list 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	ret.Flags.BoolVar(&ret.flQuiet, []string{"q", "-quiet"}, false, "Only display UUIDs")
	ret.Flags.BoolVar(&ret.flPlatform, []string{"p", "-platform"}, false, "Show the platforms")
	ret.Flags.BoolVar(&ret.flAll, []string{"a", "-all"}, false, "Show all archives (default shows !deleted) ")
	return ret
}

func (l *list) CheckFlags() (err error) {
	return
}

func (l *list) GetName() string {
	return "list"
}

func (l *list) Run(args []string) (err error) {
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
		l.OnlineAPI.FetchRessources(true, true)

		var (
			safes []api.OnlineGetSafe
		)

		if len(args) == 0 {
			if safes, err = l.OnlineAPI.GetSafes(true); err != nil {
				return
			}
		} else {
			safes = make([]api.OnlineGetSafe, len(args))

			for i, len := 0, len(args); i < len; i++ {
				if safes[i], err = l.OnlineAPI.GetSafe(args[i]); err != nil {
					return
				}
			}
		}
		l.displayArchives(safes)
	}
	return
}

func (l *list) displayArchives(val []api.OnlineGetSafe) {
	var (
		archives []api.OnlineGetArchive
		err      error
		w        *tabwriter.Writer
	)

	w = tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	defer w.Flush()
	if !l.flQuiet {
		fmt.Fprintf(w, "NAME\tSTATUS\tUUID\nna")
	}
	for i := range val {
		archives, err = l.OnlineAPI.GetArchives(val[i].UUIDRef, true)
		if err == nil {
			for j := range archives {
				if l.flQuiet {
					fmt.Fprintf(w, "%s\n", archives[j].UUIDRef)
				} else {
					if l.flAll {
						fmt.Fprintf(w, "%s\t%s\t%s\n", archives[j].Name, archives[j].Status, archives[j].UUIDRef)
					} else if archives[j].Status != "deleted" {
						fmt.Fprintf(w, "%s\t%s\t%s\n", archives[j].Name, archives[j].Status, archives[j].UUIDRef)
					}
				}
			}
		}
	}
}

func (l *list) displayPlatforms(val []api.OnlineGetPlatform) {
	for i := range val {
		if l.flQuiet {
			fmt.Println(val[i].ID)
		} else {
			fmt.Println(val[i])
		}
	}
}
