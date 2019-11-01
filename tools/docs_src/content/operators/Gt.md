---
title: "Gt"
weight: 10
---

```go
func Gt(minExpectedValue interface{}) TestDeep
```

[`Gt`]({{< ref "Gt" >}}) operator checks that data is greater than
*minExpectedValue*. *minExpectedValue* can be any numeric or
[`time.Time`](https://golang.org/pkg/time/#Time) (or assignable) value. *minExpectedValue* must be the
same kind as the compared value if numeric, and the same type if
[`time.Time`](https://golang.org/pkg/time/#Time) (or assignable).

[TypeBehind]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *minExpectedValue*.


> See also [<i class='fas fa-book'></i> Gt godoc](https://godoc.org/github.com/maxatome/go-testdeep#Gt).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Gt(155), "checks %v is > 155", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gt(156), "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := Cmp(t, got, Gt("abb"), `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gt("abc"), `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## CmpGt shortcut

```go
func CmpGt(t TestingT, got interface{}, minExpectedValue interface{}, args ...interface{}) bool
```

CmpGt is a shortcut for:

```go
Cmp(t, got, Gt(minExpectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpGt godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpGt).

### Examples

{{%expand "Int example" %}}```go
	t := &testing.T{}

	got := 156

	ok := CmpGt(t, got, 155, "checks %v is > 155", got)
	fmt.Println(ok)

	ok = CmpGt(t, got, 156, "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	got := "abc"

	ok := CmpGt(t, got, "abb", `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = CmpGt(t, got, "abc", `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## T.Gt shortcut

```go
func (t *T) Gt(got interface{}, minExpectedValue interface{}, args ...interface{}) bool
```

[`Gt`]({{< ref "Gt" >}}) is a shortcut for:

```go
t.Cmp(got, Gt(minExpectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Gt godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Gt).

### Examples

{{%expand "Int example" %}}```go
	t := NewT(&testing.T{})

	got := 156

	ok := t.Gt(got, 155, "checks %v is > 155", got)
	fmt.Println(ok)

	ok = t.Gt(got, 156, "checks %v is > 156", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := NewT(&testing.T{})

	got := "abc"

	ok := t.Gt(got, "abb", `checks "%v" is > "abb"`, got)
	fmt.Println(ok)

	ok = t.Gt(got, "abc", `checks "%v" is > "abc"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
