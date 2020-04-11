---
title: "Cap"
weight: 10
---

```go
func Cap(expectedCap interface{}) TestDeep
```

[`Cap`]({{< ref "Cap" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}). It takes data, applies `cap()` function
on it and compares its result to *expectedCap*. Of course, the
compared value must be an array, a channel or a slice.

*expectedCap* can be an `int` value:

```go
td.Cmp(t, gotSlice, td.Cap(12))
```

as well as an other operator:

```go
td.Cmp(t, gotSlice, td.Cap(td.Between(3, 4)))
```


> See also [<i class='fas fa-book'></i> Cap godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cap).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := td.Cmp(t, got, td.Cap(12), "checks %v capacity is 12", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Cap(0), "checks %v capacity is 0", got)
	fmt.Println(ok)

	got = nil

	ok = td.Cmp(t, got, td.Cap(0), "checks %v capacity is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
{{%expand "Operator example" %}}```go
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := td.Cmp(t, got, td.Cap(td.Between(10, 12)),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Cap(td.Gt(10)),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpCap shortcut

```go
func CmpCap(t TestingT, got interface{}, expectedCap interface{}, args ...interface{}) bool
```

CmpCap is a shortcut for:

```go
td.Cmp(t, got, td.Cap(expectedCap), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpCap godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpCap).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := td.CmpCap(t, got, 12, "checks %v capacity is 12", got)
	fmt.Println(ok)

	ok = td.CmpCap(t, got, 0, "checks %v capacity is 0", got)
	fmt.Println(ok)

	got = nil

	ok = td.CmpCap(t, got, 0, "checks %v capacity is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
{{%expand "Operator example" %}}```go
	t := &testing.T{}

	got := make([]int, 0, 12)

	ok := td.CmpCap(t, got, td.Between(10, 12),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	ok = td.CmpCap(t, got, td.Gt(10),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.Cap shortcut

```go
func (t *T) Cap(got interface{}, expectedCap interface{}, args ...interface{}) bool
```

[`Cap`]({{< ref "Cap" >}}) is a shortcut for:

```go
t.Cmp(got, td.Cap(expectedCap), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Cap godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Cap).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := make([]int, 0, 12)

	ok := t.Cap(got, 12, "checks %v capacity is 12", got)
	fmt.Println(ok)

	ok = t.Cap(got, 0, "checks %v capacity is 0", got)
	fmt.Println(ok)

	got = nil

	ok = t.Cap(got, 0, "checks %v capacity is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
{{%expand "Operator example" %}}```go
	t := td.NewT(&testing.T{})

	got := make([]int, 0, 12)

	ok := t.Cap(got, td.Between(10, 12),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	ok = t.Cap(got, td.Gt(10),
		"checks %v capacity is in [10 .. 12]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
