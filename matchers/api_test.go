package matchers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type decoratedString string

func (s decoratedString) String() string { return decorate(string(s)) }

func decorate(value interface{}) string { return fmt.Sprintf("Hi, I'm '%s'", value.(string)) }

func assertPasses(t *testing.T, value interface{}, m Matcher) {
	pass, desc := m.Test(value)
	assert.True(t, pass)
	assert.Equal(t, "", desc)
}

func assertFails(t *testing.T, value interface{}, m Matcher, expectedDesc string) {
	pass, desc := m.Test(value)
	assert.False(t, pass)
	assert.Equal(t, expectedDesc, desc)
}

type fakeTestScope struct {
	failures   []string
	terminated bool
}

func (t *fakeTestScope) Errorf(format string, args ...interface{}) {
	t.failures = append(t.failures, fmt.Sprintf(format, args...))
}

func (t *fakeTestScope) FailNow() {
	t.terminated = true
}

func TestUninitializedMatcher(t *testing.T) {
	m := Matcher{}
	assertPasses(t, "whatever", m)
}

func TestSimpleMatcher(t *testing.T) {
	m := New(
		func(value interface{}) bool { return value == "good" },
		func() string { return "should be good" },
		nil,
	)
	assertPasses(t, "good", m)
	assertFails(t, "bad", m, `expected: should be good`+"\n"+`full value was: "bad"`)
}

func TestSimpleMatcherWithFailureDescription(t *testing.T) {
	m := New(
		func(value interface{}) bool { return value == "good" },
		func() string { return "should be good" },
		func(interface{}) string { return "was not good" },
	)
	assertPasses(t, "good", m)
	assertFails(t, "bad", m, `was not good`+"\n"+`full value was: "bad"`)
}

func TestAssertThat(t *testing.T) {
	test1 := fakeTestScope{}
	AssertThat(&test1, 2, Equal(2))
	assert.Len(t, test1.failures, 0)
	assert.False(t, test1.terminated)

	test2 := fakeTestScope{}
	AssertThat(&test2, 3, Equal(2))
	AssertThat(&test2, 4, Equal(2))
	require.Len(t, test2.failures, 2)
	assert.False(t, test2.terminated)
	assert.Contains(t, test2.failures[0], "did not equal 2")
	assert.Contains(t, test2.failures[0], "full value was: 3")
	assert.Contains(t, test2.failures[1], "did not equal 2")
	assert.Contains(t, test2.failures[1], "full value was: 4")
}

func TestRequireThat(t *testing.T) {
	test1 := fakeTestScope{}
	RequireThat(&test1, 2, Equal(2))
	assert.Len(t, test1.failures, 0)
	assert.False(t, test1.terminated)

	test2 := fakeTestScope{}
	RequireThat(&test2, 3, Equal(2))
	assert.Len(t, test2.failures, 1)
	assert.True(t, test2.terminated)
	assert.Contains(t, test2.failures[0], "full value was: 3")
}

func TestEnsureType(t *testing.T) {
	m := New(
		func(value interface{}) bool { return value == "good" },
		func() string { return "should be good" },
		nil,
	)
	assertPasses(t, "good", m)
	assertFails(t, 3, m, "expected: should be good\nfull value was: 3")

	m1 := m.EnsureType("example string")
	assertPasses(t, "good", m1)
	assertFails(t, "bad", m1, `expected: should be good`+"\n"+`full value was: "bad"`)
	assertFails(t, 3, m1, "expected value of type string, was int\nfull value was: 3")

	m2 := m.EnsureType(nil) // no-op
	assertPasses(t, "good", m2)
	assertFails(t, 3, m2, "expected: should be good\nfull value was: 3")
}

func TestMatcherValueDescriptionForStringerType(t *testing.T) {
	m := New(
		func(value interface{}) bool { return value == decoratedString("good") },
		func() string { return "should be good" },
		nil,
	)
	assertFails(t, decoratedString("bad"), m,
		fmt.Sprintf("expected: should be good\nfull value was: %s", decorate("bad")))
}
