imirhil-go
============

[![Build Status](https://travis-ci.org/keltia/imirhil-go.svg?branch=master)](https://travis-ci.org/keltia/imirhil-go)
[![GoDoc](http://godoc.org/github.com/keltia/imirhil-go?status.svg)](http://godoc.org/github.com/keltia/imirhil-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/keltia/imirhil-go)](https://goreportcard.com/report/github.com/keltia/imirhil-go)

Go wrapper for [Imirhil/cryptcheck](https://tls.imirhil.fr/) API.  Currently v1 of the API is supported, v2 is not released or documented yet.

## Requirements

* Go >= 1.10

## API Usage


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
    
and it should be picked up.

## License

The [BSD 2-Clause license][bsd].

# Contributing

This project is an open Open Source project, please read `CONTRIBUTING.md`.

# Feedback

We welcome pull requests, bug fixes and issue reports.

Before proposing a large change, first please discuss your change by raising an issue.
