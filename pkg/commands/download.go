package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"github.com/scaleway/c14-cli/pkg/api"
	"github.com/scaleway/c14-cli/pkg/utils/ssh"
)

// TODO : flag for download in dest : ./c14 download --dest="~/Documents"

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
        $ c14 download file 83b93179-32e0-11e6-be10-10604b9b0ad9
`,
	})
	//ret.Flags.StringVar(&ret.flName, []string{"n", "-name"}, "", "Assigns a name (only with tar method)")
	return ret
}

func (d *download) CheckFlags(args []string) (err error) {
	if len(args) < 2 {
		d.PrintUsage()
		os.Exit(1)
	}
	return
}

func (d *download) GetName() string {
	return "download"
}

func getCredentials(d *download, archive string) (api.OnlineGetBucket, error) {
	var (
		bucket      api.OnlineGetBucket
		safe        api.OnlineGetSafe
		uuidArchive string
		err         error
	)

	// get UUID
	if safe, uuidArchive, err = d.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		if safe, uuidArchive, err = d.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
			return bucket, err
		}
	}

	// get bucket
	if bucket, err = d.OnlineAPI.GetBucket(safe.UUIDRef, uuidArchive); err != nil {
		return bucket, err
	}

	return bucket, err
}

func connectToSFTP(bucket api.OnlineGetBucket, sftpCred sshUtils.Credentials) (*sftp.Client, error) {

	var sftpConn *sftp.Client

	// fill credentials
	sftpCred.Host = strings.Split(bucket.Credentials[0].URI, "@")[1]
	sftpCred.Password = bucket.Credentials[0].Password
	sftpCred.User = bucket.Credentials[0].Login

	// SFTP connection
	sftpConn, err := sftpCred.NewSFTPClient()

	return sftpConn, err
}

func downloadFile(fileName string, fdRemote *sftp.File) (err error) {
	var fdLocal *os.File // file descriptor to local file

	fmt.Println(fileName)
	// Create new file
	if fdLocal, err = os.Create(fileName); err != nil {
		return
	}
	defer fdLocal.Close()

	// Copy remote file to local file
	if _, err = fdRemote.WriteTo(fdLocal); err != nil {
		return
	}

	return
}

func downloadDir(dirName string, sftpConn *sftp.Client) {
	var (
		fileName string     // Name of file to download
		fdRemote *sftp.File // file descriptor to remote file
	)

	dirName = strings.TrimSuffix(dirName, "/")
	walker := sftpConn.Walk(dirName)

	for walker.Step() {
		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		info, err := sftpConn.ReadDir(walker.Path())
		if err != nil {
			continue
		}
		for i := 0; i < len(info); i++ {
			// path of filename - "/buffer/"
			fileName = walker.Path()[len("/buffer/"):] + "/" + info[i].Name()

			if info[i].IsDir() == true {
				if err = os.MkdirAll(fileName, os.ModePerm); err != nil {
					fmt.Println(err)
				}
			} else {
				if fdRemote, err = sftpConn.Open("/buffer/" + fileName); err != nil {
					fmt.Println("err =", err)
					continue
				}
				if err = downloadFile(fileName, fdRemote); err != nil {
					fmt.Println("err download =", err)
				}
				fdRemote.Close()
			}
		}
	}
}

func (d *download) Run(args []string) (err error) {
	var (
		bucket         api.OnlineGetBucket
		sftpCred       sshUtils.Credentials
		sftpConn       *sftp.Client
		remoteFile     string     // Path to file to download
		fileName       string     // Name of file to download
		fdRemote       *sftp.File // file descriptor to remote file
		statremoteFile os.FileInfo
	)

	if err = d.InitAPI(); err != nil {
		return
	}

	archive := args[len(args)-1]
	args = args[:len(args)-1]

	// get credentials for SFTP connection
	if bucket, err = getCredentials(d, archive); err != nil {
		return
	}

	// connection in SFTP with credentials
	if sftpConn, err = connectToSFTP(bucket, sftpCred); err != nil {
		return
	}
	defer sftpCred.Close()
	defer sftpConn.Close()

	for i := 0; i < len(args); i++ {
		if args[i] == "" {
			continue
		}
		// Path of remote file
		remoteFile = "/buffer/" + args[i]

		// Open remote file
		if fdRemote, err = sftpConn.Open(remoteFile); err != nil {
			return
		}
		defer fdRemote.Close()

		// stat remote file in case file not exist
		if statremoteFile, err = fdRemote.Stat(); err != nil {
			return
		}

		if statremoteFile.IsDir() == true {
			if err = os.MkdirAll(args[i], os.ModePerm); err != nil {
				fmt.Println(err)
			}
			downloadDir(remoteFile, sftpConn)
		} else {
			// Extract name of file to download
			fileName = filepath.Base(args[i])
			if downloadFile(fileName, fdRemote); err != nil {
				return
			}
		}
	}

	return
}
