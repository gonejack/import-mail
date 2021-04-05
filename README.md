# import-mail
This command line tool imports .eml files into INBOX of IMAP account.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gonejack/import-mail)
![Build](https://github.com/gonejack/import-mail/actions/workflows/go.yml/badge.svg)
[![GitHub license](https://img.shields.io/github/license/gonejack/import-mail.svg?color=blue)](LICENSE)

### Install
```shell
> go get github.com/gonejack/import-mail
```

### Usage
```shell
> import-mail --host imap.example.com --username username --password password *.eml
```
```
Command line tool for importing .eml files to IMAP account.

Usage:
  import-mail *.eml [flags]

Flags:
      --host string       host
      --port int          port (default 993)
      --username string   username
      --password string   password
  -v, --verbose           verbose
  -h, --help              help for import-mail
```
