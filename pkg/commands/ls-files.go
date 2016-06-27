package commands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/QuentinPerez/c14-cli/pkg/api"
	"github.com/QuentinPerez/c14-cli/pkg/utils/ssh"
	"github.com/apex/log"
	"github.com/dustin/go-humanize"
	"github.com/pkg/sftp"
)

type lsFiles struct {
	Base
	lsFilesFlags
}

type lsFilesFlags struct {
}

// LsFiles returns a new command "lsFiles"
func LsFiles() Command {
	ret := &lsFiles{}
	ret.Init(Config{
		UsageLine:   "ls-files ARCHIVE",
		Description: "List the archive files",
		Help:        "List the archive files.",
		Examples: `
        $ c14 ls-files 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	return ret
}

func (l *lsFiles) GetName() string {
	return "ls-files"
}

func (l *lsFiles) Run(args []string) (err error) {
	if len(args) == 0 {
		l.PrintUsage()
		return
	}
	if err = l.InitAPI(); err != nil {
		return
	}
	l.FetchRessources(true, true)

	var (
		safe        api.OnlineGetSafe
		bucket      api.OnlineGetBucket
		sftpCred    sshUtils.Credentials
		sftpConn    *sftp.Client
		uuidArchive = args[0]
	)

	if safe, err = l.OnlineAPI.FindSafeUUIDFromArchive(uuidArchive, true); err != nil {
		return
	}
	if bucket, err = l.OnlineAPI.GetBucket(safe.UUIDRef, uuidArchive); err != nil {
		return
	}
	sftpCred.Host = strings.Split(bucket.Credentials[0].URI, "@")[1]
	sftpCred.Password = bucket.Credentials[0].Password
	sftpCred.User = bucket.Credentials[0].Login
	if sftpConn, err = sftpCred.NewSFTPClient(); err != nil {
		return
	}
	defer sftpCred.Close()
	defer sftpConn.Close()
	walker := sftpConn.Walk("/buffer")
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	fmt.Fprintf(w, "NAME\tSIZE\n")
	for walker.Step() {
		if err = walker.Err(); err != nil {
			log.Debugf("%s", err)
			continue
		}
		if walker.Stat().Mode().IsDir() {
			fmt.Fprintf(w, "%s\t\n", walker.Path())
		} else {
			fmt.Fprintf(w, "%s\t%s\n", walker.Path(), humanize.Bytes(uint64(walker.Stat().Size())))
		}
	}
	w.Flush()
	return
}
