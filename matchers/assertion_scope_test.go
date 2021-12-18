package matchers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestAssertionScopeFor(t *testing.T) {
	test1 := fakeTestScope{}
	In(&test1).For("a").Assert(2, Equal(3))
	require.Len(t, test1.failures, 1)
	assert.Regexp(t, "^a: did not equal 2", test1.failures[0])

	test2 := fakeTestScope{}
	For(&test1, "a").Assert(2, Equal(3))
	require.Len(t, test2.failures, 1)
	assert.Regexp(t, "^a: did not equal 2", test2.failures[0])

	test3 := fakeTestScope{}
	In(&test3).For("a").For("b").Assert(2, Equal(3))
	require.Len(t, test3.failures, 1)
	assert.Regexp(t, "^a: b: did not equal 2", test3.failures[0])
}

func TestAssert(t *testing.T) {
	test1 := fakeTestScope{}
	In(&test1).Assert(2, Equal(2))
	assert.Len(t, test1.failures, 0)
	assert.False(t, test1.terminated)

	test2 := fakeTestScope{}
	In(&test2).Assert(3, Equal(2))
	In(&test2).Assert(4, Equal(2))
	require.Len(t, test2.failures, 2)
	assert.False(t, test2.terminated)
	assert.Equal(t, "did not equal 2\nfull value was: 3", test2.failures[0])
	assert.Equal(t, "did not equal 2\nfull value was: 4", test2.failures[1])

	test3 := fakeTestScope{}
	For(&test3, "score").Assert(3, Equal(2))
	require.Len(t, test3.failures, 1)
	assert.False(t, test3.terminated)
	assert.Equal(t, "score did not equal 2\nfull value was: 3", test3.failures[0])
}

func TestRequire(t *testing.T) {
	test1 := fakeTestScope{}
	In(&test1).Require(2, Equal(2))
	assert.Len(t, test1.failures, 0)
	assert.False(t, test1.terminated)

	test2 := fakeTestScope{}
	In(&test2).Require(3, Equal(2))
	assert.Len(t, test2.failures, 1)
	assert.True(t, test2.terminated)
	assert.Equal(t, "did not equal 2\nfull value was: 3", test2.failures[0])

	test3 := fakeTestScope{}
	For(&test3, "score").Require(3, Equal(2))
	require.Len(t, test3.failures, 1)
	assert.True(t, test3.terminated)
	assert.Equal(t, "score did not equal 2\nfull value was: 3", test3.failures[0])
}
