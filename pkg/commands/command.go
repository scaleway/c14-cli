package commands

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/pkg/mflag"
)

type Config struct {
	UsageLine   string
	Description string
	Help        string
	Examples    string
}

type Streams struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

type Env struct {
	Streams
	Debug bool
}

type Command interface {
	GetBase() *Base
	GetName() string
	Parse(args []string) ([]string, error)
	Run(args []string) error
}

type Base struct {
	Env
	Config
	Flags  mflag.FlagSet
	flHelp *bool
}

func (b *Base) Init(c Config) {
	b.Config = c
	b.Streams.Stdout = os.Stdout
	b.Streams.Stdin = os.Stdin
	b.Streams.Stderr = os.Stderr
	b.Flags.SetOutput(ioutil.Discard)
	b.flHelp = b.Flags.Bool([]string{"h", "-help"}, false, "Print usage")
}

func (b *Base) Parse(args []string) (newArgs []string, err error) {
	if err = b.Flags.Parse(args); err != nil {
		b.Flags.PrintDefaults()
		return
	}
	if *b.flHelp {
		err = errors.New("TODO: return help message")
		return
	}
	newArgs = b.Flags.Args()
	return
}

func (b *Base) GetBase() *Base {
	return b
}

func (b *Base) PrintUsage() (err error) {
	return
}
