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

	"github.com/go-resty/resty/v2"
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

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	grade, err := c.GetScore("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, "A+", grade)
}

func TestClient_GetScoreEmpty(t *testing.T) {
	defer gock.Off()

	grade, err := NewClient().GetScore("")
	assert.Error(t, err)
	assert.NotEmpty(t, grade)
	assert.Equal(t, "Z", grade)
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

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

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

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	grade, err := c.GetScore("tls.imirhil.com")
	assert.Error(t, err)
	assert.Equal(t, "Z", grade)
}

func TestClient_GetScoreNoHosts(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.fr-empty.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, BaseURL: baseURL})

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	grade, err := c.GetScore("tls.imirhil.fr")
	assert.Error(t, err)
	assert.Equal(t, "Z", grade)
}

func TestClient_GetScoreWithError(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.fr-error.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, BaseURL: baseURL})

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	grade, err := c.GetScore("tls.imirhil.fr")
	assert.Error(t, err)
	assert.Equal(t, "Z", grade)
	assert.Equal(t, "test for error", err.Error())
}

func TestClient_GetScoreWithErrorDebug(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.fr-error.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, BaseURL: baseURL, Log: 2})

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	grade, err := c.GetScore("tls.imirhil.fr")
	assert.Error(t, err)
	assert.Equal(t, "Z", grade)
	assert.Equal(t, "test for error", err.Error())
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

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

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

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	var jr Report

	err = json.Unmarshal(ft, &jr)
	assert.NoError(t, err)

	r, err := c.GetDetailedReport("tls.imirhil.fr")
	assert.NoError(t, err)
	assert.Equal(t, jr, r)
}

func TestClient_GetDetailedReportRefresh(t *testing.T) {
	defer gock.Off()

	ft, err := ioutil.ReadFile("testdata/tls.imirhil.fr.json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr/refresh").
		Reply(200).
		BodyString(string(""))

	gock.New(baseURL).
		Get("/https/tls.imirhil.fr.json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient(Config{Timeout: 10, BaseURL: baseURL, Refresh: true})

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

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

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

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

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	var jr Report

	err = json.Unmarshal(ft, &jr)
	assert.NoError(t, err)

	r, err := c.GetDetailedReport("tls.imirhil.com")
	assert.Error(t, err)
	assert.Equal(t, jr, r)
}

func TestCallAPI(t *testing.T) {
	defer gock.Off()

	site := "tls.imirhil.fr"

	ft, err := ioutil.ReadFile("testdata/" + site + ".json")
	assert.NoError(t, err)

	gock.New(baseURL).
		Get("/https/" + site + ".json").
		Reply(200).
		BodyString(string(ft))

	c := NewClient()

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	str := "https://tls.imirhil.fr/https/" + site + ".json"
	resp, body, err := c.callAPI(str)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.IsType(t, (*resty.Response)(nil), resp)

	assert.NotNil(t, body)
	assert.NotEmpty(t, body)
}

func TestCallAPI2(t *testing.T) {
	defer gock.Off()

	site := "tls.imirhil.fr"

	gock.New(baseURL).
		Get("/https/"+site+"/refresh").
		Reply(200).
		SetHeader("Location", baseURL+"/https/"+site+".json").
		BodyString(string("HTML body we do not care about"))

	c := NewClient()

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	str := "https://tls.imirhil.fr/https/" + site + "/refresh"
	resp, body, err := c.callAPI(str)

	require.NoError(t, err)
	assert.Equal(t, "HTML body we do not care about", string(body))
	assert.NotNil(t, resp)
	assert.IsType(t, (*resty.Response)(nil), resp)

	h := resp.Header()
	assert.Equal(t, baseURL+"/https/"+site+".json", h.Get("Location"))

	assert.NotNil(t, body)
	assert.NotEmpty(t, body)
}

func TestCallAPI3(t *testing.T) {
	defer gock.Off()

	site := ""

	gock.New(baseURL).
		Get("/https/"+site+"/refresh").
		Reply(200).
		SetHeader("Location", baseURL+"/https/"+site+".json").
		BodyString(string("<!DOCTYPE html>"))

	c := NewClient()

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	str := "https://tls.imirhil.fr/https/" + site + "/refresh"
	resp, body, err := c.callAPI(str)

	require.NoError(t, err)
	assert.Equal(t, "<!DOCTYPE html>", string(body))
	assert.NotNil(t, resp)
	assert.IsType(t, (*resty.Response)(nil), resp)
	h := resp.Header()
	assert.Equal(t, baseURL+"/https/"+site+".json", h.Get("Location"))

	assert.NotNil(t, body)
	assert.NotEmpty(t, body)
}

func TestCallAPI4(t *testing.T) {
	defer gock.Off()

	site := "../../.."

	gock.New(baseURL).
		Get("/https/"+site+"/refresh").
		Reply(404).
		SetHeader("Location", baseURL+"/https/"+site+".json").
		BodyString(string("<!DOCTYPE html>"))

	c := NewClient()

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	str := "https://tls.imirhil.fr/https/" + site + "/refresh"
	resp, body, err := c.callAPI(str)

	require.Error(t, err)
	assert.Equal(t, "<!DOCTYPE html>", string(body))
	assert.NotNil(t, resp)
	assert.IsType(t, (*resty.Response)(nil), resp)
	h := resp.Header()
	assert.Equal(t, baseURL+"/https/"+site+".json", h.Get("Location"))

	assert.NotNil(t, body)
	assert.NotEmpty(t, body)
}

func TestCallAPI5(t *testing.T) {
	defer gock.Off()

	site := "../../.."

	gock.New("https://bad.imirhil.fr").
		Get("/https/" + site + "/refresh").
		Reply(404).
		BodyString(string("<HTML>"))

	c := NewClient()

	gock.InterceptClient(c.client.GetClient())
	defer gock.RestoreClient(c.client.GetClient())

	str := "https://bad.imirhil.fr/https/" + site + "/refresh"
	resp, body, err := c.callAPI(str)

	require.Error(t, err)
	assert.NotNil(t, body)
	assert.NotEmpty(t, string(body))
	assert.Contains(t, "<HTML>", string(body))
	assert.NotNil(t, resp)
	assert.IsType(t, (*resty.Response)(nil), resp)
}

func TestMyRedirect(t *testing.T) {
	err := myRedirect(nil, nil)
	require.NoError(t, err)
}

func TestVersion(t *testing.T) {
	v := Version()
	require.NotEmpty(t, v)
	assert.Equal(t, "201909", v)
}
