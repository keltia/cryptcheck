// imirhil_test.go
//
// Copyright 2018 © by Ollivier Robert <roberto@keltia.net>

package cryptcheck

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	cnfFalseZ  = Config{Log: 0}
	cnfFalseNZ = Config{Log: 1}
	cnfFalseDG = Config{Log: 2}
	cnfTrueZ   = Config{Refresh: true}
	cnfTrueZT5 = Config{Refresh: true, Timeout: 5}
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

func TestClient_GetScore(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.fr.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, BaseURL: baseURL})

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	grade, err := c.GetScore("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetScoreVerbose(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.fr.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, Log: 1, BaseURL: baseURL})

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	grade, err := c.GetScore("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetScoreNoSite(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.com.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.com.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, BaseURL: baseURL})

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	grade, err := c.GetScore("tls.imirhil.com")
	assert.Error(t, err)
	assert.Equal(t, "Z", grade)
}

func TestClient_GetScoreDebug(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.fr.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, Log: 2, BaseURL: baseURL})

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	grade, err := c.GetScore("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetDetailedReport(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.fr.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, BaseURL: baseURL})

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var jr Report

	err = json.Unmarshal(ft, &jr)
	assert.NoError(t, err)

	r, err := c.GetDetailedReport("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, jr, r)
}

func TestClient_GetDetailedVerbose(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.fr.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, Log: 1, BaseURL: baseURL})

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var jr Report

	err = json.Unmarshal(ft, &jr)
	assert.NoError(t, err)

	r, err := c.GetDetailedReport("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, jr, r)
}

func TestClient_GetDetailedNoSite(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.com.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.com.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, BaseURL: baseURL})

	gock.InterceptClient(c.client)
	defer gock.RestoreClient(c.client)

	var jr Report

	err = json.Unmarshal(ft, &jr)
	assert.NoError(t, err)

	r, err := c.GetDetailedReport("tls.imirhil.com")
	assert.Error(t, err)
	assert.Equal(t, jr, r)
}

func TestVersion(t *testing.T) {
	v := Version()
	require.NotEmpty(t, v)
	assert.Equal(t, "201809", v)
}
