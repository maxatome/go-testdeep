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

Note [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten) function can be used to group or reuse some values or
operators and so avoid boring and inefficient copies:

```go
stringOps := td.Flatten([]td.TestDeep{td.HasPrefix("fo"), td.HasSuffix("ar")})
td.Cmp(t, "foobar", td.All(
  td.Len(6),
  stringOps,
)) // succeeds
```

One can do the same with [`All`]({{< ref "All" >}}) operator itself:

```go
stringOps := td.All(td.HasPrefix("fo"), td.HasSuffix("ar"))
td.Cmp(t, "foobar", td.All(
  td.Len(6),
  stringOps,
)) // succeeds
```

but if an [`error`](https://pkg.go.dev/builtin/#error) occurs in the nested [`All`]({{< ref "All" >}}), the report is a bit more
complex to read due to the nested level. [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten) does not create a
new level, its slice is just flattened in the [`All`]({{< ref "All" >}}) parameters.

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method can return a non-`nil` [`reflect.Type`](https://pkg.go.dev/reflect/#Type) if all items
known non-interface types are equal, or if only interface types
are found (mostly issued from [`Isa()`]({{< ref "Isa" >}})) and they are equal.


> See also [<i class='fas fa-book'></i> All godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#All).

### Example

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

	// When some operators or values have to be reused and mixed between
	// several calls, Flatten can be used to avoid boring and
	// inefficient []interface{} copies:
	regOps := td.Flatten([]td.TestDeep{td.Re("o/b"), td.Re(`^fo`), td.Re(`ar$`)})
	ok = td.Cmp(t,
		got,
		td.All(td.HasPrefix("foo"), regOps, td.HasSuffix("bar")),
		"checks all operators against value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

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
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpAll godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpAll).

### Example

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

	// When some operators or values have to be reused and mixed between
	// several calls, Flatten can be used to avoid boring and
	// inefficient []interface{} copies:
	regOps := td.Flatten([]td.TestDeep{td.Re("o/b"), td.Re(`^fo`), td.Re(`ar$`)})
	ok = td.CmpAll(t, got, []interface{}{td.HasPrefix("foo"), regOps, td.HasSuffix("bar")},
		"checks all operators against value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

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
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.All godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.All).

### Example

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

	// When some operators or values have to be reused and mixed between
	// several calls, Flatten can be used to avoid boring and
	// inefficient []interface{} copies:
	regOps := td.Flatten([]td.TestDeep{td.Re("o/b"), td.Re(`^fo`), td.Re(`ar$`)})
	ok = t.All(got, []interface{}{td.HasPrefix("foo"), regOps, td.HasSuffix("bar")},
		"checks all operators against value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
