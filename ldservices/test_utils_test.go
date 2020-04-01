package ldservices

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseJSONAsMap(t *testing.T, s string) map[string]interface{} {
	var fields map[string]interface{}
	err := json.Unmarshal([]byte(s), &fields)
	require.NoError(t, err)
	return fields
}

func assertJSONEqual(t *testing.T, expected, actual string) {
	assert.Equal(t, parseJSONAsMap(t, expected), parseJSONAsMap(t, actual))
}
