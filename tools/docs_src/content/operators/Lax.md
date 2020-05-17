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
bw := td.Between(20, 30)
intValue := 21
floatValue := 21.89
td.Cmp(t, intValue, bw)           // no need to be lax here: same int types
td.Cmp(t, floatValue, td.Lax(bw)) // be lax please, as float64 ≠ int
```

Note that in the latter case, CmpLax() could be used as well:
```go
td.CmpLax(t, floatValue, bw)
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the greatest convertible or more common
[`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *expectedValue* if it is a base type (`bool`, `int*`,
`uint*`, `float*`, `complex*`, `string`), the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of
*expectedValue* otherwise, except if *expectedValue* is a [TestDeep
operator]({{< ref "operators" >}}). In this case, it delegates [`TypeBehind()`]({{< ref "operators#typebehind-method" >}}) to the operator.


> See also [<i class='fas fa-book'></i> Lax godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Lax).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := td.Between(1230, 1240) // int type here

	ok := td.Cmp(t, gotInt64, td.Lax(expected))
	fmt.Println("int64 got between ints [1230 .. 1240]:", ok)

	ok = td.Cmp(t, gotInt32, td.Lax(expected))
	fmt.Println("int32 got between ints [1230 .. 1240]:", ok)

	ok = td.Cmp(t, gotMyInt, td.Lax(expected))
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
td.Cmp(t, got, td.Lax(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpLax godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpLax).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := td.Between(1230, 1240) // int type here

	ok := td.CmpLax(t, gotInt64, expected)
	fmt.Println("int64 got between ints [1230 .. 1240]:", ok)

	ok = td.CmpLax(t, gotInt32, expected)
	fmt.Println("int32 got between ints [1230 .. 1240]:", ok)

	ok = td.CmpLax(t, gotMyInt, expected)
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
t.Cmp(got, td.Lax(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.CmpLax godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.CmpLax).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	gotInt64 := int64(1234)
	gotInt32 := int32(1235)

	type myInt uint16
	gotMyInt := myInt(1236)

	expected := td.Between(1230, 1240) // int type here

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
