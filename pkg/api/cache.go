package api

import (
	"errors"
	"sync"

	"github.com/apex/log"
	"github.com/scaleway/c14-cli/pkg/utils/configstore"
)

type cacheSafe struct {
	Safe    OnlineGetSafe
	Archive map[string]OnlineGetArchive
}

type cache struct {
	safes map[string]cacheSafe
	sync.RWMutex
}

func CleanUp() (c *cache) {
	c = &cache{safes: make(map[string]cacheSafe)}
	return
}

func NewCache() (c *cache) {
	c = CleanUp()
	_ = configStore.GetCache(&c.safes)
	return
}

func (c *cache) Save() {
	_ = configStore.SaveCache(c.safes)
}

func (c *cache) GetSafe(uuid string) (safe OnlineGetSafe, ok bool) {
	var (
		safeCache cacheSafe
	)

	c.RLock()
	if safeCache, ok = c.safes[uuid]; ok {
		safe = safeCache.Safe
	}
	c.RUnlock()
	return
}

func (c *cache) CopySafes() (safes []OnlineGetSafe, err error) {
	c.RLock()
	if length := len(c.safes); length == 0 {
		err = errors.New("No cache")
	} else {
		i := 0
		safes = make([]OnlineGetSafe, length)
		for _, val := range c.safes {
			safes[i] = val.Safe
			i++
		}
	}
	c.RUnlock()
	return
}

func (c *cache) CopyArchives(uuidSafe string) (archives []OnlineGetArchive, err error) {
	c.RLock()
	// force panic if uuidSafe do not exist
	mapArchives := c.safes[uuidSafe].Archive
	if length := len(mapArchives); length == 0 {
		err = errors.New("No cache")
	} else {
		i := 0
		archives = make([]OnlineGetArchive, length)
		for _, val := range mapArchives {
			archives[i] = val
			i++
		}
	}
	c.RUnlock()
	return
}

func (c *cache) InsertSafe(uuid string, safe OnlineGetSafe) {
	c.Lock()
	if _, found := c.safes[uuid]; !found {
		c.safes[uuid] = cacheSafe{
			Safe:    safe,
			Archive: make(map[string]OnlineGetArchive),
		}
		c.Save()
	}
	c.Unlock()
}

func (c *cache) GetArchive(uuidSafe, uuidArchive string) (archive OnlineGetArchive, ok bool) {
	c.RLock()
	// force panic if uuidSafe do not exist
	archive, ok = c.safes[uuidSafe].Archive[uuidArchive]
	c.RUnlock()
	return
}

func (c *cache) InsertArchive(uuidSafe, uuidArchive string, archive OnlineGetArchive) {
	newArchive := new(OnlineGetArchive)
	*newArchive = archive

	c.Lock()
	// force panic if uuidSafe do not exist
	log.Debugf("InsertArchive %s", uuidArchive)
	c.safes[uuidSafe].Archive[uuidArchive] = archive
	c.Save()
	c.Unlock()
}
