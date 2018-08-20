
# Main Makefile for cryptcheck
#
# Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
#

GO=		go
GOBIN=	${GOPATH}/bin

GSRCS=	cmd/cryptcheck/main.go
SRCS=	imirhil.go types.go utils.go

USRCS=	config_unix.go
WSRCS=	config_windows.go

BIN=	cryptcheck
EXE=	${BIN}.exe

OPTS=	-ldflags="-s -w" -v

all: ${BIN}

${BIN}: ${GSRCS} ${SRCS} ${USRCS}
	${GO} build ${OPTS} ./cmd/...

${EXE}: ${GSRCS} ${SRCS} ${USRCS}
	GOOS=windows ${GO} build ${OPTS} ./cmd/...

build: ${SRCS} ${USRCS}
	${GO} build ${OPTS}

test: build
	${GO} test .

install:
	${GO} install ${OPTS} ./cmd/...

lint:
	gometalinter .

clean:
	${GO} clean .
	${GO} clean ./cmd/...
	-rm -f ${BIN} ${EXE}

push:
	git push --all
	git push --tags
