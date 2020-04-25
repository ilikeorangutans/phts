package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPassword(t *testing.T) {
	p, err := NewPassword("test")
	assert.Nil(t, err)

	assert.False(t, p.Matches("foo"))
	assert.True(t, p.Matches("test"))
}

func TestScan(t *testing.T) {
	p := Password{}
	t.Logf("before %s", p)

	c, _ := NewPassword("horray")
	err := p.Scan(string(c))
	assert.Nil(t, err)
	assert.True(t, p.Matches("horray"))
}
