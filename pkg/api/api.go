package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/apex/log"
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
	verbose   bool
}

// NewC14API returns a new API
func NewC14API(client *http.Client, userAgent string, verbose bool) (api *OnlineAPI) {
	api = &OnlineAPI{
		client:    client,
		userAgent: userAgent,
		verbose:   verbose,
	}
	return
}

func (o *OnlineAPI) GetResponse(uri string) (resp *http.Response, err error) {
	var (
		req *http.Request
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
	if o.verbose {
		dump, _ := httputil.DumpRequest(req, false)
		log.Debugf("%v", string(dump))
	} else {
		log.Debugf("[GET]: %v", uri)
	}
	resp, err = o.client.Do(req)
	return
}

func (o *OnlineAPI) getWrapper(uri string, goodStatusCode []int, export interface{}) (err error) {
	var (
		resp *http.Response
		body []byte
	)

	resp, err = o.GetResponse(uri)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		err = errors.Annotatef(err, "Unable to get %v", uri)
		return
	}

	if body, err = o.handleHTTPError(goodStatusCode, resp); err != nil {
		return
	}
	err = json.Unmarshal(body, export)
	return
}

func (o *OnlineAPI) handleHTTPError(goodStatusCode []int, resp *http.Response) (content []byte, err error) {
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if o.verbose {
		dump, _ := httputil.DumpResponse(resp, false)
		log.Debugf("%v", string(dump))
	} else {
		log.Debugf("[Response]: [%v] %v", resp.StatusCode, string(content))
	}

	if resp.StatusCode >= 500 {
		err = errors.Errorf("[%v] %v", resp.StatusCode, string(content))
		return
	}
	good := false
	for _, code := range goodStatusCode {
		if code == resp.StatusCode {
			good = true
		}
	}
	if !good {
		var why OnlineError

		if err = json.Unmarshal(content, &why); err != nil {
			return
		}
		why.StatusCode = resp.StatusCode
		err = why
	}
	return
}
