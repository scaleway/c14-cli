package commands

// TODO: Refactor the shell for something more generic

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/QuentinPerez/c14-cli/pkg/api"
	"github.com/QuentinPerez/c14-cli/pkg/utils/ssh"
	"golang.org/x/crypto/ssh"
)

type shell struct {
	Base
}

// Shell returns a new command "shell"
func Shell() Command {
	ret := &shell{}
	ret.Init(Config{
		UsageLine:   "shell ARCHIVE",
		Description: "Start a shell on an active buffer",
		Help:        "Start a shell on an active buffer",
		Examples: `
        $ c14 shell myarchive
`,
	})
	return ret
}

func (u *shell) CheckFlags(args []string) (err error) {
	if len(args) < 1 {
		u.PrintUsage()
		os.Exit(1)
	}
	return
}

func (u *shell) GetName() string {
	return "shell"
}

func (u *shell) Run(args []string) (err error) {
	if err = u.InitAPI(); err != nil {
		return
	}

	var (
		safe        api.OnlineGetSafe
		bucket      api.OnlineGetBucket
		sshCred     sshUtils.Credentials
		sshConn     *ssh.Client
		uuidArchive string
		session     *ssh.Session
	)

	archive := args[0]
	if safe, uuidArchive, err = u.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		if safe, uuidArchive, err = u.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
			return
		}
	}
	if bucket, err = u.OnlineAPI.GetBucket(safe.UUIDRef, uuidArchive); err != nil {
		return
	}
	sshCred.Host = strings.Split(bucket.Credentials[0].URI, "@")[1]
	sshCred.Password = bucket.Credentials[0].Password
	sshCred.User = bucket.Credentials[0].Login
	if sshConn, err = sshCred.NewSSHClient(); err != nil {
		return
	}
	fmt.Println(sshCred)
	fmt.Println(sshConn)

	defer sshCred.Close()
	defer sshConn.Close()

	if session, err = sshConn.NewSession(); err != nil {
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	err = session.Run("/bin/sh")

	return
}
