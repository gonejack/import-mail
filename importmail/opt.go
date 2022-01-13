package importmail

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
)

type about bool

func (a about) BeforeApply() (err error) {
	fmt.Println("Visit https://github.com/gonejack/import-mail")
	os.Exit(0)
	return
}

type Options struct {
	Host      string `required:"" help:"Set IMAP host."`
	Port      int    `default:"993" help:"Set IMAP port."`
	Username  string `required:"" help:"Set IMAP username."`
	Password  string `required:"" help:"Set IMAP password."`
	RemoteDir string `name:"remote-dir" default:"INBOX" help:"Set IMAP directory."`
	SizeLimit string `name:"size-limit" default:"20M" help:"Set size limit, mail exceed this limit will be skipped."`
	About     about  `help:"Show about."`

	SaveImportedTo string `hidden:"" default:"imported"`

	Eml []string `name:".eml" arg:"" optional:"" help:"list of .eml files"`
}

func MustParseOptions() (opt Options) {
	kong.Parse(&opt,
		kong.Name("import-mail"),
		kong.Description("This command line imports .eml files into IMAP account."),
		kong.UsageOnError(),
	)
	if len(opt.Eml) == 0 || opt.Eml[0] == "*.eml" {
		opt.Eml, _ = filepath.Glob("*.eml")
	}
	return
}
