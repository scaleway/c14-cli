package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// GetVerificationURL returns an url for the verification
func GetVerificationURL(clientID string) (auth Authentication, err error) {
	var (
		requestAuthentication struct {
			ClientID       string `url:"client_id,ifStringIsNotEmpty"`
			NewCredentials string `url:"new_credentials,ifStringIsNotEmpty"`
		}
		resp *http.Response
	)

	requestAuthentication.ClientID = clientID
	requestAuthentication.NewCredentials = "yes"
	values, errs := encurl.Translate(requestAuthentication)
	if errs != nil {
		err = errors.Trace(errs[0])
		return
	}
	resp, err = http.DefaultClient.Get(fmt.Sprintf("%v?%v", deviceURL, values.Encode()))
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != 200 {
		buf, _ := ioutil.ReadAll(resp.Body)
		err = errors.Errorf(string(buf))
		return
	}
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&auth)
	return
}

type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func GetCredentials(clientID, deviceCode string) (c Credentials, err error) {
	var (
		requestCredentials struct {
			ClientID string `url:"client_id,ifStringIsNotEmpty"`
			Code     string `url:"code,ifStringIsNotEmpty"`
		}
		resp *http.Response
	)
	requestCredentials.ClientID = clientID
	requestCredentials.Code = deviceCode
	values, errs := encurl.Translate(requestCredentials)
	if errs != nil {
		err = errors.Trace(errs[0])
		return
	}
	resp, err = http.DefaultClient.Get(fmt.Sprintf("%v?%v", credentialsURL, values.Encode()))
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != 200 {
		buf, _ := ioutil.ReadAll(resp.Body)
		err = errors.Errorf(string(buf))
		return
	}
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&c)
	if err != nil {
		return
	}
	return
}
