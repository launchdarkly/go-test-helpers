package matchers

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	if strings.Contains(expectedDesc, "full value was:") {
		assert.Equal(t, expectedDesc, desc)
	} else {
		assert.Regexp(t, regexp.MustCompile("^"+regexp.QuoteMeta(expectedDesc)), desc)
	}
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

func TestMatcherUsesDescribeValue(t *testing.T) {
	m := New(
		func(value interface{}) bool { return value == decoratedString("good") },
		func() string { return "should be good" },
		nil,
	)
	assertFails(t, decoratedString("bad"), m,
		fmt.Sprintf("expected: should be good\nfull value was: %s", decorate("bad")))
}
