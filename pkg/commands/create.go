package commands

import (
	"fmt"
	"time"

	"github.com/QuentinPerez/c14-cli/pkg/api"
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/juju/errors"
)

type create struct {
	Base
	createFlags
}

type createFlags struct {
	flName string
	flDesc string
}

// Create returns a new command "create"
func Create() Command {
	ret := &create{}
	ret.Init(Config{
		UsageLine:   "",
		Description: "",
		Help:        "",
		Examples:    "",
	})
	ret.Flags.StringVar(&ret.flName, []string{"n", "-name"}, "", "Assigns a name")
	ret.Flags.StringVar(&ret.flDesc, []string{"d", "-description"}, "", "Assigns a description")
	return ret
}

func (c *create) GetName() string {
	return "create"
}

func (c *create) CheckFlags() (err error) {
	if c.flName == "" {
		c.flName = namesgenerator.GetRandomName(0)
	}
	if c.flDesc == "" {
		c.flDesc = fmt.Sprintf("Archive created at %s", time.Now())
	}
	return
}

func (c *create) Run(args []string) (err error) {
	if err = c.InitAPI(); err != nil {
		return err
	}
	var (
		uuidArchive string
		keys        []api.OnlineGetSSHKey
		bucket      api.OnlineGetBucket
	)

	if keys, err = c.OnlineAPI.GetSSHKeys(); err != nil {
		err = errors.Annotate(err, "Run:GetSSHKey")
		return
	}
	if len(keys) == 0 {
		err = errors.New("Please add an SSH Key here: https://console.online.net/en/account/ssh-keys")
		return
	}
	if _, uuidArchive, bucket, err = c.OnlineAPI.CreateSSHBucketFromScratch(api.ConfigCreateSSHBucketFromScratch{
		SafeName:    fmt.Sprintf("%s_safe", c.flName),
		ArchiveName: c.flName,
		Desc:        c.flDesc,
		UUIDSSHKeys: []string{keys[0].UUIDRef},
		Platforms:   []string{"1"},
	}); err != nil {
		err = errors.Annotate(err, "Run:CreateSSHBucketFromScratch")
		return
	}

	fmt.Println("UUID's archive:", uuidArchive)
	fmt.Println("Login:", bucket.Credentials[0].Login)
	fmt.Println("Password:", bucket.Credentials[0].Password)
	fmt.Println("Access:", bucket.Credentials[0].URI)
	return
}
