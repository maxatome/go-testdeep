---
title: "Re"
weight: 10
---

```go
func Re(reg interface{}, capture ...interface{}) TestDeep
```

[`Re`]({{< ref "Re" >}}) operator allows to apply a regexp on a `string` (or convertible),
`[]byte`, [`error`](https://pkg.go.dev/builtin/#error) or [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer) interface ([`error`](https://pkg.go.dev/builtin/#error) interface is tested
before [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer).)

*reg* is the regexp. It can be a `string` that is automatically
compiled using [`regexp.MustCompile`](https://pkg.go.dev/regexp/#MustCompile), or a [`*regexp.Regexp`](https://pkg.go.dev/regexp/#Regexp).

Optional *capture* parameter can be used to match the contents of
regexp groups. Groups are presented as a `[]string` or `[][]byte`
depending the original matched data. Note that an other operator
can be used here.

```go
td.Cmp(t, "foobar zip!", td.Re(`^foobar`)) // succeeds
td.Cmp(t, "John Doe",
  td.Re(`^(\w+) (\w+)`, []string{"John", "Doe"})) // succeeds
td.Cmp(t, "John Doe",
  td.Re(`^(\w+) (\w+)`, td.Bag("Doe", "John"))) // succeeds
```


> See also [<i class='fas fa-book'></i> Re godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Re).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foo bar"
	ok := td.Cmp(t, got, td.Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = td.Cmp(t, got, td.Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := td.Cmp(t, got, td.Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foo bar")
	ok := td.Cmp(t, got, td.Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Capture example" %}}```go
	t := &testing.T{}

	got := "foo bar biz"
	ok := td.Cmp(t, got, td.Re(`^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = td.Cmp(t, got, td.Re(`^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "Compiled example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := td.Cmp(t, got, td.Re(expected), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = td.Cmp(t, got, td.Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CompiledStringer example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := td.Cmp(t, got, td.Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "CompiledError example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := td.Cmp(t, got, td.Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "CompiledCapture example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := td.Cmp(t, got, td.Re(expected, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = td.Cmp(t, got, td.Re(expected, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## CmpRe shortcut

```go
func CmpRe(t TestingT, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool
```

CmpRe is a shortcut for:

```go
td.Cmp(t, got, td.Re(reg, capture), args...)
```

See above for details.

[`Re()`]({{< ref "Re" >}}) optional parameter *capture* is here mandatory.
`nil` value should be passed to mimic its absence in
original [`Re()`]({{< ref "Re" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpRe godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpRe).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := "foo bar"
	ok := td.CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = td.CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := td.CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := &testing.T{}

	got := errors.New("foo bar")
	ok := td.CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Capture example" %}}```go
	t := &testing.T{}

	got := "foo bar biz"
	ok := td.CmpRe(t, got, `^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = td.CmpRe(t, got, `^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "Compiled example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := td.CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = td.CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CompiledStringer example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := td.CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "CompiledError example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := td.CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "CompiledCapture example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := td.CmpRe(t, got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = td.CmpRe(t, got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## T.Re shortcut

```go
func (t *T) Re(got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool
```

[`Re`]({{< ref "Re" >}}) is a shortcut for:

```go
t.Cmp(got, td.Re(reg, capture), args...)
```

See above for details.

[`Re()`]({{< ref "Re" >}}) optional parameter *capture* is here mandatory.
`nil` value should be passed to mimic its absence in
original [`Re()`]({{< ref "Re" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Re godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Re).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := "foo bar"
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "Stringer example" %}}```go
	t := td.NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Error example" %}}```go
	t := td.NewT(&testing.T{})

	got := errors.New("foo bar")
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "Capture example" %}}```go
	t := td.NewT(&testing.T{})

	got := "foo bar biz"
	ok := t.Re(got, `^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = t.Re(got, `^(\w+) (\w+) (\w+)$`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "Compiled example" %}}```go
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CompiledStringer example" %}}```go
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "CompiledError example" %}}```go
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "CompiledCapture example" %}}```go
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := t.Re(got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = t.Re(got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
