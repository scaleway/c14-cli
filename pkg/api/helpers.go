package api

import (
	"fmt"

	"github.com/juju/errors"
)

/*
 * Get Functions
 */

// GetSafes returns a list of safe
func (o *OnlineAPI) GetSafes() (safes []OnlineGetSafe, err error) {
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe", APIUrl), &safes); err != nil {
		err = errors.Annotate(err, "GetSafes")
	}
	return
}

// GetSafe returns a safe
func (o *OnlineAPI) GetSafe(uuid string) (safe OnlineGetSafe, err error) {
	// TODO: enable to use the name instead of only the UUID
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe/%s", APIUrl, uuid), &safe); err != nil {
		err = errors.Annotate(err, "GetSafe")
	}
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
	fmt.Println(string(buff))
	return
}
