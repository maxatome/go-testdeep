---
title: "N"
weight: 10
---

```go
func N(num interface{}, tolerance ...interface{}) TestDeep
```

[`N`]({{< ref "N" >}}) operator compares a numeric data against *num* ± *tolerance*. If
*tolerance* is missing, it defaults to 0. *num* and *tolerance*
must be the same kind as the compared value.

```go
td.Cmp(t, 12.2, td.N(12., 0.3)) // succeeds
td.Cmp(t, 12.2, td.N(12., 0.1)) // fails
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *num*.


> See also [<i class='fas fa-book'></i> N godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#N).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 1.12345

	ok := td.Cmp(t, got, td.N(1.1234, 0.00006),
		"checks %v = 1.1234 ± 0.00006", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
## CmpN shortcut

```go
func CmpN(t TestingT, got interface{}, num interface{}, tolerance interface{}, args ...interface{}) bool
```

CmpN is a shortcut for:

```go
td.Cmp(t, got, td.N(num, tolerance), args...)
```

See above for details.

[`N()`]({{< ref "N" >}}) optional parameter *tolerance* is here mandatory.
0 value should be passed to mimic its absence in
original [`N()`]({{< ref "N" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpN godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpN).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 1.12345

	ok := td.CmpN(t, got, 1.1234, 0.00006,
		"checks %v = 1.1234 ± 0.00006", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
## T.N shortcut

```go
func (t *T) N(got interface{}, num interface{}, tolerance interface{}, args ...interface{}) bool
```

[`N`]({{< ref "N" >}}) is a shortcut for:

```go
t.Cmp(got, td.N(num, tolerance), args...)
```

See above for details.

[`N()`]({{< ref "N" >}}) optional parameter *tolerance* is here mandatory.
0 value should be passed to mimic its absence in
original [`N()`]({{< ref "N" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.N godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.N).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := 1.12345

	ok := t.N(got, 1.1234, 0.00006,
		"checks %v = 1.1234 ± 0.00006", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
