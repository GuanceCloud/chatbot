package config

import (
	"net"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNetjoin(t *testing.T) {
	v := net.JoinHostPort("", "1234")
	assert.Equal(t, ":1234", v)
}
