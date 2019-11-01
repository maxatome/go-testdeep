---
title: "Contains"
weight: 10
---

```go
func Contains(expectedValue interface{}) TestDeep
```

[`Contains`]({{< ref "Contains" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}) with a little convenient exception
for strings. [`Contains`]({{< ref "Contains" >}}) has to be applied on arrays, slices, maps or
strings. It compares each item of data array/slice/map/`string` (`rune`
for strings) against *expectedValue*.

```go
list := []int{12, 34, 28}
Cmp(t, list, Contains(34))              // succeeds
Cmp(t, list, Contains(Between(30, 35))) // succeeds too
Cmp(t, list, Contains(35))              // fails

hash := map[string]int{"foo": 12, "bar": 34, "zip": 28}
Cmp(t, hash, Contains(34))              // succeeds
Cmp(t, hash, Contains(Between(30, 35))) // succeeds too
Cmp(t, hash, Contains(35))              // fails

got := "foo bar"
Cmp(t, got, Contains('o'))               // succeeds
Cmp(t, got, Contains(rune('o')))         // succeeds
Cmp(t, got, Contains(Between('n', 'p'))) // succeeds
```

When [`Contains(nil)`]({{< ref "Contains" >}}) is used, `nil` is automatically converted to a
typed `nil` on the fly to avoid confusion (if the array/slice/map
item type allows it of course.) So all following Cmp calls
are equivalent (except the `(*byte)(nil)` one):

```go
num := 123
list := []*int{&num, nil}
Cmp(t, list, Contains(nil))         // succeeds → (*int)(nil)
Cmp(t, list, Contains((*int)(nil))) // succeeds
Cmp(t, list, Contains(Nil()))       // succeeds
// But...
Cmp(t, list, Contains((*byte)(nil))) // fails: (*byte)(nil) ≠ (*int)(nil)
```

As well as these ones:

```go
hash := map[string]*int{"foo": nil, "bar": &num}
Cmp(t, hash, Contains(nil))         // succeeds → (*int)(nil)
Cmp(t, hash, Contains((*int)(nil))) // succeeds
Cmp(t, hash, Contains(Nil()))       // succeeds
```

As a special case for `string` (or convertible), [`error`](https://golang.org/pkg/builtin/#error) or
[`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interface ([`error`](https://golang.org/pkg/builtin/#error) interface is tested before
[`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer)), *expectedValue* can be a `string`, a `rune` or a
`byte`. In this case, it tests if the got `string` contains this
expected `string`, `rune` or `byte`.

```go
type Foobar string
Cmp(t, Foobar("foobar"), Contains("ooba")) // succeeds

err := errors.New("error!")
Cmp(t, err, Contains("ror")) // succeeds

bstr := bytes.NewBufferString("fmt.Stringer!")
Cmp(t, bstr, Contains("String")) // succeeds
```


### Examples

{{%expand "ArraySlice example" %}}```go
	t := &testing.T{}

	ok := Cmp(t, [...]int{11, 22, 33, 44}, Contains(22))
	fmt.Println("array contains 22:", ok)

	ok = Cmp(t, [...]int{11, 22, 33, 44}, Contains(Between(20, 25)))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = Cmp(t, []int{11, 22, 33, 44}, Contains(22))
	fmt.Println("slice contains 22:", ok)

	ok = Cmp(t, []int{11, 22, 33, 44}, Contains(Between(20, 25)))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := &testing.T{}

	num := 123
	got := [...]*int{&num, nil}

	ok := Cmp(t, got, Contains(nil))
	fmt.Println("array contains untyped nil:", ok)

	ok = Cmp(t, got, Contains((*int)(nil)))
	fmt.Println("array contains *int nil:", ok)

	ok = Cmp(t, got, Contains(Nil()))
	fmt.Println("array contains Nil():", ok)

	ok = Cmp(t, got, Contains((*byte)(nil)))
	fmt.Println("array contains *byte nil:", ok) // types differ: *byte ≠ *int

	// Output:
	// array contains untyped nil: true
	// array contains *int nil: true
	// array contains Nil(): true
	// array contains *byte nil: false

```{{% /expand%}}
{{%expand "Map example" %}}```go
	t := &testing.T{}

	ok := Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33}, Contains(22))
	fmt.Println("map contains value 22:", ok)

	ok = Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33},
		Contains(Between(20, 25)))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := Cmp(t, got, Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = Cmp(t, got, Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = Cmp(t, got, Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = Cmp(t, got, Contains(Between('n', 'p')), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := Cmp(t, got, Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = Cmp(t, got, Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = Cmp(t, got, Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	// Be careful! TestDeep operators in Contains() do not work with
	// fmt.Stringer nor error interfaces
	ok = Cmp(t, got, Contains(Between('n', 'p')), "checks %s", got)
	fmt.Println("try TestDeep operator:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// try TestDeep operator: false

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := Cmp(t, got, Contains("oob"), "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = Cmp(t, got, Contains('b'), "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = Cmp(t, got, Contains(byte('a')), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	// Be careful! TestDeep operators in Contains() do not work with
	// fmt.Stringer nor error interfaces
	ok = Cmp(t, got, Contains(Between('n', 'p')), "checks %s", got)
	fmt.Println("try TestDeep operator:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// try TestDeep operator: false

```{{% /expand%}}
## CmpContains shortcut

```go
func CmpContains(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool
```

CmpContains is a shortcut for:

```go
Cmp(t, got, Contains(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "ArraySlice example" %}}```go
	t := &testing.T{}

	ok := CmpContains(t, [...]int{11, 22, 33, 44}, 22)
	fmt.Println("array contains 22:", ok)

	ok = CmpContains(t, [...]int{11, 22, 33, 44}, Between(20, 25))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = CmpContains(t, []int{11, 22, 33, 44}, 22)
	fmt.Println("slice contains 22:", ok)

	ok = CmpContains(t, []int{11, 22, 33, 44}, Between(20, 25))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := &testing.T{}

	num := 123
	got := [...]*int{&num, nil}

	ok := CmpContains(t, got, nil)
	fmt.Println("array contains untyped nil:", ok)

	ok = CmpContains(t, got, (*int)(nil))
	fmt.Println("array contains *int nil:", ok)

	ok = CmpContains(t, got, Nil())
	fmt.Println("array contains Nil():", ok)

	ok = CmpContains(t, got, (*byte)(nil))
	fmt.Println("array contains *byte nil:", ok) // types differ: *byte ≠ *int

	// Output:
	// array contains untyped nil: true
	// array contains *int nil: true
	// array contains Nil(): true
	// array contains *byte nil: false

```{{% /expand%}}
{{%expand "Map example" %}}```go
	t := &testing.T{}

	ok := CmpContains(t, map[string]int{"foo": 11, "bar": 22, "zip": 33}, 22)
	fmt.Println("map contains value 22:", ok)

	ok = CmpContains(t, map[string]int{"foo": 11, "bar": 22, "zip": 33}, Between(20, 25))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := CmpContains(t, got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = CmpContains(t, got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = CmpContains(t, got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = CmpContains(t, got, Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := CmpContains(t, got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = CmpContains(t, got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = CmpContains(t, got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	// Be careful! TestDeep operators in Contains() do not work with
	// fmt.Stringer nor error interfaces
	ok = CmpContains(t, got, Between('n', 'p'), "checks %s", got)
	fmt.Println("try TestDeep operator:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// try TestDeep operator: false

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := CmpContains(t, got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = CmpContains(t, got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = CmpContains(t, got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	// Be careful! TestDeep operators in Contains() do not work with
	// fmt.Stringer nor error interfaces
	ok = CmpContains(t, got, Between('n', 'p'), "checks %s", got)
	fmt.Println("try TestDeep operator:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// try TestDeep operator: false

```{{% /expand%}}
## T.Contains shortcut

```go
func (t *T) Contains(got interface{}, expectedValue interface{}, args ...interface{}) bool
```

[`Contains`]({{< ref "Contains" >}}) is a shortcut for:

```go
t.Cmp(got, Contains(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "ArraySlice example" %}}```go
	t := NewT(&testing.T{})

	ok := t.Contains([...]int{11, 22, 33, 44}, 22)
	fmt.Println("array contains 22:", ok)

	ok = t.Contains([...]int{11, 22, 33, 44}, Between(20, 25))
	fmt.Println("array contains at least one item in [20 .. 25]:", ok)

	ok = t.Contains([]int{11, 22, 33, 44}, 22)
	fmt.Println("slice contains 22:", ok)

	ok = t.Contains([]int{11, 22, 33, 44}, Between(20, 25))
	fmt.Println("slice contains at least one item in [20 .. 25]:", ok)

	// Output:
	// array contains 22: true
	// array contains at least one item in [20 .. 25]: true
	// slice contains 22: true
	// slice contains at least one item in [20 .. 25]: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := NewT(&testing.T{})

	num := 123
	got := [...]*int{&num, nil}

	ok := t.Contains(got, nil)
	fmt.Println("array contains untyped nil:", ok)

	ok = t.Contains(got, (*int)(nil))
	fmt.Println("array contains *int nil:", ok)

	ok = t.Contains(got, Nil())
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
	t := NewT(&testing.T{})

	ok := t.Contains(map[string]int{"foo": 11, "bar": 22, "zip": 33}, 22)
	fmt.Println("map contains value 22:", ok)

	ok = t.Contains(map[string]int{"foo": 11, "bar": 22, "zip": 33}, Between(20, 25))
	fmt.Println("map contains at least one value in [20 .. 25]:", ok)

	// Output:
	// map contains value 22: true
	// map contains at least one value in [20 .. 25]: true

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := NewT(&testing.T{})

	got := "foobar"

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = t.Contains(got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = t.Contains(got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	ok = t.Contains(got, Between('n', 'p'), "checks %s", got)
	fmt.Println("contains at least one character ['n' .. 'p']:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// contains at least one character ['n' .. 'p']: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = t.Contains(got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = t.Contains(got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	// Be careful! TestDeep operators in Contains() do not work with
	// fmt.Stringer nor error interfaces
	ok = t.Contains(got, Between('n', 'p'), "checks %s", got)
	fmt.Println("try TestDeep operator:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// try TestDeep operator: false

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.Contains(got, "oob", "checks %s", got)
	fmt.Println("contains `oob` string:", ok)

	ok = t.Contains(got, 'b', "checks %s", got)
	fmt.Println("contains 'b' rune:", ok)

	ok = t.Contains(got, byte('a'), "checks %s", got)
	fmt.Println("contains 'a' byte:", ok)

	// Be careful! TestDeep operators in Contains() do not work with
	// fmt.Stringer nor error interfaces
	ok = t.Contains(got, Between('n', 'p'), "checks %s", got)
	fmt.Println("try TestDeep operator:", ok)

	// Output:
	// contains `oob` string: true
	// contains 'b' rune: true
	// contains 'a' byte: true
	// try TestDeep operator: false

```{{% /expand%}}
