---
title: "PPtr"
weight: 10
---

```go
func PPtr(val interface{}) TestDeep
```

[`PPtr`]({{< ref "PPtr" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}). It takes the address of the address of
data and compares it to *val*.

*val* depends on data type. For example, if the compared data is an
`**int`, one can have:

```go
num := 12
pnum = &num
td.Cmp(t, &pnum, td.PPtr(12)) // succeeds
```

as well as an other operator:

```go
num := 3
pnum = &num
td.Cmp(t, &pnum, td.PPtr(td.Between(3, 4))) // succeeds
```

It is more efficient and shorter to write than:

```go
td.Cmp(t, &pnum, td.Ptr(td.Ptr(val))) // succeeds too
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of a pointer on a
pointer on *val*, except if *val* is a [TestDeep operator]({{< ref "operators" >}}). In this
case, it delegates [`TypeBehind()`]({{< ref "operators#typebehind-method" >}}) to the operator and returns the
[`reflect.Type`](https://pkg.go.dev/reflect/#Type) of a pointer on a pointer on the returned value (if
non-`nil` of course).


> See also [<i class='fas fa-book'></i> PPtr godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#PPtr).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	num := 12
	got := &num

	ok := td.Cmp(t, &got, td.PPtr(12))
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.PPtr(td.Between(4, 15)))
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpPPtr shortcut

```go
func CmpPPtr(t TestingT, got interface{}, val interface{}, args ...interface{}) bool
```

CmpPPtr is a shortcut for:

```go
td.Cmp(t, got, td.PPtr(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpPPtr godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpPPtr).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	num := 12
	got := &num

	ok := td.CmpPPtr(t, &got, 12)
	fmt.Println(ok)

	ok = td.CmpPPtr(t, &got, td.Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.PPtr shortcut

```go
func (t *T) PPtr(got interface{}, val interface{}, args ...interface{}) bool
```

[`PPtr`]({{< ref "PPtr" >}}) is a shortcut for:

```go
t.Cmp(got, td.PPtr(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.PPtr godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.PPtr).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	num := 12
	got := &num

	ok := t.PPtr(&got, 12)
	fmt.Println(ok)

	ok = t.PPtr(&got, td.Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
