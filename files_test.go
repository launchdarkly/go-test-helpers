package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTempFile(t *testing.T) {
	var filePath string
	WithTempFile(func(path string) {
		filePath = path
		assert.True(t, FilePathExists(path))
	})
	assert.False(t, FilePathExists(filePath))
}

func TestWithTempFileData(t *testing.T) {
	var filePath string
	WithTempFileData([]byte(`hello`), func(path string) {
		filePath = path
		data, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "hello", string(data))
	})
	assert.False(t, FilePathExists(filePath))
}

func TestWithTempDir(t *testing.T) {
	var path string
	WithTempDir(func(dirPath string) {
		path = dirPath
		info, err := os.Stat(path)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
		assert.NoError(t, os.WriteFile(filepath.Join(dirPath, "x"), []byte("hello"), 0600))
	})
	assert.False(t, FilePathExists(path))
}
