// types.go
//
// Copyright 2018-2019 © by Ollivier Robert <roberto@keltia.net>

// v1 on tls.imrhil.fr

// XXX Versioning of the API is nonexistent, we have to cope
// 20160510 "old" API
// 20160511 "new" API
// 20171204 add ID to struct Report
// 20180502 added two fields in the report top struct
// 20180904 added Host.Error.
// 20190914 new incompatible API but same Endpoint :(
//          data moves around, etc.

package cryptcheck

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// Key describes a single key
type Key struct {
	Type    string `json:"type"`
	Size    int    `json:"size"`
	RSASize int    `json:"rsa_size"`
}

// Cipher describes a single cipher
type Cipher struct {
	Protocol string
	Name     string
	Size     int
	DH       Key `json:"dh"`
}

// Grade aka score of the site
type Grade struct {
	Rank    string
	Details struct {
		Score           float64 `json:"score"`
		Protocol        int     `json:"protocol"`
		KeyExchange     int     `json:"key_exchange"`
		CipherStrengths int     `json:"cipher_strengths"`
	} `json:"details"`
	Error   []string
	Danger  []string
	Warning []string
	Success []string
}

// Site contains DNS site data
type Site struct {
	Name string
	IP   string `json:"ip"`
	Port int
}

// Handshake contains crypto parameters
type Handshake struct {
	Key       Key
	DH        []Key `json:"dh"`
	Protocols []string
	Ciphers   []Cipher
	HSTS      int `json:"hsts"`
}

// Host describe a single host
type Host struct {
	Host      Site      `json:"host"`
	Handshake Handshake `json:"handshake"`
	Grade     Grade
	Error     string
}

// Report describes the details for the crypto
type Report struct {
	ID struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Service   string
	Host      string
	Port      int
	Pending   bool
	Result    Result
	CreateAt  time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Result holds the actual results, host-per-host
type Result struct {
	Date  time.Time `json:"date"`
	Hosts []Host
}

// Client is used to store proxyauth & other internal state
type Client struct {
	baseurl string
	level   int
	client  *resty.Client
	timeout time.Duration
	refresh bool
}

// Config is for giving options to NewClient
type Config struct {
	BaseURL string
	Timeout int
	Refresh bool
	Log     int
}
