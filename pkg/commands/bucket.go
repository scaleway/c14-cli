package commands

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	"github.com/scaleway/c14-cli/pkg/api"
	"os"
	"strings"
	"text/tabwriter"
)

type bucket struct {
	Base
	bucketFlags
}

type bucketFlags struct {
	flPretty bool
}

func Bucket() Command {
	ret := &bucket{}
	ret.Init(Config{
		UsageLine:   "bucket [OPTIONS] [ARCHIVE]*",
		Description: "Display all information of bucket",
		Help:        "Displays protocole and connection informations of bucket.",
		Examples: `
        $ c14 bucket 83b93179-32e0-11e6-be10-10604b9b0ad9
	$ c14 bucket 83b93179-32e0-11e6-be10-10604b9b0ad9 -p`,
	})
	ret.Flags.BoolVar(&ret.flPretty, []string{"p", "-pretty"}, false, "Show all information in tab (default json output)")
	return ret
}

func (l *bucket) GetName() string {
	return "bucket"
}

func (l *bucket) CheckFlags(args []string) (err error) {
	if len(args) == 0 {
		l.PrintUsage()
		os.Exit(1)
	}
	return
}

func (l *bucket) Run(args []string) (err error) {
	var (
		safe   api.OnlineGetSafe
		bucket api.OnlineGetBucket
	)

	if err = l.InitAPI(); err != nil {
		return
	}
	l.OnlineAPI.FetchRessources()
	if safe, _, err = l.FindSafeUUIDFromArchive(args[0], true); err != nil {
		if safe, _, err = l.FindSafeUUIDFromArchive(args[0], false); err != nil {
			return
		}
	}
	if bucket, err = l.OnlineAPI.GetBucket(safe.UUIDRef, args[0]); err != nil {
		return
	}
	l.displayBucketInfo(bucket, args[0])
	return
}

func (l *bucket) displayBucketInfo(bucket api.OnlineGetBucket, archiveUUIDRef string) {
	var (
		err        error
		bucketCred []byte
	)

	if bucketCred, err = json.Marshal(bucket.Credentials); err != nil {
		err = errors.Annotate(err, "displayBucketInfo:MarshalCred")
		return
	}
	if l.flPretty {
		var w *tabwriter.Writer
		w = tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
		defer w.Flush()

		fmt.Fprintf(w, "TYPE\tID\tPASSWORD\tURI\tSSH KEYS (SSH TYPE ONLY)\n")
		for _, cred := range bucket.Credentials {
			ssh_descs := ""
			if len(cred.SSHKeys) > 0 {
				descKeys := []string{}
				for _, key := range cred.SSHKeys {
					descKeys = append(descKeys, key.Desc)
				}
				ssh_descs = strings.Join(descKeys, ",")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", cred.Protocol, cred.Login, cred.Password, cred.URI, ssh_descs)
		}
	} else {
		fmt.Println(string(bucketCred))
	}
	return
}
