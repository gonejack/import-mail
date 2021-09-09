package main

import (
	"log"
	"os"

	"github.com/gonejack/import-mail/cmd"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	c := cmd.ImportMail{
		SaveImportedTo: "imported",
	}
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
