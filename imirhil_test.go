package imirhil

import (
    "testing"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
    c := NewClient(0, false)

    require.NotNil(t, c)
    require.IsType(t, (*Client)(nil), c)
    require.NotNil(t, c.client)

    assert.Equal(t, 0, c.level)
    assert.False(t, c.refresh)
}

func TestNewClient2(t *testing.T) {
    c := NewClient(1, false)

    require.NotNil(t, c)
    require.IsType(t, (*Client)(nil), c)
    require.NotNil(t, c.client)

    assert.Equal(t, 1, c.level)
    assert.False(t, c.refresh)
}

func TestNewClient3(t *testing.T) {
    c := NewClient(0, true)

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
