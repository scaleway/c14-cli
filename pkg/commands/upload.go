package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/QuentinPerez/c14-cli/pkg/api"
	"github.com/QuentinPerez/c14-cli/pkg/utils/ssh"
	"github.com/apex/log"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/kr/fs"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
)

type upload struct {
	Base
	uploadFlags
}

type uploadFlags struct {
}

// Upload returns a new command "upload"
func Upload() Command {
	ret := &upload{}
	ret.Init(Config{
		UsageLine:   "upload [DIR|FILE]+ ARCHIVE",
		Description: "Upload your file or directory into an archive",
		Help:        "Upload your file or directory into an archive.",
		Examples: `
        $ c14 upload
        $ c14 upload test.go 83b93179-32e0-11e6-be10-10604b9b0ad9
        $ c14 upload /upload 83b93179-32e0-11e6-be10-10604b9b0ad9
`,
	})
	return ret
}

func (u *upload) CheckFlags(args []string) (err error) {
	if len(args) < 2 {
		u.PrintUsage()
		os.Exit(1)
	}
	return
}

func (u *upload) GetName() string {
	return "upload"
}

type uploadFile struct {
	FileFD *os.File
	Info   os.FileInfo
	Path   string
	Name   string
}

func (u *upload) Run(args []string) (err error) {
	if err = u.InitAPI(); err != nil {
		return
	}

	var (
		safe        api.OnlineGetSafe
		bucket      api.OnlineGetBucket
		sftpCred    sshUtils.Credentials
		sftpConn    *sftp.Client
		files       []uploadFile
		uuidArchive string
	)

	archive := args[len(args)-1]
	args = args[:len(args)-1]
	if safe, uuidArchive, err = u.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		if safe, uuidArchive, err = u.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
			return
		}
	}
	if bucket, err = u.OnlineAPI.GetBucket(safe.UUIDRef, uuidArchive); err != nil {
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

	for _, file := range args {
		var (
			f    *os.File
			info os.FileInfo
		)

		if f, err = os.Open(file); err != nil {
			log.Warnf("%s: %s", file, err)
			continue
		}
		if info, err = f.Stat(); err != nil {
			log.Warnf("%s: %s", file, err)
			f.Close()
			continue
		}
		switch mode := info.Mode(); {
		case mode.IsDir():
			walker := fs.Walk(file)
			for walker.Step() {
				if err = walker.Err(); err != nil {
					log.Warnf("%s: %s", walker.Path(), err)
					f.Close()
					continue
				}
				if walker.Stat().Mode().IsDir() {
					if err = sftpConn.Mkdir(walker.Path()); err != nil {
						continue
					}
					f.Close()
				} else if walker.Stat().Mode().IsRegular() {
					files = append(files, uploadFile{
						FileFD: f,
						Info:   info,
						Name:   walker.Path(),
						Path:   walker.Path(),
					})
				}
			}
		case mode.IsRegular():
			files = append(files, uploadFile{
				FileFD: f,
				Info:   info,
				Name:   filepath.Base(file),
				Path:   file,
			})
		}
	}
	for _, file := range files {
		var (
			info   os.FileInfo
			reader *os.File
		)

		if reader, err = os.Open(file.Path); err != nil {
			reader.Close()
			file.FileFD.Close()
			log.Warnf("%s: %s", file.Path, err)
			continue
		}
		if info, err = reader.Stat(); err != nil {
			reader.Close()
			file.FileFD.Close()
			log.Warnf("%s: %s", file.Path, err)
			continue
		}
		if err = u.uploadAFile(sftpConn, reader, file.Name, info.Size()); err != nil {
			log.Warnf("%s: %s", file.Path, err)
		}
		file.FileFD.Close()
		reader.Close()
	}
	err = nil
	return
}

func (u *upload) uploadAFile(c *sftp.Client, reader io.ReadCloser, file string, size int64) (err error) {
	log.Debugf("Upload %s -> /buffer/%s", file, file)

	var (
		buff   = make([]byte, 1<<22)
		nr, nw int
		w      *sftp.File
	)

	sf := streamformatter.NewStreamFormatter()
	progressBarOutput := sf.NewProgressOutput(os.Stdout, true)
	rc := progress.NewProgressReader(reader, progressBarOutput, size, "", file)
	defer rc.Close()
	if w, err = c.Create(fmt.Sprintf("/buffer/%s", file)); err != nil {
		return
	}
	for {
		nr, err = rc.Read(buff)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if nw, err = w.Write(buff[:nr]); err != nil {
			return
		}
		if nw != nr {
			err = errors.Errorf("Error during write")
			return
		}
	}
	return
}
