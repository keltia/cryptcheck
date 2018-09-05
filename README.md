cryptcheck
==========

[![GitHub release](https://img.shields.io/github/release/keltia/cryptcheck.svg)](https://github.com/keltia/cryptcheck/releases) 
[![GitHub issues](https://img.shields.io/github/issues/keltia/cryptcheck.svg)](https://github.com/keltia/cryptcheck/issues)
[![Go Version](https://img.shields.io/badge/go-1.10-blue.svg)](https://golang.org/dl/)
[![Build Status](https://travis-ci.org/keltia/cryptcheck.svg?branch=master)](https://travis-ci.org/keltia/cryptcheck)
[![GoDoc](http://godoc.org/github.com/keltia/cryptcheck?status.svg)](http://godoc.org/github.com/keltia/cryptcheck)
[![SemVer](http://img.shields.io/SemVer/2.0.0.png)](https://semver.org/spec/v2.0.0.html)
[![License](https://img.shields.io/pypi/l/Django.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/keltia/cryptcheck)](https://goreportcard.com/report/github.com/keltia/cryptcheck)

Go wrapper for [Imirhil/cryptcheck](https://tls.imirhil.fr/) API.  Currently v1 of the API is supported, v2 is not released or documented yet.

API v1 is now at 201809, added missing Error field in Host.

## Requirements

* Go >= 1.10

## Installation

You need to install my `proxy` module before if you are using Go 1.10.x or earlier.

    go get github.com/keltia/proxy

With Go 1.11+ and its modules support, it should work out of the box with

    go get github.com/keltia/cryptcheck/cmd/...

if you have the `GO111MODULE` environment variable set on `on`.

## USAGE

There is a small example program included in `cmd/cryptcheck` to either show the grade of a given site or JSON dump of the detailed report.

You can just get the grade like this:

    $ cryptcheck www.ssllabs.com
    cryptcheck Wrapper: 1.4.0 API version 201809
    
    Grade for 'www.ssllabs.com' is B (Date: 2018-07-30 23:52:52.494 +0200 CEST)

You can get a more detail report with `-d`:

    $ cryptcheck -d www.ssllabs.com
    cryptcheck Wrapper: 1.4.0 API version 201809
    
    {"Hosts":[{"host":{"Name":"www.ssllabs.com","ip":"64.41.200.100","Port":443},"handshake":{"Key":{"type":"rsa","size":20
    [...]

You can use `jq` to display the output of `cryptcheck -d <site>` in a colorised way (use `-raw` to remove the banner display):

    cryptcheck -raw tls.imirhil.fr | jq .

There is also a debug mode with `-D`.

By default, Cryptcheck returns the last run cached by the site, if you want to refresh, use `-R`.

## API Usage

As with many API wrappers, you will need to first create a client with some optional configuration, then there are two main functions:

``` go
    // Simplest way
    c := cryptcheck.NewClient()
    grade, err := c.GetScore("example.com")
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    
    
    // With some options, timeout at 15s and debug-like verbosity
    cnf := cryptcheck.Config{
        Timeout:15, 
        Log:2,
    }
    c := cryptcheck.NewClient(cnf)
    report, err := c.GetDetailedReport("foo.xxx")
    if err != nil {
        log.Fatalf("error: %v", err)
    }
```

OPTIONS

| Option  | Type | Description |
| ------- | ---- | ----------- |
| Timeout | int  | time for connections (default: 10s ) |
| Log     | int  | 1: verbose, 2: debug (default: 0) |
| Refresh | bool | Force refresh of the sites (default: false) |
    

## Using behind a web Proxy

Dependency: proxy support is provided by my `github.com/keltia/proxy` module.

UNIX/Linux:

```
    export HTTP_PROXY=[http://]host[:port] (sh/bash/zsh)
    setenv HTTP_PROXY [http://]host[:port] (csh/tcsh)
```

Windows:

```
    set HTTP_PROXY=[http://]host[:port]
```

The rules of Go's `ProxyFromEnvironment` apply (`HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`, lowercase variants allowed).

If your proxy requires you to authenticate, please create a file named `.netrc` in your HOME directory with permissions either `0400` or `0600` with the following data:

    machine proxy user <username> password <password>
    
and it should be picked up. On Windows, the file will be located at

    %LOCALAPPDATA%\cryptcheck\netrc

## License

The [BSD 2-Clause license](https://github.com/keltia/cryptcheck/LICENSE.md).

# Contributing

This project is an open Open Source project, please read `CONTRIBUTING.md`.

# Feedback

We welcome pull requests, bug fixes and issue reports.

Before proposing a large change, first please discuss your change by raising an issue.
