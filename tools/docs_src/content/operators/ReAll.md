---
title: "ReAll"
weight: 10
---

```go
func ReAll(reg interface{}, capture interface{}) TestDeep
```

[`ReAll`]({{< ref "ReAll" >}}) operator allows to successively apply a regexp on a `string`
(or convertible), `[]byte`, [`error`](https://pkg.go.dev/builtin/#error) or [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer) interface ([`error`](https://pkg.go.dev/builtin/#error)
interface is tested before [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer)) and to match its groups
contents.

*reg* is the regexp. It can be a `string` that is automatically
compiled using [`regexp.MustCompile`](https://pkg.go.dev/regexp/#MustCompile), or a [`*regexp.Regexp`](https://pkg.go.dev/regexp/#Regexp).

*capture* is used to match the contents of regexp groups. Groups
are presented as a `[]string` or `[][]byte` depending the original
matched data. Note that an other operator can be used here.

```go
td.Cmp(t, "John Doe",
  td.ReAll(`(\w+)(?: |\z)`, []string{"John", "Doe"})) // succeeds
td.Cmp(t, "John Doe",
  td.ReAll(`(\w+)(?: |\z)`, td.Bag("Doe", "John"))) // succeeds
```


> See also [<i class='fas fa-book'></i> ReAll godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ReAll).

### Examples

{{%expand "Capture example" %}}```go
	t := &testing.T{}

	got := "foo bar biz"
	ok := td.Cmp(t, got, td.ReAll(`(\w+)`, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = td.Cmp(t, got, td.ReAll(`(\w+)`, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CaptureComplex example" %}}```go
	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := td.Cmp(t, got,
		td.ReAll(`(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = td.Cmp(t, got,
		td.ReAll(`(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 20 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CompiledCapture example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := td.Cmp(t, got, td.ReAll(expected, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = td.Cmp(t, got, td.ReAll(expected, td.Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CompiledCaptureComplex example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := td.Cmp(t, got,
		td.ReAll(expected, td.ArrayEach(td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = td.Cmp(t, got,
		td.ReAll(expected, td.ArrayEach(td.Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 20 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## CmpReAll shortcut

```go
func CmpReAll(t TestingT, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool
```

CmpReAll is a shortcut for:

```go
td.Cmp(t, got, td.ReAll(reg, capture), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpReAll godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpReAll).

### Examples

{{%expand "Capture example" %}}```go
	t := &testing.T{}

	got := "foo bar biz"
	ok := td.CmpReAll(t, got, `(\w+)`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = td.CmpReAll(t, got, `(\w+)`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CaptureComplex example" %}}```go
	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := td.CmpReAll(t, got, `(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = td.CmpReAll(t, got, `(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CompiledCapture example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := td.CmpReAll(t, got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = td.CmpReAll(t, got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CompiledCaptureComplex example" %}}```go
	t := &testing.T{}

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := td.CmpReAll(t, got, expected, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = td.CmpReAll(t, got, expected, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
## T.ReAll shortcut

```go
func (t *T) ReAll(got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool
```

[`ReAll`]({{< ref "ReAll" >}}) is a shortcut for:

```go
t.Cmp(got, td.ReAll(reg, capture), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.ReAll godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.ReAll).

### Examples

{{%expand "Capture example" %}}```go
	t := td.NewT(&testing.T{})

	got := "foo bar biz"
	ok := t.ReAll(got, `(\w+)`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = t.ReAll(got, `(\w+)`, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CaptureComplex example" %}}```go
	t := td.NewT(&testing.T{})

	got := "11 45 23 56 85 96"
	ok := t.ReAll(got, `(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = t.ReAll(got, `(\d+)`, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CompiledCapture example" %}}```go
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := t.ReAll(got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = t.ReAll(got, expected, td.Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "CompiledCaptureComplex example" %}}```go
	t := td.NewT(&testing.T{})

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := t.ReAll(got, expected, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = t.ReAll(got, expected, td.ArrayEach(td.Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
