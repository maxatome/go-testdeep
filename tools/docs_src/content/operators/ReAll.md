---
title: "ReAll"
weight: 10
---

```go
func ReAll(reg interface{}, capture interface{}) TestDeep
```

[`ReAll`]({{< ref "ReAll" >}}) operator allows to successively apply a regexp on a `string`
(or convertible), `[]byte`, [`error`](https://golang.org/pkg/builtin/#error) or [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interface ([`error`](https://golang.org/pkg/builtin/#error)
interface is tested before [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer)) and to match its groups
contents.

*reg* is the regexp. It can be a `string` that is automatically
compiled using [`regexp.MustCompile`](https://golang.org/pkg/regexp/#MustCompile), or a [`*regexp.Regexp`](https://golang.org/pkg/regexp/#Regexp).

*capture* is used to match the contents of regexp groups. Groups
are presented as a `[]string` or `[][]byte` depending the original
matched data. Note that an other operator can be used here.

```go
Cmp(t, "John Doe",
  ReAll(`(\w+)(?: |\z)`, []string{"John", "Doe"})) // succeeds
Cmp(t, "John Doe",
  ReAll(`(\w+)(?: |\z)`, Bag("Doe", "John"))       // succeeds
```


### Examples

{{%expand "Capture example" %}}	t := &testing.T{}

	got := "foo bar biz"
	ok := Cmp(t, got, ReAll(`(\w+)`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = Cmp(t, got, ReAll(`(\w+)`, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CaptureComplex example" %}}	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := Cmp(t, got,
		ReAll(`(\d+)`, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = Cmp(t, got,
		ReAll(`(\d+)`, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 20 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CompiledCapture example" %}}	t := &testing.T{}

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := Cmp(t, got, ReAll(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = Cmp(t, got, ReAll(expected, Set("biz", "foo", "bar")),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CompiledCaptureComplex example" %}}	t := &testing.T{}

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := Cmp(t, got,
		ReAll(expected, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 10 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = Cmp(t, got,
		ReAll(expected, ArrayEach(Code(func(num string) bool {
			n, err := strconv.Atoi(num)
			return err == nil && n > 20 && n < 100
		}))),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
## CmpReAll shortcut

```go
func CmpReAll(t TestingT, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool
```

CmpReAll is a shortcut for:

```go
Cmp(t, got, ReAll(reg, capture), args...)
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

{{%expand "Capture example" %}}	t := &testing.T{}

	got := "foo bar biz"
	ok := CmpReAll(t, got, `(\w+)`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = CmpReAll(t, got, `(\w+)`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CaptureComplex example" %}}	t := &testing.T{}

	got := "11 45 23 56 85 96"
	ok := CmpReAll(t, got, `(\d+)`, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = CmpReAll(t, got, `(\d+)`, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CompiledCapture example" %}}	t := &testing.T{}

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := CmpReAll(t, got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = CmpReAll(t, got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CompiledCaptureComplex example" %}}	t := &testing.T{}

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := CmpReAll(t, got, expected, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = CmpReAll(t, got, expected, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
## T.ReAll shortcut

```go
func (t *T) ReAll(got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool
```

[`ReAll`]({{< ref "ReAll" >}}) is a shortcut for:

```go
t.Cmp(got, ReAll(reg, capture), args...)
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

{{%expand "Capture example" %}}	t := NewT(&testing.T{})

	got := "foo bar biz"
	ok := t.ReAll(got, `(\w+)`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = t.ReAll(got, `(\w+)`, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CaptureComplex example" %}}	t := NewT(&testing.T{})

	got := "11 45 23 56 85 96"
	ok := t.ReAll(got, `(\d+)`, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = t.ReAll(got, `(\d+)`, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CompiledCapture example" %}}	t := NewT(&testing.T{})

	expected := regexp.MustCompile(`(\w+)`)

	got := "foo bar biz"
	ok := t.ReAll(got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but all catured groups do not match Set
	got = "foo BAR biz"
	ok = t.ReAll(got, expected, Set("biz", "foo", "bar"),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
{{%expand "CompiledCaptureComplex example" %}}	t := NewT(&testing.T{})

	expected := regexp.MustCompile(`(\d+)`)

	got := "11 45 23 56 85 96"
	ok := t.ReAll(got, expected, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 10 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Matches, but 11 is not greater than 20
	ok = t.ReAll(got, expected, ArrayEach(Code(func(num string) bool {
		n, err := strconv.Atoi(num)
		return err == nil && n > 20 && n < 100
	})),
		"checks value %s", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
{{% /expand%}}
