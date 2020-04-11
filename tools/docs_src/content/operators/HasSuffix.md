---
title: "HasSuffix"
weight: 10
---

```go
func HasSuffix(expected string) TestDeep
```

[`HasSuffix`]({{< ref "HasSuffix" >}}) operator allows to compare the suffix of a `string` (or
convertible), `[]byte` (or convertible), [`error`](https://pkg.go.dev/builtin/#error) or [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer)
interface ([`error`](https://pkg.go.dev/builtin/#error) interface is tested before [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer)).

```go
td.Cmp(t, []byte("foobar"), td.HasSuffix("bar")) // succeeds

type Foobar string
td.Cmp(t, Foobar("foobar"), td.HasSuffix("bar")) // succeeds

err := errors.New("error!")
td.Cmp(t, err, td.HasSuffix("!")) // succeeds

bstr := bytes.NewBufferString("fmt.Stringer!")
td.Cmp(t, bstr, td.HasSuffix("!")) // succeeds
```


> See also [<i class='fas fa-book'></i> HasSuffix godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#HasSuffix).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := td.Cmp(t, got, td.HasSuffix("bar"), "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.HasSuffix("bar"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.HasSuffix("bar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.Cmp(t, got, td.HasSuffix("bar"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
## CmpHasSuffix shortcut

```go
func CmpHasSuffix(t TestingT, got interface{}, expected string, args ...interface{}) bool
```

CmpHasSuffix is a shortcut for:

```go
td.Cmp(t, got, td.HasSuffix(expected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpHasSuffix godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpHasSuffix).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := td.CmpHasSuffix(t, got, "bar", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.HasSuffix("bar"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.CmpHasSuffix(t, got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.CmpHasSuffix(t, got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
## T.HasSuffix shortcut

```go
func (t *T) HasSuffix(got interface{}, expected string, args ...interface{}) bool
```

[`HasSuffix`]({{< ref "HasSuffix" >}}) is a shortcut for:

```go
t.Cmp(got, td.HasSuffix(expected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.HasSuffix godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.HasSuffix).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := "foobar"

	ok := t.HasSuffix(got, "bar", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = t.Cmp([]byte(got), td.HasSuffix("bar"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.HasSuffix(got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := td.NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.HasSuffix(got, "bar", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
