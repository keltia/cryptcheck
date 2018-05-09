// config_unix.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
// +build !windows

/*
File location: $HOME/.netrc
*/
package cryptcheck

import (
	"os"
	"path/filepath"
)

var (
	netrcFile = filepath.Join(os.Getenv("HOME"), ".netrc")
)
