---
title: "String"
weight: 10
---

```go
func String(expected string) TestDeep
```

[`String`]({{< ref "String" >}}) operator allows to compare a `string` (or convertible), [`error`](https://golang.org/pkg/builtin/#error)
or [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interface ([`error`](https://golang.org/pkg/builtin/#error) interface is tested before
[`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer).)

```go
err := errors.New("error!")
Cmp(t, err, String("error!")) // succeeds

bstr := bytes.NewBufferString("fmt.Stringer!")
Cmp(t, bstr, String("fmt.Stringer!")) // succeeds
```


### Examples

{{%expand "Base example" %}}	t := &testing.T{}

	got := "foobar"

	ok := Cmp(t, got, String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Stringer example" %}}	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := Cmp(t, got, String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Error example" %}}	t := &testing.T{}

	got := errors.New("foobar")

	ok := Cmp(t, got, String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
## CmpString shortcut

```go
func CmpString(t TestingT, got interface{}, expected string, args ...interface{}) bool
```

CmpString is a shortcut for:

```go
Cmp(t, got, String(expected), args...)
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

{{%expand "Base example" %}}	t := &testing.T{}

	got := "foobar"

	ok := CmpString(t, got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Stringer example" %}}	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := CmpString(t, got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Error example" %}}	t := &testing.T{}

	got := errors.New("foobar")

	ok := CmpString(t, got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
## T.String shortcut

```go
func (t *T) String(got interface{}, expected string, args ...interface{}) bool
```

[`String`]({{< ref "String" >}}) is a shortcut for:

```go
t.Cmp(got, String(expected), args...)
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

{{%expand "Base example" %}}	t := NewT(&testing.T{})

	got := "foobar"

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Stringer example" %}}	t := NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Error example" %}}	t := NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
