package commands

// TODO: Refactor the upload for something more generic

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/apex/log"
	"github.com/docker/docker/pkg/progress"
	"github.com/docker/docker/pkg/streamformatter"
	"github.com/dustin/go-humanize"
	"github.com/kr/fs"
	"github.com/scaleway/c14-cli/pkg/api"
	"github.com/scaleway/c14-cli/pkg/utils/ssh"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
)

type upload struct {
	Base
	isPiped bool
	uploadFlags
}

type uploadFlags struct {
	flName string
}

// Upload returns a new command "upload"
func Upload() Command {
	ret := &upload{}
	ret.Init(Config{
		UsageLine:   "upload [DIR|FILE]* ARCHIVE",
		Description: "Upload your file or directory into an archive",
		Help:        "Upload your file or directory into an archive, use SFTP protocol.",
		Examples: `
        $ c14 upload
        $ c14 upload test.go 83b93179-32e0-11e6-be10-10604b9b0ad9
        $ c14 upload /upload 83b93179-32e0-11e6-be10-10604b9b0ad9
        $ tar cvf - /upload 2> /dev/null | ./c14 upload --name "file.tar.gz" fervent_austin
`,
	})
	ret.Flags.StringVar(&ret.flName, []string{"n", "-name"}, "", "Assigns a name (only with tar method)")
	return ret
}

func (u *upload) CheckFlags(args []string) (err error) {
	u.isPiped = !terminal.IsTerminal(int(os.Stdin.Fd()))
	nbArgs := 1

	if !u.isPiped {
		nbArgs = 2
	} else {
		if u.flName == "" && len(args) == 1 {
			err = errors.Errorf("You need to specified a name")
			return
		}
	}
	if len(args) < nbArgs {
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
		padding     int
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

	if u.isPiped && u.flName != "" {
		return u.pipedUpload(sftpConn)
	}
	for _, file := range args {
		var (
			f    *os.File
			info os.FileInfo
		)

		if f, err = os.Open(file); err != nil {
			log.Warnf("Open %s: %s", file, err)
			continue
		}
		if info, err = f.Stat(); err != nil {
			log.Warnf("Stat %s: %s", file, err)
			f.Close()
			continue
		}
		switch mode := info.Mode(); {
		case mode.IsDir():
			walker := fs.Walk(file)
			for walker.Step() {
				if err = walker.Err(); err != nil {
					log.Warnf("Walker %s: %s", walker.Path(), err)
					f.Close()
					continue
				}
				name := walker.Path()
				for name[0] == '/' {
					name = name[1:]
				}
				if walker.Stat().Mode().IsDir() {
					if err = sftpConn.Mkdir("/buffer/" + name); err != nil {
						if err.Error() == "file does not exist" { // bad :/
							sp := strings.Split(name, string(os.PathSeparator))
							path := sp[0]
							for i, n := range sp {
								if i != 0 {
									path = path + "/" + n
								}
								sftpConn.Mkdir("/buffer/" + path)
							}
						}
						continue
					}
					f.Close()
				} else if walker.Stat().Mode().IsRegular() {
					if len(name) > padding {
						padding = len(name)
					}
					files = append(files, uploadFile{
						FileFD: f,
						Info:   info,
						Name:   name,
						Path:   walker.Path(),
					})
				}
			}
		case mode.IsRegular():
			name := filepath.Base(file)
			if len(name) > padding {
				padding = len(name)
			}
			files = append(files, uploadFile{
				FileFD: f,
				Info:   info,
				Name:   name,
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
			log.Warnf("reader Open %s: %s", file.Path, err)
			continue
		}
		if info, err = reader.Stat(); err != nil {
			reader.Close()
			file.FileFD.Close()
			log.Warnf("reader Stat %s: %s", file.Path, err)
			continue
		}
		if err = u.uploadAFile(sftpConn, reader, file.Name, info.Size(), padding); err != nil {
			log.Warnf("upload %s: %s", file.Path, err)
		}
		file.FileFD.Close()
		reader.Close()
	}
	err = nil
	return
}

func (u *upload) uploadAFile(c *sftp.Client, reader io.ReadCloser, file string, size int64, padding int) (err error) {
	log.Debugf("Upload %s -> /buffer/%s", file, file)

	var (
		buff   = make([]byte, 1<<23)
		nr, nw int
		w      *sftp.File
	)
	if w, err = c.Create(fmt.Sprintf("/buffer/%s", file)); err != nil {
		return
	}
	defer w.Close()
	if size == 0 {
		log.Warnf("upload %s is empty", file)
		return
	}
	sf := streamformatter.NewStreamFormatter()
	progressBarOutput := sf.NewProgressOutput(os.Stdout, true)
	rc := progress.NewProgressReader(reader, progressBarOutput, size, "", fmt.Sprintf("%-*s", padding, file))
	defer rc.Close()
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

func (u *upload) pipedUpload(c *sftp.Client) (err error) {
	var (
		buff   = make([]byte, 1<<23)
		nr, nw int
		total  uint64
		w      *sftp.File
	)

	if w, err = c.Create(fmt.Sprintf("/buffer/%s", u.flName)); err != nil {
		return
	}
	for {
		nr, err = os.Stdin.Read(buff)
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
		total += uint64(nr)
		fmt.Printf("\rUploading \t%s", humanize.Bytes(total))

	}
	return
}
