package main

import (
	"log"
	"os"

	"github.com/gonejack/import-mail/importmail"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	cmd := importmail.Import{
		Options: importmail.MustParseOptions(),
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
