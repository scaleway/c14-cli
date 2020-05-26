package rclone

import (
	"fmt"
	"os"
	"text/template"

	fsConfig "github.com/rclone/rclone/fs/config/obscure"
)

type Config struct {
	SafeUUID    string
	ArchiveUUID string
	C14Host     string
	C14Port     string
	C14User     string
	C14Password string
	S3AccessKey string
	S3SecretKey string
	S3Profile   string
}

// GenerateConfig : generate rclone config wit bot remotes from a template
func GenerateConfig(config Config) (err error) {
	fmt.Println("Converting SFTP password for rclone config")

	config.C14Password, err = fsConfig.Obscure(config.C14Password)
	if err != nil {
		return
	}

	tmpl, err := template.New("conf").Parse(
		`[c14]
type = sftp
host = {{ .C14Host }}
port = {{ .C14Port }}
user = {{ .C14User }}
pass = {{ .C14Password }}
md5sum_command = md5sum
sha1sum_command = sha1sum

[{{ if ne .S3Profile "" }}{{ .S3Profile }}{{ else }}s3{{ end }}]
type = s3
provider = Scaleway
{{ if ne .S3Profile "" -}}
env_auth =  true
{{ else -}}
access_key_id = {{ .S3AccessKey }}
secret_access_key = {{ .S3SecretKey }}
{{ end -}}
region = fr-par
endpoint = https://s3.fr-par.scw.cloud
storage_class = GLACIER
`)
	if err != nil {
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configPath := home + "/rclone-c14-migration_" + config.SafeUUID + "_" + config.ArchiveUUID + ".conf"

	f, err := os.Create(configPath)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Println("Writing config file to", configPath)

	err = tmpl.Execute(f, config)
	if err != nil {
		return
	}

	fmt.Println("The following config file has been generated: ")
	err = tmpl.Execute(os.Stdout, config)
	if err != nil {
		return
	}

	return
}
