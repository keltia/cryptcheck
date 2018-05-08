package imirhil

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	cnfFalseZ  = Config{Log: 0}
	cnfFalseNZ = Config{Log: 1}
	cnfTrueZ   = Config{Refresh: true}
)

func TestNewClient(t *testing.T) {
	c := NewClient(cnfFalseZ)

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 0, c.level)
	assert.False(t, c.refresh)
}

func TestNewClient2(t *testing.T) {
	c := NewClient(cnfFalseNZ)

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 1, c.level)
	assert.False(t, c.refresh)
}

func TestNewClient3(t *testing.T) {
	c := NewClient(cnfTrueZ)

	require.NotNil(t, c)
	require.IsType(t, (*Client)(nil), c)
	require.NotNil(t, c.client)

	assert.Equal(t, 0, c.level)
	assert.True(t, c.refresh)
}

func TestClient_GetScore(t *testing.T) {

}

func TestClient_GetDetailedReport(t *testing.T) {

}
