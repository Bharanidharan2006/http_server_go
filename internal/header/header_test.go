package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Test: Valid single header
headers := NewHeaders()
data := []byte("Host: localhost:42069\r\nHost: localhost:2309\r\n\r\n")
n, done, err := headers.Parse(data)
require.NoError(t, err)
require.NotNil(t, headers)
assert.Equal(t, "localhost:42069, localhost:2309", headers["host"])
assert.Equal(t, 47, n)
assert.True(t, done)

// Test: Invalid spacing header
headers = NewHeaders()
data = []byte("       Host : localhost:42069       \r\n\r\n")
n, done, err = headers.Parse(data)
require.Error(t, err)
assert.Equal(t, 0, n)
assert.False(t, done)

headers = NewHeaders()
data = []byte("H©st: localhost:42069       \r\n\r\n")
n, done, err = headers.Parse(data)
require.Error(t, err)
assert.Equal(t, 0, n)
assert.False(t, done)

}