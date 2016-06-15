package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/juju/errors"
)

var (
	// APIUrl represents Online's endpoint
	APIUrl = "https://api.online.net/api/v1"
)

// OnlineAPI is used to communicate with Online API
type OnlineAPI struct {
	client    *http.Client
	userAgent string
}

// NewC14API returns a new API
func NewC14API(client *http.Client, userAgent string) (api *OnlineAPI) {
	api = &OnlineAPI{
		client:    client,
		userAgent: userAgent,
	}
	return
}

func (o *OnlineAPI) GetResponse(apiURL, resource string) (resp *http.Response, err error) {
	var (
		req *http.Request
		uri = fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource)
	)

	req, err = http.NewRequest("GET", uri, nil)
	if err != nil {
		err = errors.Annotatef(err, "NewRequest Get %v", uri)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", o.userAgent)

	// curl, err := http2curl.GetCurlCommand(req)
	// if err != nil {
	// 	return nil, err
	// }
	resp, err = o.client.Do(req)
	return
}
