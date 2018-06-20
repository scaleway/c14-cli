package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"golang.org/x/oauth2"

	"github.com/cocooma/mflag"
	"github.com/scaleway/c14-cli/pkg/api"
	"github.com/scaleway/c14-cli/pkg/api/auth"
	"github.com/scaleway/c14-cli/pkg/version"
)

// Config represents the informations on the usages
type Config struct {
	UsageLine   string
	Description string
	Help        string
	Examples    string
}

// Streams allows to overload the output and input
type Streams struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

// Command is the interface that is used to handle the commands
type Command interface {
	GetName() string
	Parse(args []string) ([]string, error)
	Run(args []string) error
	CheckFlags(args []string) error
	PrintUsage()
}

// Base must be embedded in the commands
type Base struct {
	Config
	Flags mflag.FlagSet
	*api.OnlineAPI
	flHelp *bool
}

// Init initiates the Base structure
func (b *Base) Init(c Config) {
	b.Config = c
	b.Flags.SetOutput(ioutil.Discard)
	b.flHelp = b.Flags.Bool([]string{"h", "-help"}, false, "Print usage")
}

// InitAPI initiates the Online API with the credentials
func (b *Base) InitAPI() (err error) {
	var (
		c            *auth.Credentials
		privateToken string
	)

	if privateToken = os.Getenv("C14_PRIVATE_TOKEN"); privateToken != "" {
		c = &auth.Credentials{
			AccessToken: privateToken,
		}
	} else {
		if c, err = auth.GetCredentials(); err != nil {
			return
		}
	}
	b.OnlineAPI = api.NewC14API(oauth2.NewClient(oauth2.NoContext, c), version.UserAgent, Root.Verbose)
	return
}

// Parse parses the argurments
func (b *Base) Parse(args []string) (newArgs []string, err error) {
	if err = b.Flags.Parse(args); err != nil {
		err = fmt.Errorf("usage: c14 %s: %s", b.UsageLine, err)
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

// CheckFlags can be overloaded by the commands to check the flags
func (b *Base) CheckFlags(args []string) (err error) {
	return
}

// PrintUsage print on Stdout the usage message
func (b *Base) PrintUsage() {
	var usageTemplate = `Usage: c14 {{.UsageLine}}

{{.Help}}

{{.Options}}
{{.ExamplesHelp}}
`

	t := template.New("full")
	template.Must(t.Parse(usageTemplate))
	_ = t.Execute(os.Stdout, b)
}

// Options returns the options available, it used by PrintUsage
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

// ExamplesHelp returns the examples, it used by PrintUsage
func (b *Base) ExamplesHelp() string {
	if b.Examples == "" {
		return ""
	}
	return fmt.Sprintf("Examples:\n%s", strings.Trim(b.Examples, "\n"))
}
