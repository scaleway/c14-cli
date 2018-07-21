package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"github.com/scaleway/c14-cli/pkg/api"
	"github.com/scaleway/c14-cli/pkg/utils/ssh"
)

type download struct {
	Base
	isPiped bool
	downloadFlags
}

type downloadFlags struct {
	flName string
}

// Upload returns a new command "upload"
func Download() Command {
	ret := &download{}
	ret.Init(Config{
		UsageLine:   "download [DIR|FILE]* ARCHIVE",
		Description: "Download your file or directory into an archive",
		Help:        "Download your file or directory into an archive, use SFTP protocol.",
		Examples: `
$ c14 download
$ c14 download test.go 83b93179-32e0-11e6-be10-10604b9b0ad9
$ c14 download /download 83b93179-32e0-11e6-be10-10604b9b0ad9
$ tar cvf - /download 2> /dev/null | ./c14 download --name "file.tar.gz" fervent_austin
`,
	})
	ret.Flags.StringVar(&ret.flName, []string{"n", "-name"}, "", "Assigns a name (only with tar method)")
	return ret
}

func (u *download) GetName() string {
	return "download"
}

func (d *download) Run(args []string) (err error) {
	var (
		safe     api.OnlineGetSafe
		bucket   api.OnlineGetBucket
		sftpCred sshUtils.Credentials
		sftpConn *sftp.Client
		//files       []uploadFile
		uuidArchive string
		//padding     int
	)

	if err = d.InitAPI(); err != nil {
		return
	}

	archive := args[len(args)-1]
	args = args[:len(args)-1]
	fmt.Println("args =", args[0])
	fmt.Println("archive =", archive)

	// get UUID
	if safe, uuidArchive, err = d.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		if safe, uuidArchive, err = d.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
			return
		}
	}

	// get bucket
	if bucket, err = d.OnlineAPI.GetBucket(safe.UUIDRef, uuidArchive); err != nil {
		return
	}
	fmt.Println("bucket =", bucket)

	// fill credentials
	sftpCred.Host = strings.Split(bucket.Credentials[0].URI, "@")[1]
	sftpCred.Password = bucket.Credentials[0].Password
	sftpCred.User = bucket.Credentials[0].Login

	// connect
	if sftpConn, err = sftpCred.NewSFTPClient(); err != nil {
		return
	}

	defer sftpCred.Close()
	defer sftpConn.Close()

	fmt.Println("Host =", sftpCred.Host)
	fmt.Println("Password =", sftpCred.Password)
	fmt.Println("User =", sftpCred.User)
	//=======================Connection end====================================

	var fileName string

	splittedString := strings.Split(args[0], "/")
	if splittedString != nil {
		fileName = splittedString[len(splittedString)-1]
	} else {
		fileName = args[0]
	}

	fdLocal, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer fdLocal.Close()
	fmt.Println("file created")

	file := "/buffer/" + args[0]
	fmt.Println("file =", file)
	fdRemote, err := sftpConn.Open(file)
	if err != nil {
		return
	}
	defer fdRemote.Close()

	fdRemote.WriteTo(fdLocal)

	return
}
