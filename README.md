# Affectum

![affectum](./assets/affectum.png)

> Simple application to periodically save email attachments. You enter your mail configuration and set a folder, then affectum will scan that folder every minute and save any new mail attachments.

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

Copyright (c) 2021 nkoporec
Licensed under the MIT license.

