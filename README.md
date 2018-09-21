# C14 CLI

Interact with [Online C14](https://www.online.net/en/c14) API from the command line.

![Online C14 logo](https://raw.githubusercontent.com/scaleway/c14-cli/master/assets/logo.png)

#### Table of Contents

1. [Overview](#overview)
2. [Setup](#setup)
3. [Usage](#usage)
  * [Login](#login)
  * [Commands](#commands)
    * [`create`](#c14-create)
    * [`freeze`](#c14-freeze)
    * [`help`](#c14-help)
    * [`ls`](#c14-ls)
    * [`login`](#c14-login)
    * [`files`](#c14-files)
    * [`rename`](#c14-rename)
    * [`remove`](#c14-remove)
    * [`unfreeze`](#c14-unfreeze)
    * [`upload`](#c14-upload)
    * [`version`](#c14-version)
    * [`download`](#c14-download)
  * [Examples](#examples)
4. [Changelog](#changelog)
5. [Development](#development)
  * [Hack](#hack)
6. [License](#license)

## Overview

A command-line tool to manage your C14 storage easily

## Setup

⚠️ Ensure you have a go version >= 1.6

```shell
go get -u github.com/scaleway/c14-cli/cmd/c14
```

## Usage


```console
$ c14
Usage: c14 [OPTIONS] COMMAND [arg...]

Interact with C14 from the command line.

Options:
  -D, --debug        Enable debug mode
  -V, --verbose      Enable verbose mode

Commands:
    create    Create a new archive
    files     List the files of an archive
    freeze    Lock an archive
    help      Help of the c14 command line
    login     Log in to Online API
    ls        List the archives
    rename    Rename an archive
    remove    Remove an archive
    unfreeze  Unlock an archive
    upload    Upload your file or directory into an archive
    bucket    Displays all information of bucket
    version   Show the version information
    download  Download your file or directory into an archive

Run 'c14 COMMAND --help' for more information on a command.
```

### Login

```console
$ c14 login
Please opens this link with your browser: https://console.online.net/oauth/v2/device/usercode
Then copy paste the code XXXXXX
$
```

### Commands

#### `c14 create`

```console
Usage: c14 create [OPTIONS]

Create a new archive, by default with a random name, standard storage (0.0002€/GB/month), automatic locked in 7 days and your datas will be stored in the choosen platform (by default at DC4).

Options:
  -c, --crypto=true     Enable aes-256-bc cryptography, enabled by default.
  -d, --description=""  Assigns a description
  -h, --help=false      Print usage
  -k, --sshkey=""       A list of UUIDs corresponding to the SSH keys (separated by a comma) that will be used for the connection.
  -l, --large=false     Ask for a large bucket
  -n, --name=""         Assigns a name
  -P, --platform=2      Select the platform (by default at DC4)
  -p, --parity=standard Specify a parity to use
  -q, --quiet=false     Don't display the waiting loop
  -s, --safe=""         Name of the safe to use. If it doesn't exists it will be created.

Examples:
        $ c14 create
        $ c14 create --name "MyBooks" --description "hardware books" -P 1
        $ c14 create --name "MyBooks" --description "hardware books" --safe "Bookshelf" --platform 2
        $ c14 create --sshkey "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx,xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

```

#### `c14 freeze`

```console
Usage: c14 freeze [OPTIONS] [ARCHIVE]+

Lock an archive, your archive will be stored in highly secure Online data centers and will stay available On Demand (0.01€/GB).

Options:
  -h, --help=false      Print usage
  --nowait=false
  -q, --quiet=false

Examples:
        $ c14 freeze 83b93179-32e0-11e6-be10-10604b9b0ad9
```


#### `c14 ls`

```console
Usage: c14 ls [OPTIONS] [ARCHIVE]*

Displays the archives, by default only the NAME, STATUS, UUID.

Options:
  -a, --all=false       Show all information on archives (size,parity,creationDate,description)
  -h, --help=false      Print usage
  -p, --platform=false  Show the platforms
  -q, --quiet=false     Only display UUIDs

Examples:
        $ c14 ls
        $ c14 ls -a
```

#### `c14 help`

```console
Usage: c14 help [COMMAND]

Help prints help information about c14 and its commands.
By default, help lists available commands.
When invoked with a command name, it prints the usage and the help of
the command.

Options:
  -h, --help=false      Print usage

Examples:
        $ c14 help
        $ c14 help create
```


#### `c14 login`

```console
Usage: c14 login

Generates a credentials file in $CONFIG/c14-cli/c14rc.json
containing informations to generate a token.

Options:
  -h, --help=false      Print usage

Examples:
    $ c14 login
```


#### `c14 files`

```console
Usage: c14 files ARCHIVE

List the files of an archive, displays the name and size of files

Options:
  -h, --help=false      Print usage

Examples:
        $ c14 files 83b93179-32e0-11e6-be10-10604b9b0ad9
```


#### `c14 rename`

```console
Usage: c14 rename ARCHIVE new_name

Rename an archive.

Options:
  -h, --help=false      Print usage

Examples:
        $ c14 rename 83b93179-32e0-11e6-be10-10604b9b0ad9 new_name
        $ c14 rename old_name new_name
```


#### `c14 remove`

```console
Usage: c14 remove [ARCHIVE]+

Remove an archive

Options:
  -h, --help=false      Print usage

Examples:
        $ c14 remove 83b93179-32e0-11e6-be10-10604b9b0ad9 2d752399-429f-447f-85cd-c6104dfed5db
```


#### `c14 unfreeze`

```console
Usage: c14 unfreeze [OPTIONS] [ARCHIVE]+

Unlock an archive, extraction of the archive's data (0.01€/GB).

Options:
  -h, --help=false      Print usage
  --nowait=false
  -q, --quiet=false

Examples:
        $ c14 unfreeze 83b93179-32e0-11e6-be10-10604b9b0ad9
```


#### `c14 upload`

```console
Usage: c14 upload [DIR|FILE]* ARCHIVE

Upload your file or directory into an archive, use SFTP protocol.

Options:
  -h, --help=false      Print usage
  -n, --name=""         Assigns a name (only with tar method)

Examples:
        $ c14 upload
        $ c14 upload test.go 83b93179-32e0-11e6-be10-10604b9b0ad9
        $ c14 upload /upload 83b93179-32e0-11e6-be10-10604b9b0ad9
        $ tar cvf - /upload 2> /dev/null | ./c14 upload --name "file.tar.gz" fervent_austin
```

#### `c14 download`

```console
Usage: c14 download [DIR|FILE]* ARCHIVE

Download your file or directory into an archive, use SFTP protocol.

Options:
  -h, --help=false      Print usage

Examples:
        $ c14 download
        $ c14 download file 83b93179-32e0-11e6-be10-10604b9b0ad9
```

#### `c14 bucket`

```console
Usage: c14 bucket [OPTIONS] [ARCHIVE]*

Displays (JSON or tab output) all information of bucket.

Options:
  -h, --help=false      Print usage
  -p, --pretty=""       Show all information in tab (default json output)

Examples:
        $ c14 bucket
        $ c14 bucket 83b93179-32e0-11e6-be10-10604b9b0ad9
        $ c14 bucket -p 83b93179-32e0-11e6-be10-10604b9b0ad9
```

#### `c14 version`

```console
Usage: c14 version

Show the version information.

Options:
  -h, --help=false      Print usage

Examples:
        $ c14 version
```



### Examples

Soon

---

## Changelog

### master (unreleased)

 * Support of `verify` command
 * Support of `create` command
 * Support of `freeze` command
 * Support of `help` command
 * Support of `ls` command
 * Support of `login` command
 * Support of `files` command
 * Support of `rename` command
 * Support of `remove` command
 * Support of `unfreeze` command
 * Support of `upload` command
 * Support of `version` command
 * Support of `download` command

---

## Development

Feel free to contribute :smiley::beers:


### Hack

1. [Install go](https://golang.org/doc/install)
2. Ensure you have `$GOPATH` and `$PATH` well configured, something like:
  * `export GOPATH=$HOME/go`
  * `export PATH=$PATH:$GOPATH/bin`
3. Fetch the project: `go get -u github.com/scaleway/c14-cli`
4. Go to c14-cli directory: `cd $GOPATH/src/github.com/scaleway/c14-cli`
5. Hack: `vim`
6. Build: `make`
7. Run: `./c14`

## License

[MIT](https://github.com/scaleway/c14-cli/blob/master/LICENSE.md)

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
