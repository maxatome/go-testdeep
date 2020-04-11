---
title: "HasPrefix"
weight: 10
---

```go
func HasPrefix(expected string) TestDeep
```

[`HasPrefix`]({{< ref "HasPrefix" >}}) operator allows to compare the prefix of a `string` (or
convertible), `[]byte` (or convertible), [`error`](https://pkg.go.dev/builtin/#error) or [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer)
interface ([`error`](https://pkg.go.dev/builtin/#error) interface is tested before [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer)).

```go
td.Cmp(t, []byte("foobar"), td.HasPrefix("foo")) // succeeds

type Foobar string
td.Cmp(t, Foobar("foobar"), td.HasPrefix("foo")) // succeeds

err := errors.New("error!")
td.Cmp(t, err, td.HasPrefix("err")) // succeeds

bstr := bytes.NewBufferString("fmt.Stringer!")
td.Cmp(t, bstr, td.HasPrefix("fmt")) // succeeds
```


> See also [<i class='fas fa-book'></i> HasPrefix godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#HasPrefix).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := td.Cmp(t, got, td.HasPrefix("foo"), "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.HasPrefix("foo"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.Cmp(t, got, td.HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
## CmpHasPrefix shortcut

```go
func CmpHasPrefix(t TestingT, got interface{}, expected string, args ...interface{}) bool
```

CmpHasPrefix is a shortcut for:

```go
td.Cmp(t, got, td.HasPrefix(expected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpHasPrefix godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpHasPrefix).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := td.CmpHasPrefix(t, got, "foo", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = td.Cmp(t, []byte(got), td.HasPrefix("foo"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := td.CmpHasPrefix(t, got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := td.CmpHasPrefix(t, got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
## T.HasPrefix shortcut

```go
func (t *T) HasPrefix(got interface{}, expected string, args ...interface{}) bool
```

[`HasPrefix`]({{< ref "HasPrefix" >}}) is a shortcut for:

```go
t.Cmp(got, td.HasPrefix(expected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.HasPrefix godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.HasPrefix).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := "foobar"

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println("using string:", ok)

	ok = t.Cmp([]byte(got), td.HasPrefix("foo"), "checks %s", got)
	fmt.Println("using []byte:", ok)

	// Output:
	// using string: true
	// using []byte: true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := td.NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
