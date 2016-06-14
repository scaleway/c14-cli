package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/QuentinPerez/c14-cli/pkg/api/oauth2"
	"github.com/QuentinPerez/c14-cli/pkg/config"
	"github.com/juju/errors"

	"golang.org/x/crypto/ssh/terminal"
)

type login struct {
	Base
	loginFlags
}

type loginFlags struct {
	clientID string
}

// Login returns a new command "login"
func Login() Command {
	ret := &login{}
	ret.Init(Config{
		UsageLine:   "login",
		Description: "Log in to Online API",
		Help: `Generates a credentials file in $HOME/.c14rc
containing informations to generate a token`,
		Examples: `
    $ c14 login`,
	})
	ret.Flags.StringVar(&ret.clientID, []string{"-client-id"}, "", "Online's ClientID")
	return ret
}

func (l *login) GetName() string {
	return "login"
}

func (l *login) Run(args []string) (err error) {
	var (
		empty      string
		auth       oauth2.Authentication
		c          oauth2.Credentials
		configFile config.Credentials
	)

	if l.clientID == "" {
		if err = promptUser(`! You must provide a ClientID !

Please opens this link with your web browser: https://console.online.net/en/api/apps
You can create a new app with the following fields:
---
    App name:    c14
    Description: c14 cli
    Permissions:
        (storage) Read/Write
---
Then copy the client_id here : `, &l.clientID, true); err != nil {
			return
		}
	}
	l.clientID = strings.Trim(l.clientID, "\r\n")
	if auth, err = oauth2.GetVerificationURL(l.clientID); err != nil {
		return
	}
	if err = promptUser(fmt.Sprintf(`
Please opens this link with your browser: %v
Then copy paste the code %v (enter when is done) : `, auth.VerficationURL, auth.UserCode), &empty, true); err != nil {
		return
	}
	if c, err = oauth2.GetCredentials(l.clientID, auth.DeviceCode); err != nil {
		return
	}
	configFile.ClientID = c.ClientID
	configFile.ClientSecret = c.ClientSecret
	err = configFile.Save()
	return
}

func promptUser(prompt string, output *string, echo bool) (err error) {
	var b []byte

	fmt.Fprintf(os.Stdout, prompt)
	os.Stdout.Sync()
	if !echo {
		if b, err = terminal.ReadPassword(int(os.Stdin.Fd())); err != nil {
			err = errors.Annotate(err, "Unable to prompt for password")
			return
		}
		*output = string(b)
		fmt.Fprintf(os.Stdout, "\n")
	} else {
		reader := bufio.NewReader(os.Stdin)
		*output, err = reader.ReadString('\n')
	}
	return
}
