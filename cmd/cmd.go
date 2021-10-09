package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
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
	About     bool   `help:"Show about."`

	Eml []string `arg:"" optional:""`
}

type Import struct {
	options
	SaveImportedTo string

	limit  int
	buf    bytes.Buffer
	client *client.Client
}

func (c *Import) Run() (err error) {
	kong.Parse(&c.options,
		kong.Name("import-mail"),
		kong.Description("Command line tool for importing .eml files to IMAP account."),
		kong.UsageOnError(),
	)
	if c.About {
		fmt.Println("Visit https://github.com/gonejack/import-mail")
		return
	}
	if len(c.Eml) == 0 {
		c.Eml, _ = filepath.Glob("*.eml")
	}
	if len(c.Eml) == 0 {
		return errors.New("no .eml file found")
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
	c.limit = int(localLimit)

	remoteLimit, err := c.queryAppendLimit()
	if err == nil && remoteLimit != 0 {
		c.limit = humanize.IByte * int(remoteLimit)
	}
	log.Printf("APPENDLIMIT is %s", humanize.Bytes(uint64(c.limit)))

	return c.doAppend()
}
func (c *Import) doAppend() error {
	for _, eml := range c.Eml {
		log.Printf("process %s", eml)

		mail, err := os.Open(eml)
		if err != nil {
			return err
		}

		if c.limit > 0 {
			stat, err := mail.Stat()
			if err != nil {
				return err
			}
			size := int(stat.Size())
			if size > c.limit {
				log.Printf("skipped, %s's size %s is larger than APPENDLIMIT %s", eml, humanize.Bytes(uint64(size)), humanize.Bytes(uint64(c.limit)))
				continue
			}
		}

		err = c.doAppendOne(mail)
		if err != nil {
			return err
		}

		_ = mail.Close()
		err = os.Rename(eml, filepath.Join(c.SaveImportedTo, filepath.Base(eml)))
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Import) doAppendOne(mail io.Reader) (err error) {
	defer c.buf.Reset()

	scan := bufio.NewScanner(mail)
	for scan.Scan() {
		c.buf.WriteString(scan.Text())
		c.buf.WriteString("\r\n")
	}
	err = scan.Err()
	if err != nil {
		return err
	}

	return c.client.Append(c.RemoteDir, nil, time.Time{}, &c.buf)
}
func (c *Import) queryAppendLimit() (size uint32, err error) {
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
func (c *Import) connect() (err error) {
	c.client, err = client.DialTLS(fmt.Sprintf("%s:%d", c.Host, c.Port), nil)
	if err == nil {
		err = c.client.Login(c.Username, c.Password)
	}
	return
}
func (c *Import) disconnect() (err error) {
	return c.client.Logout()
}
