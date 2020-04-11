---
title: "NaN"
weight: 10
---

```go
func NaN() TestDeep
```

[`NaN`]({{< ref "NaN" >}}) operator checks that data is a float and is not-a-number.

```go
got := math.NaN()
td.Cmp(t, got, td.NaN()) // succeeds
td.Cmp(t, 4.2, td.NaN()) // fails
```


> See also [<i class='fas fa-book'></i> NaN godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NaN).

### Examples

{{%expand "Float32 example" %}}```go
	t := &testing.T{}

	got := float32(math.NaN())

	ok := td.Cmp(t, got, td.NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is float32 not-a-number:", ok)

	got = 12

	ok = td.Cmp(t, got, td.NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is float32 not-a-number: true
	// float32(12) is float32 not-a-number: false

```{{% /expand%}}
{{%expand "Float64 example" %}}```go
	t := &testing.T{}

	got := math.NaN()

	ok := td.Cmp(t, got, td.NaN(),
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = td.Cmp(t, got, td.NaN(),
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
td.Cmp(t, got, td.NaN(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNaN godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNaN).

### Examples

{{%expand "Float32 example" %}}```go
	t := &testing.T{}

	got := float32(math.NaN())

	ok := td.CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is float32 not-a-number:", ok)

	got = 12

	ok = td.CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is float32 not-a-number: true
	// float32(12) is float32 not-a-number: false

```{{% /expand%}}
{{%expand "Float64 example" %}}```go
	t := &testing.T{}

	got := math.NaN()

	ok := td.CmpNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = td.CmpNaN(t, got,
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
t.Cmp(got, td.NaN(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.NaN godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.NaN).

### Examples

{{%expand "Float32 example" %}}```go
	t := td.NewT(&testing.T{})

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
	t := td.NewT(&testing.T{})

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
