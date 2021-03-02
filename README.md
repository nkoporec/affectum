# Affectum ![Affectum](./assets/affectum.png)

> Simple application to periodically save email attachments.

## About

Affectum is an application that will periodically scan your email folder and save any email attachments (files, images ...). On startup the application will scan the folder and save any existing attachments, then every minute a re-scan will occur and save any new attachments if found.

This is useful if you don't want to manually save your attachments or if you want to automate your email workflow.

It implements a local IMAP client which will connect to your email provider and scan your emails. It runs locally, so it's private, safe, and secure

## Install/run

### macOS

Install from [.dmg](https://github.com/nkoporec/affectum/releases) and run as any other application.

### Windows

Just run `affectum.exe`.

### Linux

```bash
go get github.com/nkoporec/affectum
nohup affectum >/dev/null 2>&1 &
```

## Configuration

The configuration file is located at ~/.affectum/affectum.env

```env
# Email host
# eq: imap.gmail.com
MAIL_HOST=

# Email port
# eq: 993 (for gmail.com)
MAIL_PORT=

# Email username
# eq: "user"
MAIL_USERNAME=

# Email password
# eq: "password"
MAIL_PASSWORD=

# Email folder
# eq: "Bills"
MAIL_FOLDER=

# Full path to a folder where the attachments will be saved.
# Defaults to ~/.affectum/files/
ATTACHMENT_FOLDER_PATH=

# Connect with STARTTLS.
# Defaults to false.
STARTTLS= false
```

Make sure you have this file created, before running Affectum. If the file is not available or the configuration is not correct the application will not start.

### Debugging

If for some reason the Affectum is not working, please check the log located at ~/.affectum/affectum.log and see if there are any errors.

## Local development

### Build command

```bash
make build-darwin # only in macOS
make build-win
make build-linux # only on linux
```

**Prerequisites**

The project uses these major dependencies and inherits their prerequisites:

- [systray](https://github.com/getlantern/systray)

Due to the `systray` package, the build for macOS can be done only in Mac, a linux build only on a Linux machine. Platform specific prerequisites required.

## Credits ##

Affectum icon is part of [Streamline](https://streamlinehq.com) icon pack.

## License ##

Copyright (c) 2021 nkoporec <br>
Licensed under the MIT license.

