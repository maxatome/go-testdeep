---
title: "Between"
weight: 10
---

```go
func Between(from interface{}, to interface{}, bounds ...BoundsKind) TestDeep
```

[`Between`]({{< ref "Between" >}}) operator checks that data is between *from* and
*to*. *from* and *to* can be any numeric, `string` or [`time.Time`](https://golang.org/pkg/time/#Time) (or
assignable) value. *from* and *to* must be the same kind as the
compared value if numeric, and the same type if `string` or [`time.Time`](https://golang.org/pkg/time/#Time) (or
assignable). *bounds* allows to specify whether bounds are included
or not:

- [`BoundsInIn`](https://godoc.org/github.com/maxatome/go-testdeep#BoundsKind) (default): between *from* and *to* both included
- [`BoundsInOut`](https://godoc.org/github.com/maxatome/go-testdeep#BoundsKind): between *from* included and *to* excluded
- [`BoundsOutIn`](https://godoc.org/github.com/maxatome/go-testdeep#BoundsKind): between *from* excluded and *to* included
- [`BoundsOutOut`](https://godoc.org/github.com/maxatome/go-testdeep#BoundsKind): between *from* and *to* both excluded


If *bounds* is missing, it defaults to [`BoundsInIn`](https://godoc.org/github.com/maxatome/go-testdeep#BoundsKind).

[TypeBehind]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *from* (same as the *to* one.)


### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Between(154, 156),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = Cmp(t, got, Between(154, 156, BoundsInIn),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between(154, 156, BoundsInOut),
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between(154, 156, BoundsOutIn),
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between(154, 156, BoundsOutOut),
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := Cmp(t, got, Between("aaa", "abc"),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = Cmp(t, got, Between("aaa", "abc", BoundsInIn),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between("aaa", "abc", BoundsInOut),
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between("aaa", "abc", BoundsOutIn),
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Between("aaa", "abc", BoundsOutOut),
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
## CmpBetween shortcut

```go
func CmpBetween(t TestingT, got interface{}, from interface{}, to interface{}, bounds BoundsKind, args ...interface{}) bool
```

CmpBetween is a shortcut for:

```go
Cmp(t, got, Between(from, to, bounds), args...)
```

See above for details.

[`Between()`]({{< ref "Between" >}}) optional parameter *bounds* is here mandatory.
[`BoundsInIn`](https://godoc.org/github.com/maxatome/go-testdeep#BoundsKind) value should be passed to mimic its absence in
original [`Between()`]({{< ref "Between" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := CmpBetween(t, got, 154, 156, BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = CmpBetween(t, got, 154, 156, BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = CmpBetween(t, got, 154, 156, BoundsInOut,
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = CmpBetween(t, got, 154, 156, BoundsOutIn,
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = CmpBetween(t, got, 154, 156, BoundsOutOut,
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := CmpBetween(t, got, "aaa", "abc", BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = CmpBetween(t, got, "aaa", "abc", BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = CmpBetween(t, got, "aaa", "abc", BoundsInOut,
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = CmpBetween(t, got, "aaa", "abc", BoundsOutIn,
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = CmpBetween(t, got, "aaa", "abc", BoundsOutOut,
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
## T.Between shortcut

```go
func (t *T) Between(got interface{}, from interface{}, to interface{}, bounds BoundsKind, args ...interface{}) bool
```

[`Between`]({{< ref "Between" >}}) is a shortcut for:

```go
t.Cmp(got, Between(from, to, bounds), args...)
```

See above for details.

[`Between()`]({{< ref "Between" >}}) optional parameter *bounds* is here mandatory.
[`BoundsInIn`](https://godoc.org/github.com/maxatome/go-testdeep#BoundsKind) value should be passed to mimic its absence in
original [`Between()`]({{< ref "Between" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Int example" %}}```go
	t := NewT(&testing.T{})

	got := 156

	ok := t.Between(got, 154, 156, BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = t.Between(got, 154, 156, BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, BoundsInOut,
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, BoundsOutIn,
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, BoundsOutOut,
		"checks %v is in ]154 .. 156[", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := NewT(&testing.T{})

	got := "abc"

	ok := t.Between(got, "aaa", "abc", BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = t.Between(got, "aaa", "abc", BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", BoundsInOut,
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", BoundsOutIn,
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", BoundsOutOut,
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
