package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoolPtr(t *testing.T) {
	p := BoolPtr(true)
	require.NotNil(t, p)
	assert.True(t, *p)
}

func TestIntPtr(t *testing.T) {
	p := IntPtr(2)
	require.NotNil(t, p)
	assert.Equal(t, 2, *p)
}

func TestFloat64Ptr(t *testing.T) {
	p := Float64Ptr(2.5)
	require.NotNil(t, p)
	assert.Equal(t, float64(2.5), *p)
}

func TestStrPtr(t *testing.T) {
	p := StrPtr("x")
	require.NotNil(t, p)
	assert.Equal(t, "x", *p)
}

func TestUint64Ptr(t *testing.T) {
	p := Uint64Ptr(uint64(2))
	require.NotNil(t, p)
	assert.Equal(t, uint64(2), *p)
}
