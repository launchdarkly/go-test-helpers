package matchers

import (
	"fmt"
	"reflect"
	"strings"
)

// Equal is a matcher that tests whether the input value matches the expected value according
// to reflect.DeepEqual, except in the case of numbers where an exact type match is not needed.
func Equal(expectedValue interface{}) Matcher {
	return New(
		func(value interface{}) bool {
			return reflect.DeepEqual(canonicalizeValue(value), canonicalizeValue(expectedValue))
		},
		func() string {
			return fmt.Sprintf("equal to %s", DescribeValue(expectedValue))
		},
		func(value interface{}) string {
			return fmt.Sprintf("did not equal %s", DescribeValue(expectedValue))
		},
	)
}

// Automatic numeric type conversion for use with Equals(), to avoid the common problem of
// expecting for instance int(1) in a JSON data structure which parsed it as float64(1).
func canonicalizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case uint8:
		return uint64(v)
	case uint:
		return uint64(v)
	case int8:
		return float64(v)
	case int:
		return float64(v)
	case float32:
		return float64(v)
	}
	return value
}

// StringContains is a matcher for string values that tests for the presence of a substring,
// case-sensitively.
func StringContains(substring string) Matcher {
	return New(
		func(value interface{}) bool {
			return strings.Contains(value.(string), substring)
		},
		func() string {
			return fmt.Sprintf("contains %q", substring)
		},
		func(interface{}) string {
			return fmt.Sprintf("did not contain %q", substring)
		},
	).EnsureType("")
}
