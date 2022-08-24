package helpers

import (
	"testing"
	"time"

	"github.com/launchdarkly/go-test-helpers/v2/testbox"

	"github.com/stretchr/testify/assert"
)

func TestTryReceive(t *testing.T) {
	ch := make(chan string, 1)
	v, ok, closed := TryReceive(ch, time.Millisecond)
	assert.False(t, ok)
	assert.False(t, closed)
	assert.Equal(t, "", v)

	ch <- "a"
	v, ok, closed = TryReceive(ch, time.Millisecond)
	assert.True(t, ok)
	assert.False(t, closed)
	assert.Equal(t, "a", v)

	go func() {
		close(ch)
	}()
	v, ok, closed = TryReceive(ch, time.Second)
	assert.False(t, ok)
	assert.True(t, closed)
	assert.Equal(t, "", v)
}

func TestRequireValue(t *testing.T) {
	testbox.ShouldFailAndExitEarly(t, func(t testbox.TestingT) {
		ch := make(chan string, 1)
		_ = RequireValue(t, ch, time.Millisecond)
	})

	ch := make(chan string, 1)
	go func() {
		ch <- "a"
	}()
	v := RequireValue(t, ch, time.Second)
	assert.Equal(t, "a", v)

	testbox.ShouldFailAndExitEarly(t, func(t testbox.TestingT) {
		ch := make(chan string, 1)
		go func() {
			close(ch)
		}()
		_ = RequireValue(t, ch, time.Second)
	})
}

func TestAssertNoMoreValues(t *testing.T) {
	ch := make(chan string, 1)
	AssertNoMoreValues(t, ch, time.Millisecond)

	testbox.ShouldFail(t, func(t testbox.TestingT) {
		ch := make(chan string, 1)
		go func() {
			ch <- "a"
		}()
		AssertNoMoreValues(t, ch, time.Second)
	})

	testbox.ShouldFail(t, func(t testbox.TestingT) {
		ch := make(chan string, 1)
		go func() {
			close(ch)
		}()
		AssertNoMoreValues(t, ch, time.Second)
	})
}

func TestAssertChannelClosed(t *testing.T) {
	ch := make(chan string, 1)
	go func() {
		close(ch)
	}()
	AssertChannelClosed(t, ch, time.Second)

	testbox.ShouldFail(t, func(t testbox.TestingT) {
		ch := make(chan string, 1)
		AssertChannelClosed(t, ch, time.Millisecond)
	})

	testbox.ShouldFail(t, func(t testbox.TestingT) {
		ch := make(chan string, 1)
		ch <- "a"
		AssertChannelClosed(t, ch, time.Millisecond)
	})
}

func TestAssertChannelNotClosed(t *testing.T) {
	testbox.ShouldFail(t, func(t testbox.TestingT) {
		ch := make(chan string, 1)
		go func() {
			close(ch)
		}()
		AssertChannelNotClosed(t, ch, time.Second)
	})

	ch := make(chan string, 1)
	AssertChannelNotClosed(t, ch, time.Millisecond)

	ch <- "a"
	AssertChannelNotClosed(t, ch, time.Millisecond)
}

func TestFailureMessages(t *testing.T) {
	result := testbox.SandboxTest(func(t testbox.TestingT) {
		ch := make(chan string, 1)
		_ = RequireValue(t, ch, time.Millisecond, "sorry%s", ".")
	})
	if assert.Len(t, result.Failures, 2) {
		assert.Equal(t, "expected a string value from channel but did not receive one in 1ms", result.Failures[0].Message)
		assert.Equal(t, "sorry.", result.Failures[1].Message)
	}

	result = testbox.SandboxTest(func(t testbox.TestingT) {
		ch := make(chan string, 1)
		go func() {
			ch <- "a"
		}()
		AssertNoMoreValues(t, ch, time.Second, "sorry%s", ".")
	})
	if assert.Len(t, result.Failures, 2) {
		assert.Equal(t, "expected no more string values from channel but got one: a", result.Failures[0].Message)
		assert.Equal(t, "sorry.", result.Failures[1].Message)
	}

	result = testbox.SandboxTest(func(t testbox.TestingT) {
		ch := make(chan string, 1)
		go func() {
			close(ch)
		}()
		AssertNoMoreValues(t, ch, time.Second, "sorry%s", ".")
	})
	if assert.Len(t, result.Failures, 2) {
		assert.Equal(t, "channel was unexpectedly closed", result.Failures[0].Message)
		assert.Equal(t, "sorry.", result.Failures[1].Message)
	}

	result = testbox.SandboxTest(func(t testbox.TestingT) {
		ch := make(chan string, 1)
		AssertChannelClosed(t, ch, time.Millisecond, "sorry%s", ".")
	})
	if assert.Len(t, result.Failures, 2) {
		assert.Equal(t, "expected channel to be closed within 1ms but it was not", result.Failures[0].Message)
		assert.Equal(t, "sorry.", result.Failures[1].Message)
	}

	result = testbox.SandboxTest(func(t testbox.TestingT) {
		ch := make(chan string, 1)
		go func() {
			close(ch)
		}()
		AssertChannelNotClosed(t, ch, time.Second, "sorry%s", ".")
	})
	if assert.Len(t, result.Failures, 2) {
		assert.Equal(t, "channel was unexpectedly closed", result.Failures[0].Message)
		assert.Equal(t, "sorry.", result.Failures[1].Message)
	}
}
