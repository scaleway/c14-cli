package commands

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

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
		err = fmt.Errorf("usage: c14 %v", b.UsageLine)
		return
	}
	if *b.flHelp {
		b.PrintUsage()
		os.Exit(1)
		return
	}
	newArgs = b.Flags.Args()
	return
}

func (b *Base) GetBase() *Base {
	return b
}

func (b *Base) PrintUsage() {
	var usageTemplate = `Usage: c14 {{.UsageLine}}

{{.Help}}

{{.Options}}
{{.ExamplesHelp}}
`

	t := template.New("full")
	template.Must(t.Parse(usageTemplate))
	t.Execute(os.Stdout, b)
}

func (b *Base) Options() string {
	var options string

	visitor := func(flag *mflag.Flag) {
		var optionUsage string

		name := strings.Join(flag.Names, ", -")
		if flag.DefValue == "" {
			optionUsage = fmt.Sprintf("%s=\"\"", name)
		} else {
			optionUsage = fmt.Sprintf("%s=%s", name, flag.DefValue)
		}
		options += fmt.Sprintf("  -%-20s %s\n", optionUsage, flag.Usage)
	}
	b.Flags.VisitAll(visitor)
	if len(options) == 0 {
		return ""
	}
	return fmt.Sprintf("Options:\n%s", options)
}

func (b *Base) ExamplesHelp() string {
	if b.Examples == "" {
		return ""
	}
	return fmt.Sprintf("Examples:\n%s", strings.Trim(b.Examples, "\n"))
}
