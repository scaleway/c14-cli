package commands

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/QuentinPerez/c14-cli/pkg/api/auth"
	"github.com/juju/errors"

	"golang.org/x/crypto/ssh/terminal"
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
		Help: `Generates a credentials file in $HOME/.c14rc
containing informations to generate a token`,
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
Then copy paste the code %v
`, authentication.VerficationURL, authentication.UserCode)
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
