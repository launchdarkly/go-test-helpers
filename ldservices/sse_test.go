package ldservices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSSEEvent(t *testing.T) {
	e := NewSSEEvent("my-id", "my-event", "my-data")
	assert.Equal(t, "my-id", e.Id())
	assert.Equal(t, "my-event", e.Event())
	assert.Equal(t, "my-data", e.Data())
}
