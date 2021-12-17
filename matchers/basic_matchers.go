package matchers

import (
	"fmt"
	"reflect"
	"strings"
)

// Equal is a matcher that tests whether the input value matches the expected value according
// to reflect.DeepEqual.
func Equal(expectedValue interface{}) Matcher {
	return New(
		func(value interface{}) bool {
			return reflect.DeepEqual(value, expectedValue)
		},
		func() string {
			return fmt.Sprintf("equal to %s", DescribeValue(expectedValue))
		},
		func(value interface{}) string {
			return fmt.Sprintf("did not equal %s", DescribeValue(expectedValue))
		},
	)
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
