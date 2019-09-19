// imirhil.go
//
// Copyright 2018-1019 © by Ollivier Robert <roberto@keltia.net>

/*
  This file contains the datatypes used by tls.imirhil.fr
*/

package cryptcheck // import "github.com/keltia/cryptcheck"

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	baseURL = "https://tls.imirhil.fr"
	typeURL = "https"
	ext     = ".json"

	// DefaultWait is the timeout
	DefaultWait = 10 * time.Second

	// APIVersion is the cryptcheck API v1 as observed
	APIVersion = "201909"

	// MyVersion is the API version
	MyVersion = "1.5.2"

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

	c.client = resty.New().SetTimeout(DefaultWait)
	proxy := os.Getenv("http_proxy")
	if proxy != "" {
		c.client.SetProxy(proxy)
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
	if len(full.Result.Hosts) == 0 {
		return "Z", fmt.Errorf("empty hosts")
	}

	c.debug("fullscore=%#v", full)
	score = full.Result.Hosts[0].Grade.Rank
	return
}

// GetDetailedReport retrieve the full data
func (c *Client) GetDetailedReport(site string) (report Report, err error) {
	var (
		retry = 0
		body  []byte
		str   string
	)

	if site == "" {
		return Report{}, errors.Wrap(err, "empty site")
	}

	if c.refresh {
		str = fmt.Sprintf("%s/%s/%s/%s", c.baseurl, typeURL, site, "refresh")
	} else {
		str = fmt.Sprintf("%s/%s/%s%s", c.baseurl, typeURL, site, ext)
	}

	resp, body, err := c.callAPI(str)
	if err != nil {
		return Report{}, errors.Wrap(err, "err resp")
	}

	for {
		if retry >= DefaultRetry {
			return Report{}, errors.Wrap(err, "retry expired")
		}

		if resp == nil {
			retry++
			c.debug("nil resp/loop")
			continue
		}
		if resp.StatusCode() == http.StatusOK {

			c.debug("status OK")

			// If refreshing, we discard the body
			if c.refresh {
				c.debug("refresh requested")
				str = fmt.Sprintf("%s/%s/%s%s", c.baseurl, typeURL, site, ext)
				resp, body, err = c.callAPI(str)
				if err != nil {
					return Report{}, errors.Wrap(err, "refresh error")
				}

				// Reset it otherwise we loop forever-ish
				c.refresh = false
				continue
			}

			var r Report

			if err := json.Unmarshal(resp.Body(), &r); err != nil {
				return Report{}, errors.Wrapf(err, "bad json: %s", body)
			}

			if r.Pending {
				retry++
				time.Sleep(10 * time.Second)

				resp, body, err = c.callAPI(str)
				if err != nil {
					return Report{}, errors.Wrap(err, "pending error")
				}
				c.debug("resp was %v", resp)
			} else {
				// Next call succeed
				break
			}
		} else {
			return Report{}, errors.Wrapf(err, "bad status %v body %v", resp.Status, body)
		}
	}
	c.debug("success: %s", string(body))
	err = json.Unmarshal(body, &report)

	if len(report.Result.Hosts) != 0 {
		if report.Result.Hosts[0].Error != "" {
			c.debug("got errors")
			err = errors.New(fmt.Sprintf("%v", report.Result.Hosts[0].Error))
			return
		}
	}

	return
}

// callAPI does the main chunk of a call
func (c *Client) callAPI(strURL string) (*resty.Response, []byte, error) {

	c.debug("strURL=%s", strURL)
	c.debug("clt=%#v", c.client)

	resp, err := c.client.R().
		SetHeader("User-Agent", fmt.Sprintf("%s/%s", MyName, MyVersion)).
		Get(strURL)

	if err != nil {
		c.debug("err=%s", err.Error())
		return nil, nil, errors.Wrap(err, "client.Do")
	}

	body := resp.Body()

	c.debug("resp=%#v, body=%s", resp, string(body))

	if resp.StatusCode() != http.StatusOK {
		return resp, body, fmt.Errorf("NOK code=%d", resp.StatusCode)
	}

	c.debug("success: %s", string(body))
	return resp, body, nil
}

// Version returns our internal API version
func Version() string {
	return APIVersion
}
