package commands

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/juju/errors"
	"github.com/online-net/c14-cli/pkg/api"
)

type create struct {
	Base
	createFlags
}

type createFlags struct {
	flName  string
	flDesc  string
	flSafe  string
	flQuiet bool
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
`,
	})
	ret.Flags.StringVar(&ret.flName, []string{"n", "-name"}, "", "Assigns a name")
	ret.Flags.StringVar(&ret.flDesc, []string{"d", "-description"}, "", "Assigns a description")
	ret.Flags.BoolVar(&ret.flQuiet, []string{"q", "-quiet"}, false, "Don't display the waiting loop")
	ret.Flags.StringVar(&ret.flSafe, []string{"s", "-safe"}, "", "Name of the safe to use. If it doesn't exists it will be created.")
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
	)

	if keys, err = c.OnlineAPI.GetSSHKeys(); err != nil {
		err = errors.Annotate(err, "Run:GetSSHKey")
		return
	}
	if len(keys) == 0 {
		err = errors.New("Please add an SSH Key here: https://console.online.net/en/account/ssh-keys")
		return
	}

	safeName = c.flSafe

	if safeName == "" {
		safeName = fmt.Sprintf("%s_safe", c.flName)
	}

	if _, uuidArchive, _, err = c.OnlineAPI.CreateSSHBucketFromScratch(api.ConfigCreateSSHBucketFromScratch{
		SafeName:    safeName,
		ArchiveName: c.flName,
		Desc:        c.flDesc,
		UUIDSSHKeys: []string{keys[0].UUIDRef},
		Platforms:   []string{"1"},
		Days:        7,
		Quiet:       c.flQuiet,
	}); err != nil {
		err = errors.Annotate(err, "Run:CreateSSHBucketFromScratch")
		return
	}
	fmt.Printf("%s\n", uuidArchive)
	return
}
