package matchers

import (
	"fmt"
	"reflect"
	"strings"
)

// Items is a matcher for a slice value. It tests that the slice contains the same number of elements
// as the number of parameters, and that each parameter matches the corresponding matcher in order.
//
//     s := []int{6,2}
//     matchers.Items(matchers.Equal(6), matchers.Equal(2)).Test(s) // pass
//     matchers.Items(matchers.Equal(2), matchers.Equal(6)).Test(s) // fail
func Items(matchers ...Matcher) Matcher {
	return New(
		func(value interface{}) bool {
			v := reflect.ValueOf(value)
			if v.Type().Kind() != reflect.Slice {
				return false
			}
			if v.Len() != len(matchers) {
				return false
			}
			for i, m := range matchers {
				elementValue := v.Index(i).Interface()
				if !m.test(elementValue) {
					return false
				}
			}
			return true
		},
		func() string {
			return "items: " + describeMatchers(matchers, "")
		},
		func(value interface{}) string {
			v := reflect.ValueOf(value)
			if v.Type().Kind() != reflect.Slice {
				return "a slice"
			}
			if v.Len() != len(matchers) {
				return fmt.Sprintf("expected slice with %d item(s), got %d item(s)", len(matchers), v.Len())
			}
			parts := make([]string, 0, len(matchers))
			for i, m := range matchers {
				elementValue := v.Index(i).Interface()
				if !m.test(elementValue) {
					parts = append(parts, fmt.Sprintf("item[%d] %s", i, m.describeFailure(elementValue)))
				}
			}
			return strings.Join(parts, ", ")
		},
	)
}

// ItemsInAnyOrder is a matcher for a slice value. It tests that the slice contains the same number of
// elements as the number of parameters, and that each parameter is a matcher that matches one item in
// the slice.
//
//     s := []int{6,2}
//     matchers.ItemsInAnyOrder(matchers.Equal(2), matchers.Equal(6)).Test(s) // pass
func ItemsInAnyOrder(matchers ...Matcher) Matcher {
	return New(
		func(value interface{}) bool {
			v := reflect.ValueOf(value)
			if v.Type().Kind() != reflect.Slice {
				return false
			}
			if v.Len() != len(matchers) {
				return false
			}
			foundCount := 0
			for _, m := range matchers {
				for j := 0; j < v.Len(); j++ {
					elementValue := v.Index(j).Interface()
					if m.test(elementValue) {
						foundCount++
						break
					}
				}
			}
			return foundCount == len(matchers)
		},
		func() string {
			return "items in any order: " + describeMatchers(matchers, ", ")
		},
		nil,
		// func(value interface{}) string {
		// 	// It should be possible to make a better failure message where it lists the specific
		// 	// matchers that weren't found, and/or the non-matched items. That will be particularly
		// 	// helpful for lists of events. For now, it's just spitting out the whole condition.
		// 	v := reflect.ValueOf(value)
		// 	if v.Type().Kind() != reflect.Slice {
		// 		return "a slice"
		// 	}
		// 	if v.Len() != len(matchers) {
		// 		return fmt.Sprintf("should have %d item(s) (had %d)", len(matchers), v.Len())
		// 	}
		// 	return "contains in any order: " + desc(matchers, value, ", ")
		// },
	)
}
