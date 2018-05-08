package imirhil

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

const (
	// GoodAuth is test:test
	GoodAuth = "Basic dGVzdDp0ZXN0"
)

// --- setupProxyAuth
func TestSetupProxyAuthNoNetrc(t *testing.T) {
	client := &Client{}

	f := filepath.Join(".", "test/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	_, err = setupProxyAuth(client)
	assert.Error(t, err, "should be an error")
	assert.Equal(t, ErrNoAuth, err)
}

func TestSetupProxyAuthVerboseNoNetrc(t *testing.T) {
	client := &Client{level: 1}

	f := filepath.Join(".", "test/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	_, err = setupProxyAuth(client)
	assert.Error(t, err, "should be an error")
	assert.Equal(t, ErrNoAuth, err)
}

func TestSetupProxyAuth(t *testing.T) {
	client := &Client{}

	f := filepath.Join(".", "test/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	auth, err := setupProxyAuth(client)
	assert.NoError(t, err, "no error")
	assert.Equal(t, GoodAuth, auth)
}

func TestSetupProxyAuthVerbose(t *testing.T) {
	client := &Client{level: 1}

	f := filepath.Join(".", "test/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	auth, err := setupProxyAuth(client)
	assert.NoError(t, err, "no error")
	assert.Equal(t, GoodAuth, auth)
}

// -- loadNetrc
func TestLoadNetrcNoFile(t *testing.T) {
	client := &Client{}

	f := filepath.Join(".", "test/no-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	user, password := loadNetrc(client)
	assert.EqualValues(t, "", user, "null user")
	assert.EqualValues(t, "", password, "null password")
}

func TestLoadNetrcZero(t *testing.T) {
	client := &Client{}

	err := os.Setenv("NETRC", filepath.Join(".", "test/zero-netrc"))
	require.NoError(t, err)

	user, password := loadNetrc(client)
	assert.EqualValues(t, "", user, "test user")
	assert.EqualValues(t, "", password, "test password")
}

func TestLoadNetrcPerms(t *testing.T) {
	client := &Client{}

	f := filepath.Join(".", "test/perms-netrc")
	err := os.Setenv("NETRC", f)
	assert.NoError(t, err)

	err = os.Chmod(f, 0644)
	require.NoError(t, err)

	user, password := loadNetrc(client)
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	assert.EqualValues(t, "", user, "test user")
	assert.EqualValues(t, "", password, "test password")
}

func TestLoadNetrcGood(t *testing.T) {
	client := &Client{}

	f := filepath.Join(".", "test/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	user, password := loadNetrc(client)
	assert.EqualValues(t, "test", user, "test user")
	assert.EqualValues(t, "test", password, "test password")
}

func TestLoadNetrcGoodVerbose(t *testing.T) {
	client := &Client{level: 1}

	f := filepath.Join(".", "test/test-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	user, password := loadNetrc(client)
	assert.EqualValues(t, "test", user, "test user")
	assert.EqualValues(t, "test", password, "test password")
}

func TestLoadNetrcBad(t *testing.T) {
	client := &Client{}

	f := filepath.Join(".", "test/bad-netrc")
	err := os.Setenv("NETRC", f)
	require.NoError(t, err)

	// We must ensure propre perms
	err = os.Chmod(f, 0600)
	require.NoError(t, err)

	user, password := loadNetrc(client)
	assert.EqualValues(t, "", user, "test user")
	assert.EqualValues(t, "", password, "test password")
}
