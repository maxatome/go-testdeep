---
title: "SuperBagOf"
weight: 10
---

```go
func SuperBagOf(expectedItems ...interface{}) TestDeep
```

[`SuperBagOf`]({{< ref "SuperBagOf" >}}) operator compares the contents of an array or a slice (or a
pointer on array/slice) without taking care of the order of items.

During a match, each expected item should match in the compared
array/slice. But some items in the compared array/slice may not be
expected.

```go
Cmp(t, []int{1, 1, 2}, SuperBagOf(1))       // succeeds
Cmp(t, []int{1, 1, 2}, SuperBagOf(1, 1, 1)) // fails, one 1 is missing
```


> See also [<i class='fas fa-book'></i> SuperBagOf godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#SuperBagOf).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.Cmp(t, got, td.SuperBagOf(8, 5, 8),
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.SuperBagOf(td.Gt(5), td.Lte(2)),
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpSuperBagOf shortcut

```go
func CmpSuperBagOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

CmpSuperBagOf is a shortcut for:

```go
td.Cmp(t, got, td.SuperBagOf(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSuperBagOf godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#CmpSuperBagOf).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.CmpSuperBagOf(t, got, []interface{}{8, 5, 8},
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = td.CmpSuperBagOf(t, got, []interface{}{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.SuperBagOf shortcut

```go
func (t *T) SuperBagOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

[`SuperBagOf`]({{< ref "SuperBagOf" >}}) is a shortcut for:

```go
t.Cmp(got, td.SuperBagOf(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SuperBagOf godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#T.SuperBagOf).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SuperBagOf(got, []interface{}{8, 5, 8},
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = t.SuperBagOf(got, []interface{}{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
