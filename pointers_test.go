package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsPtr(t *testing.T) {
	boolP := AsPointer(true)
	require.NotNil(t, boolP)
	assert.True(t, *boolP)

	intP := AsPointer(2)
	require.NotNil(t, intP)
	assert.Equal(t, 2, *intP)
}
