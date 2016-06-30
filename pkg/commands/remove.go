package commands

import (
	"sync"

	"github.com/QuentinPerez/c14-cli/pkg/api"
	"github.com/apex/log"
)

type remove struct {
	Base
	removeFlags
}

type removeFlags struct {
}

// Remove returns a new command "remove"
func Remove() Command {
	ret := &remove{}
	ret.Init(Config{
		UsageLine:   "remove [ARCHIVE]+",
		Description: "",
		Help:        "",
		Examples: `
        $ c14 remove 83b93179-32e0-11e6-be10-10604b9b0ad9 2d752399-429f-447f-85cd-c6104dfed5db`,
	})
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

	for _, uuid := range args {
		wait.Add(1)
		go r.remove(&wait, uuid)
	}
	wait.Wait()
	return
}

func (r *remove) remove(wait *sync.WaitGroup, uuid string) (err error) {
	defer wait.Done()

	var (
		safe        api.OnlineGetSafe
		uuidArchive string
	)

	if safe, uuidArchive, err = r.OnlineAPI.FindSafeUUIDFromArchive(uuid, true); err != nil {
		log.Warnf("%s: %s", uuid, err)
		return
	}
	if err = r.OnlineAPI.DeleteArchive(safe.UUIDRef, uuidArchive); err != nil {
		log.Warnf("%s: %s", uuidArchive, err)
	}
	return
}
