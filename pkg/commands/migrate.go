package commands

import (
	"fmt"
	"os"
	"strings"

	"net/url"

	"github.com/juju/errors"
	"github.com/pkg/sftp"
	"github.com/scaleway/c14-cli/pkg/api"
	"github.com/scaleway/c14-cli/pkg/utils/rclone"
	"github.com/scaleway/c14-cli/pkg/utils/sis"
	sshUtils "github.com/scaleway/c14-cli/pkg/utils/ssh"
)

type migrate struct {
	Base
	migrateFlags
}

type migrateFlags struct {
	fls3AccessKey    string
	fls3SecretKey    string
	fls3Profile      string
	fls3Bucket       string
	fls3Prefix       string
	fls3CreateBucket bool
}

// Migrate returns a new command "migrate"
func Migrate() Command {
	ret := &migrate{}
	ret.Init(Config{
		UsageLine:   "migrate [OPTIONS] [ACTION] [ARCHIVE]",
		Description: "Migration helper to S3 Cold Storage",
		Help:        "Migrate an archive to Cold Storage\n\n[ACTION] is one of [precheck, generate-rclone-config, rclone-sync]",
		Examples: `
        $ c14 migrate --s3-access-key xxx --s3-secret-key yyy precheck d28d0f7b-4524-4f7c-a7a3-7341503e9110"
        $ c14 migrate --s3-profile scw-par generate-rclone-config d28d0f7b-4524-4f7c-a7a3-7341503e9110`,
	})

	ret.Flags.StringVar(&ret.fls3AccessKey, []string{"-s3-access-key"}, "", "aws_access_key_id")
	ret.Flags.StringVar(&ret.fls3SecretKey, []string{"-s3-secret-key"}, "", "aws_secret_access_key")
	ret.Flags.StringVar(&ret.fls3Profile, []string{"-s3-profile"}, "", "aws_profile")
	ret.Flags.StringVar(&ret.fls3Bucket, []string{"-s3-bucket"}, "", "Destination bucket name")
	ret.Flags.StringVar(&ret.fls3Prefix, []string{"-s3-prefix"}, "", "Prefix in destination bucket")
	ret.Flags.BoolVar(&ret.fls3CreateBucket, []string{"-s3-create-bucket"}, false, "Prefix in destination bucket")

	return ret
}

func (m *migrate) GetName() string {
	return "migrate"
}

func (m *migrate) CheckFlags(args []string) (err error) {
	if len(args) != 2 {
		m.PrintUsage()
		os.Exit(1)
	}
	return
}

func (m *migrate) Run(args []string) (err error) {

	if m.fls3Profile != "" {
		fmt.Println("Using AWS profile " + m.fls3Profile)
	} else {
		if m.fls3AccessKey == "" || m.fls3SecretKey == "" {
			return errors.New("Please set --s3-access-key and --s3-access-key or use --s3-aws-profile")
		}
	}

	if err = m.InitAPI(); err != nil {
		return
	}

	switch args[0] {
	case "precheck":
		if err = m.precheck(args); err != nil {
			return
		}
	case "generate-rclone-config":
		if err = m.generateRcloneConfig(args); err != nil {
			return
		}
	case "rclone-sync":
		if err = m.runRcloneSync(args); err != nil {
			return
		}
	default:
		return errors.New("invalid action")
	}

	return
}

