package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"

	"github.com/apex/log"
	"github.com/juju/errors"
)

type cacheArchive struct {
	Archive *OnlineGetArchive
}

type cacheSafe struct {
	Safe    OnlineGetSafe
	Archive map[string]cacheArchive
}

type cache struct {
	safes map[string]cacheSafe
	sync.RWMutex
}

func getCachePath() (path string, err error) {
	homeDir := os.Getenv("HOME") // *nix
	if homeDir == "" {           // Windows
		homeDir = os.Getenv("USERPROFILE")
	}
	if homeDir == "" {
		return "", errors.New("user home directory not found")
	}
	path = fmt.Sprintf("%s/.c14-cache", homeDir)
	return
}

func CleanUp() (c *cache) {
	c = &cache{safes: make(map[string]cacheSafe)}
	return
}

func NewCache() (c *cache) {
	var (
		path        string
		fileContent []byte
		err         error
	)
	c = CleanUp()
	if path, err = getCachePath(); err == nil {
		// Don't check permissions on Windows
		if runtime.GOOS != "windows" {
			stat, errStat := os.Stat(path)
			if errStat == nil {
				perm := stat.Mode().Perm()
				if perm&0066 != 0 {
					log.Debugf("Permissions %#o for %v are too open", perm, path)
					return
				}
			} else {
				return
			}
		}
		if fileContent, err = ioutil.ReadFile(path); err != nil {
			return
		}
		json.Unmarshal(fileContent, &c.safes)
	}
	return
}

func (c *cache) Save() {
	var (
		path string
		err  error
		data []byte
	)

	if path, err = getCachePath(); err == nil {
		if data, err = json.Marshal(c.safes); err == nil {
			_ = ioutil.WriteFile(path, data, 0600)
		}
	}
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
			archives[i] = *val.Archive
			i++
		}
	}
	c.RUnlock()
	return
}

func (c *cache) InsertSafe(uuid string, safe OnlineGetSafe) {
	c.Lock()
	c.safes[uuid] = cacheSafe{
		Safe:    safe,
		Archive: make(map[string]cacheArchive),
	}
	c.Save()
	c.Unlock()
}

func (c *cache) GetArchive(uuidSafe, uuidArchive string) (archive OnlineGetArchive, ok bool) {
	var (
		archiveCache cacheArchive
	)

	c.RLock()
	// force panic if uuidSafe do not exist
	if archiveCache, ok = c.safes[uuidSafe].Archive[uuidArchive]; ok && archiveCache.Archive != nil {
		archive = *archiveCache.Archive
	} else {
		ok = false
	}
	c.RUnlock()
	return
}

func (c *cache) InsertArchive(uuidSafe, uuidArchive string, archive OnlineGetArchive) {
	newArchive := new(OnlineGetArchive)
	*newArchive = archive

	c.Lock()
	// force panic if uuidSafe do not exist
	log.Debugf("InsertArchive %s", uuidArchive)
	c.safes[uuidSafe].Archive[uuidArchive] = cacheArchive{
		Archive: newArchive,
	}
	c.Save()
	c.Unlock()
}

func (c *cache) DisplayCache() {
	c.RLock()
	for key, val := range c.safes {
		fmt.Println(key)
		for key2, val2 := range val.Archive {
			fmt.Println("    ", key2, val2.Archive)
		}
	}
	c.RUnlock()
}
