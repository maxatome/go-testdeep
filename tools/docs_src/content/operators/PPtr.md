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
PPtr(12)
```
as well as an other operator:
```go
PPtr(Between(3, 4))
```

It is more efficient and shorter to write than:
```go
Ptr(Ptr(val))
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of a pointer on a
pointer on *val*, except if *val* is a [TestDeep operator]({{< ref "operators" >}}). In this
case, it delegates [`TypeBehind()`]({{< ref "operators#typebehind-method" >}}) to the operator and returns the
[`reflect.Type`](https://golang.org/pkg/reflect/#Type) of a pointer on a pointer on the returned value (if
non-`nil` of course).


> See also [<i class='fas fa-book'></i> PPtr godoc](https://godoc.org/github.com/maxatome/go-testdeep#PPtr).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	num := 12
	got := &num

	ok := Cmp(t, &got, PPtr(12))
	fmt.Println(ok)

	ok = Cmp(t, &got, PPtr(Between(4, 15)))
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
Cmp(t, got, PPtr(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpPPtr godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpPPtr).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	num := 12
	got := &num

	ok := CmpPPtr(t, &got, 12)
	fmt.Println(ok)

	ok = CmpPPtr(t, &got, Between(4, 15))
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
t.Cmp(got, PPtr(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.PPtr godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.PPtr).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	num := 12
	got := &num

	ok := t.PPtr(&got, 12)
	fmt.Println(ok)

	ok = t.PPtr(&got, Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
