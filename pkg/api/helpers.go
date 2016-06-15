package api

import (
	"fmt"

	"github.com/juju/errors"
)

func (o *OnlineAPI) GetSafes() (safes []OnlineGetSafe, err error) {
	if err = o.getWrapper(fmt.Sprintf("%s/storage/c14/safe", APIUrl), []int{200}, &safes); err != nil {
		err = errors.Annotate(err, "GetSafes")
	}
	return
}
