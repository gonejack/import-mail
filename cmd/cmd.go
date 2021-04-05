package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/emersion/go-imap/client"
)

var zeroTime time.Time

type ImportMail struct {
	Host        string
	Port        int
	Username    string
	Password    string
	ImportedDir string

	client *client.Client
}

func (c *ImportMail) Execute(emails []string) (err error) {
	if len(emails) == 0 {
		emails, _ = filepath.Glob("*.eml")
	}
	if len(emails) == 0 {
		return errors.New("no eml given")
	}

	err = c.mkdir()
	if err != nil {
		return
	}

	err = c.connect()
	if err != nil {
		return
	}
	defer c.disconnect()

	for _, eml := range emails {
		log.Printf("procssing %s", eml)

		mail, err := os.Open(eml)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		scan := bufio.NewScanner(mail)
		for scan.Scan() {
			buf.WriteString(scan.Text())
			buf.WriteString("\r\n")
		}
		err = scan.Err()
		if err != nil {
			return err
		}

		err = c.client.Append("INBOX", nil, zeroTime, &buf)
		if err != nil {
			return err
		}

		err = os.Rename(eml, filepath.Join(c.ImportedDir, filepath.Base(eml)))
		if err != nil {
			return err
		}
	}

	return
}

func (c *ImportMail) mkdir() error {
	err := os.MkdirAll(c.ImportedDir, 0777)
	if err != nil {
		return fmt.Errorf("cannot make images dir %s", err)
	}

	return nil
}
func (c *ImportMail) connect() (err error) {
	c.client, err = client.DialTLS(fmt.Sprintf("%s:%d", c.Host, c.Port), nil)
	if err == nil {
		err = c.client.Login(c.Username, c.Password)
	}
	return
}
func (c *ImportMail) disconnect() (err error) {
	return c.client.Logout()
}
