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

## Requirements

* Go >= 1.10

## USAGE

There is a small example program included in `cmd/getgrade` to either show the grade of a given site or JSON dump of the detailed report.

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
