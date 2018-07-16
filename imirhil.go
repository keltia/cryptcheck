// imirhil.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

/*
  This file contains the datatypes used by tls.imirhil.fr
*/

package cryptcheck // import "github.com/keltia/cryptcheck"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "https://tls.imirhil.fr"
	typeURL = "https"
	ext     = ".json"

	// DefaultWait is the timeout
	DefaultWait = 10 * time.Second

	// APIVersion is the cryptcheck API v1 as observed
	APIVersion = "201805"

	// MyVersion is the API version
	MyVersion = "1.2.1"

	// MyName is the name used for the configuration
	MyName = "cryptcheck"
)

// Private area

func myRedirect(req *http.Request, via []*http.Request) error {
	return nil
}

// Public functions

// NewClient setups proxy authentication
func NewClient(cnf ...Config) *Client {
	var c *Client

	// Set default
	if len(cnf) == 0 {
		c = &Client{
			baseurl: baseURL,
			timeout: DefaultWait,
		}
	} else {
		c = &Client{
			baseurl: cnf[0].BaseURL,
			level:   cnf[0].Log,
			refresh: cnf[0].Refresh,
		}

		if cnf[0].Timeout == 0 {
			c.timeout = DefaultWait
		} else {
			c.timeout = time.Duration(cnf[0].Timeout) * time.Second
		}

		// Ensure we have the API endpoint right
		if c.baseurl == "" {
			c.baseurl = baseURL
		}

		c.verbose("got cnf: %#v", cnf[0])
	}

	proxyauth, err := setupProxyAuth(c)
	if err != nil {
		c.proxyauth = proxyauth
	}

	_, trsp := c.setupTransport(c.baseurl)
	c.client = &http.Client{
		Transport:     trsp,
		Timeout:       c.timeout,
		CheckRedirect: myRedirect,
	}
	c.debug("cryptcheck: c=%#v", c)
	return c
}

// GetScore retrieves the current score from tls.imirhil.fr
func (c *Client) GetScore(site string) (score string, err error) {
	full, err := c.GetDetailedReport(site)
	if err != nil {
		score = "Z"
		return
	}
	score = full.Hosts[0].Grade.Rank
	return
}

// GetDetailedReport retrieve the full data
func (c *Client) GetDetailedReport(site string) (report Report, err error) {
	var body []byte

	str := fmt.Sprintf("%s/%s/%s%s", c.baseurl, typeURL, site, ext)

	if c.refresh {
		str = str + "/refresh"
	}

	c.debug("str=%s", str)
	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
		log.Printf("error: req is nil: %v", err)
		return Report{}, nil
	}

	c.debug("req=%#v", req)
	c.debug("clt=%#v", c.client)

	resp, err := c.client.Do(req)
	if err != nil {
		c.verbose("err=%#v", err)
		return
	}
	c.debug("resp=%#v", resp)
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusOK {

		c.debug("status OK")

		if string(body) == "pending" {
			time.Sleep(10 * time.Second)
			resp, err = c.client.Do(req)
			if err != nil {
				return
			}
			c.verbose("resp was %v", resp)
		}
	} else if resp.StatusCode == http.StatusFound {
		str := resp.Header["Location"][0]

		c.debug("Got 302 to %s", str)

		req, err = http.NewRequest("GET", str, nil)
		if err != nil {
			err = fmt.Errorf("Cannot handle redirect: %v", err)
			return
		}

		resp, err = c.client.Do(req)
		if err != nil {
			return
		}
		c.verbose("resp was %v", resp)
	} else {
		err = fmt.Errorf("did not get acceptable status code: %v body: %q", resp.Status, body)
		return
	}

	err = json.Unmarshal(body, &report)
	return
}

// Version returns our internal API version
func Version() string {
	return APIVersion
}
