---
title: "NaN"
weight: 10
---

```go
func NaN() TestDeep
```

[`NaN`]({{< ref "NaN" >}}) operator checks that data is a float and is not-a-number.


### Examples

{{%expand "Float32 example" %}}```go
	t := &testing.T{}

	got := float32(math.NaN())

	ok := Cmp(t, got, NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is float32 not-a-number:", ok)

	got = 12

	ok = Cmp(t, got, NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is float32 not-a-number: true
	// float32(12) is float32 not-a-number: false

```{{% /expand%}}
{{%expand "Float64 example" %}}```go
	t := &testing.T{}

	got := math.NaN()

	ok := Cmp(t, got, NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = Cmp(t, got, NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is not-a-number: true
	// float64(12) is not-a-number: false

```{{% /expand%}}
## CmpNaN shortcut

```go
func CmpNaN(t TestingT, got interface{}, args ...interface{}) bool
```

CmpNaN is a shortcut for:

```go
Cmp(t, got, NaN(), args...)
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

{{%expand "Float32 example" %}}```go
	t := &testing.T{}

	got := float32(math.NaN())

	ok := CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is float32 not-a-number:", ok)

	got = 12

	ok = CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is float32 not-a-number: true
	// float32(12) is float32 not-a-number: false

```{{% /expand%}}
{{%expand "Float64 example" %}}```go
	t := &testing.T{}

	got := math.NaN()

	ok := CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is not-a-number: true
	// float64(12) is not-a-number: false

```{{% /expand%}}
## T.NaN shortcut

```go
func (t *T) NaN(got interface{}, args ...interface{}) bool
```

[`NaN`]({{< ref "NaN" >}}) is a shortcut for:

```go
t.Cmp(got, NaN(), args...)
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

{{%expand "Float32 example" %}}```go
	t := NewT(&testing.T{})

	got := float32(math.NaN())

	ok := t.NaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is float32 not-a-number:", ok)

	got = 12

	ok = t.NaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is float32 not-a-number: true
	// float32(12) is float32 not-a-number: false

```{{% /expand%}}
{{%expand "Float64 example" %}}```go
	t := NewT(&testing.T{})

	got := math.NaN()

	ok := t.NaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = t.NaN(got,
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is not-a-number: true
	// float64(12) is not-a-number: false

```{{% /expand%}}
