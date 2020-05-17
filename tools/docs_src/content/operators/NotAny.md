---
title: "NotAny"
weight: 10
---

```go
func NotAny(notExpectedItems ...interface{}) TestDeep
```

[`NotAny`]({{< ref "NotAny" >}}) operator checks that the contents of an array or a slice (or
a pointer on array/slice) does not contain any of *notExpectedItems*.

```go
td.Cmp(t, []int{1}, td.NotAny(1, 2, 3)) // fails
td.Cmp(t, []int{5}, td.NotAny(1, 2, 3)) // succeeds

// works with slices/arrays of any type
td.Cmp(t, personSlice, td.NotAny(
  Person{Name: "Bob", Age: 32},
  Person{Name: "Alice", Age: 26},
))
```

To flatten a non-`[]interface{}` slice/array, use [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten) function
and so avoid boring and inefficient copies:

```go
notExpected := []int{2, 1}
td.Cmp(t, []int{4, 4, 3, 8}, td.NotAny(td.Flatten(notExpected))) // succeeds
// = td.Cmp(t, []int{4, 4, 3, 8}, td.NotAny(2, 1))

notExp1 := []int{2, 1}
notExp2 := []int{5, 8}
td.Cmp(t, []int{4, 4, 42, 8},
  td.NotAny(td.Flatten(notExp1), 3, td.Flatten(notExp2))) // succeeds
// = td.Cmp(t, []int{4, 4, 42, 8}, td.NotAny(2, 1, 3, 5, 8))
```

Beware that [`NotAny(…)`]({{< ref "NotAny" >}}) is not equivalent to [`Not(Any(…)`]({{< ref "Not" >}})) but is like
[`Not(SuperSet(…)`]({{< ref "Not" >}})).


> See also [<i class='fas fa-book'></i> NotAny godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotAny).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{4, 5, 9, 42}

	ok := td.Cmp(t, got, td.NotAny(3, 6, 8, 41, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.NotAny(3, 6, 8, 42, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using notExpected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	notExpected := []int{3, 6, 8, 41, 43}
	ok = td.Cmp(t, got, td.NotAny(td.Flatten(notExpected)),
		"checks %v contains no item listed in notExpected", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
## CmpNotAny shortcut

```go
func CmpNotAny(t TestingT, got interface{}, notExpectedItems []interface{}, args ...interface{}) bool
```

CmpNotAny is a shortcut for:

```go
td.Cmp(t, got, td.NotAny(notExpectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNotAny godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNotAny).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{4, 5, 9, 42}

	ok := td.CmpNotAny(t, got, []interface{}{3, 6, 8, 41, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = td.CmpNotAny(t, got, []interface{}{3, 6, 8, 42, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using notExpected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	notExpected := []int{3, 6, 8, 41, 43}
	ok = td.CmpNotAny(t, got, []interface{}{td.Flatten(notExpected)},
		"checks %v contains no item listed in notExpected", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
## T.NotAny shortcut

```go
func (t *T) NotAny(got interface{}, notExpectedItems []interface{}, args ...interface{}) bool
```

[`NotAny`]({{< ref "NotAny" >}}) is a shortcut for:

```go
t.Cmp(got, td.NotAny(notExpectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.NotAny godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.NotAny).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{4, 5, 9, 42}

	ok := t.NotAny(got, []interface{}{3, 6, 8, 41, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = t.NotAny(got, []interface{}{3, 6, 8, 42, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// When expected is already a non-[]interface{} slice, it cannot be
	// flattened directly using notExpected... without copying it to a new
	// []interface{} slice, then use td.Flatten!
	notExpected := []int{3, 6, 8, 41, 43}
	ok = t.NotAny(got, []interface{}{td.Flatten(notExpected)},
		"checks %v contains no item listed in notExpected", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
