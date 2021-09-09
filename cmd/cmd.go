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

	"github.com/alecthomas/kong"
	"github.com/dustin/go-humanize"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap-appendlimit"
	"github.com/emersion/go-imap/client"
)

type options struct {
	Host      string `required:"" help:"Set IMAP host."`
	Port      int    `default:"993" help:"Set IMAP port."`
	Username  string `required:"" help:"Set IMAP username."`
	Password  string `required:"" help:"Set IMAP password."`
	RemoteDir string `name:"remote-dir" default:"INBOX" help:"Set IMAP directory."`
	SizeLimit string `name:"size-limit" default:"20M" help:"Set size limit, mail exceed this limit will be skipped."`

	Eml []string `arg:"" optional:""`
}

type ImportMail struct {
	options
	SaveImportedTo string
	sizeLimit      int
	client         *client.Client
}

func (c *ImportMail) Run() (err error) {
	kong.Parse(&c.options,
		kong.Name("import-mail"),
		kong.Description("Command line tool for importing .eml files to IMAP account."),
		kong.UsageOnError(),
	)

	if len(c.Eml) == 0 {
		c.Eml, _ = filepath.Glob("*.eml")
	}
	if len(c.Eml) == 0 {
		return errors.New("not .eml file found")
	}

	err = os.MkdirAll(c.SaveImportedTo, 0777)
	if err != nil {
		return fmt.Errorf("cannot make dir %s: %s", c.SaveImportedTo, err)
	}

	err = c.connect()
	if err != nil {
		return
	}
	defer c.disconnect()

	localLimit, err := humanize.ParseBytes(c.SizeLimit)
	if err != nil {
		return
	}
	c.sizeLimit = int(localLimit)

	remoteLimit, err := c.queryAppendLimit()
	if err != nil && remoteLimit != 0 {
		c.sizeLimit = humanize.IByte * int(remoteLimit)
	}
	log.Printf("APPENDLIMIT is %s", humanize.Bytes(uint64(c.sizeLimit)))

	return c.doAppend(c.Eml)
}
func (c *ImportMail) doAppend(emails []string) error {
	for _, eml := range emails {
		log.Printf("process %s", eml)

		mail, err := os.Open(eml)
		if err != nil {
			return err
		}

		if c.sizeLimit > 0 {
			stat, err := mail.Stat()
			if err != nil {
				return err
			}
			if size := int(stat.Size()); size > c.sizeLimit {
				log.Printf("skipped %s's size %s larger than APPENDLIMIT %s", eml, humanize.Bytes(uint64(size)), humanize.Bytes(uint64(c.sizeLimit)))
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
		err = c.client.Append(c.RemoteDir, nil, time.Time{}, &buf)
		if err != nil {
			return err
		}
		err = os.Rename(eml, filepath.Join(c.SaveImportedTo, filepath.Base(eml)))
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *ImportMail) queryAppendLimit() (size uint32, err error) {
	status, err := c.client.Status(c.RemoteDir, []imap.StatusItem{appendlimit.Capability})
	if err != nil {
		return
	}
	val := status.Items[appendlimit.StatusAppendLimit]
	if val == nil {
		return
	}
	return imap.ParseNumber(val)
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
