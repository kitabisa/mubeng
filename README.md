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
  <a href="#contributors">Contributing</a> •
  <a href="#changes">What's new</a> •
  <a href="https://pkg.go.dev/ktbs.dev/mubeng">Documentation</a> •
  <a href="https://github.com/kitabisa/mubeng/issues/new/choose">Report Issues</a>
</p>

---

- [Features](#features)
- [Why mubeng?](#why-mubeng)
- [Demo](#demo)
- [Installation](#installation)
  - [Binary](#binary)
  - [Docker](#docker)
  - [Source](#source)
- [Usage](#usage)
  - [Basic](#basic)
  - [Options](#options)
  - [Examples](#examples)
    - [Proxy checker](#proxy-checker)
    - [Proxy IP rotator](#proxy-ip-rotator)
- [Limitations](#limitations)
- [Roadmap](#roadmap)
- [Contributors](#contributors)
- [Pronunciation](#pronunciation)
- [Changes](#changes)
- [License](#license)

---

# Features

- **Proxy IP rotator**: Rotates your IP address for every specific request.
- **Proxy checker**: Check your proxy IP which is still alive.
- **All HTTP methods** are supported.
- **HTTP/S & SOCKSv5** proxy protocols apply.
- **All parameters & URIs are passed**.
- **Minimal configuration**: You can just run it against your proxy file, and choose the action you want!

# Why mubeng?

`mubeng` has 2 core functionality:

- Run proxy server as proxy IP rotation, or
- Perform proxy checks.

So, you don't need any extra proxy checking tools out there if you want to do that, and we also leave it entirely up to user to provide the proxy resources, therefore it's **NOT** under any policy for any use of this tool, but developers assume no liability and are not responsible for any misuse or damage. Be responsible for your actions!

# Demo

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

```bash
▶ mubeng -h
```

| **Flag**                      	| **Description**                                              	|
|-------------------------------	|--------------------------------------------------------------	|
| -f, --file `<FILE>`           	| Proxy file.                                                  	|
| -a, --address `<ADDR>:<PORT>` 	| Run proxy server.                                            	|
| -c, --check                   	| To perform proxy live check.                                 	|
| -t, --timeout                 	| Max. time allowed for proxy server/live check (default: 30s). |
| -r, --rotate `<AFTER>`        	| Rotate proxy IP for every `AFTER` request.                    |
| -v, --verbose                 	| Dump HTTP request/responses or show died proxy checks.        |
| -o, --output <FILE>           	| Log output from proxy servers or live proxy checks.          	|


## Examples

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

Furthermore, if you wish to do IP rotator from proxies that are still alive earlier from the results of checking `(live.txt)` or if you have your own list, you must use `-a` _(--address)_ flag instead of `-c` _(--check)_ to run proxy server:

```bash
▶ mubeng -a localhost:8080 -f live.txt -r 10
```

The `-r` _(--rotate)_ flag works to rotate your IP for every _N_ request value you provide `(10)`.

# Limitations

Currently IP rotation runs the proxy server only as an HTTP protocol, not a SOCKSv5 protocol, even though the resource you have is SOCKSv5. In other words, the SOCKSv5 resource that you provide is used properly because it uses auto-switch transport on the client, but this proxy server **DOES NOT** switch to anything other than HTTP protocol.

# Roadmap

- [ ] Support HTTPS protocol for proxy server.
- [ ] Rotate IP proxy `AFTER` request.
- [ ] `mubeng` proxy server as service, daemonize it!

# Contributors

[![contributions](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/kitabisa/mubeng/issues)

This project exists thanks to all the people who contribute. To learn how to setup a development environment and for contribution guidelines, see [CONTRIBUTING.md](https://github.com/kitabisa/mubeng/blob/master/.github/CONTRIBUTING.md).

<a href="https://github.com/kitabisa/mubeng/graphs/contributors">
	<img src=".github/CONTRIBUTORS.svg">
</a>

# Pronunciation

`jv_ID` • **/mo͞oˌbēNG/** — mubeng nganti mumet. (ꦩꦸꦧꦺꦁ​ꦔꦤ꧀ꦠꦶ​ꦩꦸꦩꦺꦠ꧀)

# Changes

For changes, see [CHANGELOG.md](https://github.com/kitabisa/mubeng/blob/master/CHANGELOG.md).

# License

This program is free software: you can redistribute it and/or modify it under the terms of the [Apache license](https://github.com/kitabisa/mubeng/blob/master/LICENSE). Kitabisa mubeng and any contributions are Copyright © by Dwi Siswanto 2021.