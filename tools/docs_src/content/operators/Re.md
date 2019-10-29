---
title: "Re"
weight: 10
---

```go
func Re(reg interface{}, capture ...interface{}) TestDeep
```

[`Re`]({{< ref "Re" >}}) operator allows to apply a regexp on a `string` (or convertible),
`[]byte`, [`error`](https://golang.org/pkg/builtin/#error) or [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interface ([`error`](https://golang.org/pkg/builtin/#error) interface is tested
before [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer).)

*reg* is the regexp. It can be a `string` that is automatically
compiled using [`regexp.MustCompile`](https://golang.org/pkg/regexp/#MustCompile), or a [`*regexp.Regexp`](https://golang.org/pkg/regexp/#Regexp).

Optional *capture* parameter can be used to match the contents of
regexp groups. Groups are presented as a `[]string` or `[][]byte`
depending the original matched data. Note that an other operator
can be used here.

```go
Cmp(t, "foobar zip!", Re(`^foobar`))     // succeeds
Cmp(t, "John Doe",
  Re(`^(\w+) (\w+)`, []string{"John", "Doe"})) // succeeds
Cmp(t, "John Doe",
  Re(`^(\w+) (\w+)`, Bag("Doe", "John"))       // succeeds
```


### Examples

{{%expand "Base example" %}}	t := &testing.T{}

	got := "foo bar"
	ok := Cmp(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = Cmp(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "Stringer example" %}}	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := Cmp(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Error example" %}}	t := &testing.T{}

	got := errors.New("foo bar")
	ok := Cmp(t, got, Re("(zip|bar)$"), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Capture example" %}}	t := &testing.T{}

	got := "foo bar biz"
	ok := Cmp(t, got, Re(`^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = Cmp(t, got, Re(`^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "Compiled example" %}}	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := Cmp(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = Cmp(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CompiledStringer example" %}}	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := Cmp(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "CompiledError example" %}}	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := Cmp(t, got, Re(expected), "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "CompiledCapture example" %}}	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := Cmp(t, got, Re(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = Cmp(t, got, Re(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
## CmpRe shortcut

```go
func CmpRe(t TestingT, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool
```

CmpRe is a shortcut for:

```go
Cmp(t, got, Re(reg, capture), args...)
```

See above for details.

[`Re()`]({{< ref "Re" >}}) optional parameter *capture* is here mandatory.
`nil` value should be passed to mimic its absence in
original [`Re()`]({{< ref "Re" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Base example" %}}	t := &testing.T{}

	got := "foo bar"
	ok := CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "Stringer example" %}}	t := &testing.T{}

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Error example" %}}	t := &testing.T{}

	got := errors.New("foo bar")
	ok := CmpRe(t, got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Capture example" %}}	t := &testing.T{}

	got := "foo bar biz"
	ok := CmpRe(t, got, `^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = CmpRe(t, got, `^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "Compiled example" %}}	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := "foo bar"
	ok := CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CompiledStringer example" %}}	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "CompiledError example" %}}	t := &testing.T{}

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := CmpRe(t, got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "CompiledCapture example" %}}	t := &testing.T{}

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := CmpRe(t, got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = CmpRe(t, got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
## T.Re shortcut

```go
func (t *T) Re(got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool
```

[`Re`]({{< ref "Re" >}}) is a shortcut for:

```go
t.Cmp(got, Re(reg, capture), args...)
```

See above for details.

[`Re()`]({{< ref "Re" >}}) optional parameter *capture* is here mandatory.
`nil` value should be passed to mimic its absence in
original [`Re()`]({{< ref "Re" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Base example" %}}	t := NewT(&testing.T{})

	got := "foo bar"
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	got = "bar foo"
	ok = t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "Stringer example" %}}	t := NewT(&testing.T{})

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Error example" %}}	t := NewT(&testing.T{})

	got := errors.New("foo bar")
	ok := t.Re(got, "(zip|bar)$", nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "Capture example" %}}	t := NewT(&testing.T{})

	got := "foo bar biz"
	ok := t.Re(got, `^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = t.Re(got, `^(\w+) (\w+) (\w+)$`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "Compiled example" %}}	t := NewT(&testing.T{})

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
{{% /expand%}}
{{%expand "CompiledStringer example" %}}	t := NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	// bytes.Buffer implements fmt.Stringer
	got := bytes.NewBufferString("foo bar")
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "CompiledError example" %}}	t := NewT(&testing.T{})

	expected := regexp.MustCompile("(zip|bar)$")

	got := errors.New("foo bar")
	ok := t.Re(got, expected, nil, "checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
{{% /expand%}}
{{%expand "CompiledCapture example" %}}	t := NewT(&testing.T{})

	expected := regexp.MustCompile(`^(\w+) (\w+) (\w+)$`)

	got := "foo bar biz"
	ok := t.Re(got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	got = "foo bar! biz"
	ok = t.Re(got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
