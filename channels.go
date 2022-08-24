package helpers

import (
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TryReceive returns the next value from the channel and true if successful, or returns an empty
// value and false if the timeout expires first.
func TryReceive[V any](ch <-chan V, timeout time.Duration) (V, bool) {
	select {
	case v := <-ch:
		return v, true
	case <-time.After(timeout):
		var empty V
		return empty, false
	}
}

// RequireValue returns the next value from the channel, or forces an immediate test failure
// and exit if the timeout expires first.
func RequireValue[V any](t require.TestingT, ch <-chan V, timeout time.Duration) V {
	if v, ok := TryReceive(ch, timeout); ok {
		return v
	}
	var empty V
	t.Errorf("expected a %T value from channel but did not receive one in %s", empty, timeout)
	t.FailNow()
	return empty // never reached
}

// AssertNoMoreValues asserts that no value is available from the channel within the timeout.
func AssertNoMoreValues[V any](t assert.TestingT, ch <-chan V, timeout time.Duration) bool {
	if v, ok := TryReceive(ch, timeout); ok {
		t.Errorf("expected no more %T values from channel but got one: %+v", v, v)
		return false
	}
	return true
}

// RequireNoMoreValues is equivalent to AssertNoMoreValues except that it forces an immediate
// test exit on failure.
func RequireNoMoreValues[V any](t require.TestingT, ch <-chan V, timeout time.Duration) {
	if !AssertNoMoreValues(t, ch, timeout) {
		t.FailNow()
	}
}
