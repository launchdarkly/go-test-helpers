package matchers

import "testing"

func TestItems(t *testing.T) {
	slice := []string{"y", "z", "x"}
	array := [3]string{"y", "z", "x"}

	assertPasses(t, slice, Items(Equal("y"), Equal("z"), Equal("x")))

	assertPasses(t, array, Items(Equal("y"), Equal("z"), Equal("x")))

	assertFails(t, slice, Items(Equal("a"), Equal("b"), Equal("c")),
		`item[0] did not equal "a", item[1] did not equal "b", item[2] did not equal "c"`)

	assertFails(t, slice, Items(Equal("y"), Equal("b"), Equal("x")),
		`item[1] did not equal "b"`)

	assertFails(t, slice, Items(Equal("x"), Equal("y")),
		"expected slice with 2 item(s), got 3 item(s)")

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
		`got more items than expected: [1]: "z"`)

	assertFails(t, slice, ItemsInAnyOrder(Equal("x"), Equal("a"), Equal("z")),
		`failed expectation for one item [0] with value: "y"`+"\n"+
			`failure was: did not equal "a"`)

	assertFails(t, slice, ItemsInAnyOrder(Equal("x"), Equal("a"), Equal("b")),
		`no items were found to match: (equal to "a"), (equal to "b")`)

	assertFails(t, slice, ItemsInAnyOrder(Equal("a"), Equal("b"), Equal("c")),
		`no items were found to match: (equal to "a"), (equal to "b"), (equal to "c")`)

	assertFails(t, 2, ItemsInAnyOrder(Equal("x"), Equal("y")),
		"expected slice or array value but got int\nfull value was: 2")
}

func TestItemsInclude(t *testing.T) {
	slice := []string{"y", "z", "x"}
	array := [3]string{"y", "z", "x"}

	assertPasses(t, slice, ItemsInclude(Equal("x"), Equal("y"), Equal("z")))
	assertPasses(t, array, ItemsInclude(Equal("x"), Equal("y"), Equal("z")))

	assertPasses(t, slice, ItemsInclude(Equal("y")))

	assertPasses(t, slice, ItemsInclude(Equal("y"), Equal("y")))

	assertFails(t, slice, ItemsInclude(Equal("x"), Equal("q"), Equal("r"), Equal("y")),
		`no items were found to match: (equal to "q"), (equal to "r")`)

	assertFails(t, 2, ItemsInclude(Equal("x"), Equal("y")),
		"expected slice or array value but got int\nfull value was: 2")
}

func TestMapOf(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	assertPasses(t, m, MapOf(
		KV("b", Equal(2)),
		KV("a", Equal(1)),
	))

	assertFails(t, m, MapOf(
		KV("b", Equal(3)),
		KV("a", Equal(1)),
	), `key [b] did not equal 3`)

	assertFails(t, m, MapOf(
		KV("b", Equal(2)),
	), `expected map keys [b] but got map keys [a b]`)
}

func TestMapIncluding(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	assertPasses(t, m, MapIncluding(
		KV("b", Equal(2)),
		KV("a", Equal(1)),
	))

	assertPasses(t, m, MapIncluding(
		KV("b", Equal(2)),
	))

	assertFails(t, m, MapIncluding(
		KV("b", Equal(3)),
		KV("a", Equal(1)),
	), `key [b] did not equal 3`)

	assertFails(t, m, MapIncluding(
		KV("c", Equal(3)),
		KV("b", Equal(2)),
		KV("a", Equal(1)),
	), `key [c] not found`)
}

func TestValueForKey(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}

	assertPasses(t, m, ValueForKey("b").Should(Equal(2)))

	assertFails(t, m, ValueForKey("c").Should(Equal(2)), `map key "c" not found`)

	assertFails(t, []int{}, ValueForKey("c").Should(Equal(2)), `expected a map but got []int`)

	assertFails(t, nil, ValueForKey("c").Should(Equal(2)), `map was nil`)
}

func TestOptValueForKey(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]interface{}{"a": 1, "b": 2}

	assertPasses(t, m1, OptValueForKey("b").Should(Equal(2)))

	assertPasses(t, m1, OptValueForKey("c").Should(Equal(0)))

	assertFails(t, []int{}, OptValueForKey("c").Should(Equal(2)), `expected a map but got []int`)

	assertPasses(t, m2, OptValueForKey("c").Should(BeNil()))
}
