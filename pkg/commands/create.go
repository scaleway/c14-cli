package commands

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/juju/errors"
	"github.com/online-net/c14-cli/pkg/api"
	"strings"
)

type create struct {
	Base
	createFlags
}

type createFlags struct {
	flName    string
	flDesc    string
	flSafe    string
	flQuiet   bool
	flParity  string
	flLarge   bool
	flCrypto  bool
	flSshKeys string
}

// Create returns a new command "create"
func Create() Command {
	ret := &create{}
	ret.Init(Config{
		UsageLine:   "create [OPTIONS]",
		Description: "Create a new archive",
		Help:        "Create a new archive, by default with a random name, standard storage (0.0002€/GB/month), automatic locked in 7 days and your datas will be stored at DC2.",
		Examples: `
        $ c14 create
        $ c14 create --name "MyBooks" --description "hardware books"
        $ c14 create --name "MyBooks" --description "hardware books" --safe "Bookshelf"
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
	ret.Flags.StringVar(&ret.flSshKeys, []string{"k", "-sshkey"}, "", "UUID of ssh keys use for ssh connections (separate by comma).")
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
			err = errors.New("Please add an SSH Key here: https://console.online.net/en/account/ssh-keys")
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
		Platforms:   []string{"1"},
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
