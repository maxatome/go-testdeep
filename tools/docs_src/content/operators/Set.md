---
title: "Set"
weight: 10
---

```go
func Set(expectedItems ...interface{}) TestDeep
```

[`Set`]({{< ref "Set" >}}) operator compares the contents of an array or a slice (or a
pointer on array/slice) ignoring duplicates and without taking care
of the order of items.

During a match, each expected item should match in the compared
array/slice, and each array/slice item should be matched by an
expected item to succeed.

```go
td.Cmp(t, []int{1, 1, 2}, td.Set(1, 2))    // succeeds
td.Cmp(t, []int{1, 1, 2}, td.Set(2, 1))    // succeeds
td.Cmp(t, []int{1, 1, 2}, td.Set(1, 2, 3)) // fails, 3 is missing

// works with slices/arrays of any type
td.Cmp(t, personSlice, td.Set(
  Person{Name: "Bob", Age: 32},
  Person{Name: "Alice", Age: 26},
))
```

To flatten a non-`[]interface{}` slice/array, use [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten) function
and so avoid boring and inefficient copies:

```go
expected := []int{2, 1}
td.Cmp(t, []int{1, 1, 2}, td.Set(td.Flatten(expected))) // succeeds
// = td.Cmp(t, []int{1, 1, 2}, td.Set(2, 1))

exp1 := []int{2, 1}
exp2 := []int{5, 8}
td.Cmp(t, []int{1, 5, 1, 2, 8, 3, 3},
  td.Set(td.Flatten(exp1), 3, td.Flatten(exp2))) // succeeds
// = td.Cmp(t, []int{1, 5, 1, 2, 8, 3, 3}, td.Set(2, 1, 3, 5, 8))
```


> See also [<i class='fas fa-book'></i> Set godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Set).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present, ignoring duplicates
	ok := td.Cmp(t, got, td.Set(1, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = td.Cmp(t, got, td.Set(1, 2, 2, 2, 2, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = td.Cmp(t, got, td.Set(td.Between(1, 4), 3, td.Between(2, 10)),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5, 8}
	ok = td.Cmp(t, got, td.Set(td.Flatten(expected)),
		"checks all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## CmpSet shortcut

```go
func CmpSet(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

CmpSet is a shortcut for:

```go
td.Cmp(t, got, td.Set(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSet godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSet).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present, ignoring duplicates
	ok := td.CmpSet(t, got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = td.CmpSet(t, got, []interface{}{1, 2, 2, 2, 2, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = td.CmpSet(t, got, []interface{}{td.Between(1, 4), 3, td.Between(2, 10)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5, 8}
	ok = td.CmpSet(t, got, []interface{}{td.Flatten(expected)},
		"checks all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## T.Set shortcut

```go
func (t *T) Set(got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

[`Set`]({{< ref "Set" >}}) is a shortcut for:

```go
t.Cmp(got, td.Set(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Set godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Set).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present, ignoring duplicates
	ok := t.Set(got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = t.Set(got, []interface{}{1, 2, 2, 2, 2, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = t.Set(got, []interface{}{td.Between(1, 4), 3, td.Between(2, 10)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using expected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	expected := []int{1, 2, 3, 5, 8}
	ok = t.Set(got, []interface{}{td.Flatten(expected)},
		"checks all expected items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
