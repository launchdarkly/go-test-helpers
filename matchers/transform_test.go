package matchers

import (
	"errors"
	"strings"
	"testing"
)

func stringLength() MatcherTransform {
	return Transform(
		"string length",
		func(value any) (any, error) { return len(value.(string)), nil },
	)
}

func TestTransform(t *testing.T) {
	m := stringLength().Should(Equal(3))

	assertPasses(t, "abc", m)
	assertFails(t, "abcd", m, `string length did not equal 3`+"\n"+`full value was: "abcd"`)
}

func TestTransformEnsureType(t *testing.T) {
	m := stringLength().EnsureInputValueType("example string").
		Should(Equal(3))

	assertPasses(t, "abc", m)
	assertFails(t, "abcd", m, `string length did not equal 3`+"\n"+`full value was: "abcd"`)
	assertFails(t, 3, m, "expected value of type string, was int\nfull value was: 3")
}

func TestTransformError(t *testing.T) {
	stringLengthForLowercaseStringsOnly := Transform(
		"string length",
		func(value any) (any, error) {
			if strings.ToLower(value.(string)) == value.(string) {
				return len(value.(string)), nil
			}
			return 0, errors.New("was not lowercase")
		},
	)
	m := stringLengthForLowercaseStringsOnly.Should(Equal(3))

	assertPasses(t, "abc", m)
	assertFails(t, "Abc", m, `was not lowercase`+"\n"+`full value was: "Abc"`)
}

func TestLength(t *testing.T) {
	assertPasses(t, [3]int{7, 8, 9}, Length().Should(Equal(3)))

	assertPasses(t, []int{7, 8, 9}, Length().Should(Equal(3)))

	assertPasses(t, "abc", Length().Should(Equal(3)))

	assertPasses(t, map[string]int{"a": 1, "b": 2}, Length().Should(Equal(2)))

	assertFails(t, 3, Length().Should(Equal(3)),
		"matchers.Length() was used for an inapplicable type (int)\nfull value was: 3")
}
