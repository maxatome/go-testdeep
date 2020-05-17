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
td.Cmp(t, []int{1, 1, 2}, td.SuperBagOf(1))       // succeeds
td.Cmp(t, []int{1, 1, 2}, td.SuperBagOf(1, 1, 1)) // fails, one 1 is missing

// works with slices/arrays of any type
td.Cmp(t, personSlice, td.SuperBagOf(
  Person{Name: "Bob", Age: 32},
  Person{Name: "Alice", Age: 26},
))
```

To flatten a non-`[]interface{}` slice/array, use [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten) function
and so avoid boring and inefficient copies:

```go
expected := []int{1, 2, 1}
td.Cmp(t, []int{1}, td.SuperBagOf(td.Flatten(expected))) // succeeds
// = td.Cmp(t, []int{1}, td.SuperBagOf(1, 2, 1))

exp1 := []int{5, 1, 1}
exp2 := []int{8, 42}
td.Cmp(t, []int{1, 5, 1, 8, 42, 3, 3, 6},
  td.SuperBagOf(td.Flatten(exp1), 3, td.Flatten(exp2))) // succeeds
// = td.Cmp(t, []int{1, 5, 1, 8, 42, 3, 3, 6}, td.SuperBagOf(5, 1, 1, 3, 8, 42))
```


> See also [<i class='fas fa-book'></i> SuperBagOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperBagOf).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.Cmp(t, got, td.SuperBagOf(8, 5, 8),
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.SuperBagOf(td.Gt(5), td.Lte(2)),
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{8, 5, 8}
	ok = td.Cmp(t, got, td.SuperBagOf(td.Flatten(expected)),
		"checks the expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
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
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSuperBagOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSuperBagOf).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := td.CmpSuperBagOf(t, got, []interface{}{8, 5, 8},
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = td.CmpSuperBagOf(t, got, []interface{}{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{8, 5, 8}
	ok = td.CmpSuperBagOf(t, got, []interface{}{td.Flatten(expected)},
		"checks the expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
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
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SuperBagOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SuperBagOf).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SuperBagOf(got, []interface{}{8, 5, 8},
		"checks the items are present, in any order")
	fmt.Println(ok)

	ok = t.SuperBagOf(got, []interface{}{td.Gt(5), td.Lte(2)},
		"checks at least 2 items of %v match", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{8, 5, 8}
	ok = t.SuperBagOf(got, []interface{}{td.Flatten(expected)},
		"checks the expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
