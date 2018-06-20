package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/QuentinPerez/go-encodeUrl"
	"github.com/juju/errors"
	"github.com/scaleway/c14-cli/pkg/utils/configstore"
)

var (
	deviceURL      = "https://console.online.net/oauth/v2/device/code"
	credentialsURL = "https://console.online.net/oauth/v2/device/credentials"
	tokenURL       = "https://console.online.net/oauth/v2/token"
)

// Authentication is used to exchange the keys with the oauth2 API
type Authentication struct {
	DeviceCode     string `json:"device_code"`
	UserCode       string `json:"user_code"`
	VerficationURL string `json:"verification_url"`
}

// Credentials represents the informations to get a token, there are saved in `~/$CONFIG/c14-cli/c14rc.json`
type Credentials struct {
	ClientID     string `json:"client_id" url:"client_id,ifStringIsNotEmpty"`
	ClientSecret string `json:"client_secret" url:"client_secret,ifStringIsNotEmpty"`
	Code         string `json:"-" url:"code,ifStringIsNotEmpty"`
	GrantType    string `json:"-" url:"grant_type,ifStringIsNotEmpty"`
	AccessToken  string `json:"access_token" url:"-"`
}

func getURL(url string, encode, decode interface{}) (err error) {
	var resp *http.Response

	values, errs := encurl.Translate(encode)
	if errs != nil {
		err = errors.Trace(errs[0])
		return
	}
	resp, err = http.DefaultClient.Get(fmt.Sprintf("%v?%v", url, values.Encode()))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buf, _ := ioutil.ReadAll(resp.Body)
		err = errors.Errorf(string(buf))
		return
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(decode)
	return
}

// GetVerificationURL returns an url for the verification
func GetVerificationURL(clientID string) (auth Authentication, err error) {
	var (
		requestAuthentication struct {
			ClientID       string `url:"client_id,ifStringIsNotEmpty"`
			NewCredentials string `url:"new_credentials,ifStringIsNotEmpty"`
		}
	)

	requestAuthentication.ClientID = clientID
	requestAuthentication.NewCredentials = "yes"
	err = getURL(deviceURL, requestAuthentication, &auth)
	return
}

// GenerateCredentials calls the oauth2 API to get the credentials
func GenerateCredentials(clientID, deviceCode string) (c Credentials, err error) {
	var (
		requestCredentials struct {
			ClientID string `url:"client_id,ifStringIsNotEmpty"`
			Code     string `url:"code,ifStringIsNotEmpty"`
		}
	)
	requestCredentials.ClientID = clientID
	requestCredentials.Code = deviceCode
	if err = getURL(credentialsURL, requestCredentials, &c); err != nil {
		return
	}
	c.Code = deviceCode
	c.GrantType = "http://oauth.net/grant_type/device/1.0"
	err = getURL(tokenURL, c, &c)
	return
}

// Token oauth2.TokenSource implementation
func (c *Credentials) Token() (t *oauth2.Token, err error) {
	t = &oauth2.Token{
		AccessToken: c.AccessToken,
		TokenType:   "Bearer",
	}
	return
}

// Save writes the credentials file
func (c *Credentials) Save() (err error) {
	err = configStore.SaveRC(c)
	return
}

// GetCredentials returns the C14 credentials file
func GetCredentials() (c *Credentials, err error) {
	c = &Credentials{}
	if err = configStore.GetRC(c); err != nil {
		return
	}
	if c.ClientSecret == "" || c.ClientID == "" || c.AccessToken == "" {
		err = errors.Errorf("You need to login first: c14 login")
	}
	return
}
