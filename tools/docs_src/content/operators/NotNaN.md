---
title: "NotNaN"
weight: 10
---

```go
func NotNaN() TestDeep
```

[`NotNaN`]({{< ref "NotNaN" >}}) operator checks that data is a float and is not not-a-number.


### Examples

{{%expand "Float32 example" %}}	t := &testing.T{}

	got := float32(math.NaN())

	ok := Cmp(t, got, NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is NOT float32 not-a-number:", ok)

	got = 12

	ok = Cmp(t, got, NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is NOT float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is NOT float32 not-a-number: false
	// float32(12) is NOT float32 not-a-number: true
{{% /expand%}}
{{%expand "Float64 example" %}}	t := &testing.T{}

	got := math.NaN()

	ok := Cmp(t, got, NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = Cmp(t, got, NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is NOT not-a-number: false
	// float64(12) is NOT not-a-number: true
{{% /expand%}}
## CmpNotNaN shortcut

```go
func CmpNotNaN(t TestingT, got interface{}, args ...interface{}) bool
```

CmpNotNaN is a shortcut for:

```go
Cmp(t, got, NotNaN(), args...)
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

{{%expand "Float32 example" %}}	t := &testing.T{}

	got := float32(math.NaN())

	ok := CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is NOT float32 not-a-number:", ok)

	got = 12

	ok = CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is NOT float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is NOT float32 not-a-number: false
	// float32(12) is NOT float32 not-a-number: true
{{% /expand%}}
{{%expand "Float64 example" %}}	t := &testing.T{}

	got := math.NaN()

	ok := CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is NOT not-a-number: false
	// float64(12) is NOT not-a-number: true
{{% /expand%}}
## T.NotNaN shortcut

```go
func (t *T) NotNaN(got interface{}, args ...interface{}) bool
```

[`NotNaN`]({{< ref "NotNaN" >}}) is a shortcut for:

```go
t.Cmp(got, NotNaN(), args...)
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

{{%expand "Float32 example" %}}	t := NewT(&testing.T{})

	got := float32(math.NaN())

	ok := t.NotNaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is NOT float32 not-a-number:", ok)

	got = 12

	ok = t.NotNaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is NOT float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is NOT float32 not-a-number: false
	// float32(12) is NOT float32 not-a-number: true
{{% /expand%}}
{{%expand "Float64 example" %}}	t := NewT(&testing.T{})

	got := math.NaN()

	ok := t.NotNaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = t.NotNaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is NOT not-a-number: false
	// float64(12) is NOT not-a-number: true
{{% /expand%}}
