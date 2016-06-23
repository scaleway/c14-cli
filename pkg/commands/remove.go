package commands

import (
	"sync"

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
		UsageLine:   "remove [ARGS]+",
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
	wait.Add(1)
	go func() {
		defer wait.Done()
		var (
			errSafe error
		)
		if _, errSafe = r.OnlineAPI.GetSafe(uuid); errSafe == nil {
			if errSafe = r.OnlineAPI.DeleteSafe(uuid); errSafe != nil {
				log.Errorf("%s: %s", uuid, errSafe)
			}
		} else {
			log.Errorf("%s: %s", uuid, errSafe)
		}
	}()
	return
}
