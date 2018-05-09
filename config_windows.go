// config_windows.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
// +build windows

package cryptcheck

import (
	"os"
	"path/filepath"
)

/*
File location: %LOCALAPPDATA%\cryptcheck\netrc
*/
var (
	netrcFile = filepath.Join(os.Getenv("%LOCALAPPDATA%"), MyName, "netrc")
)
