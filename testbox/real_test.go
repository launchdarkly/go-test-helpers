package testbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRealTest(t *testing.T) {
	// can't test failure cases, since then *this* test would fail

	t.Run("success", func(t *testing.T) {
		rt := RealTest(t)
		assert.True(rt, true)

		assert.False(t, rt.Failed())
		assert.False(t, t.Failed())
		assert.False(t, t.Skipped())
	})

	t.Run("subtest success", func(t *testing.T) {
		ran := false

		rt := RealTest(t)
		rt.Run("sub", func(u TestingT) {
			ran = true
			assert.True(u, true)
		})

		assert.True(t, ran)

		assert.False(t, rt.Failed())
		assert.False(t, t.Failed())
		assert.False(t, t.Skipped())
	})

	t.Run("skip", func(t *testing.T) { // this test will always be reported as skipped
		rt := RealTest(t)
		rt.Skip()

		assert.True(t, false) // won't execute because we exited early on Skip
	})

	t.Run("subtest skip", func(t *testing.T) {
		ran := false
		continued := false

		rt := RealTest(t)
		rt.Run("sub", func(u TestingT) {
			ran = true
			u.Skip("let's skip this")
			continued = true
		})

		assert.True(t, ran)
		assert.False(t, continued)

		assert.False(t, rt.Failed())
		assert.False(t, t.Failed())
		assert.False(t, t.Skipped())
	})

	t.Run("subtest SkipNow", func(t *testing.T) {
		ran := false
		continued := false

		rt := RealTest(t)
		rt.Run("sub", func(u TestingT) {
			ran = true
			u.SkipNow()
			continued = true
		})

		assert.True(t, ran)
		assert.False(t, continued)

		assert.False(t, rt.Failed())
		assert.False(t, t.Failed())
		assert.False(t, t.Skipped())
	})
}
