package jsonhelpers

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONDiff(t *testing.T) {
	diffResult := func(t *testing.T, value1, value2 []byte) JSONDiffResult {
		diff, err := JSONDiff(value1, value2)
		assert.NoError(t, err)
		return diff
	}

	t.Run("equality and inequality without detailed diff", func(t *testing.T) {
		values := []interface{}{
			nil,
			true,
			false,
			3,
			3.5,
			"x",
			[]string{"a", "b"},
			map[string]interface{}{"a": []int{1, 2}},
		}
		for i, value1 := range values {
			jsonValue1, _ := json.Marshal(value1)
			t.Run(fmt.Sprintf("%s == %s", string(jsonValue1), string(jsonValue1)), func(t *testing.T) {
				assert.Nil(t, diffResult(t, jsonValue1, jsonValue1))
			})
			for j, value2 := range values {
				if j == i {
					continue
				}
				jsonValue2, _ := json.Marshal(value2)
				t.Run(fmt.Sprintf("%s != %s", string(jsonValue1), string(jsonValue2)), func(t *testing.T) {
					diff := diffResult(t, jsonValue1, jsonValue2)
					assert.Len(t, diff, 1)
					assert.Nil(t, diff[0].Path)
					assert.Equal(t, string(jsonValue1), diff[0].Value1)
					assert.Equal(t, string(jsonValue2), diff[0].Value2)
				})
			}
		}
	})

	t.Run("inequality with object diff", func(t *testing.T) {
		assert.Equal(t, JSONDiffResult{
			{Path: JSONPath{{Property: "b"}}, Value1: "2", Value2: "3"},
		}, diffResult(t, []byte(`{"a":1,"b":2}`), []byte(`{"a":1,"b":3}`)))

		assert.Equal(t, JSONDiffResult{
			{Path: JSONPath{{Property: "b"}}, Value1: "2", Value2: ""},
		}, diffResult(t, []byte(`{"a":1,"b":2}`), []byte(`{"a":1}`)))

		assert.Equal(t, JSONDiffResult{
			{Path: JSONPath{{Property: "b"}}, Value1: "", Value2: "2"},
		}, diffResult(t, []byte(`{"a":1}`), []byte(`{"a":1,"b":2}`)))

		assert.Equal(t, JSONDiffResult{
			{Path: JSONPath{{Property: "b"}, {Property: "c"}}, Value1: "2", Value2: "3"},
		}, diffResult(t, []byte(`{"a":1,"b":{"c":2}}`), []byte(`{"a":1,"b":{"c":3}}`)))

		assert.Equal(t, JSONDiffResult{
			{Path: JSONPath{{Property: "b"}, {Index: 1}}, Value1: `"d"`, Value2: `"e"`},
		}, diffResult(t, []byte(`{"a":1,"b":["c","d"]}`), []byte(`{"a":1,"b":["c","e"]}`)))
	})

	t.Run("inequality with array diff", func(t *testing.T) {
		assert.Equal(t, JSONDiffResult{
			{Path: JSONPath{{Index: 1}}, Value1: `"b"`, Value2: `"c"`},
		}, diffResult(t, []byte(`["a","b"]`), []byte(`["a","c"]`)))

		assert.Equal(t, JSONDiffResult{
			{Path: JSONPath{{Index: 1}, {Property: "b"}}, Value1: `2`, Value2: `3`},
		}, diffResult(t, []byte(`["a",{"b":2}]`), []byte(`["a",{"b":3}]`)))

		assert.Equal(t, JSONDiffResult{
			{Path: JSONPath{{Index: 2}}, Value1: ``, Value2: `,"c"`},
		}, diffResult(t, []byte(`["a","b"]`), []byte(`["a","b","c"]`)))

		assert.Equal(t, JSONDiffResult{
			{Path: JSONPath{{Index: 2}}, Value1: `,"c"`, Value2: ``},
		}, diffResult(t, []byte(`["a","b","c"]`), []byte(`["a","b"]`)))

		assert.Equal(t, JSONDiffResult{
			{Path: nil, Value1: `["a","d"]`, Value2: `["a","b","c"]`},
		}, diffResult(t, []byte(`["a","d"]`), []byte(`["a","b","c"]`)))
	})
}

func TestJSONDiffResultStrings(t *testing.T) {
	assert.Equal(t, "x = abc, y = def",
		JSONDiffElement{Value1: "abc", Value2: "def"}.Describe("x", "y"))

	assert.Equal(t, `at "prop1": x = abc, y = def`,
		JSONDiffElement{Path: JSONPath{{Property: "prop1"}}, Value1: "abc", Value2: "def"}.
			Describe("x", "y"))

	assert.Equal(t, `at "prop1"."prop2": x = abc, y = def`,
		JSONDiffElement{Path: JSONPath{{Property: "prop1"}, {Property: "prop2"}}, Value1: "abc", Value2: "def"}.
			Describe("x", "y"))

	assert.Equal(t, "at [1]: x = abc, y = def",
		JSONDiffElement{Path: JSONPath{{Index: 1}}, Value1: "abc", Value2: "def"}.
			Describe("x", "y"))

	assert.Equal(t, `at [1]."prop2": x = abc, y = def`,
		JSONDiffElement{Path: JSONPath{{Index: 1}, {Property: "prop2"}}, Value1: "abc", Value2: "def"}.
			Describe("x", "y"))
}
