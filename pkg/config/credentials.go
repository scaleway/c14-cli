package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"runtime"

	"github.com/juju/errors"
)

type Credentials struct {
	DeviceCode   string `json:"device_code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func getCredentialsFile() (path string, err error) {
	var u *user.User

	u, err = user.Current()
	if err != nil {
		err = errors.Trace(err)
		return
	}
	path = fmt.Sprintf("%s/.c14rc", u.HomeDir)
	return
}

// Save writes the credentials file
func (c *Credentials) Save() (err error) {
	var (
		path  string
		c14rc *os.File
	)

	path, err = getCredentialsFile()
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

	path, err = getCredentialsFile()
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
