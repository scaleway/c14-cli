package sshUtils

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Credentials struct {
	Password  string
	Host      string
	User      string
	sshClient *ssh.Client
}

func (c *Credentials) NewSFTPClient() (client *sftp.Client, err error) {
	if _, err = c.NewSSHClient(); err != nil {
		return
	}

	client, err = sftp.NewClient(c.sshClient)
	return
}

func (c *Credentials) NewSSHClient() (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
	}

	var err error
	c.sshClient, err = ssh.Dial("tcp", c.Host, sshConfig)
	return c.sshClient, err
}

func (c *Credentials) Close() (err error) {
	if c.sshClient != nil {
		err = c.sshClient.Close()
	}
	return
}
