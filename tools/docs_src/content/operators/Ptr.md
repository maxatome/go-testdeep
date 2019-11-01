---
title: "Ptr"
weight: 10
---

```go
func Ptr(val interface{}) TestDeep
```

[`Ptr`]({{< ref "Ptr" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}). It takes the address of data and
compares it to *val*.

*val* depends on data type. For example, if the compared data is an
`*int`, one can have:
```go
Ptr(12)
```
as well as an other operator:
```go
Ptr(Between(3, 4))
```

[TypeBehind]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of a pointer on *val*,
except if *val* is a [TestDeep operator]({{< ref "operators" >}}). In this case, it delegates
[TypeBehind()]({{< ref "operators#typebehind-method" >}}) to the operator and returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of a
pointer on the returned value (if non-`nil` of course).


### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 12

	ok := Cmp(t, &got, Ptr(12))
	fmt.Println(ok)

	ok = Cmp(t, &got, Ptr(Between(4, 15)))
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpPtr shortcut

```go
func CmpPtr(t TestingT, got interface{}, val interface{}, args ...interface{}) bool
```

CmpPtr is a shortcut for:

```go
Cmp(t, got, Ptr(val), args...)
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

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 12

	ok := CmpPtr(t, &got, 12)
	fmt.Println(ok)

	ok = CmpPtr(t, &got, Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.Ptr shortcut

```go
func (t *T) Ptr(got interface{}, val interface{}, args ...interface{}) bool
```

[`Ptr`]({{< ref "Ptr" >}}) is a shortcut for:

```go
t.Cmp(got, Ptr(val), args...)
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

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	got := 12

	ok := t.Ptr(&got, 12)
	fmt.Println(ok)

	ok = t.Ptr(&got, Between(4, 15))
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
