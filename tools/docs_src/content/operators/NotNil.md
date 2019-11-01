---
title: "NotNil"
weight: 10
---

```go
func NotNil() TestDeep
```

[`NotNil`]({{< ref "NotNil" >}}) operator checks that data is not `nil` (or is a non-`nil`
interface, containing a non-`nil` pointer.)


### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := Cmp(t, got, Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = Cmp(t, got, NotNil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = Cmp(t, got, Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = Cmp(t, got, NotNil())
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false

```{{% /expand%}}
## CmpNotNil shortcut

```go
func CmpNotNil(t TestingT, got interface{}, args ...interface{}) bool
```

CmpNotNil is a shortcut for:

```go
Cmp(t, got, NotNil(), args...)
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

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := Cmp(t, got, Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = CmpNotNil(t, got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = Cmp(t, got, Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = CmpNotNil(t, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false

```{{% /expand%}}
## T.NotNil shortcut

```go
func (t *T) NotNil(got interface{}, args ...interface{}) bool
```

[`NotNil`]({{< ref "NotNil" >}}) is a shortcut for:

```go
t.Cmp(got, NotNil(), args...)
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

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := t.Cmp(got, Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = t.NotNil(got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = t.Cmp(got, Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = t.NotNil(got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false

```{{% /expand%}}
