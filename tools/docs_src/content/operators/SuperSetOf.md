---
title: "SuperSetOf"
weight: 10
---

```go
func SuperSetOf(expectedItems ...interface{}) TestDeep
```

[`SuperSetOf`]({{< ref "SuperSetOf" >}}) operator compares the contents of an array or a slice (or
a pointer on array/slice) ignoring duplicates and without taking
care of the order of items.

During a match, each expected item should match in the compared
array/slice. But some items in the compared array/slice may not be
expected.

```go
td.Cmp(t, []int{1, 1, 2}, td.SuperSetOf(1))    // succeeds
td.Cmp(t, []int{1, 1, 2}, td.SuperSetOf(1, 3)) // fails, 3 is missing

// works with slices/arrays of any type
td.Cmp(t, personSlice, td.SuperSetOf(
  Person{Name: "Bob", Age: 32},
  Person{Name: "Alice", Age: 26},
))
```

To flatten a non-`[]interface{}` slice/array, use [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten) function
and so avoid boring and inefficient copies:

```go
expected := []int{2, 1}
td.Cmp(t, []int{1, 1, 2, 8}, td.SuperSetOf(td.Flatten(expected))) // succeeds
// = td.Cmp(t, []int{1, 1, 2, 8}, td.SubSetOf(2, 1))

exp1 := []int{2, 1}
exp2 := []int{5, 8}
td.Cmp(t, []int{1, 5, 1, 8, 42, 3, 3},
  td.SuperSetOf(td.Flatten(exp1), 3, td.Flatten(exp2))) // succeeds
// = td.Cmp(t, []int{1, 5, 1, 8, 42, 3, 3}, td.SuperSetOf(2, 1, 3, 5, 8))
```


> See also [<i class='fas fa-book'></i> SuperSetOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperSetOf).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.Cmp(t, got, td.SuperSetOf(1, 2, 3),
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.SuperSetOf(td.Gt(5), td.Lte(2)),
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3}
	ok = td.Cmp(t, got, td.SuperSetOf(td.Flatten(expected)),
		"checks the expected items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
## CmpSuperSetOf shortcut

```go
func CmpSuperSetOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

CmpSuperSetOf is a shortcut for:

```go
td.Cmp(t, got, td.SuperSetOf(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSuperSetOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSuperSetOf).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.CmpSuperSetOf(t, got, []interface{}{1, 2, 3},
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = td.CmpSuperSetOf(t, got, []interface{}{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3}
	ok = td.CmpSuperSetOf(t, got, []interface{}{td.Flatten(expected)},
		"checks the expected items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
## T.SuperSetOf shortcut

```go
func (t *T) SuperSetOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

[`SuperSetOf`]({{< ref "SuperSetOf" >}}) is a shortcut for:

```go
t.Cmp(got, td.SuperSetOf(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SuperSetOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SuperSetOf).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SuperSetOf(got, []interface{}{1, 2, 3},
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = t.SuperSetOf(got, []interface{}{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3}
	ok = t.SuperSetOf(got, []interface{}{td.Flatten(expected)},
		"checks the expected items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
