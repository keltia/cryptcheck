// imirhil_test.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

package cryptcheck

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/goware/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testURL = "http://localhost:10000"

var (
	cnfFalseZ  = Config{Log: 0}
	cnfFalseNZ = Config{Log: 1}
	cnfFalseDG = Config{Log: 2}
	cnfTrueZ   = Config{Refresh: true}
	cnfTrueZT5 = Config{Refresh: true, Timeout: 5}

	mockService *httpmock.MockHTTPServer
)

func TestNewClientDefault(t *testing.T) {
	f := filepath.Join(".", "testdata/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	c := NewClient()

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 0, c.level)
	assert.Equal(t, DefaultWait, c.timeout)
	assert.False(t, c.refresh)
}

func TestNewClient(t *testing.T) {
	f := filepath.Join(".", "testdata/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	c := NewClient(cnfFalseZ)

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 0, c.level)
	assert.Equal(t, DefaultWait, c.timeout)
	assert.False(t, c.refresh)
}

func TestNewClient2(t *testing.T) {
	f := filepath.Join(".", "testdata/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	c := NewClient(cnfFalseNZ)

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 1, c.level)
	assert.Equal(t, DefaultWait, c.timeout)
	assert.False(t, c.refresh)
}

func TestNewClient3(t *testing.T) {
	f := filepath.Join(".", "testdata/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	c := NewClient(cnfTrueZ)

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 0, c.level)
	assert.Equal(t, DefaultWait, c.timeout)
	assert.True(t, c.refresh)
}

func TestNewClient4(t *testing.T) {
	f := filepath.Join(".", "testdata/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	c := NewClient(cnfTrueZT5)

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 0, c.level)
	assert.Equal(t, 5*time.Second, c.timeout)
	assert.True(t, c.refresh)
}

func TestNewClient5(t *testing.T) {
	f := filepath.Join(".", "testdata/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	c := NewClient(cnfFalseDG)

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 2, c.level)
	assert.Equal(t, DefaultWait, c.timeout)
	assert.False(t, c.refresh)
}

func TestNewClientNoProxy(t *testing.T) {
	f := filepath.Join(".", "testdata/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	c := NewClient(cnfFalseZ)

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 0, c.level)
	assert.Equal(t, DefaultWait, c.timeout)
	assert.False(t, c.refresh)
}

func BeforeAPI(t *testing.T) {
	if mockService == nil {
		// new mocking server
		mockService = httpmock.NewMockHTTPServer("127.0.0.1:10000")
	}

	// define request->response pairs
	request1, _ := url.Parse("http://127.0.0.1:10000/https/tls.imirhil.fr.json")
	request2, _ := url.Parse("http://127.0.0.1:10000/https/tls.imirhil.com.json")
	ft, err := ioutil.ReadFile("test/tls.imirhil.fr.json")
	assert.NoError(t, err)

	mockService.AddResponses([]httpmock.MockResponse{
		{
			Request: http.Request{
				Method: "GET",
				URL:    request1,
			},
			Response: httpmock.Response{
				StatusCode: 200,
				Body:       string(ft),
			},
		},
		{
			Request: http.Request{
				Method: "GET",
				URL:    request2,
			},
			Response: httpmock.Response{
				StatusCode: 302,
				Body:       "no site",
			},
		},
	})
}

func TestClient_GetScore(t *testing.T) {
	ct := NewClient(Config{Timeout: 10, BaseURL: testURL})
	BeforeAPI(t)

	t.Logf("ct=%#v", ct)
	grade, err := ct.GetScore("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetScoreVerbose(t *testing.T) {
	ct := NewClient(Config{Timeout: 10, Log: 1, BaseURL: testURL})
	BeforeAPI(t)

	t.Logf("ct=%#v", ct)
	grade, err := ct.GetScore("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetScoreNoSite(t *testing.T) {
	ct := NewClient(Config{Timeout: 10, BaseURL: testURL, Log: 2})
	BeforeAPI(t)

	t.Logf("ct=%#v", ct)
	grade, err := ct.GetScore("tls.imirhil.com")
	assert.Error(t, err)
	t.Logf("error=%v", err)
	assert.Equal(t, "Z", grade)
}

func TestClient_GetScoreDebug(t *testing.T) {
	ct := NewClient(Config{Timeout: 10, Log: 2, BaseURL: testURL})
	BeforeAPI(t)

	t.Logf("ct=%#v", ct)
	grade, err := ct.GetScore("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetDetailedReport(t *testing.T) {
	ct := NewClient(Config{BaseURL: testURL})
	BeforeAPI(t)

	var jr Report

	ft, err := ioutil.ReadFile("test/tls.imirhil.fr.json")
	require.NoError(t, err)

	err = json.Unmarshal(ft, &jr)
	assert.NoError(t, err)

	r, err := ct.GetDetailedReport("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, jr, r)
}

func TestClient_GetDetailedVerbose(t *testing.T) {
	ct := NewClient(Config{Log: 1, BaseURL: testURL})
	BeforeAPI(t)

	var jr Report

	ft, err := ioutil.ReadFile("test/tls.imirhil.fr.json")
	require.NoError(t, err)

	err = json.Unmarshal(ft, &jr)
	assert.NoError(t, err)

	r, err := ct.GetDetailedReport("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, jr, r)
}

func TestClient_GetDetailedNoSite(t *testing.T) {
	ct := NewClient(Config{BaseURL: testURL})
	BeforeAPI(t)

	r, err := ct.GetDetailedReport("tls.imirhil.com")
	assert.Error(t, err)
	assert.Equal(t, Report{}, r)
}

func TestVersion(t *testing.T) {
	v := Version()
	require.NotEmpty(t, v)
	assert.Equal(t, "201805", v)
}
