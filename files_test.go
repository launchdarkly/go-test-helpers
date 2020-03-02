package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithTempFile(t *testing.T) {
	var filePath string
	WithTempFile(func(path string) {
		filePath = path
		assert.True(t, FilePathExists(path))
	})
	assert.False(t, FilePathExists(filePath))
}
