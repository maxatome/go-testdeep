---
title: "Between"
weight: 10
---

```go
func Between(from interface{}, to interface{}, bounds ...BoundsKind) TestDeep
```

[`Between`]({{< ref "Between" >}}) operator checks that data is between *from* and
*to*. *from* and *to* can be any numeric, `string` or [`time.Time`](https://pkg.go.dev/time/#Time) (or
assignable) value. *from* and *to* must be the same kind as the
compared value if numeric, and the same type if `string` or [`time.Time`](https://pkg.go.dev/time/#Time) (or
assignable). *bounds* allows to specify whether bounds are included
or not:

- [`BoundsInIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind) (default): between *from* and *to* both included
- [`BoundsInOut`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind): between *from* included and *to* excluded
- [`BoundsOutIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind): between *from* excluded and *to* included
- [`BoundsOutOut`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind): between *from* and *to* both excluded


If *bounds* is missing, it defaults to [`BoundsInIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind).

```go
tc.Cmp(t, 17, td.Between(17, 20))               // succeeds, BoundsInIn by default
tc.Cmp(t, 17, td.Between(10, 17, BoundsInOut))  // fails
tc.Cmp(t, 17, td.Between(10, 17, BoundsOutIn))  // succeeds
tc.Cmp(t, 17, td.Between(17, 20, BoundsOutOut)) // fails
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *from* (same as the *to* one.)


> See also [<i class='fas fa-book'></i> Between godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Between).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.Cmp(t, got, td.Between(154, 156),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsInIn),
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsInOut),
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsOutIn),
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between(154, 156, td.BoundsOutOut),
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

	ok := td.Cmp(t, got, td.Between("aaa", "abc"),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsInIn),
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsInOut),
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsOutIn),
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Between("aaa", "abc", td.BoundsOutOut),
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
td.Cmp(t, got, td.Between(from, to, bounds), args...)
```

See above for details.

[`Between()`]({{< ref "Between" >}}) optional parameter *bounds* is here mandatory.
td.[`BoundsInIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind) value should be passed to mimic its absence in
original [`Between()`]({{< ref "Between" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpBetween godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpBetween).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := td.CmpBetween(t, got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.CmpBetween(t, got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, 154, 156, td.BoundsInOut,
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, 154, 156, td.BoundsOutIn,
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, 154, 156, td.BoundsOutOut,
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

	ok := td.CmpBetween(t, got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsInOut,
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsOutIn,
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = td.CmpBetween(t, got, "aaa", "abc", td.BoundsOutOut,
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
t.Cmp(got, td.Between(from, to, bounds), args...)
```

See above for details.

[`Between()`]({{< ref "Between" >}}) optional parameter *bounds* is here mandatory.
td.[`BoundsInIn`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#BoundsKind) value should be passed to mimic its absence in
original [`Between()`]({{< ref "Between" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Between godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Between).

### Examples

{{%expand "Int example" %}}```go
	t := td.NewT(&testing.T{})

	got := 156

	ok := t.Between(got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = t.Between(got, 154, 156, td.BoundsInIn,
		"checks %v is in [154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, td.BoundsInOut,
		"checks %v is in [154 .. 156[", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, td.BoundsOutIn,
		"checks %v is in ]154 .. 156]", got)
	fmt.Println(ok)

	ok = t.Between(got, 154, 156, td.BoundsOutOut,
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
	t := td.NewT(&testing.T{})

	got := "abc"

	ok := t.Between(got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	// BoundsInIn is implicit
	ok = t.Between(got, "aaa", "abc", td.BoundsInIn,
		`checks "%v" is in ["aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", td.BoundsInOut,
		`checks "%v" is in ["aaa" .. "abc"[`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", td.BoundsOutIn,
		`checks "%v" is in ]"aaa" .. "abc"]`, got)
	fmt.Println(ok)

	ok = t.Between(got, "aaa", "abc", td.BoundsOutOut,
		`checks "%v" is in ]"aaa" .. "abc"[`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true
	// false

```{{% /expand%}}
