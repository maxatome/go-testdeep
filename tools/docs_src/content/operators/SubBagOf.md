---
title: "SubBagOf"
weight: 10
---

```go
func SubBagOf(expectedItems ...interface{}) TestDeep
```

[`SubBagOf`]({{< ref "SubBagOf" >}}) operator compares the contents of an array or a slice (or a
pointer on array/slice) without taking care of the order of items.

During a match, each array/slice item should be matched by an
expected item to succeed. But some expected items can be missing
from the compared array/slice.

```go
Cmp(t, []int{1}, SubBagOf(1, 1, 2))       // succeeds
Cmp(t, []int{1, 1, 1}, SubBagOf(1, 1, 2)) // fails, one 1 is an extra item
```


> See also [<i class='fas fa-book'></i> SubBagOf godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#SubBagOf).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.Cmp(t, got, td.SubBagOf(0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 8, 9, 9),
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	// got contains one 8 too many
	ok = td.Cmp(t, got, td.SubBagOf(0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 9, 9),
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 2}

	ok = td.Cmp(t, got, td.SubBagOf(
		td.Between(0, 3),
		td.Between(0, 3),
		td.Between(0, 3),
		td.Between(0, 3),
		td.Gt(4),
		td.Gt(4)),
		"checks at least all items match, in any order with TestDeep operators")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
## CmpSubBagOf shortcut

```go
func CmpSubBagOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

CmpSubBagOf is a shortcut for:

```go
td.Cmp(t, got, td.SubBagOf(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSubBagOf godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#CmpSubBagOf).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.CmpSubBagOf(t, got, []interface{}{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	// got contains one 8 too many
	ok = td.CmpSubBagOf(t, got, []interface{}{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 2}

	ok = td.CmpSubBagOf(t, got, []interface{}{td.Between(0, 3), td.Between(0, 3), td.Between(0, 3), td.Between(0, 3), td.Gt(4), td.Gt(4)},
		"checks at least all items match, in any order with TestDeep operators")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
## T.SubBagOf shortcut

```go
func (t *T) SubBagOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

[`SubBagOf`]({{< ref "SubBagOf" >}}) is a shortcut for:

```go
t.Cmp(got, td.SubBagOf(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SubBagOf godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#T.SubBagOf).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SubBagOf(got, []interface{}{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	// got contains one 8 too many
	ok = t.SubBagOf(got, []interface{}{0, 0, 1, 1, 2, 2, 3, 3, 5, 5, 8, 9, 9},
		"checks at least all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 2}

	ok = t.SubBagOf(got, []interface{}{td.Between(0, 3), td.Between(0, 3), td.Between(0, 3), td.Between(0, 3), td.Gt(4), td.Gt(4)},
		"checks at least all items match, in any order with TestDeep operators")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
