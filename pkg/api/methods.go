package api

import (
	"fmt"

	"github.com/juju/errors"
)

/*
 * Get Functions
 */

// GetSafes returns a list of safe
func (o *OnlineAPI) GetSafes(useCache bool) (safes []OnlineGetSafe, err error) {
	if safes, err = o.cache.CopySafes(); err == nil {
		return
	}
	if !useCache {
		if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe", APIUrl), &safes); err != nil {
			err = errors.Annotate(err, "GetSafes")
			return
		}
		for _, safe := range safes {
			if _, ok := o.cache.GetSafe(safe.UUIDRef); !ok {
				o.cache.InsertSafe(safe.UUIDRef, safe)
			}
		}
	}
	return
}

// GetSafe returns a safe
func (o *OnlineAPI) GetSafe(uuid string) (safe OnlineGetSafe, err error) {
	// TODO: enable to use the name instead of only the UUID
	var (
		ok bool
	)

	if safe, ok = o.cache.GetSafe(uuid); ok {
		return
	}
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe/%s", APIUrl, uuid), &safe); err != nil {
		err = errors.Annotate(err, "GetSafe")
		return
	}
	o.cache.InsertSafe(safe.UUIDRef, safe)
	return
}

// GetPlatforms returns a list of platform
func (o *OnlineAPI) GetPlatforms() (platform []OnlineGetPlatform, err error) {
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/platform", APIUrl), &platform); err != nil {
		err = errors.Annotate(err, "GetPlatforms")
	}
	return
}

// GetPlatform returns a platform
func (o *OnlineAPI) GetPlatform(uuid string) (platform OnlineGetPlatform, err error) {
	// TODO: enable to use the name instead of only the UUID
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/platform/%s", APIUrl, uuid), &platform); err != nil {
		err = errors.Annotate(err, "GetPlatform")
	}
	return
}

func (o *OnlineAPI) GetSSHKeys() (keys []OnlineGetSSHKey, err error) {
	if err = o.getWrapper(fmt.Sprintf("%s/user/key/ssh", APIUrl), &keys); err != nil {
		err = errors.Annotate(err, "GetSSHKeys")
	}
	return
}

func (o *OnlineAPI) GetSSHKey(uuid string) (key OnlineGetSSHKey, err error) {
	if err = o.getWrapper(fmt.Sprintf("%s/user/key/ssh/%s", APIUrl, uuid), &key); err != nil {
		err = errors.Annotate(err, "GetSSHKey")
	}
	return
}

func (o *OnlineAPI) GetArchives(uuidSafe string, useCache bool) (archives []OnlineGetArchive, err error) {
	if archives, err = o.cache.CopyArchives(uuidSafe); err == nil {
		return
	}
	if !useCache {
		if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe/%s/archive", APIUrl, uuidSafe), &archives); err != nil {
			err = errors.Annotate(err, "GetArchives")
			return
		}
		for _, archive := range archives {
			if _, ok := o.cache.GetArchive(uuidSafe, archive.UUIDRef); !ok {
				o.cache.InsertArchive(uuidSafe, archive.UUIDRef, archive)
			}
		}
	}
	return
}

func (o *OnlineAPI) GetArchive(uuidSafe, uuidArchive string) (archive OnlineGetArchive, err error) {
	var (
		ok bool
	)

	if archive, ok = o.cache.GetArchive(uuidSafe, uuidArchive); ok {
		return
	}
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe/%s/archive/%s", APIUrl, uuidSafe, uuidArchive), &archive); err != nil {
		err = errors.Annotate(err, "GetArchive")
		return
	}
	o.cache.InsertArchive(uuidSafe, uuidArchive, archive)
	return
}

func (o *OnlineAPI) GetBucket(uuidSafe, uuidArchive string) (bucket OnlineGetBucket, err error) {
	var (
		ok bool
	)

	if bucket, ok = o.cache.GetBucket(uuidSafe, uuidArchive); ok {
		return
	}
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe/%s/archive/%s/bucket", APIUrl, uuidSafe, uuidArchive), &bucket); err != nil {
		err = errors.Annotate(err, "GetBucket")
		return
	}
	o.cache.InsertBucket(uuidSafe, uuidArchive, bucket)
	return
}

/*
 * Create Functions
 */

func (o *OnlineAPI) CreateSafe(name, desc string) (uuid string, err error) {
	var (
		buff []byte
	)

	if buff, err = o.postWrapper(fmt.Sprintf("%s/storage/c14/safe", APIUrl), OnlinePostSafe{
		Name:        name,
		Description: desc,
	}); err != nil {
		err = errors.Annotate(err, "CreateSafe")
		return
	}
	uuid = string(buff)
	_, err = o.GetSafe(uuid)
	return
}

type ConfigCreateArchive struct {
	UUIDSafe  string
	Name      string
	Desc      string
	Parity    string
	Protocols []string
	SSHKeys   []string
	Platforms []string
	Days      int
}

func (o *OnlineAPI) CreateArchive(config ConfigCreateArchive) (uuid string, err error) {
	var (
		buff []byte
	)

	if buff, err = o.postWrapper(fmt.Sprintf("%s/storage/c14/safe/%s/archive", APIUrl, config.UUIDSafe), OnlinePostArchive{
		Name:        config.Name,
		Description: config.Desc,
		Protocols:   config.Protocols,
		SSHKeys:     config.SSHKeys,
		Platforms:   config.Platforms,
		Days:        config.Days,
	}); err != nil {
		err = errors.Annotate(err, "CreateArchive")
		return
	}
	uuid = string(buff)
	_, err = o.GetArchive(config.UUIDSafe, uuid)
	return
}

/*
 * Delete Functions
 */

func (o *OnlineAPI) DeleteSafe(uuid string) (err error) {
	// TODO: remove from cache
	if err = o.deleteWrapper(fmt.Sprintf("%s/storage/c14/safe/%s", APIUrl, uuid)); err != nil {
		err = errors.Annotate(err, "DeleteSafe")
	}
	return
}

func (o *OnlineAPI) DeleteArchive(uuidSafe, uuidArchive string) (err error) {
	// TODO: remove from cache
	if err = o.deleteWrapper(fmt.Sprintf("%s/storage/c14/safe/%s/archive/%s", APIUrl, uuidSafe, uuidArchive)); err != nil {
		err = errors.Annotate(err, "DeleteArchive")
	}
	return
}
