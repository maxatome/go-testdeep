---
title: "NotNaN"
weight: 10
---

```go
func NotNaN() TestDeep
```

[`NotNaN`]({{< ref "NotNaN" >}}) operator checks that data is a float and is not not-a-number.

```go
got := math.NaN()
td.Cmp(t, got, td.NotNaN()) // fails
td.Cmp(t, 4.2, td.NotNaN()) // succeeds
td.Cmp(t, 4, td.NotNaN())   // fails, as 4 is not a float
```


> See also [<i class='fas fa-book'></i> NotNaN godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotNaN).

### Examples

{{%expand "Float32 example" %}}```go
	t := &testing.T{}

	got := float32(math.NaN())

	ok := td.Cmp(t, got, td.NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is NOT float32 not-a-number:", ok)

	got = 12

	ok = td.Cmp(t, got, td.NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is NOT float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is NOT float32 not-a-number: false
	// float32(12) is NOT float32 not-a-number: true

```{{% /expand%}}
{{%expand "Float64 example" %}}```go
	t := &testing.T{}

	got := math.NaN()

	ok := td.Cmp(t, got, td.NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = td.Cmp(t, got, td.NotNaN(),
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is NOT not-a-number: false
	// float64(12) is NOT not-a-number: true

```{{% /expand%}}
## CmpNotNaN shortcut

```go
func CmpNotNaN(t TestingT, got interface{}, args ...interface{}) bool
```

CmpNotNaN is a shortcut for:

```go
td.Cmp(t, got, td.NotNaN(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNotNaN godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNotNaN).

### Examples

{{%expand "Float32 example" %}}```go
	t := &testing.T{}

	got := float32(math.NaN())

	ok := td.CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(math.NaN()) is NOT float32 not-a-number:", ok)

	got = 12

	ok = td.CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float32(12) is NOT float32 not-a-number:", ok)

	// Output:
	// float32(math.NaN()) is NOT float32 not-a-number: false
	// float32(12) is NOT float32 not-a-number: true

```{{% /expand%}}
{{%expand "Float64 example" %}}```go
	t := &testing.T{}

	got := math.NaN()

	ok := td.CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("math.NaN() is not-a-number:", ok)

	got = 12

	ok = td.CmpNotNaN(t, got,
		"checks %v is not-a-number", got)

	fmt.Println("float64(12) is not-a-number:", ok)

	// math.NaN() is NOT not-a-number: false
	// float64(12) is NOT not-a-number: true

```{{% /expand%}}
## T.NotNaN shortcut

```go
func (t *T) NotNaN(got interface{}, args ...interface{}) bool
```

[`NotNaN`]({{< ref "NotNaN" >}}) is a shortcut for:

```go
t.Cmp(got, td.NotNaN(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.NotNaN godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.NotNaN).

### Examples

{{%expand "Float32 example" %}}```go
	t := td.NewT(&testing.T{})

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

```{{% /expand%}}
{{%expand "Float64 example" %}}```go
	t := td.NewT(&testing.T{})

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

```{{% /expand%}}
