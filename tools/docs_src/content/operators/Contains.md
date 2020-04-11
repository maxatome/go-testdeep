---
title: "Contains"
weight: 10
---

```go
func Contains(expectedValue interface{}) TestDeep
```

[`Contains`]({{< ref "Contains" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}) to check if something is contained
in another thing. [`Contains`]({{< ref "Contains" >}}) has to be applied on arrays, slices, maps or
strings. It tries to be as smarter as possible.

If *expectedValue* is a [TestDeep operator]({{< ref "operators" >}}), each item of data
array/slice/map/`string` (`rune` for strings) is compared to it. The
use of a [TestDeep operator]({{< ref "operators" >}}) as *expectedValue* works only in this
way: item per item.

If data is a slice, and *expectedValue* has the same type, then
*expectedValue* is searched as a sub-slice, otherwise
*expectedValue* is compared to each slice value.

```go
list := []int{12, 34, 28}
td.Cmp(t, list, td.Contains(34))                 // succeeds
td.Cmp(t, list, td.Contains(td.Between(30, 35))) // succeeds too
td.Cmp(t, list, td.Contains(35))                 // fails
td.Cmp(t, list, td.Contains([]int{34, 28}))      // succeeds
```

If data is an array or a map, each value is compared to
*expectedValue*. [`Map`]({{< ref "Map" >}}) keys are not checked: see [`ContainsKey`]({{< ref "ContainsKey" >}}) to check
map keys existence.

```go
hash := map[string]int{"foo": 12, "bar": 34, "zip": 28}
td.Cmp(t, hash, td.Contains(34))                 // succeeds
td.Cmp(t, hash, td.Contains(td.Between(30, 35))) // succeeds too
td.Cmp(t, hash, td.Contains(35))                 // fails

array := [...]int{12, 34, 28}
td.Cmp(t, array, td.Contains(34))                 // succeeds
td.Cmp(t, array, td.Contains(td.Between(30, 35))) // succeeds too
td.Cmp(t, array, td.Contains(35))                 // fails
```

If data is a `string` (or convertible), `[]byte` (or convertible),
[`error`](https://pkg.go.dev/builtin/#error) or [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer) interface ([`error`](https://pkg.go.dev/builtin/#error) interface is tested before
[`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer)), *expectedValue* can be a `string`, a `[]byte`, a `rune` or
a `byte`. In this case, it tests if the got `string` contains this
expected `string`, `[]byte`, `rune` or `byte`.

```go
got := "foo bar"
td.Cmp(t, got, td.Contains('o'))                  // succeeds
td.Cmp(t, got, td.Contains(rune('o')))            // succeeds
td.Cmp(t, got, td.Contains(td.Between('n', 'p'))) // succeeds
td.Cmp(t, got, td.Contains("bar"))                // succeeds
td.Cmp(t, got, td.Contains([]byte("bar")))        // succeeds

td.Cmp(t, []byte("foobar"), td.Contains("ooba")) // succeeds

type Foobar string
td.Cmp(t, Foobar("foobar"), td.Contains("ooba")) // succeeds

err := errors.New("error!")
td.Cmp(t, err, td.Contains("ror")) // succeeds

bstr := bytes.NewBufferString("fmt.Stringer!")
td.Cmp(t, bstr, td.Contains("String")) // succeeds
```

Pitfall: if you want to check if 2 words are contained in got, don't do:

```go
td.Cmp(t, "foobar", td.Contains(td.All("foo", "bar"))) // Bad!
```

as [TestDeep operator]({{< ref "operators" >}}) [`All`]({{< ref "All" >}}) in [`Contains`]({{< ref "Contains" >}}) operates on each `rune`, so it
does not work as expected, but do::

```go
td.Cmp(t, "foobar", td.All(td.Contains("foo"), td.Contains("bar")))
```

When [`Contains(nil)`]({{< ref "Contains" >}}) is used, `nil` is automatically converted to a
typed `nil` on the fly to avoid confusion (if the array/slice/map
item type allows it of course.) So all following Cmp calls
are equivalent (except the `(*byte)(nil)` one):

```go
num := 123
list := []*int{&num, nil}
td.Cmp(t, list, td.Contains(nil))         // succeeds → (*int)(nil)
td.Cmp(t, list, td.Contains((*int)(nil))) // succeeds
td.Cmp(t, list, td.Contains(td.Nil()))    // succeeds
// But...
td.Cmp(t, list, td.Contains((*byte)(nil))) // fails: (*byte)(nil) ≠ (*int)(nil)
```

As well as these ones:

```go
hash := map[string]*int{"foo": nil, "bar": &num}
td.Cmp(t, hash, td.Contains(nil))         // succeeds → (*int)(nil)
td.Cmp(t, hash, td.Contains((*int)(nil))) // succeeds
td.Cmp(t, hash, td.Contains(td.Nil()))    // succeeds
```


> See also [<i class='fas fa-book'></i> Contains godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Contains).

### Examples

{{%expand "ArraySlice example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t, [...]int{11, 22, 33, 44}, td.Contains(22))
	fmt.Println("array contains 22:", ok)

	ok = td.Cmp(t, [...]int{11, 22, 33, 44}, td.Contains(td.Between(20, 25)))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = td.Cmp(t, []int{11, 22, 33, 44}, td.Contains(22))
	fmt.Println("slice contains 22:", ok)

	ok = td.Cmp(t, []int{11, 22, 33, 44}, td.Contains(td.Between(20, 25)))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	ok = td.Cmp(t, []int{11, 22, 33, 44}, td.Contains([]int{22, 33}))
	fmt.Println("slice contains the sub-slice [22, 33]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true
	// slice contains the sub-slice [22, 33]: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := &testing.T{}

	num := 123
	got := [...]*int{&num, nil}

	ok := td.Cmp(t, got, td.Contains(nil))
	fmt.Println("array contains untyped nil:", ok)

	ok = td.Cmp(t, got, td.Contains((*int)(nil)))
	fmt.Println("array contains *int nil:", ok)

	ok = td.Cmp(t, got, td.Contains(td.Nil()))
	fmt.Println("array contains Nil():", ok)

	ok = td.Cmp(t, got, td.Contains((*byte)(nil)))
	fmt.Println("array contains *byte nil:", ok) // types differ: *byte ≠ *int

	// Output:
	// array contains untyped nil: true
	// array contains *int nil: true
	// array contains Nil(): true
	// array contains *byte nil: false

```{{% /expand%}}
{{%expand "Map example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33}, td.Contains(22))
	fmt.Println("map contains value 22:", ok)

	ok = td.Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33},
		td.Contains(td.Between(20, 25)))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := td.Cmp(t, got, td.Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.Cmp(t, got, td.Contains([]byte("oob")), "checks %s", got)
	fmt.Println("contains `oob` []byte:", ok)

	ok = td.Cmp(t, got, td.Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.Cmp(t, got, td.Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.Cmp(t, got, td.Contains(td.Between('n', 'p')), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains `oob` []byte: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.Cmp(t, got, td.Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.Cmp(t, got, td.Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.Cmp(t, got, td.Contains(td.Between('n', 'p')), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.Cmp(t, got, td.Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.Cmp(t, got, td.Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.Cmp(t, got, td.Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.Cmp(t, got, td.Contains(td.Between('n', 'p')), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
## CmpContains shortcut

```go
func CmpContains(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool
```

CmpContains is a shortcut for:

```go
td.Cmp(t, got, td.Contains(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpContains godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpContains).

### Examples

{{%expand "ArraySlice example" %}}```go
	t := &testing.T{}

	ok := td.CmpContains(t, [...]int{11, 22, 33, 44}, 22)
	fmt.Println("array contains 22:", ok)

	ok = td.CmpContains(t, [...]int{11, 22, 33, 44}, td.Between(20, 25))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = td.CmpContains(t, []int{11, 22, 33, 44}, 22)
	fmt.Println("slice contains 22:", ok)

	ok = td.CmpContains(t, []int{11, 22, 33, 44}, td.Between(20, 25))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	ok = td.CmpContains(t, []int{11, 22, 33, 44}, []int{22, 33})
	fmt.Println("slice contains the sub-slice [22, 33]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true
	// slice contains the sub-slice [22, 33]: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := &testing.T{}

	num := 123
	got := [...]*int{&num, nil}

	ok := td.CmpContains(t, got, nil)
	fmt.Println("array contains untyped nil:", ok)

	ok = td.CmpContains(t, got, (*int)(nil))
	fmt.Println("array contains *int nil:", ok)

	ok = td.CmpContains(t, got, td.Nil())
	fmt.Println("array contains Nil():", ok)

	ok = td.CmpContains(t, got, (*byte)(nil))
	fmt.Println("array contains *byte nil:", ok) // types differ: *byte ≠ *int

	// Output:
	// array contains untyped nil: true
	// array contains *int nil: true
	// array contains Nil(): true
	// array contains *byte nil: false

```{{% /expand%}}
{{%expand "Map example" %}}```go
	t := &testing.T{}

	ok := td.CmpContains(t, map[string]int{"foo": 11, "bar": 22, "zip": 33}, 22)
	fmt.Println("map contains value 22:", ok)

	ok = td.CmpContains(t, map[string]int{"foo": 11, "bar": 22, "zip": 33}, td.Between(20, 25))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := td.CmpContains(t, got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.CmpContains(t, got, []byte("oob"), "checks %s", got)
	fmt.Println("contains `oob` []byte:", ok)

	ok = td.CmpContains(t, got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.CmpContains(t, got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.CmpContains(t, got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains `oob` []byte: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.CmpContains(t, got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.CmpContains(t, got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.CmpContains(t, got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.CmpContains(t, got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.CmpContains(t, got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = td.CmpContains(t, got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = td.CmpContains(t, got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = td.CmpContains(t, got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
## T.Contains shortcut

```go
func (t *T) Contains(got interface{}, expectedValue interface{}, args ...interface{}) bool
```

[`Contains`]({{< ref "Contains" >}}) is a shortcut for:

```go
t.Cmp(got, td.Contains(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Contains godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Contains).

### Examples

{{%expand "ArraySlice example" %}}```go
	t := td.NewT(&testing.T{})

	ok := t.Contains([...]int{11, 22, 33, 44}, 22)
	fmt.Println("array contains 22:", ok)

	ok = t.Contains([...]int{11, 22, 33, 44}, td.Between(20, 25))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = t.Contains([]int{11, 22, 33, 44}, 22)
	fmt.Println("slice contains 22:", ok)

	ok = t.Contains([]int{11, 22, 33, 44}, td.Between(20, 25))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	ok = t.Contains([]int{11, 22, 33, 44}, []int{22, 33})
	fmt.Println("slice contains the sub-slice [22, 33]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true
	// slice contains the sub-slice [22, 33]: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := td.NewT(&testing.T{})

	num := 123
	got := [...]*int{&num, nil}

	ok := t.Contains(got, nil)
	fmt.Println("array contains untyped nil:", ok)

	ok = t.Contains(got, (*int)(nil))
	fmt.Println("array contains *int nil:", ok)

	ok = t.Contains(got, td.Nil())
	fmt.Println("array contains Nil():", ok)

	ok = t.Contains(got, (*byte)(nil))
	fmt.Println("array contains *byte nil:", ok) // types differ: *byte ≠ *int

	// Output:
	// array contains untyped nil: true
	// array contains *int nil: true
	// array contains Nil(): true
	// array contains *byte nil: false

```{{% /expand%}}
{{%expand "Map example" %}}```go
	t := td.NewT(&testing.T{})

	ok := t.Contains(map[string]int{"foo": 11, "bar": 22, "zip": 33}, 22)
	fmt.Println("map contains value 22:", ok)

	ok = t.Contains(map[string]int{"foo": 11, "bar": 22, "zip": 33}, td.Between(20, 25))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := td.NewT(&testing.T{})

	got := "foobar"

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = t.Contains(got, []byte("oob"), "checks %s", got)
	fmt.Println("contains `oob` []byte:", ok)

	ok = t.Contains(got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = t.Contains(got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = t.Contains(got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains `oob` []byte: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = t.Contains(got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = t.Contains(got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = t.Contains(got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := td.NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = t.Contains(got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = t.Contains(got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = t.Contains(got, td.Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
