package commands

import "github.com/scaleway/c14-cli/pkg/api"

type rename struct {
	Base
	renameFlags
}

type renameFlags struct {
}

// Rename returns a new command "rename"
func Rename() Command {
	ret := &rename{}
	ret.Init(Config{
		UsageLine:   "rename ARCHIVE new_name",
		Description: "Rename an archive",
		Help:        "Rename an archive.",
		Examples: `
        $ c14 rename 83b93179-32e0-11e6-be10-10604b9b0ad9 new_name
        $ c14 rename old_name new_name`,
	})
	return ret
}

func (r *rename) GetName() string {
	return "rename"
}

func (r *rename) Run(args []string) (err error) {
	if len(args) != 2 {
		r.PrintUsage()
		return
	}

	if err = r.InitAPI(); err != nil {
		return
	}
	var (
		safe        api.OnlineGetSafe
		uuidArchive string
	)

	archive := args[0]
	if safe, uuidArchive, err = r.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		if safe, uuidArchive, err = r.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
			return
		}
	}
	err = r.OnlineAPI.PatchArchive(safe.UUIDRef, uuidArchive, api.OnlinePatchArchive{
		Name: args[1],
	})
	return
}
