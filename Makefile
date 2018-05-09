
# Main Makefile for cryptcheck
#
# Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
#

GOBIN=   ${GOPATH}/bin

SRCS= imirhil.go proxy.go types.go utils.go

USRCS=	config_unix.go
WSRCS=	config_windows.go

OPTS=	-ldflags="-s -w" -v

all: build

build: ${SRCS} ${USRCS}
	go build ${OPTS}

test: build
	go test

windows: ${SRCS} ${WSRCS}
	GOOS=windows go build ${OPTS} .

install:
	go install ${OPTS}

lint:
	gometalinter .

clean:
	go clean

push:
	git push --all
	git push --tags
