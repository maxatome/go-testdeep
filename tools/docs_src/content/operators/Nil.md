---
title: "Nil"
weight: 10
---

```go
func Nil() TestDeep
```

[`Nil`]({{< ref "Nil" >}}) operator checks that data is `nil` (or is a non-`nil` interface,
but containing a `nil` pointer.)


### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := Cmp(t, got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = Cmp(t, got, Nil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = Cmp(t, got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = Cmp(t, got, Nil())
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true

```{{% /expand%}}
## CmpNil shortcut

```go
func CmpNil(t TestingT, got interface{}, args ...interface{}) bool
```

CmpNil is a shortcut for:

```go
Cmp(t, got, Nil(), args...)
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

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := Cmp(t, got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = CmpNil(t, got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = Cmp(t, got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = CmpNil(t, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true

```{{% /expand%}}
## T.Nil shortcut

```go
func (t *T) Nil(got interface{}, args ...interface{}) bool
```

[`Nil`]({{< ref "Nil" >}}) is a shortcut for:

```go
t.Cmp(got, Nil(), args...)
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

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := t.Cmp(got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = t.Nil(got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = t.Cmp(got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = t.Nil(got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true

```{{% /expand%}}
