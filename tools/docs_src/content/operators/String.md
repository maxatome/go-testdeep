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


> See also [<i class='fas fa-book'></i> String godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#String).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := td.Cmp(t, got, td.String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.Cmp(t, got, td.String("foobar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
## CmpString shortcut

```go
func CmpString(t TestingT, got interface{}, expected string, args ...interface{}) bool
```

CmpString is a shortcut for:

```go
td.Cmp(t, got, td.String(expected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpString godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#CmpString).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := td.CmpString(t, got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.CmpString(t, got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.CmpString(t, got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
## T.String shortcut

```go
func (t *T) String(got interface{}, expected string, args ...interface{}) bool
```

[`String`]({{< ref "String" >}}) is a shortcut for:

```go
t.Cmp(got, td.String(expected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.String godoc](https://godoc.org/github.com/maxatome/go-testdeep/td#T.String).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := "foobar"

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := td.NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.String(got, "foobar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
