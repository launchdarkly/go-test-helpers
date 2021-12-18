package matchers

import "testing"

func TestStringContains(t *testing.T) {
	assertPasses(t, "abc", StringContains("b"))
	assertFails(t, "abc", StringContains("x"), `did not contain "x"`+"\n"+`full value was: "abc"`)
	assertFails(t, "abc", StringContains("B"), `did not contain "B"`+"\n"+`full value was: "abc"`)
}

func TestStringHasPrefix(t *testing.T) {
	assertPasses(t, "abc", StringHasPrefix("a"))
	assertFails(t, "abc", StringHasPrefix("x"), `did not start with "x"`+"\n"+`full value was: "abc"`)
	assertFails(t, "abc", StringHasPrefix("A"), `did not start with "A"`+"\n"+`full value was: "abc"`)
}

func TestStringHasSuffix(t *testing.T) {
	assertPasses(t, "abc", StringHasSuffix("c"))
	assertFails(t, "abc", StringHasSuffix("x"), `did not end with "x"`+"\n"+`full value was: "abc"`)
	assertFails(t, "abc", StringHasSuffix("C"), `did not end with "C"`+"\n"+`full value was: "abc"`)
}
