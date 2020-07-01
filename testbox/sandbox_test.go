package testbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSandboxTest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := SandboxTest(func(u TestingT) {
			assert.True(u, true)
		})

		assert.False(t, r.Failed)
		assert.Len(t, r.Failures, 0)

		assert.False(t, r.Skipped)
		assert.Len(t, r.Skips, 0)
	})

	t.Run("failure", func(t *testing.T) {
		r := SandboxTest(func(u TestingT) {
			assert.False(t, u.Failed())
			assert.Equal(u, "abc", "def")
			assert.Fail(u, "another message")
			assert.True(t, u.Failed())
		})

		assert.True(t, r.Failed)
		if assert.Len(t, r.Failures, 2) {
			assert.Nil(t, r.Failures[0].Path)
			assert.Contains(t, r.Failures[0].Message, "abc")
			assert.Nil(t, r.Failures[1].Path)
			assert.Contains(t, r.Failures[1].Message, "another")
		}

		assert.False(t, r.Skipped)
		assert.Len(t, r.Skips, 0)
	})

	t.Run("FailNow", func(t *testing.T) {
		continued := false

		r := SandboxTest(func(u TestingT) {
			u.FailNow()
			continued = true
		})

		assert.True(t, r.Failed)
		assert.Len(t, r.Failures, 0)

		assert.False(t, r.Skipped)
		assert.Len(t, r.Skips, 0)

		assert.False(t, continued)
	})

	t.Run("skip", func(t *testing.T) {
		ran := false
		continued := false

		r := SandboxTest(func(u TestingT) {
			ran = true
			u.Skip("please", "skip")
			continued = true
		})

		assert.True(t, ran)
		assert.False(t, continued)

		assert.False(t, r.Failed)
		assert.Len(t, r.Failures, 0)

		assert.True(t, r.Skipped)
		if assert.Len(t, r.Skips, 1) {
			assert.Nil(t, r.Skips[0].Path)
			assert.Equal(t, "please skip", r.Skips[0].Message)
		}
	})

	t.Run("SkipNow", func(t *testing.T) {
		ran := false
		continued := false

		r := SandboxTest(func(u TestingT) {
			ran = true
			u.SkipNow()
			continued = true
		})

		assert.True(t, ran)
		assert.False(t, continued)

		assert.False(t, r.Failed)
		assert.Len(t, r.Failures, 0)

		assert.True(t, r.Skipped)
		if assert.Len(t, r.Skips, 1) {
			assert.Nil(t, r.Skips[0].Path)
			assert.Equal(t, "", r.Skips[0].Message)
		}
	})
}

func TestSandboxTestSubtests(t *testing.T) {
	t.Run("successes", func(t *testing.T) {
		r := SandboxTest(func(u TestingT) {
			u.Run("sub1", func(uu TestingT) {
				assert.True(uu, true)
			})
			u.Run("sub2", func(uu TestingT) {
				assert.True(uu, true)
			})
		})

		assert.False(t, r.Failed)
		assert.Len(t, r.Failures, 0)

		assert.False(t, r.Skipped)
		assert.Len(t, r.Skips, 0)
	})

	t.Run("failures", func(t *testing.T) {
		r := SandboxTest(func(u TestingT) {
			u.Run("sub1", func(uu TestingT) {
				assert.Equal(uu, "abc", "def")
			})
			u.Run("sub2", func(uu TestingT) {
				assert.Equal(uu, "ghi", "jkl")
			})
		})

		assert.True(t, r.Failed)
		if assert.Len(t, r.Failures, 2) {
			assert.Equal(t, TestPath{"sub1"}, r.Failures[0].Path)
			assert.Contains(t, r.Failures[0].Message, "abc")
			assert.Equal(t, TestPath{"sub2"}, r.Failures[1].Path)
			assert.Contains(t, r.Failures[1].Message, "ghi")
		}

		assert.False(t, r.Skipped)
		assert.Len(t, r.Skips, 0)
	})

	t.Run("successes", func(t *testing.T) {
		r := SandboxTest(func(u TestingT) {
			u.Run("sub1", func(uu TestingT) {
				assert.True(uu, true)
			})
			u.Run("sub2", func(uu TestingT) {
				assert.True(uu, true)
			})
		})

		assert.False(t, r.Failed)
		assert.Len(t, r.Failures, 0)

		assert.False(t, r.Skipped)
		assert.Len(t, r.Skips, 0)
	})

	t.Run("FailNow", func(t *testing.T) {
		continued1 := false
		ran2 := false

		r := SandboxTest(func(u TestingT) {
			u.Run("sub1", func(uu TestingT) {
				assert.True(uu, false)
				uu.FailNow()           // equivalent to require.True(uu, false)
				assert.False(uu, true) // we shouldn't get here
				continued1 = true
			})
			u.Run("sub2", func(uu TestingT) {
				ran2 = true
			})
		})

		assert.False(t, continued1)
		assert.True(t, ran2)

		assert.True(t, r.Failed)
		if assert.Len(t, r.Failures, 1) {
			assert.Equal(t, TestPath{"sub1"}, r.Failures[0].Path)
		}

		assert.False(t, r.Skipped)
		assert.Len(t, r.Skips, 0)
	})

	t.Run("skip", func(t *testing.T) {
		continued1 := false
		ran2 := false

		r := SandboxTest(func(u TestingT) {
			u.Run("sub1", func(uu TestingT) {
				uu.Skip("please", "skip")
				continued1 = true
			})
			u.Run("sub2", func(uu TestingT) {
				ran2 = true
			})
		})

		assert.False(t, continued1)
		assert.True(t, ran2)

		assert.False(t, r.Failed)

		assert.False(t, r.Skipped)
		if assert.Len(t, r.Skips, 1) {
			assert.Equal(t, TestPath{"sub1"}, r.Skips[0].Path)
			assert.Equal(t, "please skip", r.Skips[0].Message)
		}
	})

	t.Run("SkipNow", func(t *testing.T) {
		continued1 := false
		ran2 := false

		r := SandboxTest(func(u TestingT) {
			u.Run("sub1", func(uu TestingT) {
				uu.SkipNow()
				continued1 = true // we shouldn't get here
			})
			u.Run("sub2", func(uu TestingT) {
				ran2 = true
			})
		})

		assert.False(t, continued1)
		assert.True(t, ran2)

		assert.False(t, r.Failed)

		assert.False(t, r.Skipped)
		if assert.Len(t, r.Skips, 1) {
			assert.Equal(t, TestPath{"sub1"}, r.Skips[0].Path)
			assert.Equal(t, "", r.Skips[0].Message)
		}
	})
}
