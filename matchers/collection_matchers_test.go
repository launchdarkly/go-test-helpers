package matchers

import "testing"

func TestItems(t *testing.T) {
	slice := []string{"y", "z", "x"}
	array := [3]string{"y", "z", "x"}

	assertPasses(t, slice, Items(Equal("y"), Equal("z"), Equal("x")))

	assertPasses(t, array, Items(Equal("y"), Equal("z"), Equal("x")))

	assertFails(t, slice, Items(Equal("a"), Equal("b"), Equal("c")),
		`item[0] did not equal "a", item[1] did not equal "b", item[2] did not equal "c"`+
			"\n"+`full value was: ["y", "z", "x"]`)

	assertFails(t, slice, Items(Equal("y"), Equal("b"), Equal("x")),
		`item[1] did not equal "b"`+
			"\n"+`full value was: ["y", "z", "x"]`)

	assertFails(t, slice, Items(Equal("x"), Equal("y")),
		"expected slice with 2 item(s), got 3 item(s)\n"+`full value was: ["y", "z", "x"]`)

	assertFails(t, 2, Items(Equal("x"), Equal("y")),
		"expected slice or array value but got int\nfull value was: 2")
}

func TestItemsInAnyOrder(t *testing.T) {
	slice := []string{"y", "z", "x"}
	array := [3]string{"y", "z", "x"}

	assertPasses(t, slice, ItemsInAnyOrder(Equal("x"), Equal("y"), Equal("x")))
	assertPasses(t, slice, ItemsInAnyOrder(Equal("y"), Equal("z"), Equal("x")))

	assertPasses(t, array, ItemsInAnyOrder(Equal("x"), Equal("y"), Equal("x")))

	assertFails(t, slice, ItemsInAnyOrder(Equal("x"), Equal("y")),
		`got more items than expected: [1]: "z"`+"\n"+`full value was: ["y", "z", "x"]`)

	assertFails(t, slice, ItemsInAnyOrder(Equal("x"), Equal("a"), Equal("z")),
		`failed expectation for one item [0] with value: "y"`+"\n"+
			`failure was: did not equal "a"`+
			"\n"+`full value was: ["y", "z", "x"]`)

	assertFails(t, slice, ItemsInAnyOrder(Equal("x"), Equal("a"), Equal("b")),
		`no items were found to match: (equal to "a"), (equal to "b")`+
			"\n"+`full value was: ["y", "z", "x"]`)

	assertFails(t, slice, ItemsInAnyOrder(Equal("a"), Equal("b"), Equal("c")),
		`no items were found to match: (equal to "a"), (equal to "b"), (equal to "c")`+
			"\n"+`full value was: ["y", "z", "x"]`)
}
