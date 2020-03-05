package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type myCloser struct {
	closed bool
}

func (m *myCloser) Close() error {
	m.closed = true
	return nil
}

func TestWithCloser(t *testing.T) {
	c := &myCloser{}
	WithCloser(c, func() {
		assert.False(t, c.closed)
	})
	assert.True(t, c.closed)
}
