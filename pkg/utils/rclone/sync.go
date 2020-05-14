package rclone

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// CheckRcloneExists : check if rclone is in PATH
func CheckRcloneExists() (err error) {
	path, err := exec.LookPath("rclone")
	if err != nil {
		return
	}
	fmt.Printf("rclone executable is in %s\n", path)

	return
}

// Sync : run rclone sync command
func Sync(safeUUID string, archiveUUID string, s3Profile string) (err error) {

	// rclone destination will be named after awscli profile if --s3-profile is used
	destRemote := "s3"
	if s3Profile != "" {
		destRemote = s3Profile
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configPath := home + "/rclone-c14-migration_" + safeUUID + "_" + archiveUUID + ".conf"

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return
	}

	bucketName := fmt.Sprintf("c14-%s", safeUUID)

	app := "rclone"

	config := "--config=" + configPath
	action := "sync"
	src := "c14:/buffer/"
	dest := destRemote + ":" + bucketName + "/" + archiveUUID
	logLevel := "--log-level=INFO"

	rcloneCmd := fmt.Sprintf("%s %s %s %s %s %s", app, config, action, src, dest, logLevel)

	fmt.Println("Running command: " + rcloneCmd)

	cmd := exec.Command(app, config, action, src, dest, logLevel)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		return
	}
	outStr, _ := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Println(outStr)

	return
}
