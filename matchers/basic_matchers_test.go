package matchers

import "testing"

func TestEqual(t *testing.T) {
	assertPasses(t, 3, Equal(3))
	assertFails(t, 4, Equal(3), "did not equal 3\nfull value was: 4")

	assertPasses(t, 3, Equal(float64(3)))
	assertPasses(t, float64(3), Equal(3))

	assertPasses(t, map[string]any{"a": []int{1, 2}},
		Equal(map[string]any{"a": []int{1, 2}}))
}
