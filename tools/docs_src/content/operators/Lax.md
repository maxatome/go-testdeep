---
title: "Lax"
weight: 10
---

```go
func Lax(expectedValue interface{}) TestDeep
```

[`Lax`]({{< ref "Lax" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}), it temporarily enables the BeLax config
flag before letting the comparison process continue its course.

It is more commonly used as CmpLax function than as an operator. It
could be used when, for example, an operator is constructed once
but applied to different, but compatible types as in:

```go
bw := Between(20, 30)
intValue := 21
floatValue := 21.89
Cmp(t, intValue, bw)        // no need to be lax here: same int types
Cmp(t, floatValue, Lax(bw)) // be lax please, as float64 ≠ int
```

Note that in the latter case, CmpLax() could be used as well:
```go
CmpLax(t, floatValue, bw)
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the greatest convertible or more common
[`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *expectedValue* if it is a base type (`bool`, `int*`,
`uint*`, `float*`, `complex*`, `string`), the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of
*expectedValue* otherwise, except if *expectedValue* is a [TestDeep
operator]({{< ref "operators" >}}). In this case, it delegates [`TypeBehind()`]({{< ref "operators#typebehind-method" >}}) to the operator.


> See also [<i class='fas fa-book'></i> Lax godoc](https://godoc.org/github.com/maxatome/go-testdeep#Lax).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := Between(1230, 1240) // int type here

	ok := Cmp(t, gotInt64, Lax(expected))
	fmt.Println("int64 got between ints [1230 .. 1240]:", ok)

	ok = Cmp(t, gotInt32, Lax(expected))
	fmt.Println("int32 got between ints [1230 .. 1240]:", ok)

	ok = Cmp(t, gotMyInt, Lax(expected))
	fmt.Println("myInt got between ints [1230 .. 1240]:", ok)

	// Output:
	// int64 got between ints [1230 .. 1240]: true
	// int32 got between ints [1230 .. 1240]: true
	// myInt got between ints [1230 .. 1240]: true

```{{% /expand%}}
## CmpLax shortcut

```go
func CmpLax(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool
```

CmpLax is a shortcut for:

```go
Cmp(t, got, Lax(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpLax godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpLax).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := Between(1230, 1240) // int type here

	ok := CmpLax(t, gotInt64, expected)
	fmt.Println("int64 got between ints [1230 .. 1240]:", ok)

	ok = CmpLax(t, gotInt32, expected)
	fmt.Println("int32 got between ints [1230 .. 1240]:", ok)

	ok = CmpLax(t, gotMyInt, expected)
	fmt.Println("myInt got between ints [1230 .. 1240]:", ok)

	// Output:
	// int64 got between ints [1230 .. 1240]: true
	// int32 got between ints [1230 .. 1240]: true
	// myInt got between ints [1230 .. 1240]: true

```{{% /expand%}}
## T.CmpLax shortcut

```go
func (t *T) CmpLax(got interface{}, expectedValue interface{}, args ...interface{}) bool
```

CmpLax is a shortcut for:

```go
t.Cmp(got, Lax(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.CmpLax godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpLax).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := Between(1230, 1240) // int type here

	ok := t.CmpLax(gotInt64, expected)
	fmt.Println("int64 got between ints [1230 .. 1240]:", ok)

	ok = t.CmpLax(gotInt32, expected)
	fmt.Println("int32 got between ints [1230 .. 1240]:", ok)

	ok = t.CmpLax(gotMyInt, expected)
	fmt.Println("myInt got between ints [1230 .. 1240]:", ok)

	// Output:
	// int64 got between ints [1230 .. 1240]: true
	// int32 got between ints [1230 .. 1240]: true
	// myInt got between ints [1230 .. 1240]: true

```{{% /expand%}}
