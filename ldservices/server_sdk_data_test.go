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

func TestFlagOrSegment(t *testing.T) {
	f := FlagOrSegment("my-key", 2)
	assert.Equal(t, "my-key", f.GetKey())

	bytes, err := json.Marshal(f)
	assert.NoError(t, err)
	assertJSONEqual(t, `{"key":"my-key","version":2}`, string(bytes))
}

func TestEmptyServerSDKData(t *testing.T) {
	expectedJSON := `{"flags":{},"segments":{}}`
	data := NewServerSDKData()
	bytes, err := json.Marshal(data)
	assert.NoError(t, err)
	assertJSONEqual(t, expectedJSON, string(bytes))
}

func TestSDKDataWithFlagsAndSegments(t *testing.T) {
	flag1 := FlagOrSegment("flagkey1", 1)
	flag2 := FlagOrSegment("flagkey2", 2)
	segment1 := FlagOrSegment("segkey1", 3)
	segment2 := FlagOrSegment("segkey2", 4)
	data := NewServerSDKData().Flags(flag1, flag2).Segments(segment1, segment2)

	expectedJSON := `{
		"flags": {
			"flagkey1": {
				"key": "flagkey1",
				"version": 1
			},
			"flagkey2": {
				"key": "flagkey2",
				"version": 2
			}
		},
		"segments": {
			"segkey1": {
				"key": "segkey1",
				"version": 3
			},
			"segkey2": {
				"key": "segkey2",
				"version": 4
			}
		}
	}`
	bytes, err := json.Marshal(data)
	assert.NoError(t, err)
	assertJSONEqual(t, expectedJSON, string(bytes))
}
