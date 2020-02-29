---
title: "All"
weight: 10
---

```go
func All(expectedValues ...interface{}) TestDeep
```

[`All`]({{< ref "All" >}}) operator compares data against several expected values. During
a match, all of them have to match to succeed. Consider it
as a "AND" logical operator.

```go
td.Cmp(t, "foobar", td.All(
  td.Len(6),
  td.HasPrefix("fo"),
  td.HasSuffix("ar"),
)) // succeeds
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method can return a non-`nil` [`reflect.Type`](https://golang.org/pkg/reflect/#Type) if all items
known non-interface types are equal, or if only interface types
are found (mostly issued from [`Isa()`]({{< ref "Isa" >}})) and they are equal.


> See also [<i class='fas fa-book'></i> All godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#All).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "foo/bar" string
	ok := td.Cmp(t,
		got,
		td.All(td.Re("o/b"), td.HasSuffix("bar"), "foo/bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = td.Cmp(t,
		got,
		td.All(td.Re("o/b"), td.HasSuffix("bar"), "fooX/Ybar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## CmpAll shortcut

```go
func CmpAll(t TestingT, got interface{}, expectedValues []interface{}, args ...interface{}) bool
```

CmpAll is a shortcut for:

```go
td.Cmp(t, got, td.All(expectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpAll godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#CmpAll).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "foo/bar" string
	ok := td.CmpAll(t, got, []interface{}{td.Re("o/b"), td.HasSuffix("bar"), "foo/bar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = td.CmpAll(t, got, []interface{}{td.Re("o/b"), td.HasSuffix("bar"), "fooX/Ybar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## T.All shortcut

```go
func (t *T) All(got interface{}, expectedValues []interface{}, args ...interface{}) bool
```

[`All`]({{< ref "All" >}}) is a shortcut for:

```go
t.Cmp(got, td.All(expectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.All godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#T.All).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := "foo/bar"

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "foo/bar" string
	ok := t.All(got, []interface{}{td.Re("o/b"), td.HasSuffix("bar"), "foo/bar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "o/b" regexp *AND* "bar" suffix *AND* exact "fooX/Ybar" string
	ok = t.All(got, []interface{}{td.Re("o/b"), td.HasSuffix("bar"), "fooX/Ybar"},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
