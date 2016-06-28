package api

import (
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
}

// CreateSSHBucketFromScratch creates a safe, an archive and returns the bucket available over SSH
func (o *OnlineAPI) CreateSSHBucketFromScratch(c ConfigCreateSSHBucketFromScratch) (uuidSafe, uuidArchive string, bucket OnlineGetBucket, err error) {
	if uuidSafe, err = o.CreateSafe(c.SafeName, ""); err != nil {
		err = errors.Annotate(err, "CreateSSHBucketFromScratch:CreateSafe")
		return
	}
	if uuidArchive, err = o.CreateArchive(ConfigCreateArchive{
		UUIDSafe:  uuidSafe,
		Name:      c.ArchiveName,
		Desc:      c.Desc,
		Protocols: []string{"SSH"},
		Platforms: c.Platforms,
		SSHKeys:   c.UUIDSSHKeys,
		Days:      c.Days,
	}); err != nil {
		o.DeleteSafe(uuidSafe)
		err = errors.Annotate(err, "CreateSSHBucketFromScratch:CreateArchive")
		return
	}
	for i := 0; i < 60; i++ {
		err = nil
		if bucket, err = o.GetBucket(uuidSafe, uuidArchive); err == nil {
			break
		}
		if onlineError, ok := errors.Cause(err).(*OnlineError); ok && onlineError.StatusCode != 404 {
			return
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		o.DeleteArchive(uuidSafe, uuidArchive)
		o.DeleteSafe(uuidSafe)
		err = errors.Annotate(err, "CreateSSHBucketFromScratch:GetBucket")
		return
	}
	return
}

// FetchRessources get the ressources to fill the cache
func (o *OnlineAPI) FetchRessources(archive, bucket bool) (err error) {
	var (
		wgSafe sync.WaitGroup
		safes  []OnlineGetSafe
	)

	if safes, err = o.GetSafes(false); err != nil {
		err = errors.Annotate(err, "FetchRessources")
		return
	}
	if archive {
		for indexSafe := range safes {
			if safes[indexSafe].Status == "deleted" {
				continue
			}
			wgSafe.Add(1)
			go func(uuidSafe string, wgSafe *sync.WaitGroup) {
				var (
					archives   []OnlineGetArchive
					wgArchive  sync.WaitGroup
					errArchive error
				)

				archives, errArchive = o.GetArchives(uuidSafe, false)
				if bucket && errArchive == nil {
					for indexArchive := range archives {
						wgArchive.Add(1)
						go func(uuidSafe, uuidArchive string, wgArchive *sync.WaitGroup) {
							_, _ = o.GetBucket(uuidSafe, uuidArchive)
							wgArchive.Done()
						}(uuidSafe, archives[indexArchive].UUIDRef, &wgArchive)
					}
				}
				wgArchive.Wait()
				wgSafe.Done()
			}(safes[indexSafe].UUIDRef, &wgSafe)
		}
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
		if safes[indexSafe].Status != "deleted" {
			if archives, err = o.GetArchives(safes[indexSafe].UUIDRef, useCache); err == nil {
				for indexArchive := range archives {
					if archive == archives[indexArchive].UUIDRef || archive == archives[indexArchive].Name {
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
	}
	switch len(ret) {
	case 0:
		err = errors.Errorf("Archive %s not found", archive)
	case 1:
		safe = ret[0].safe
		uuidArchive = ret[0].uuid
	default:
		err = errors.Errorf("Multiple candidate for %s", archive)
	}
	return
}
