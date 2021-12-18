package matchers

import (
	"fmt"
	"reflect"
	"strings"
)

// Items is a matcher for a slice or array value. It tests that the number of elements is equal to
// the number of matchers, and that each element matches the corresponding matcher in order.
//
//     s := []int{6,2}
//     matchers.Items(matchers.Equal(6), matchers.Equal(2)).Test(s) // pass
//     matchers.Items(matchers.Equal(2), matchers.Equal(6)).Test(s) // fail
func Items(matchers ...Matcher) Matcher {
	return New(
		func(value interface{}) bool {
			elements, err := getSliceOrArrayElementValues(value)
			if err != nil || len(elements) != len(matchers) {
				return false
			}
			for i, m := range matchers {
				elementValue := elements[i]
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
			elements, err := getSliceOrArrayElementValues(value)
			if err != nil {
				return err.Error()
			}
			if len(elements) != len(matchers) {
				return fmt.Sprintf("expected slice with %d item(s), got %d item(s)", len(matchers), len(elements))
			}
			parts := make([]string, 0, len(matchers))
			for i, m := range matchers {
				elementValue := elements[i]
				if !m.test(elementValue) {
					parts = append(parts, fmt.Sprintf("item[%d] %s", i, m.describeFailure(elementValue)))
				}
			}
			return strings.Join(parts, ", ")
		},
	)
}

// ItemsInAnyOrder is a matcher for a slice or array value. It tests that the number of elements is
// equal to the number of matchers, and that each matcher matches an element.
//
//     s := []int{6,2}
//     matchers.ItemsInAnyOrder(matchers.Equal(2), matchers.Equal(6)).Test(s) // pass
func ItemsInAnyOrder(matchers ...Matcher) Matcher {
	return New(
		func(value interface{}) bool {
			elements, err := getSliceOrArrayElementValues(value)
			if err != nil || len(elements) != len(matchers) {
				return false
			}
			foundCount := 0
			for _, m := range matchers {
				for _, elementValue := range elements {
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
		func(value interface{}) string {
			// Describing a failure for ItemsInAnyOrder requires us to repeat the matching logic we
			// previously executed, but in a bit more detail. For any matcher that successfully
			// matched an item, we don't need to describe that matcher or that item.
			elements, err := getSliceOrArrayElementValues(value)
			if err != nil {
				return err.Error()
			}
			type unmatchedElement struct {
				index int
				value interface{}
			}
			unmatchedElements := make([]unmatchedElement, 0)
			for i, e := range elements {
				unmatchedElements = append(unmatchedElements, unmatchedElement{index: i, value: e})
			}
			unmatchedMatchers := append([]Matcher(nil), matchers...)
			for i := 0; i < len(unmatchedMatchers); i++ {
				m := unmatchedMatchers[i]
				for j := 0; j < len(unmatchedElements); j++ {
					if m.test(unmatchedElements[j].value) {
						unmatchedElements = append(unmatchedElements[0:j], unmatchedElements[j+1:]...)
						unmatchedMatchers = append(unmatchedMatchers[0:i], unmatchedMatchers[i+1:]...)
						j--
						i--
						break
					}
				}
			}
			if len(unmatchedMatchers) == 0 {
				// Every matcher matched a value but there were some values left over.
				parts := make([]string, 0)
				for _, e := range unmatchedElements {
					parts = append(parts, fmt.Sprintf("[%d]: %s", e.index, DescribeValue(e.value)))
				}
				return fmt.Sprintf("got more items than expected: %s", strings.Join(parts, ", "))
			}
			if len(unmatchedMatchers) == 1 && len(unmatchedElements) == 1 {
				// In this case we'll assume that this was the element that was supposed to be matched
				// by this matcher, but its value was wrong somehow. So we'll print the failure message
				// from the matcher.
				return fmt.Sprintf("failed expectation for one item [%d] with value: %s\nfailure was: %s",
					unmatchedElements[0].index, DescribeValue(unmatchedElements[0].value),
					unmatchedMatchers[0].describeFailure(unmatchedElements[0].value))
			}
			// If there was more than one unmatched matcher and/or element, we can't really guess
			// which one was supposed to go with which, so we'll just print all the unmet conditions.
			return fmt.Sprintf("no items were found to match: %s", describeMatchers(unmatchedMatchers, ", "))
		},
	)
}

func getSliceOrArrayElementValues(sliceValue interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(sliceValue)
	if v.Type().Kind() != reflect.Slice && v.Type().Kind() != reflect.Array {
		return nil, fmt.Errorf("expected slice or array value but got %T", sliceValue)
	}
	ret := make([]interface{}, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		ret = append(ret, v.Index(i).Interface())
	}
	return ret, nil
}
