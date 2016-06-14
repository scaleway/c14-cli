package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/QuentinPerez/go-encodeUrl"
	"github.com/juju/errors"
)

var (
	deviceURL      = "https://console.online.net/oauth/v2/device/code"
	usercodeURL    = "https://console.online.net/oauth/v2/device/usercode"
	credentialsURL = "https://console.online.net/oauth/v2/device/credentials"
	tokenURL       = "https://console.online.net/oauth/v2/token"
)

type Authentication struct {
	DeviceCode     string `json:"device_code"`
	UserCode       string `json:"user_code"`
	VerficationURL string `json:"verification_url"`
}

var auth Authentication

func newAuthentication(clientID string) (err error) {
	var (
		requestCredentials struct {
			ClientID       string `url:"client_id,ifStringIsNotEmpty"`
			NewCredentials string `url:"new_credentials,ifStringIsNotEmpty"`
		}
		resp *http.Response
	)

	requestCredentials.ClientID = clientID
	requestCredentials.NewCredentials = "yes"
	values, errs := encurl.Translate(requestCredentials)
	if errs != nil {
		err = errors.Trace(errs[0])
		return
	}
	resp, err = http.DefaultClient.Get(fmt.Sprintf("%v?%v", deviceURL, values.Encode()))
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != 200 {
		err = errors.Errorf("Invalid ClienID ?")
		return
	}
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&auth)
	return
}

func newCredentials(clientID string) (err error) {
	return
}

// GetVerificationURL returns an url for the verification
func GetVerificationURL(clientID string) (copy Authentication, err error) {
	if auth.DeviceCode == "" {
		if err = newAuthentication(clientID); err != nil {
			return
		}
	}
	copy = auth
	return
}

func GetToken(clientID string) (err error) {
	return
}
