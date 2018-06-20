package commands

import (
	"fmt"
	"os"

	"strings"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/juju/errors"
	"github.com/scaleway/c14-cli/pkg/api"
)

type create struct {
	Base
	createFlags
}

type createFlags struct {
	flName     string
	flDesc     string
	flSafe     string
	flQuiet    bool
	flParity   string
	flLarge    bool
	flCrypto   bool
	flSshKeys  string
	flPlatform string
}

// Create returns a new command "create"
func Create() Command {
	ret := &create{}
	ret.Init(Config{
		UsageLine:   "create [OPTIONS]",
		Description: "Create a new archive",
		Help:        "Create a new archive, by default with a random name, standard storage (0.0002€/GB/month), automatic locked in 7 days and your datas will be stored in the choosen platform (by default at DC4).",
		Examples: `
        $ c14 create
        $ c14 create --name "MyBooks" --description "hardware books" -P 1
        $ c14 create --name "MyBooks" --description "hardware books" --safe "Bookshelf" --platform 2
        $ c14 create --sshkey "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx,xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
`,
	})
	ret.Flags.StringVar(&ret.flName, []string{"n", "-name"}, "", "Assigns a name")
	ret.Flags.StringVar(&ret.flDesc, []string{"d", "-description"}, "", "Assigns a description")
	ret.Flags.BoolVar(&ret.flQuiet, []string{"q", "-quiet"}, false, "Don't display the waiting loop")
	ret.Flags.StringVar(&ret.flSafe, []string{"s", "-safe"}, "", "Name of the safe to use. If it doesn't exists it will be created.")
	ret.Flags.StringVar(&ret.flParity, []string{"p", "-parity"}, "standard", "Specify a parity to use")
	ret.Flags.BoolVar(&ret.flLarge, []string{"l", "-large"}, false, "Ask for a large bucket")
	ret.Flags.BoolVar(&ret.flCrypto, []string{"c", "-crypto"}, true, "Enable aes-256-bc cryptography, enabled by default.")
	ret.Flags.StringVar(&ret.flSshKeys, []string{"k", "-sshkey"}, "", "A list of UUIDs corresponding to the SSH keys (separated by a comma) that will be used for the connection.")
	ret.Flags.StringVar(&ret.flPlatform, []string{"P", "-platform"}, "2", "Select the platform (by default at DC4)")
	return ret
}

func (c *create) GetName() string {
	return "create"
}

func (c *create) CheckFlags(args []string) (err error) {
	if len(args) != 0 {
		c.PrintUsage()
		os.Exit(1)
	}

	if c.flName == "" {
		c.flName = namesgenerator.GetRandomName(0)
	}
	if c.flDesc == "" {
		c.flDesc = " "
	}
	return
}

func (c *create) Run(args []string) (err error) {
	if err = c.InitAPI(); err != nil {
		return
	}
	var (
		uuidArchive string
		safeName    string
		keys        []api.OnlineGetSSHKey
		crypto      string
		UuidSshKeys []string
	)

	if len(c.flSshKeys) == 0 {
		if keys, err = c.OnlineAPI.GetSSHKeys(); err != nil {
			err = errors.Annotate(err, "Run:GetSSHKey")
			return
		}
		if len(keys) == 0 {
			err = errors.New("Please add your SSH key here: https://console.online.net/en/account/ssh-keys")
			return
		}
		UuidSshKeys = append(UuidSshKeys, keys[0].UUIDRef)
	} else {
		UuidSshKeys = strings.Split(c.flSshKeys, ",")
		for _, keyArg := range UuidSshKeys {
			_, checkErr := c.OnlineAPI.GetSSHKey(keyArg)
			if checkErr != nil {
				err = errors.New(fmt.Sprintf("%s : %s", checkErr.Error(), keyArg))
				return
			}
		}
	}

	safeName = c.flSafe

	if safeName == "" {
		safeName = fmt.Sprintf("%s_safe", c.flName)
	}

	if c.flCrypto == false {
		crypto = "none"
	} else {
		crypto = "aes-256-cbc"
	}

	if _, uuidArchive, _, err = c.OnlineAPI.CreateSSHBucketFromScratch(api.ConfigCreateSSHBucketFromScratch{
		SafeName:    safeName,
		ArchiveName: c.flName,
		Desc:        c.flDesc,
		UUIDSSHKeys: UuidSshKeys,
		Platforms:   []string{c.flPlatform},
		Days:        7,
		Quiet:       c.flQuiet,
		Parity:      c.flParity,
		LargeBucket: c.flLarge,
		Crypto:      crypto,
	}); err != nil {
		err = errors.Annotate(err, "Run:CreateSSHBucketFromScratch")
		return
	}
	fmt.Printf("%s\n", uuidArchive)
	return
}
