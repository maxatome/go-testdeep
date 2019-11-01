---
title: "Lt"
weight: 10
---

```go
func Lt(maxExpectedValue interface{}) TestDeep
```

[`Lt`]({{< ref "Lt" >}}) operator checks that data is lesser than
*maxExpectedValue*. *maxExpectedValue* can be any numeric or
[`time.Time`](https://golang.org/pkg/time/#Time) (or assignable) value. *maxExpectedValue* must be the
same kind as the compared value if numeric, and the same type if
[`time.Time`](https://golang.org/pkg/time/#Time) (or assignable).

[TypeBehind]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *maxExpectedValue*.


### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Lt(157), "checks %v is < 157", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lt(156), "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := Cmp(t, got, Lt("abd"), `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lt("abc"), `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## CmpLt shortcut

```go
func CmpLt(t TestingT, got interface{}, maxExpectedValue interface{}, args ...interface{}) bool
```

CmpLt is a shortcut for:

```go
Cmp(t, got, Lt(maxExpectedValue), args...)
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

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := CmpLt(t, got, 157, "checks %v is < 157", got)
	fmt.Println(ok)

	ok = CmpLt(t, got, 156, "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := CmpLt(t, got, "abd", `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = CmpLt(t, got, "abc", `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## T.Lt shortcut

```go
func (t *T) Lt(got interface{}, maxExpectedValue interface{}, args ...interface{}) bool
```

[`Lt`]({{< ref "Lt" >}}) is a shortcut for:

```go
t.Cmp(got, Lt(maxExpectedValue), args...)
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

{{%expand "Int example" %}}```go
	t := NewT(&testing.T{})

	got := 156

	ok := t.Lt(got, 157, "checks %v is < 157", got)
	fmt.Println(ok)

	ok = t.Lt(got, 156, "checks %v is < 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := NewT(&testing.T{})

	got := "abc"

	ok := t.Lt(got, "abd", `checks "%v" is < "abd"`, got)
	fmt.Println(ok)

	ok = t.Lt(got, "abc", `checks "%v" is < "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
