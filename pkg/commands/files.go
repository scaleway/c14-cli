package commands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/apex/log"
	"github.com/dustin/go-humanize"
	"github.com/scaleway/c14-cli/pkg/api"
	"github.com/scaleway/c14-cli/pkg/utils/ssh"
	"github.com/pkg/sftp"
)

type files struct {
	Base
	filesFlags
}

type filesFlags struct {
}

// Files returns a new command "files"
func Files() Command {
	ret := &files{}
	ret.Init(Config{
		UsageLine:   "files ARCHIVE",
		Description: "List the files of an archive",
		Help:        "List the files of an archive, displays the name and size of files",
		Examples: `
        $ c14 files 83b93179-32e0-11e6-be10-10604b9b0ad9`,
	})
	return ret
}

func (l *files) CheckFlags(args []string) (err error) {
	if len(args) == 0 {
		l.PrintUsage()
		os.Exit(1)
	}
	return
}

func (l *files) GetName() string {
	return "files"
}

func (l *files) Run(args []string) (err error) {
	if err = l.InitAPI(); err != nil {
		return
	}

	var (
		safe                 api.OnlineGetSafe
		bucket               api.OnlineGetBucket
		sftpCred             sshUtils.Credentials
		sftpConn             *sftp.Client
		uuidArchive, archive = args[0], args[0]
	)

	if safe, uuidArchive, err = l.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		if safe, uuidArchive, err = l.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
			return
		}
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
			if walker.Path() != "/buffer" {
				fmt.Fprintf(w, "%s/\t\n", walker.Path()[8:])
			}
		} else {
			fmt.Fprintf(w, "%s\t%s\n", walker.Path()[8:], humanize.Bytes(uint64(walker.Stat().Size())))
		}
	}
	w.Flush()
	return
}
