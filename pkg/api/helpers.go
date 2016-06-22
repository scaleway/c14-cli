package api

import (
	"time"

	"github.com/juju/errors"
)

type ConfigCreateSSHBucketFromScratch struct {
	SafeName    string
	ArchiveName string
	Desc        string
	UUIDSSHKeys []string
	Platforms   []string
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
		Days:      7,
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
