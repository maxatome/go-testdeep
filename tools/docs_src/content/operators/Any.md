---
title: "Any"
weight: 10
---

```go
func Any(expectedValues ...interface{}) TestDeep
```

[`Any`]({{< ref "Any" >}}) operator compares data against several expected values. During
a match, at least one of them has to match to succeed.

[TypeBehind]({{< ref "operators#typebehind-method" >}}) method can return a non-`nil` [`reflect.Type`](https://golang.org/pkg/reflect/#Type) if all items
known non-interface types are equal, or if only interface types
are found (mostly issued from [`Isa()`]({{< ref "Isa" >}})) and they are equal.


### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := Cmp(t, got, Any(Re("zip"), HasSuffix("bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = Cmp(t, got, Any(Re("zip"), HasSuffix("foo")),
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
Cmp(t, got, Any(expectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := CmpAny(t, got, []interface{}{Re("zip"), HasSuffix("bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = CmpAny(t, got, []interface{}{Re("zip"), HasSuffix("foo")},
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
t.Cmp(got, Any(expectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	got := "foo/bar"

	// Checks got string against:
	//   "zip" regexp *OR* "bar" suffix
	ok := t.Any(got, []interface{}{Re("zip"), HasSuffix("bar")},
		"checks value %s", got)
	fmt.Println(ok)

	// Checks got string against:
	//   "zip" regexp *OR* "foo" suffix
	ok = t.Any(got, []interface{}{Re("zip"), HasSuffix("foo")},
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
