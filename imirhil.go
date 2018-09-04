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

	"github.com/keltia/proxy"
	"github.com/pkg/errors"
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
	MyVersion = "1.3.3"

	// MyName is the name used for the configuration
	MyName = "cryptcheck"

	// DefaultRetry is the number of times we try hard to get an answer
	DefaultRetry = 5
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
			level:   0,
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

		c.debug("got cnf: %#v", cnf[0])
	}

	// for informational purposes, ErrNoAuth is a ignorable error
	c.proxyauth, _ = proxy.SetupProxyAuth()

	_, trsp := proxy.SetupTransport(c.baseurl)
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
	c.debug("fullscore=%#v", full)
	if full.Hosts[0].Error != "" {
		c.debug("got errors")
		score = "Z"
		err = errors.New(fmt.Sprintf("unknown site: %v", full.Hosts[0].Error))
		return
	}
	score = full.Hosts[0].Grade.Rank
	return
}

// GetDetailedReport retrieve the full data
func (c *Client) GetDetailedReport(site string) (report Report, err error) {
	var (
		retry = 0
		body  []byte
	)

	str := fmt.Sprintf("%s/%s/%s%s", c.baseurl, typeURL, site, ext)

	if c.refresh {
		str = str + "/refresh"
	}

	c.debug("str=%s", str)
	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
		log.Printf("error: req is nil: %v", err)
		return Report{}, errors.Wrap(err, "http.newrequest")
	}

	c.debug("req=%#v", req)
	c.debug("clt=%#v", c.client)

	resp, err := c.client.Do(req)
	if err != nil {
		c.verbose("err=%s", err.Error())
		return Report{}, errors.Wrap(err, "1st call")
	}
	c.debug("resp=%#v, body=%s", resp, string(body))

	for {
		if retry == DefaultRetry {
			return Report{}, errors.Wrap(err, "retry expired")
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return Report{}, errors.Wrap(err, "read body")
		}

		if resp.StatusCode == http.StatusOK {

			c.debug("status OK")

			if string(body) == "pending" {
				retry++
				time.Sleep(10 * time.Second)
				resp, err = c.client.Do(req)
				if err != nil {
					return Report{}, errors.Wrap(err, "pending error")
				}
				c.verbose("resp was %v", resp)
			} else {
				// Next call succeed
				break
			}
		} else if resp.StatusCode == http.StatusFound {
			str := resp.Header["Location"][0]

			c.debug("Got 302 to %s", str)

			req, err = http.NewRequest("GET", str, nil)
			if err != nil {
				return Report{}, errors.Wrap(err, "bad redirect")
			}

			resp, err = c.client.Do(req)
			if err != nil {
				return Report{}, errors.Wrap(err, "redirect")
			}
			c.verbose("resp was %v", resp)
		} else {
			return Report{}, errors.Wrapf(err, "bad status %v body %v", resp.Status, body)
		}
	}
	c.debug("success: %s", string(body))
	err = json.Unmarshal(body, &report)
	return
}

// Version returns our internal API version
func Version() string {
	return APIVersion
}
