---
title: "NotAny"
weight: 10
---

```go
func NotAny(expectedItems ...interface{}) TestDeep
```

[`NotAny`]({{< ref "NotAny" >}}) operator checks that the contents of an array or a slice (or
a pointer on array/slice) does not contain any of *expectedItems*.

```go
Cmp(t, []int{1}, NotAny(1, 2, 3)) // fails
Cmp(t, []int{5}, NotAny(1, 2, 3)) // succeeds
```

Beware that [`NotAny(…)`]({{< ref "NotAny" >}}) is not equivalent to [`Not(Any(…)`]({{< ref "Not" >}})).


> See also [<i class='fas fa-book'></i> NotAny godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#NotAny).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{4, 5, 9, 42}

	ok := td.Cmp(t, got, td.NotAny(3, 6, 8, 41, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.NotAny(3, 6, 8, 42, 43),
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## CmpNotAny shortcut

```go
func CmpNotAny(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

CmpNotAny is a shortcut for:

```go
td.Cmp(t, got, td.NotAny(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNotAny godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#CmpNotAny).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{4, 5, 9, 42}

	ok := td.CmpNotAny(t, got, []interface{}{3, 6, 8, 41, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = td.CmpNotAny(t, got, []interface{}{3, 6, 8, 42, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## T.NotAny shortcut

```go
func (t *T) NotAny(got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

[`NotAny`]({{< ref "NotAny" >}}) is a shortcut for:

```go
t.Cmp(got, td.NotAny(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.NotAny godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#T.NotAny).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{4, 5, 9, 42}

	ok := t.NotAny(got, []interface{}{3, 6, 8, 41, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	ok = t.NotAny(got, []interface{}{3, 6, 8, 42, 43},
		"checks %v contains no item listed in NotAny()", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
