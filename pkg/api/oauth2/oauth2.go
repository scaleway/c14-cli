package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"runtime"

	"golang.org/x/oauth2"

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

type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func getURL(url string, encode, decode interface{}) (err error) {
	var resp *http.Response

	values, errs := encurl.Translate(encode)
	if errs != nil {
		err = errors.Trace(errs[0])
		return
	}
	resp, err = http.DefaultClient.Get(fmt.Sprintf("%v?%v", url, values.Encode()))
	if resp != nil {
		defer resp.Body.Close()
	}
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

func GenerateCredentials(clientID, deviceCode string) (c Credentials, err error) {
	var (
		requestCredentials struct {
			ClientID string `url:"client_id,ifStringIsNotEmpty"`
			Code     string `url:"code,ifStringIsNotEmpty"`
		}
	)
	requestCredentials.ClientID = clientID
	requestCredentials.Code = deviceCode
	err = getURL(credentialsURL, requestCredentials, &c)
	return
}

func getCredentialsPath() (path string, err error) {
	var u *user.User

	u, err = user.Current()
	if err != nil {
		err = errors.Trace(err)
		return
	}
	path = fmt.Sprintf("%s/.c14rc", u.HomeDir)
	return
}

func (c *Credentials) Token() (t *oauth2.Token, err error) {
	err = errors.New("Not implemented yet")
	return
}

// Save writes the credentials file
func (c *Credentials) Save() (err error) {
	var (
		path  string
		c14rc *os.File
	)

	path, err = getCredentialsPath()
	if err != nil {
		err = errors.Annotate(err, "Unable to get credentials file")
		return
	}
	c14rc, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		err = errors.Annotatef(err, "Unable to create %v config file", path)
		return
	}
	defer c14rc.Close()
	encoder := json.NewEncoder(c14rc)
	err = encoder.Encode(c)
	if err != nil {
		err = errors.Annotatef(err, "Unable to encode %v", path)
		return
	}
	return
}

// GetCredentials returns the C14 credentials file
func GetCredentials() (c *Credentials, err error) {
	var (
		path        string
		fileContent []byte
	)

	path, err = getCredentialsPath()
	if err != nil {
		err = errors.Annotate(err, "Unable to get credentials file")
		return
	}
	// Don't check permissions on Windows
	if runtime.GOOS != "windows" {
		stat, errStat := os.Stat(path)
		if errStat == nil {
			perm := stat.Mode().Perm()
			if perm&0066 != 0 {
				err = errors.Errorf("Permissions %#o for %v are too open", perm, path)
				return
			}
		} else {
			err = errors.Errorf("You need to login first: c14 login")
			return
		}
	}
	fileContent, err = ioutil.ReadFile(path)
	if err != nil {
		err = errors.Annotatef(err, "Unable to read %v", path)
		return
	}
	c = &Credentials{}
	err = json.Unmarshal(fileContent, c)
	return
}
