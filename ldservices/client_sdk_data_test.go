package ldservices

import (
	"encoding/json"
	"testing"

	"gopkg.in/launchdarkly/go-sdk-common.v1/ldvalue"

	"github.com/stretchr/testify/assert"
)

func TestEmptyClientSDKData(t *testing.T) {
	expectedJSON := `{}`
	data := NewClientSDKData()
	bytes, err := json.Marshal(data)
	assert.NoError(t, err)
	assertJSONEqual(t, expectedJSON, string(bytes))
}

func TestClientSDKDataWithFlags(t *testing.T) {
	flag1 := FlagValueData{
		Key:                  "flagkey1",
		Version:              1,
		FlagVersion:          1000,
		Value:                ldvalue.String("a"),
		VariationIndex:       2,
		Reason:               ldvalue.ObjectBuild().Set("kind", ldvalue.String("FALLTHROUGH")).Build(),
		TrackEvents:          true,
		DebugEventsUntilDate: uint64(3000),
	}
	flag2 := FlagValueData{
		Key:            "flagkey2",
		Version:        2,
		FlagVersion:    2000,
		Value:          ldvalue.String("b"),
		VariationIndex: -1,
	}
	data := NewClientSDKData().Flags(flag1, flag2)

	expectedJSON := `{
		"flagkey1": {
			"version": 1,
			"flagVersion": 1000,
			"value": "a",
			"variation": 2,
			"reason": { "kind": "FALLTHROUGH" },
			"trackEvents": true,
			"debugEventsUntilDate": 3000
		},
		"flagkey2": {
			"version": 2,
			"flagVersion": 2000,
			"value": "b"
		}
	}`
	bytes, err := json.Marshal(data)
	assert.NoError(t, err)
	assertJSONEqual(t, expectedJSON, string(bytes))
}
