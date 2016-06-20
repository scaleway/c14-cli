package api

import (
	"fmt"

	"github.com/juju/errors"
)

// GetSafes returns a list of safe
func (o *OnlineAPI) GetSafes() (safes []OnlineGetSafe, err error) {
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe", APIUrl), []int{200}, &safes); err != nil {
		err = errors.Annotate(err, "GetSafes")
	}
	return
}

// GetSafe returns a safe
func (o *OnlineAPI) GetSafe(uuid string) (safe OnlineGetSafe, err error) {
	// TODO: enable to use the name instead of only the UUID
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe/%s", APIUrl, uuid), []int{200}, &safe); err != nil {
		err = errors.Annotate(err, "GetSafe")
	}
	return
}
