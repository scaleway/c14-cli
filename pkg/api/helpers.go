package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/juju/errors"
)

type ConfigCreateSSHBucketFromScratch struct {
	SafeName    string
	ArchiveName string
	Desc        string
	UUIDSSHKeys []string
	Platforms   []string
	Days        int
	Quiet       bool
	Parity      string
	LargeBucket bool
}

// CreateSSHBucketFromScratch creates a safe, an archive and returns the bucket available over SSH
func (o *OnlineAPI) CreateSSHBucketFromScratch(c ConfigCreateSSHBucketFromScratch) (uuidSafe, uuidArchive string, bucket OnlineGetBucket, err error) {

	var (
		safes []OnlineGetSafe
		found = false
	)

	if safes, err = o.GetSafes(true); err != nil {
		err = errors.Annotate(err, "CreateSSHBucketFromScratch:GetSafes")
		return
	}

	for idxSafe := range safes {
		if safes[idxSafe].Name == c.SafeName {
			uuidSafe = safes[idxSafe].UUIDRef
			found = true
		}
	}
	if !found {
		if uuidSafe, err = o.CreateSafe(c.SafeName, ""); err != nil {
			err = errors.Annotate(err, "CreateSSHBucketFromScratch:CreateSafe")
			return
		}
	}

	if uuidArchive, err = o.CreateArchive(ConfigCreateArchive{
		UUIDSafe:    uuidSafe,
		Name:        c.ArchiveName,
		Desc:        c.Desc,
		Protocols:   []string{"SSH"},
		Platforms:   c.Platforms,
		SSHKeys:     c.UUIDSSHKeys,
		Days:        c.Days,
		Parity:      c.Parity,
		LargeBucket: c.LargeBucket,
	}); err != nil {
		err = errors.Annotate(err, "CreateSSHBucketFromScratch:CreateArchive")
		return
	}
	if !c.Quiet {
		defer fmt.Printf("\r \r")
	}
	errChan := make(chan error)
	go func() {
		var errGoRoutine error

		for i := 0; i < 120; i++ {
			errGoRoutine = nil
			if bucket, errGoRoutine = o.GetBucket(uuidSafe, uuidArchive); errGoRoutine == nil {
				break
			}
			if onlineError, ok := errors.Cause(errGoRoutine).(*OnlineError); ok && onlineError.StatusCode != 404 {
				break
			}
			time.Sleep(1 * time.Second)
		}
		errChan <- errGoRoutine
	}()
	loop := 0
	for {
		select {
		case err = <-errChan:
			goto OUT
		case <-time.After(100 * time.Millisecond):
			if !c.Quiet {
				fmt.Printf("\r%c\r", "-\\|/"[loop%4])
				loop++
				if loop == 5 {
					loop = 0
				}
			}
		}
	}
OUT:
	if err != nil {
		err = errors.Annotate(err, "CreateSSHBucketFromScratch:GetBucket")
	}
	return
}

// FetchRessources get the ressources to fill the cache
func (o *OnlineAPI) FetchRessources() (err error) {
	var (
		wgSafe sync.WaitGroup
		safes  []OnlineGetSafe
	)

	if safes, err = o.GetSafes(false); err != nil {
		err = errors.Annotate(err, "FetchRessources")
		return
	}
	for indexSafe := range safes {
		wgSafe.Add(1)
		go func(uuidSafe string, wgSafe *sync.WaitGroup) {
			_, _ = o.GetArchives(uuidSafe, false)
			wgSafe.Done()
		}(safes[indexSafe].UUIDRef, &wgSafe)
	}
	wgSafe.Wait()
	return
}

func (o *OnlineAPI) FindSafeUUIDFromArchive(archive string, useCache bool) (safe OnlineGetSafe, uuidArchive string, err error) {
	var (
		safes []OnlineGetSafe
		ret   []struct {
			safe OnlineGetSafe
			uuid string
		}
	)

	if safes, err = o.GetSafes(useCache); err != nil {
		err = errors.Annotate(err, "FindArchiveFromCache:GetSafes")
		return
	}
	for indexSafe := range safes {
		var (
			archives []OnlineGetArchive
		)
		if archives, err = o.GetArchives(safes[indexSafe].UUIDRef, useCache); err == nil {
			for indexArchive := range archives {
				if archive == archives[indexArchive].UUIDRef {
					safe = safes[indexSafe]
					uuidArchive = archives[indexArchive].UUIDRef
					return
				}
				if archive == archives[indexArchive].Name {
					ret = append(ret, struct {
						safe OnlineGetSafe
						uuid string
					}{
						safes[indexSafe],
						archives[indexArchive].UUIDRef,
					})
				}
			}
		}
	}
	switch len(ret) {
	case 0:
		err = errors.Errorf("Archive %s not found", archive)
	case 1:
		safe = ret[0].safe
		uuidArchive = ret[0].uuid
	default:
		err = errors.Errorf("Too many candidate for %s", archive)
	}
	return
}

func (o *OnlineAPI) CleanUpCache() {
	o.cache = CleanUp()
}
