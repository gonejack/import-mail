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
Flags:
  -h, --help                  Show context-sensitive help.
      --host=STRING           Set IMAP host.
      --port=993              Set IMAP port.
      --username=STRING       Set IMAP username.
      --password=STRING       Set IMAP password.
      --remote-dir="INBOX"    Set IMAP directory.
      --size-limit="20M"      Set size limit, mail exceed this limit will be skipped.
      --about                 Show about.
```
