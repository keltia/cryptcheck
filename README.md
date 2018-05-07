imirhil-go
============

[![Build Status](https://travis-ci.org/keltia/imirhil-go.svg?branch=master)](https://travis-ci.org/keltia/imirhil-go)
[![GoDoc](http://godoc.org/github.com/keltia/imirhil-go?status.svg)](http://godoc.org/github.com/keltia/imirhil-go)

Go wrapper for [Imirhil/cryptcheck](https://tls.imirhil.fr/) API.  Currently v1 of the API is supported, v2 is not released or dicumented yet.

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

## License

The [BSD 2-Clause license][bsd].

# Feedback

We welcome pull requests, bug fixes and issue reports.

Before proposing a large change, first please discuss your change by raising an issue.
