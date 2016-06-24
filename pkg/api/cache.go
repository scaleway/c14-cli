package api

import (
	"errors"
	"fmt"
	"sync"

	"github.com/apex/log"
)

type cacheArchive struct {
	archive *OnlineGetArchive
	bucket  *OnlineGetBucket
}

type cacheSafe struct {
	safe    OnlineGetSafe
	archive map[string]cacheArchive
}

type cache struct {
	safes map[string]cacheSafe
	sync.RWMutex
}

func NewCache() (c *cache) {
	c = &cache{
		safes: make(map[string]cacheSafe),
	}
	return
}

func (c *cache) GetSafe(uuid string) (safe OnlineGetSafe, ok bool) {
	var (
		safeCache cacheSafe
	)

	c.RLock()
	if safeCache, ok = c.safes[uuid]; ok {
		safe = safeCache.safe
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
			safes[i] = val.safe
			i++
		}
	}
	c.RUnlock()
	return
}

func (c *cache) CopyArchives(uuidSafe string) (archives []OnlineGetArchive, err error) {
	c.RLock()
	// force panic if uuidSafe do not exist
	mapArchives := c.safes[uuidSafe].archive
	if length := len(mapArchives); length == 0 {
		err = errors.New("No cache")
	} else {
		i := 0
		archives = make([]OnlineGetArchive, length)
		for _, val := range mapArchives {
			archives[i] = *val.archive
			i++
		}
	}
	c.RUnlock()
	return
}

func (c *cache) InsertSafe(uuid string, safe OnlineGetSafe) {
	c.Lock()
	c.safes[uuid] = cacheSafe{
		safe:    safe,
		archive: make(map[string]cacheArchive),
	}
	c.Unlock()
}

func (c *cache) GetArchive(uuidSafe, uuidArchive string) (archive OnlineGetArchive, ok bool) {
	var (
		archiveCache cacheArchive
	)

	c.RLock()
	// force panic if uuidSafe do not exist
	if archiveCache, ok = c.safes[uuidSafe].archive[uuidArchive]; ok && archiveCache.archive != nil {
		archive = *archiveCache.archive
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
	val := c.safes[uuidSafe].archive[uuidArchive]
	log.Debugf("InsertArchive %s", uuidArchive)
	c.safes[uuidSafe].archive[uuidArchive] = cacheArchive{
		archive: newArchive,
		bucket:  val.bucket,
	}
	c.Unlock()
}

func (c *cache) GetBucket(uuidSafe, uuidArchive string) (bucket OnlineGetBucket, ok bool) {
	var (
		archiveCache cacheArchive
	)

	c.RLock()
	// force panic if uuidSafe do not exist
	if archiveCache, ok = c.safes[uuidSafe].archive[uuidArchive]; ok && archiveCache.bucket != nil {
		bucket = *archiveCache.bucket
	} else {
		ok = false
	}
	c.RUnlock()
	return
}

func (c *cache) InsertBucket(uuidSafe, uuidArchive string, bucket OnlineGetBucket) {
	newBucket := new(OnlineGetBucket)
	*newBucket = bucket

	c.Lock()
	// force panic if uuidSafe do not exist
	val := c.safes[uuidSafe].archive[uuidArchive]
	log.Debugf("InsertBucket %s:%s", uuidSafe, uuidArchive)
	c.safes[uuidSafe].archive[uuidArchive] = cacheArchive{
		archive: val.archive,
		bucket:  newBucket,
	}
	c.Unlock()
}

func (c *cache) DisplayCache() {
	c.RLock()
	for key, val := range c.safes {
		fmt.Println(key)
		for key2, val2 := range val.archive {
			fmt.Println("    ", key2, val2.archive, val2.bucket)
		}
	}
	c.RUnlock()
}
