package matchers

import "testing"

func TestItems(t *testing.T) {
	slice := []string{"y", "z", "x"}

	assertPasses(t, slice, Items(Equal("y"), Equal("z"), Equal("x")))

	assertFails(t, slice, Items(Equal("a"), Equal("b"), Equal("c")),
		`item[0] did not equal "a", item[1] did not equal "b", item[2] did not equal "c"`+
			"\nfull value was: [y z x]")

	assertFails(t, slice, Items(Equal("y"), Equal("b"), Equal("x")),
		`item[1] did not equal "b"`+
			"\nfull value was: [y z x]")

	assertFails(t, slice, Items(Equal("x"), Equal("y")),
		"expected slice with 2 item(s), got 3 item(s)\nfull value was: [y z x]")
}

func TestItemsInAnyOrder(t *testing.T) {
	slice := []string{"y", "z", "x"}

	assertPasses(t, slice, ItemsInAnyOrder(Equal("x"), Equal("y"), Equal("x")))
	assertPasses(t, slice, ItemsInAnyOrder(Equal("y"), Equal("z"), Equal("x")))

	assertFails(t, slice, ItemsInAnyOrder(Equal("x"), Equal("y")),
		`expected: items in any order: (equal to "x"), (equal to "y")`+"\nfull value was: [y z x]")

	assertFails(t, slice, ItemsInAnyOrder(Equal("x"), Equal("a"), Equal("z")),
		`expected: items in any order: (equal to "x"), (equal to "a"), (equal to "z")`+
			"\nfull value was: [y z x]")
}
