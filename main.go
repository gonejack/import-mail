package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gonejack/import-mail/cmd"
	"github.com/spf13/cobra"
)

var (
	host     string
	port     int
	username string
	password string

	prog = &cobra.Command{
		Use:   "import-mail *.eml",
		Short: "Command line tool for importing .eml files to IMAP account.",
		Run: func(c *cobra.Command, args []string) {
			err := run(c, args)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	log.SetOutput(os.Stdout)

	prog.Flags().SortFlags = false
	prog.PersistentFlags().SortFlags = false

	flags := prog.PersistentFlags()
	{
		flags.StringVarP(&host, "host", "", "", "host")
		flags.IntVarP(&port, "port", "", 993, "port")
		flags.StringVarP(&username, "username", "", "", "username")
		flags.StringVarP(&password, "password", "", "", "password")
	}
}

func run(c *cobra.Command, args []string) error {
	switch "" {
	case host:
		return fmt.Errorf("argument --host is required")
	case username:
		return fmt.Errorf("argument --username is required")
	case password:
		return fmt.Errorf("argument --password is required")
	}

	exec := cmd.ImportMail{
		Host:        host,
		Port:        port,
		Username:    username,
		Password:    password,
		ImportedDir: "imported",
	}
	return exec.Execute(args)
}

func main() {
	_ = prog.Execute()
}