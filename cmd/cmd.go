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

	"github.com/dustin/go-humanize"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap-appendlimit"
	"github.com/emersion/go-imap/client"
)

var appendLimitSize int
var zeroTime time.Time

type ImportMail struct {
	Host           string
	Port           int
	Username       string
	Password       string
	ImportedDir    string
	ArgAppendLimit string

	client *client.Client
}

func (c *ImportMail) Execute(emails []string) (err error) {
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

	limit, err := c.getAppendLimit()
	if err != nil {
		return
	}
	appendLimitSize = humanize.IByte * int(limit)
	if appendLimitSize == 0 {
		parsed, perr := humanize.ParseBytes(c.ArgAppendLimit)
		if perr != nil {
			return perr
		}
		appendLimitSize = int(parsed)
	}
	log.Printf("APPENDLIMIT is %s", humanize.Bytes(uint64(appendLimitSize)))

	return c.appendMails(emails)
}

func (c *ImportMail) appendMails(emails []string) error {
	for _, eml := range emails {
		log.Printf("procssing %s", eml)

		mail, err := os.Open(eml)
		if err != nil {
			return err
		}

		if appendLimitSize > 0 {
			stat, err := mail.Stat()
			if err != nil {
				return err
			}
			if size := int(stat.Size()); size > appendLimitSize {
				log.Printf("mail size %s larger than %s", humanize.Bytes(uint64(size)), humanize.Bytes(uint64(appendLimitSize)))
				continue
			}
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

	return nil
}

func (c *ImportMail) getAppendLimit() (size uint32, err error) {
	status, err := c.client.Status("INBOX", []imap.StatusItem{appendlimit.Capability})
	if err != nil {
		return
	}
	val := status.Items[appendlimit.StatusAppendLimit]
	if val == nil {
		return
	}
	return imap.ParseNumber(val)
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
