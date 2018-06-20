package commands

import (
	"sync"

	"github.com/apex/log"
	"github.com/scaleway/c14-cli/pkg/api"
)

type remove struct {
	Base
	removeFlags
}

type removeFlags struct {
	flForce bool
}

// Remove returns a new command "remove"
func Remove() Command {
	ret := &remove{}
	ret.Init(Config{
		UsageLine:   "remove [ARCHIVE]+",
		Description: "Remove an archive",
		Help:        "Remove an archive.",
		Examples: `
        $ c14 remove 83b93179-32e0-11e6-be10-10604b9b0ad9 2d752399-429f-447f-85cd-c6104dfed5db`,
	})
	ret.Flags.BoolVar(&ret.flForce, []string{"f", "-force"}, false, "Remove the archive and the safe")
	return ret
}

func (r *remove) GetName() string {
	return "remove"
}

func (r *remove) Run(args []string) (err error) {
	if len(args) == 0 {
		r.PrintUsage()
		return
	}

	if err = r.InitAPI(); err != nil {
		return
	}

	var wait sync.WaitGroup

	for _, archive := range args {
		wait.Add(1)
		go r.remove(&wait, archive)
	}
	wait.Wait()
	return
}

func (r *remove) remove(wait *sync.WaitGroup, archive string) (err error) {
	defer wait.Done()

	var (
		safe        api.OnlineGetSafe
		uuidArchive string
	)

	if safe, uuidArchive, err = r.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		log.Warnf("%s: %s", archive, err)
		return
	}
	if err = r.OnlineAPI.DeleteArchive(safe.UUIDRef, uuidArchive); err != nil {
		log.Warnf("%s: %s", uuidArchive, err)
		return
	}
	if r.flForce {
		if err = r.OnlineAPI.DeleteSafe(safe.UUIDRef); err != nil {
			log.Warnf("%s: %s", safe.UUIDRef, err)
			return
		}
	}
	return
}
