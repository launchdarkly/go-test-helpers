// Package matchers provides a flexible test assertion API similar to Java's Hamcrest. Matchers are
// constructed separately from the values being tested, and can then be applied to any value, or
// negated, or combined in various ways.
//
// This implementation is for Go 1.17 so it does not yet have generics. Instead, all matchers take
// values of type interface{} and must explicitly cast the type if needed. The simplest way to
// provide type safety is to use Matcher.EnsureType().
//
// Examples of syntax:
//
//     import m "github.com/launchdarkly/go-test-helpers/matchers"
//
//     func TestSomething(t *T) {
//         m.AssertThat(t, "")
//     }
package matchers
