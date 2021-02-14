<h1 align="center">
  <br>
  <a href="#"><img src="https://user-images.githubusercontent.com/25837540/107883163-dbf70380-6f1f-11eb-856f-e69e733313e5.png" width="400px" alt="mubeng"></a>
</h1>

<h4 align="center">An proxy checker, and HTTP proxy server as IP rotator, all applicable HTTP/S & SOCKS5 protocols with ease.</h4>

<p align="center">
	<a href="#"><img src="https://img.shields.io/badge/kitabisa-security%20project-blue"></a>
	<a href="https://golang.org"><img src="https://img.shields.io/badge/made%20with-Go-brightgreen"></a>
	<!-- <a href="https://goreportcard.com/report/ktbs.dev/mubeng"><img src="https://goreportcard.com/badge/ktbs.dev/mubeng"></a> -->
	<a href="https://github.com/kitabisa/mubeng/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-yellowgreen"></a>
	<a href="#"><img src="https://img.shields.io/badge/platform-osx%2Flinux%2Fwindows-green"></a>
	<!-- <a href="https://github.com/kitabisa/mubeng/releases"><img src="https://img.shields.io/github/release/kitabisa/mubeng"></a> -->
	<!-- <a href="https://github.com/kitabisa/mubeng/issues"><img src="https://img.shields.io/github/issues/kitabisa/mubeng"></a> -->
</p>

<p align="center">
  <a href="https://github.com/kitabisa/mubeng/blob/master/.github/CONTRIBUTING.md">Contributing</a> •
  <a href="https://github.com/kitabisa/mubeng/blob/master/CHANGELOG.md">What's new</a> •
  <a href="https://pkg.go.dev/ktbs.dev/mubeng">Documentation</a> •
  <a href="https://github.com/kitabisa/mubeng/issues/new/choose">Report Issues</a>
</p>

---

# Table of Contents

_TODO_

# Installation

## Binary

Simply, download a pre-built binary from [releases page](https://github.com/kitabisa/mubeng/releases), unpack and run!

## Docker

Pull the [Docker](https://docs.docker.com/get-docker/) image by running:

```bash
▶ docker pull kitabisa/mubeng
```

## Source

Using [Go _(v1.15+)_](https://golang.org/doc/install) compiler:

```bash
▶ GO111MODULE=on go get -v -u ktbs.dev/mubeng/cmd/mubeng
```

> **NOTE:** The same command above also works for updating.

### — or

Manual building executable from source code:

```bash
▶ git clone https://github.com/kitabisa/mubeng
▶ cd mubeng
▶ make build
▶ (sudo) mv ./bin/mubeng /usr/local/bin
▶ make clean
```

# Usage

For usage, it's always required to provide your proxy list, whether it is used to check or run as a proxy server for your proxy IP rotation.

## Basic

```bash
▶ mubeng [-c|-a :8080] -f file.txt [options...]
```

## Options

Here are all the options it supports.

| **Flag**                      	| **Description**                                              	|
|-------------------------------	|--------------------------------------------------------------	|
| -f, --file `<FILE>`           	| Proxy file.                                                  	|
| -a, --address `<ADDR>:<PORT>` 	| Run proxy server.                                            	|
| -c, --check                   	| To perform proxy live check.                                 	|
| -t, --timeout                 	| Max. time allowed for proxy server/live check (default: 30s). |
| -r, --rotate `<AFTER>`        	| Rotate proxy IP for every `AFTER` request.                    |
| -v, --verbose                 	| Dump HTTP request/responses show died proxy checks.          	|
| -o, --output <FILE>           	| Log output from proxy servers or live proxy checks.          	|


## Examples

`mubeng` has 2 core functionality:

- Perform proxy checks, or
- Run proxy server as proxy IP rotation.

For example, you've proxy list `(proxies.txt)` as:

<table>
	<td>
		<pre>http://127.0.0.1:8080
https://127.0.0.1:3128
socks5://127.0.0.1:2121
...
...</pre>
	</td>
</table>

> Because we use auto-switch transport, `mubeng` can accept multiple proxy protocol schemes at once.<br>
> Please refer to [documentation](https://pkg.go.dev/ktbs.dev/mubeng/pkg/mubeng#Transport) for this package.

### Proxy checker

Pass `--check` flag to command to perform proxy checks:

```bash
▶ mubeng -f proxies.txt --check --output live.txt
```

The above case also uses `--output` flag to save a live proxy into file `(live.txt)` from checking result.

### Proxy IP rotator

Furthermore, if you wish to do IP rotator from proxies that are still alive earlier from the results of checking `(live.txt)`, you must use `-a` _(--address)_ flag instead of `-c` _(--check)_ to run proxy server:

```bash
▶ mubeng -a localhost:8080 -f live.txt -r 10
```

The `-r` _(--rotate)_ flag works to rotate your IP for every _N_ request value you provide `(10)`.

# Limitations

Currently IP rotation runs the proxy server only as an HTTP protocol, not a SOCKSv5 protocol, even though the resource you have is SOCKSv5. In other words, the SOCKSv5 resource that you provide is used properly because it uses auto-switch transport on the client, but this proxy server **DOES NOT** switch to anything other than HTTP protocol.

# TODOs

- [ ] TLS support. _Built-in certs maker?_
- [ ] Rotate IP proxy `AFTER` request.
- [ ] `mubeng` proxy server as service, daemonize it!

## Pronunciation

_TODO_

## Changes

For changes, see [CHANGELOG.md](https://github.com/kitabisa/mubeng/blob/master/CHANGELOG.md).

## License

This program is free software: you can redistribute it and/or modify it under the terms of the Apache license. Kitabisa mubeng and any contributions are Copyright © by Dwi Siswanto 2021.