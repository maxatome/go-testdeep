---
title: "Any"
weight: 10
---

```go
func Any(expectedValues ...interface{}) TestDeep
```

[`Any`]({{< ref "Any" >}}) operator compares data against several expected values. During
a match, at least one of them has to match to succeed. Consider it
as a "OR" logical operator.

```go
td.Cmp(t, "foo", td.Any("bar", "foo", "zip")) // succeeds
td.Cmp(t, "foo", td.Any(
  td.Len(4),
  td.HasPrefix("f"),
  td.HasSuffix("z"),
)) // succeeds coz "f" prefix
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method can return a non-`nil` [`reflect.Type`](https://golang.org/pkg/reflect/#Type) if all items
known non-interface types are equal, or if only interface types
are found (mostly issued from [`Isa()`]({{< ref "Isa" >}})) and they are equal.


> See also [<i class='fas fa-book'></i> Any godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#Any).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := td.Cmp(t, got, td.Any(td.Re("zip"), td.HasSuffix("bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = td.Cmp(t, got, td.Any(td.Re("zip"), td.HasSuffix("foo")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## CmpAny shortcut

```go
func CmpAny(t TestingT, got interface{}, expectedValues []interface{}, args ...interface{}) bool
```

CmpAny is a shortcut for:

```go
td.Cmp(t, got, td.Any(expectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpAny godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#CmpAny).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := td.CmpAny(t, got, []interface{}{td.Re("zip"), td.HasSuffix("bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = td.CmpAny(t, got, []interface{}{td.Re("zip"), td.HasSuffix("foo")},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## T.Any shortcut

```go
func (t *T) Any(got interface{}, expectedValues []interface{}, args ...interface{}) bool
```

[`Any`]({{< ref "Any" >}}) is a shortcut for:

```go
t.Cmp(got, td.Any(expectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Any godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#T.Any).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := t.Any(got, []interface{}{td.Re("zip"), td.HasSuffix("bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = t.Any(got, []interface{}{td.Re("zip"), td.HasSuffix("foo")},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
