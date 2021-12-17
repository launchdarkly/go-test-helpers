package matchers

import (
	"testing"
)

func TestNot(t *testing.T) {
	assertPasses(t, "bad", Not(Equal("good")))
	assertFails(t, "good", Not(Equal("good")), `expected: not (equal to "good")`+"\n"+`full value was: "good"`)
}

func TestAllOf(t *testing.T) {
	hasA := StringContains("A")
	hasB := StringContains("B")
	assertPasses(t, "an A and a B", AllOf(hasA, hasB))
	assertFails(t, "a B", AllOf(hasA, hasB), `did not contain "A"`+"\n"+`full value was: "a B"`)
	assertFails(t, "an A", AllOf(hasA, hasB), `did not contain "B"`+"\n"+`full value was: "an A"`)
	assertFails(t, "a C", AllOf(hasA, hasB), `did not contain "A", did not contain "B"`+"\n"+`full value was: "a C"`)
}

func TestAnyOf(t *testing.T) {
	hasA := StringContains("A")
	hasB := StringContains("B")
	assertPasses(t, "an A and a B", AnyOf(hasA, hasB))
	assertPasses(t, "a B", AnyOf(hasA, hasB))
	assertPasses(t, "an A", AnyOf(hasA, hasB))
	assertFails(t, "a C", AnyOf(hasA, hasB), `did not contain "A", did not contain "B"`+"\n"+`full value was: "a C"`)
}
