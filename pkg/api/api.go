package api

import (
	"bytes"
	"encoding/json"
	"io"
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
	cache     *cache
}

// NewC14API returns a new API
func NewC14API(client *http.Client, userAgent string, verbose bool) (api *OnlineAPI) {
	api = &OnlineAPI{
		client:    client,
		userAgent: userAgent,
		verbose:   verbose,
		cache:     NewCache(),
	}
	return
}

func (o *OnlineAPI) response(method, uri string, content io.Reader) (resp *http.Response, err error) {
	var (
		req *http.Request
	)

	req, err = http.NewRequest(method, uri, content)
	if err != nil {
		err = errors.Annotatef(err, "response %s %s", method, uri)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", o.userAgent)

	// curl, err := http2curl.GetCurlCommand(req)
	// if err != nil {
	// 	return nil, err
	// }
	if o.verbose {
		dump, _ := httputil.DumpRequest(req, true)
		log.Debugf("%v", string(dump))
	} else {
		log.Debugf("[%s]: %v", method, uri)
	}
	resp, err = o.client.Do(req)
	return
}

func (o *OnlineAPI) getWrapper(uri string, export interface{}) (err error) {
	var (
		resp *http.Response
		body []byte
	)

	resp, err = o.response("GET", uri, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		err = errors.Annotatef(err, "Unable to get %s", uri)
		return
	}

	if body, err = o.handleHTTPError([]int{200}, resp); err != nil {
		return
	}
	err = json.Unmarshal(body, export)
	return
}

func (o *OnlineAPI) deleteWrapper(uri string) (err error) {
	var (
		resp *http.Response
	)

	resp, err = o.response("DELETE", uri, nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		err = errors.Annotatef(err, "Unable to delete %s", uri)
		return
	}

	if _, err = o.handleHTTPError([]int{204}, resp); err != nil {
		return
	}
	return
}

func (o *OnlineAPI) postWrapper(uri string, content interface{}, goodStatusCode ...[]int) (body []byte, err error) {
	var (
		resp    *http.Response
		payload = new(bytes.Buffer)
	)

	encoder := json.NewEncoder(payload)
	if content != nil {
		if err = encoder.Encode(content); err != nil {
			return
		}
	}
	resp, err = o.response("POST", uri, payload)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		err = errors.Annotatef(err, "Unable to post %s", uri)
		return
	}
	goodStatus := []int{201}
	if len(goodStatusCode) > 0 {
		goodStatus = goodStatusCode[0]
	}
	body, err = o.handleHTTPError(goodStatus, resp)
	return
}

func (o *OnlineAPI) handleHTTPError(goodStatusCode []int, resp *http.Response) (content []byte, err error) {
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if o.verbose {
		dump, _ := httputil.DumpResponse(resp, true)
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
