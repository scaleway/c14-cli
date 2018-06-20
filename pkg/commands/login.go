package commands

import (
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/scaleway/c14-cli/pkg/api/auth"
)

type login struct {
	Base
	loginFlags
}

type loginFlags struct {
}

// Login returns a new command "login"
func Login() Command {
	ret := &login{}
	ret.Init(Config{
		UsageLine:   "login",
		Description: "Log in to Online API",
		Help: `Generates a credentials file in $CONFIG/c14-cli/c14rc.json
containing informations to generate a token.`,
		Examples: `
    $ c14 login`,
	})
	return ret
}

func (l *login) GetName() string {
	return "login"
}

func (l *login) Run(args []string) (err error) {
	var (
		authentication auth.Authentication
		credentials    auth.Credentials
		clientID       = "38320_2wln446j992cgo0088g04coo8gswkcws0c4sww0oo0ggs8kos8"
	)

	if authentication, err = auth.GetVerificationURL(clientID); err != nil {
		return
	}
	fmt.Printf(`Please opens this link with your browser: %v
Then copy paste the code %v`, authentication.VerficationURL, authentication.UserCode)
	for i := 0; i < 1500; i++ {
		if credentials, err = auth.GenerateCredentials(clientID, authentication.DeviceCode); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		err = errors.Annotate(err, "Timeout")
		return
	}
	err = credentials.Save()
	return
}
