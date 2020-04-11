---
title: "SubSetOf"
weight: 10
---

```go
func SubSetOf(expectedItems ...interface{}) TestDeep
```

[`SubSetOf`]({{< ref "SubSetOf" >}}) operator compares the contents of an array or a slice (or a
pointer on array/slice) ignoring duplicates and without taking care
of the order of items.

During a match, each array/slice item should be matched by an
expected item to succeed. But some expected items can be missing
from the compared array/slice.

```go
td.Cmp(t, []int{1, 1}, td.SubSetOf(1, 2))    // succeeds
td.Cmp(t, []int{1, 1, 2}, td.SubSetOf(1, 3)) // fails, 2 is an extra item
```


> See also [<i class='fas fa-book'></i> SubSetOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubSetOf).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are expected, ignoring duplicates
	ok := td.Cmp(t, got, td.SubSetOf(1, 2, 3, 4, 5, 6, 7, 8),
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several SubSetOf entries
	ok = td.Cmp(t, got, td.SubSetOf(td.Between(1, 4), 3, td.Between(2, 10), td.Gt(100)),
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpSubSetOf shortcut

```go
func CmpSubSetOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

CmpSubSetOf is a shortcut for:

```go
td.Cmp(t, got, td.SubSetOf(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSubSetOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSubSetOf).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are expected, ignoring duplicates
	ok := td.CmpSubSetOf(t, got, []interface{}{1, 2, 3, 4, 5, 6, 7, 8},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several SubSetOf entries
	ok = td.CmpSubSetOf(t, got, []interface{}{td.Between(1, 4), 3, td.Between(2, 10), td.Gt(100)},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.SubSetOf shortcut

```go
func (t *T) SubSetOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

[`SubSetOf`]({{< ref "SubSetOf" >}}) is a shortcut for:

```go
t.Cmp(got, td.SubSetOf(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SubSetOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SubSetOf).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are expected, ignoring duplicates
	ok := t.SubSetOf(got, []interface{}{1, 2, 3, 4, 5, 6, 7, 8},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several SubSetOf entries
	ok = t.SubSetOf(got, []interface{}{td.Between(1, 4), 3, td.Between(2, 10), td.Gt(100)},
		"checks at least all items are present, in any order, ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
