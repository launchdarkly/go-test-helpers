package helpers

import (
	"testing"
	"time"

	"github.com/launchdarkly/go-test-helpers/v2/testbox"
	"github.com/stretchr/testify/assert"
)

func TestTryReceive(t *testing.T) {
	ch := make(chan string, 1)
	v, ok := TryReceive(ch, time.Millisecond)
	assert.False(t, ok)
	assert.Equal(t, "", v)

	ch <- "a"
	v, ok = TryReceive(ch, time.Millisecond)
	assert.True(t, ok)
	assert.Equal(t, "a", v)
}

func TestRequireValue(t *testing.T) {
	result := testbox.SandboxTest(func(t1 testbox.TestingT) {
		ch := make(chan string, 1)
		_ = RequireValue(t1, ch, time.Millisecond)
		t.Errorf("test should have exited early but did not")
	})
	assert.True(t, result.Failed)

	ch := make(chan string, 1)
	go func() {
		ch <- "a"
	}()
	v := RequireValue(t, ch, time.Second)
	assert.Equal(t, "a", v)
}

func TestAssertNoMoreValues(t *testing.T) {
	ch := make(chan string, 1)
	AssertNoMoreValues(t, ch, time.Millisecond)

	result := testbox.SandboxTest(func(t testbox.TestingT) {
		ch := make(chan string, 1)
		go func() {
			ch <- "a"
		}()
		AssertNoMoreValues(t, ch, time.Second)
	})
	assert.True(t, result.Failed)
}

func TestRequireNoMoreValues(t *testing.T) {
	ch := make(chan string, 1)
	AssertNoMoreValues(t, ch, time.Millisecond)

	result := testbox.SandboxTest(func(t1 testbox.TestingT) {
		ch := make(chan string, 1)
		go func() {
			ch <- "a"
		}()
		RequireNoMoreValues(t1, ch, time.Second)
		t.Errorf("test should have exited early but did not")
	})
	assert.True(t, result.Failed)
}
