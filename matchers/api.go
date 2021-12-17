package matchers

import (
	"fmt"
	"reflect"
)

// TestingT is an interface for any test scope type that has an Errorf method for reporting
// failures. This is compatible with Go's testing.T, and with assert.TestingT. See
// AssertThat.
type TestingT interface {
	Errorf(format string, args ...interface{})
}

// RequireT is an interface for any test scope type that has an Errorf method for reporting
// failures and also a FailNow method for stopping the test immediately. This is compatible
// with Go's testing.T, and with require.TestingT. See RequireThat.
type RequireT interface {
	TestingT
	FailNow()
}

// TestFunc is a function used in defining a new Matcher. It returns true if the value passes
// the test or false for failure.
type TestFunc func(value interface{}) bool

// DescribeTestFunc is a function used in defining a new Matcher. It returns a description of
// the test expectation.
type DescribeTestFunc func() string

// DescribeFailureFunc is a function used in defining a new Matcher. Given the value that was
// tested, and assuming that the test failed, it returns a descriptive string.
//
// For simple conditions, this function can be omitted or can return an empty string, in which
// case the failure description will be produced from only the DescribeTestFunc and a
// description of the test value
//
// The second parameter is the function to use for making a string description of a value of
// the expected type.
type DescribeFailureFunc func(value interface{}) string

// Matcher is a general mechanism for declaring expectations about a value. Expectations can be combined,
// and they self-describe on failure.
type Matcher struct {
	testFn            TestFunc
	describeTestFn    DescribeTestFunc
	describeFailureFn DescribeFailureFunc
}

// New creates a Matcher.
func New(
	test TestFunc,
	describeTest DescribeTestFunc,
	describeFailure DescribeFailureFunc,
) Matcher {
	return Matcher{testFn: test, describeTestFn: describeTest, describeFailureFn: describeFailure}
}

// Test executes the expectation for a specific value. It returns true if the value passes the
// test or false for failure, plus a string describing the expectation that failed.
func (m Matcher) Test(value interface{}) (pass bool, failDescription string) {
	if m.test(value) {
		return true, ""
	}
	var failureDesc string
	if m.describeFailureFn != nil {
		failureDesc = m.describeFailureFn(value)
	}
	if failureDesc == "" {
		failureDesc = fmt.Sprintf("expected: %s", m.describeTest())
	}
	return false, fmt.Sprintf("%s\nfull value was: %s", failureDesc, DescribeValue(value))
}

func (m Matcher) test(value interface{}) bool {
	if m.testFn == nil {
		return true
	}
	return m.testFn(value)
}

func (m Matcher) describeTest() string {
	if m.describeTestFn == nil {
		return "[no description given for assertion]"
	}
	return m.describeTestFn()
}

func (m Matcher) describeFailure(value interface{}) string {
	if m.describeFailureFn != nil {
		return m.describeFailureFn(value)
	}
	return m.describeTest()
}

// AssertThat is for use with any test framework that has a test scope type with the same Errorf
// method as Go's testing.T. It tests a value against a matcher and, on failure, calls the test
// scope's Errorf method. This logs a failure but does not stop the test.
func AssertThat(t TestingT, value interface{}, matcher Matcher) bool {
	if pass, desc := matcher.Test(value); !pass {
		t.Errorf("%s", desc)
		return false
	}
	return true
}

// RequireThat is for use with any test framework that has a test scope type with the same Errorf
// and FailNow methods as Go's testing.T. It tests a value against a matcher and, on failure, calls
// the test scope's Errorf method and then FailNow. This logs a failure and immediately terminates
// the test.
func RequireThat(t RequireT, value interface{}, matcher Matcher) bool {
	if pass, desc := matcher.Test(value); !pass {
		t.Errorf("%s", desc)
		t.FailNow()
		return false // does not return since FailNow() will force an early exit
	}
	return true
}

// EnsureType adds type safety to a matcher. The valueOfType parameter should be any value of the
// expected type. The returned Matcher will guarantee that the value is of that type before calling
// the original test function, so it is safe for the test function to cast the value.
func (m Matcher) EnsureType(valueOfType interface{}) Matcher {
	return New(
		func(value interface{}) bool {
			if valueOfType != nil && (reflect.TypeOf(value) != reflect.TypeOf(valueOfType)) {
				return false
			}
			return m.test(value)
		},
		m.describeTest,
		func(value interface{}) string {
			if valueOfType != nil && reflect.TypeOf(value) != reflect.TypeOf(valueOfType) {
				return fmt.Sprintf("expected value of type %T, was %T", valueOfType, value)
			}
			if m.describeFailureFn == nil {
				return ""
			}
			return m.describeFailure(value)
		},
	)
}
