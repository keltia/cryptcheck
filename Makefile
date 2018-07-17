
# Main Makefile for cryptcheck
#
# Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
#

.PATH=	cmd/getgrade:.
GOBIN=	${GOPATH}/bin

GSRCS=	cmd/getgrade/main.go
SRCS=	imirhil.go types.go utils.go

USRCS=	config_unix.go
WSRCS=	config_windows.go

BIN=	getgrade
EXE=	${BIN}.exe

OPTS=	-ldflags="-s -w" -v

all: ${BIN}

${BIN}: ${GSRCS} ${SRCS} ${USRCS}
	go build ${OPTS} ./cmd/...

${EXE}: ${GSRCS} ${SRCS} ${USRCS}
	GOOS=windows go build ${OPTS} ./cmd/...

build: ${SRCS} ${USRCS}
	go build ${OPTS}

test: build
	go test

windows: ${EXE}
	GOOS=windows go build ${OPTS} .

install:
	go install ${OPTS} ./cmd/...

lint:
	gometalinter .

clean:
	go clean .
	go clean ./cmd/...
	-rm -f ${BIN} ${EXE}

push:
	git push --all
	git push --tags
