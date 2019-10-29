---
title: "Lte"
weight: 10
---

```go
func Lte(maxExpectedValue interface{}) TestDeep
```

[`Lte`]({{< ref "Lte" >}}) operator checks that data is lesser or equal than
*maxExpectedValue*. *maxExpectedValue* can be any numeric or
[`time.Time`](https://golang.org/pkg/time/#Time) (or assignable) value. *maxExpectedValue* must be the
same kind as the compared value if numeric, and the same type if
[`time.Time`](https://golang.org/pkg/time/#Time) (or assignable).

[TypeBehind]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *maxExpectedValue*.


### Examples

{{%expand "Int example" %}}	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Lte(156), "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lte(157), "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lte(155), "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
{{%expand "String example" %}}	t := &testing.T{}

	got := "abc"

	ok := Cmp(t, got, Lte("abc"), `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lte("abd"), `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Lte("abb"), `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
## CmpLte shortcut

```go
func CmpLte(t TestingT, got interface{}, maxExpectedValue interface{}, args ...interface{}) bool
```

CmpLte is a shortcut for:

```go
Cmp(t, got, Lte(maxExpectedValue), args...)
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

{{%expand "Int example" %}}	t := &testing.T{}

	got := 156

	ok := CmpLte(t, got, 156, "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = CmpLte(t, got, 157, "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = CmpLte(t, got, 155, "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
{{%expand "String example" %}}	t := &testing.T{}

	got := "abc"

	ok := CmpLte(t, got, "abc", `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = CmpLte(t, got, "abd", `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = CmpLte(t, got, "abb", `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
## T.Lte shortcut

```go
func (t *T) Lte(got interface{}, maxExpectedValue interface{}, args ...interface{}) bool
```

[`Lte`]({{< ref "Lte" >}}) is a shortcut for:

```go
t.Cmp(got, Lte(maxExpectedValue), args...)
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

{{%expand "Int example" %}}	t := NewT(&testing.T{})

	got := 156

	ok := t.Lte(got, 156, "checks %v is ≤ 156", got)
	fmt.Println(ok)

	ok = t.Lte(got, 157, "checks %v is ≤ 157", got)
	fmt.Println(ok)

	ok = t.Lte(got, 155, "checks %v is ≤ 155", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
{{%expand "String example" %}}	t := NewT(&testing.T{})

	got := "abc"

	ok := t.Lte(got, "abc", `checks "%v" is ≤ "abc"`, got)
	fmt.Println(ok)

	ok = t.Lte(got, "abd", `checks "%v" is ≤ "abd"`, got)
	fmt.Println(ok)

	ok = t.Lte(got, "abb", `checks "%v" is ≤ "abb"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
