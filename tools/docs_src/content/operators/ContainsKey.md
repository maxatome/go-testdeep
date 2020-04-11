---
title: "ContainsKey"
weight: 10
---

```go
func ContainsKey(expectedValue interface{}) TestDeep
```

[`ContainsKey`]({{< ref "ContainsKey" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}) and works on maps only. It
compares each key of map against *expectedValue*.

```go
hash := map[string]int{"foo": 12, "bar": 34, "zip": 28}
td.Cmp(t, hash, td.ContainsKey("foo"))             // succeeds
td.Cmp(t, hash, td.ContainsKey(td.HasPrefix("z"))) // succeeds
td.Cmp(t, hash, td.ContainsKey(td.HasPrefix("x"))) // fails

hnum := map[int]string{1: "foo", 42: "bar"}
td.Cmp(t, hash, td.ContainsKey(42))                 // succeeds
td.Cmp(t, hash, td.ContainsKey(td.Between(40, 45))) // succeeds
```

When [`ContainsKey(nil)`]({{< ref "ContainsKey" >}}) is used, `nil` is automatically converted to a
typed `nil` on the fly to avoid confusion (if the map key type allows
it of course.) So all following Cmp calls are equivalent
(except the `(*byte)(nil)` one):

```go
num := 123
hnum := map[*int]bool{&num: true, nil: true}
td.Cmp(t, hnum, td.ContainsKey(nil))         // succeeds → (*int)(nil)
td.Cmp(t, hnum, td.ContainsKey((*int)(nil))) // succeeds
td.Cmp(t, hnum, td.ContainsKey(td.Nil()))    // succeeds
// But...
td.Cmp(t, hnum, td.ContainsKey((*byte)(nil))) // fails: (*byte)(nil) ≠ (*int)(nil)
```


> See also [<i class='fas fa-book'></i> ContainsKey godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ContainsKey).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t,
		map[string]int{"foo": 11, "bar": 22, "zip": 33}, td.ContainsKey("foo"))
	fmt.Println(`map contains key "foo":`, ok)

	ok = td.Cmp(t,
		map[int]bool{12: true, 24: false, 42: true, 51: false},
		td.ContainsKey(td.Between(40, 50)))
	fmt.Println("map contains at least a key in [40 .. 50]:", ok)

	ok = td.Cmp(t,
		map[string]int{"FOO": 11, "bar": 22, "zip": 33},
		td.ContainsKey(td.Smuggle(strings.ToLower, "foo")))
	fmt.Println(`map contains key "foo" without taking case into account:`, ok)

	// Output:
	// map contains key "foo": true
	// map contains at least a key in [40 .. 50]: true
	// map contains key "foo" without taking case into account: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := &testing.T{}

	num := 1234
	got := map[*int]bool{&num: false, nil: true}

	ok := td.Cmp(t, got, td.ContainsKey(nil))
	fmt.Println("map contains untyped nil key:", ok)

	ok = td.Cmp(t, got, td.ContainsKey((*int)(nil)))
	fmt.Println("map contains *int nil key:", ok)

	ok = td.Cmp(t, got, td.ContainsKey(td.Nil()))
	fmt.Println("map contains Nil() key:", ok)

	ok = td.Cmp(t, got, td.ContainsKey((*byte)(nil)))
	fmt.Println("map contains *byte nil key:", ok) // types differ: *byte ≠ *int

	// Output:
	// map contains untyped nil key: true
	// map contains *int nil key: true
	// map contains Nil() key: true
	// map contains *byte nil key: false

```{{% /expand%}}
## CmpContainsKey shortcut

```go
func CmpContainsKey(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool
```

CmpContainsKey is a shortcut for:

```go
td.Cmp(t, got, td.ContainsKey(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpContainsKey godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpContainsKey).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.CmpContainsKey(t, map[string]int{"foo": 11, "bar": 22, "zip": 33}, "foo")
	fmt.Println(`map contains key "foo":`, ok)

	ok = td.CmpContainsKey(t, map[int]bool{12: true, 24: false, 42: true, 51: false}, td.Between(40, 50))
	fmt.Println("map contains at least a key in [40 .. 50]:", ok)

	ok = td.CmpContainsKey(t, map[string]int{"FOO": 11, "bar": 22, "zip": 33}, td.Smuggle(strings.ToLower, "foo"))
	fmt.Println(`map contains key "foo" without taking case into account:`, ok)

	// Output:
	// map contains key "foo": true
	// map contains at least a key in [40 .. 50]: true
	// map contains key "foo" without taking case into account: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := &testing.T{}

	num := 1234
	got := map[*int]bool{&num: false, nil: true}

	ok := td.CmpContainsKey(t, got, nil)
	fmt.Println("map contains untyped nil key:", ok)

	ok = td.CmpContainsKey(t, got, (*int)(nil))
	fmt.Println("map contains *int nil key:", ok)

	ok = td.CmpContainsKey(t, got, td.Nil())
	fmt.Println("map contains Nil() key:", ok)

	ok = td.CmpContainsKey(t, got, (*byte)(nil))
	fmt.Println("map contains *byte nil key:", ok) // types differ: *byte ≠ *int

	// Output:
	// map contains untyped nil key: true
	// map contains *int nil key: true
	// map contains Nil() key: true
	// map contains *byte nil key: false

```{{% /expand%}}
## T.ContainsKey shortcut

```go
func (t *T) ContainsKey(got interface{}, expectedValue interface{}, args ...interface{}) bool
```

[`ContainsKey`]({{< ref "ContainsKey" >}}) is a shortcut for:

```go
t.Cmp(got, td.ContainsKey(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.ContainsKey godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.ContainsKey).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	ok := t.ContainsKey(map[string]int{"foo": 11, "bar": 22, "zip": 33}, "foo")
	fmt.Println(`map contains key "foo":`, ok)

	ok = t.ContainsKey(map[int]bool{12: true, 24: false, 42: true, 51: false}, td.Between(40, 50))
	fmt.Println("map contains at least a key in [40 .. 50]:", ok)

	ok = t.ContainsKey(map[string]int{"FOO": 11, "bar": 22, "zip": 33}, td.Smuggle(strings.ToLower, "foo"))
	fmt.Println(`map contains key "foo" without taking case into account:`, ok)

	// Output:
	// map contains key "foo": true
	// map contains at least a key in [40 .. 50]: true
	// map contains key "foo" without taking case into account: true

```{{% /expand%}}
{{%expand "Nil example" %}}```go
	t := td.NewT(&testing.T{})

	num := 1234
	got := map[*int]bool{&num: false, nil: true}

	ok := t.ContainsKey(got, nil)
	fmt.Println("map contains untyped nil key:", ok)

	ok = t.ContainsKey(got, (*int)(nil))
	fmt.Println("map contains *int nil key:", ok)

	ok = t.ContainsKey(got, td.Nil())
	fmt.Println("map contains Nil() key:", ok)

	ok = t.ContainsKey(got, (*byte)(nil))
	fmt.Println("map contains *byte nil key:", ok) // types differ: *byte ≠ *int

	// Output:
	// map contains untyped nil key: true
	// map contains *int nil key: true
	// map contains Nil() key: true
	// map contains *byte nil key: false

```{{% /expand%}}
