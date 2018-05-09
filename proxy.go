// proxy.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

package cryptcheck

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	proxyTag = "proxy"
)

var ErrNoAuth = fmt.Errorf("no proxy auth")

// Private functions

func getProxy(req *http.Request) (uri *url.URL, err error) {
	uri, err = http.ProxyFromEnvironment(req)
	if err != nil {
		log.Printf("no proxy in environment")
		uri = &url.URL{}
	} else if uri == nil {
		log.Println("No proxy configured or url excluded")
	}
	return
}

func setupProxyAuth(c *Client) (proxyauth string, err error) {
	// Try to load $HOME/.netrc or file pointed at by $NETRC
	user, password := loadNetrc(c)

	if user != "" {
		c.verbose("Proxy user %s found.", user)
	}

	err = ErrNoAuth

	// Do we have a proxy user/password?
	if user != "" && password != "" {
		auth := fmt.Sprintf("%s:%s", user, password)
		proxyauth = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		err = nil
	}
	return
}

// loadNetrc supports a subset of the original ftp(1) .netrc file.
/*
We support:

  machine
  default
  login
  password

Format:
  machine <host> login <user> password <pass>
*/
func loadNetrc(c *Client) (user, password string) {
	var dnetrc string

	// is $NETRC defined?
	dnetVar := os.Getenv("NETRC")

	// Allow override
	if dnetVar == "" {
		dnetrc = netrcFile
	} else {
		dnetrc = dnetVar
	}

	c.verbose("NETRC=%s", dnetVar)

	// First check for permissions
	fh, err := os.Open(dnetrc)
	if err != nil {
		c.verbose("warning: can not find/read %s: %v", dnetrc, err)
		return "", ""
	}
	defer fh.Close()

	// Now check permissions
	st, err := fh.Stat()
	if err != nil {
		c.verbose("unable to stat: %v", err)
		return "", ""
	}

	if (st.Mode() & 077) != 0 {
		c.verbose("invalid permissions, must be 0400/0600")
		return "", ""
	}

	c.verbose("now parsing")
	user, password = parseNetrc(c, fh)
	return
}

/*
   Format:
   machine proxy|default login <user> password <pass>
*/
// parseDbrc loads the file format historically defined by DBI::Dbrc
func parseNetrc(c *Client, r io.Reader) (user, password string) {
	c.verbose("found netrc")

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break
		}

		flds := strings.Split(line, " ")
		c.verbose("%s: %d fields", line, len(flds))

		if flds[0] != "machine" {
			c.verbose("machine is not the first word")
			continue
		}

		// Check what we need
		if len(flds) != 6 {
			c.verbose("bad format")
			continue
		}

		if flds[1] == proxyTag || flds[1] == "default" {

			if flds[2] == "login" && flds[4] == "password" {
				user = flds[3]
				password = flds[5]
				c.verbose("got %s/default entry for user %s", proxyTag, user)
			}
			break
		}
	}
	if err := s.Err(); err != nil {
		c.verbose("error reading netrc: %v", err)
		return "", ""
	}

	if user == "" {
		c.verbose("no user/password for %s/default in netrc", proxyTag)
	}

	return
}

func (c *Client) setupTransport(str string) (*http.Request, *http.Transport) {
	/*
	   Proxy code taken from https://github.com/LeoCBS/poc-proxy-https/blob/master/main.go
	*/
	myurl, err := url.Parse(str)
	if err != nil {
		log.Printf("error parsing %s: %v", str, err)
		return nil, nil
	}

	req, err := http.NewRequest("GET", str, nil)
	if err != nil {
		c.debug("error: req is nil: %v", err)
		return nil, nil
	}
	req.Header.Set("Host", myurl.Host)
	req.Header.Add("User-Agent", fmt.Sprintf("cryptcheck/%s", MyVersion))

	// Get proxy URL
	proxyURL, err := getProxy(req)
	if c.proxyauth != "" {
		req.Header.Add("Proxy-Authorization", c.proxyauth)
	}

	transport := &http.Transport{
		Proxy:              http.ProxyURL(proxyURL),
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		ProxyConnectHeader: req.Header,
	}
	c.debug("transport=%#v", transport)
	return req, transport
}
