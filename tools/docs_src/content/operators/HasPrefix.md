---
title: "HasPrefix"
weight: 10
---

```go
func HasPrefix(expected string) TestDeep
```

[`HasPrefix`]({{< ref "HasPrefix" >}}) operator allows to compare the prefix of a `string` (or
convertible), [`error`](https://golang.org/pkg/builtin/#error) or [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interface ([`error`](https://golang.org/pkg/builtin/#error) interface is
tested before [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer).)

```go
type Foobar string
Cmp(t, Foobar("foobar"), HasPrefix("foo")) // succeeds

err := errors.New("error!")
Cmp(t, err, HasPrefix("err")) // succeeds

bstr := bytes.NewBufferString("fmt.Stringer!")
Cmp(t, bstr, HasPrefix("fmt")) // succeeds
```


> See also [<i class='fas fa-book'></i> HasPrefix godoc](https://godoc.org/github.com/maxatome/go-testdeep#HasPrefix).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := Cmp(t, got, HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := Cmp(t, got, HasPrefix("foo"), "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := Cmp(t, got, HasPrefix("foo"), "checks %s", got)
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
Cmp(t, got, HasPrefix(expected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpHasPrefix godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpHasPrefix).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foobar"

	ok := CmpHasPrefix(t, got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := CmpHasPrefix(t, got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foobar")

	ok := CmpHasPrefix(t, got, "foo", "checks %s", got)
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
t.Cmp(got, HasPrefix(expected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.HasPrefix godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.HasPrefix).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	got := "foobar"

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foobar")

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := NewT(&testing.T{})

	got := errors.New("foobar")

	ok := t.HasPrefix(got, "foo", "checks %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
