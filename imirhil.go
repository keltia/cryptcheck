// imirhil.go

/*
  This file contains the datatypes used by tls.imirhil.fr
*/

package imirhil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "https://tls.imirhil.fr/"
	typeURL = "https/"
	ext     = ".json"

	DefaultWait = 10 * time.Second
	Version     = "201805"
)

// Private area

func myRedirect(req *http.Request, via []*http.Request) error {
	return nil
}

// Public functions

// NewClient setups proxy authentication
func NewClient(logLevel int, refresh bool) *Client {
	c := &Client{refresh: false}

	if logLevel >= 0 {
		c.level = logLevel
	}

	proxyauth, err := setupProxyAuth(c)
	if err != nil {
		c.proxyauth = proxyauth
	}

	if refresh {
		c.refresh = refresh
	}

	_, trsp := c.setupTransport(baseURL)
	c.client = &http.Client{
		Transport:     trsp,
		Timeout:       DefaultWait,
		CheckRedirect: myRedirect,
	}
	c.debug("imirhil: c=%#v", c)
	return c
}

// GetScore retrieves the current score for tls.imirhil.fr
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

	str := fmt.Sprintf("%s/%s/%s.%s", baseURL, typeURL, site, ext)

	if c.refresh {
		str = str + "/refresh"
	}

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
	if resp.StatusCode == http.StatusOK {

		c.debug("status OK")

		if string(body) == "pending" {
			time.Sleep(10 * time.Second)
			resp, err = c.client.Do(req)
			if err != nil {
				return
			}
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
	} else {
		err = fmt.Errorf("did not get acceptable status code: %v body: %q", resp.Status, body)
		return
	}

	err = json.Unmarshal(body, &report)
	return
}
