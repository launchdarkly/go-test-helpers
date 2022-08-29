package jsonhelpers

import (
	"testing"

	"github.com/launchdarkly/go-test-helpers/v2/testbox"
	"github.com/stretchr/testify/assert"
)

func TestAssertEqual(t *testing.T) {
	AssertEqual(t, `{"a":true,"b":false}`, `{"b":false,"a":true}`)

	AssertEqual(t, JValueOf(`{"a":true,"b":false}`), JValueOf(`{"b":false,"a":true}`))

	result := testbox.SandboxTest(func(t testbox.TestingT) {
		AssertEqual(t, `{"a":true,"b":false}`, `{"a":false,"b":false}`)
	})
	assert.True(t, result.Failed)
	if assert.Len(t, result.Failures, 1) {
		assert.Equal(t, `incorrect JSON value: {"a":false,"b":false}
at "a": expected = true, actual = false`, result.Failures[0].Message)
	}

	result = testbox.SandboxTest(func(t testbox.TestingT) {
		AssertEqual(t, `{"a":true,"b":false}`, `{`)
	})
	assert.True(t, result.Failed)
	if assert.Len(t, result.Failures, 1) {
		assert.Equal(t, `invalid actual value (JSON unmarshaling error: unexpected end of JSON input): {`,
			result.Failures[0].Message)
	}

	result = testbox.SandboxTest(func(t testbox.TestingT) {
		AssertEqual(t, `{`, `{"a":true,"b":false}`)
	})
	assert.True(t, result.Failed)
	if assert.Len(t, result.Failures, 1) {
		assert.Equal(t, `invalid expected value (JSON unmarshaling error: unexpected end of JSON input): {`,
			result.Failures[0].Message)
	}
}
