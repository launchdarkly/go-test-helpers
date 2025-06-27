package matchers

import (
	"encoding/json"
	"testing"
)

func TestJSONEqual(t *testing.T) {
	t.Run("simple values", func(t *testing.T) {
		for _, value := range []any{nil, true, false, 3, 3.5, "x"} {
			jsonValue, _ := json.Marshal(value)
			t.Run(string(jsonValue), func(t *testing.T) {
				assertPasses(t, value, JSONEqual(value))
				assertPasses(t, jsonValue, JSONEqual(value))
				assertPasses(t, value, JSONEqual(jsonValue))
				assertPasses(t, string(jsonValue), JSONStrEqual(string(jsonValue)))
			})
		}
	})

	t.Run("deep equality", func(t *testing.T) {
		for _, value := range []any{
			[]string{"a", "b"},
			map[string]any{"a": []int{1, 2}},
		} {
			jsonValue, _ := json.Marshal(value)
			t.Run(string(jsonValue), func(t *testing.T) {
				assertPasses(t, value, JSONEqual(value))
				assertPasses(t, jsonValue, JSONEqual(value))
				assertPasses(t, value, JSONEqual(jsonValue))
				assertPasses(t, value, JSONStrEqual(string(jsonValue)))
				assertPasses(t, string(jsonValue), JSONStrEqual(string(jsonValue)))
			})
		}

		assertPasses(t, []byte(`{"a": true, "b": false}`),
			JSONEqual([]byte(`{"b": false, "a": true}`)))
	})

	t.Run("inequality with basic message", func(t *testing.T) {
		assertFails(t, true, JSONEqual(3), "expected: JSON equal to 3\nfull value was: true")
		assertFails(t, []byte("[1,2]"), JSONEqual(3), "expected: JSON equal to 3\nfull value was: [1,2]")
	})

	t.Run("inequality with detailed diff", func(t *testing.T) {
		assertFails(t, `{"a":1,"b":3}`, JSONStrEqual(`{"a":1,"b":2}`),
			`JSON values at "b": expected = 2, actual = 3`+
				"\n"+`full value was: {"a":1,"b":3}`)
	})
}

func TestJSONProperty(t *testing.T) {
	assertPasses(t, []byte(`{"a":1,"b":2}`), JSONProperty("b").Should(Equal(2)))

	assertFails(t, []byte(`{"a":1,"b":2}`), JSONProperty("b").Should(Equal(3)),
		`JSON property "b" did not equal 3`+"\n"+`full value was: {"a":1,"b":2}`)

	assertFails(t, []byte(`{"a":1,"b":2}`), JSONProperty("c").Should(Equal(3)),
		`JSON property "c" not found`+"\n"+`full value was: {"a":1,"b":2}`)
}

func TestJSONOptProperty(t *testing.T) {
	assertPasses(t, []byte(`{"a":1,"b":2}`), JSONOptProperty("b").Should(Equal(2)))

	assertFails(t, []byte(`{"a":1,"b":2}`), JSONOptProperty("b").Should(Equal(3)),
		`JSON property "b" did not equal 3`+"\n"+`full value was: {"a":1,"b":2}`)

	assertFails(t, []byte(`{"a":1,"b":2}`), JSONOptProperty("c").Should(Equal(3)),
		`JSON property "c" did not equal 3`+"\n"+`full value was: {"a":1,"b":2}`)

	assertPasses(t, []byte(`{"a":1,"b":2}`), JSONOptProperty("c").Should(Equal(nil)))
}

func TestJSONArray(t *testing.T) {
	assertPasses(t, []byte(`[true, false]`), JSONArray().Should(Length().Should(Equal(2))))

	assertPasses(t, []byte(`[true, false]`), JSONArray().Should(Items(Equal(true), Equal(false))))

	assertFails(t, []byte(`{"a":1,"b":2}`), JSONArray().Should(Length().Should(Equal(2))),
		`wanted a JSON array but found a different type`)
}

func TestJSONMap(t *testing.T) {
	assertPasses(t, []byte(`{"a":1,"b":2}`), JSONMap().Should(Length().Should(Equal(2))))

	assertPasses(t, []byte(`{"a":1,"b":2}`), JSONMap().Should(ValueForKey("a").Should(Equal(1))))

	assertFails(t, []byte(`[true, false]`), JSONMap().Should(Length().Should(Equal(2))),
		`wanted a JSON object but found a different type`)
}