func (m *migrate) precheck(args []string) (err error) {
	var (
		safe                 api.OnlineGetSafe
		bucket               api.OnlineGetBucket
		sftpCred             sshUtils.Credentials
		sftpConn             *sftp.Client
		archiveUUID, archive = args[1], args[1]
	)

	if safe, archiveUUID, err = m.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		if safe, archiveUUID, err = m.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
			return
		}
	}
	if bucket, err = m.OnlineAPI.GetBucket(safe.UUIDRef, archiveUUID); err != nil {
		return
	}

	fmt.Println("Making sure all files are < 5 TB for S3 compatibility")
	sftpCred.Host = strings.Split(bucket.Credentials[0].URI, "@")[1]
	sftpCred.Password = bucket.Credentials[0].Password
	sftpCred.User = bucket.Credentials[0].Login
	if sftpConn, err = sftpCred.NewSFTPClient(); err != nil {
		return
	}
	defer sftpCred.Close()
	defer sftpConn.Close()
	w := sftpConn.Walk("/buffer")
	for w.Step() {
		if w.Err() != nil {
			continue
		}

		if w.Stat().Size() > 5497559962838 {
			return errors.New("File '" + w.Path() + "' is bigger than 5 TB")
		}
	}

	err = sis.CheckAPI(m.fls3AccessKey, m.fls3SecretKey, m.fls3Profile)
	if err != nil {
		return
	}

	bucketName := m.fls3Bucket
	if bucketName == "" {
		bucketName = fmt.Sprintf("c14-%s", safe.UUIDRef)
	}

	fmt.Printf("Checking if S3 migration destination bucket %s exists...\n", bucketName)
	bucketExists, err := sis.CheckBucket(bucketName, m.fls3AccessKey, m.fls3SecretKey, m.fls3Profile)
	if err != nil {
		return
	}

	if !bucketExists {
		if m.fls3CreateBucket {
			fmt.Println("Creating bucket...")
			err = sis.CreateBucket(bucketName, m.fls3AccessKey, m.fls3SecretKey, m.fls3Profile)
			if err != nil {
				return
			}
		} else {
			fmt.Println("You can use --s3-create-bucket to automatically create the bucket")
		}
	}

	fmt.Println("All good!")

	return
}

func (m *migrate) generateRcloneConfig(args []string) (err error) {
	var (
		safe                 api.OnlineGetSafe
		bucket               api.OnlineGetBucket
		archiveUUID, archive = args[1], args[1]
	)

	if safe, archiveUUID, err = m.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		if safe, archiveUUID, err = m.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
			return
		}
	}
	if bucket, err = m.OnlineAPI.GetBucket(safe.UUIDRef, archiveUUID); err != nil {
		return
	}
	u, err := url.Parse(bucket.Credentials[0].URI)
	if err != nil {
		return
	}

	err = rclone.GenerateConfig(rclone.Config{
		SafeUUID:    safe.UUIDRef,
		ArchiveUUID: archiveUUID,
		C14Host:     u.Hostname(),
		C14Port:     u.Port(),
		C14User:     u.User.Username(),
		C14Password: bucket.Credentials[0].Password,
		S3AccessKey: m.fls3AccessKey,
		S3SecretKey: m.fls3SecretKey,
		S3Profile:   m.fls3Profile,
	})

	return
}

func (m *migrate) runRcloneSync(args []string) (err error) {
	var (
		safe                 api.OnlineGetSafe
		archiveUUID, archive = args[1], args[1]
	)

	if safe, archiveUUID, err = m.OnlineAPI.FindSafeUUIDFromArchive(archive, true); err != nil {
		if safe, archiveUUID, err = m.OnlineAPI.FindSafeUUIDFromArchive(archive, false); err != nil {
			return
		}
	}
	fmt.Println("Checking if rclone is installed")
	err = rclone.CheckRcloneExists()
	if err != nil {
		return
	}

	bucketName := m.fls3Bucket
	if bucketName == "" {
		bucketName = fmt.Sprintf("c14-%s", safe.UUIDRef)
	}

	bucketPrefix := m.fls3Prefix
	if bucketPrefix == "" {
		bucketPrefix = archiveUUID
	}

	fmt.Println("Running sync")
	err = rclone.Sync(safe.UUIDRef, archiveUUID, m.fls3Profile, bucketName, bucketPrefix)
	if err != nil {
		return
	}

	fmt.Println("Sync done.")
	fmt.Println("Please freeze or delete your archive once you made sure it migrated properly.")
	fmt.Println("To freeze the archive, run: c14 freeze " + archiveUUID)
	return
}
