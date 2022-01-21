package matchers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueFormatting(t *testing.T) {
	assert.Equal(t, `"abc"`, DescribeValue("abc"))

	assert.Equal(t, `{abc}`, DescribeValue("{abc}"))

	assert.Equal(t, `[abc]`, DescribeValue("[abc]"))

	assert.Equal(t, decorate("abc"), DescribeValue(decoratedString("abc")))

	assert.Equal(t, "abc", DescribeValue([]byte("abc")))

	assert.Equal(t, `{"a":1,"b":2}`, DescribeValue([]byte(`{"b":2,"a":1}`)))

	assert.Equal(t, `{"a":1,"b":2}`, DescribeValue(json.RawMessage(`{"b":2,"a":1}`)))

	taggedStruct := struct {
		Name   string `json:"name"`
		Values []int  `json:"values"`
	}{"Lucy", []int{1, 2}}

	untaggedStruct := struct {
		Name   string
		Values []int
	}{"Mina", []int{1, 2}}

	assert.Equal(t, `{"name":"Lucy","values":[1,2]}`, DescribeValue(taggedStruct))

	assert.Equal(t, `{Name:Mina Values:[1 2]}`, DescribeValue(untaggedStruct))

	assert.Equal(t, `[1, 2]`, DescribeValue([]int{1, 2}))
	assert.Equal(t, `[1, 2]`, DescribeValue([2]int{1, 2}))
	assert.Equal(t, `["a", "b"]`, DescribeValue([]string{"a", "b"}))
}

func TestSomething(t *testing.T) {
	eventData := []string{
		`{"kind": "feature", "value": true}`,
		`{"key": "x", "kind": "custom"}`,
	}
	For(t, "event data").Assert(eventData, ItemsInAnyOrder(
		JSONStrEqual(`{"kind": "custom", "key": "x"}`),
		JSONStrEqual(`{"kind": "feature", "value": true}`),
	))
	For(t, "first event").Assert(eventData[0],
		JSONProperty("kind").Should(Not(Equal("summary"))))
}
