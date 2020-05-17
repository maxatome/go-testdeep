---
title: "Nil"
weight: 10
---

```go
func Nil() TestDeep
```

[`Nil`]({{< ref "Nil" >}}) operator checks that data is `nil` (or is a non-`nil` interface,
but containing a `nil` pointer.)

```go
var got *int
td.Cmp(t, got, td.Nil())    // succeeds
td.Cmp(t, got, nil)         // fails as (*int)(nil) ≠ untyped nil
td.Cmp(t, got, (*int)(nil)) // succeeds
```

but:

```go
var got fmt.Stringer = (*bytes.Buffer)(nil)
td.Cmp(t, got, td.Nil()) // succeeds
td.Cmp(t, got, nil)      // fails, as the interface is not nil
got = nil
td.Cmp(t, got, nil) // succeeds
```


> See also [<i class='fas fa-book'></i> Nil godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Nil).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := td.Cmp(t, got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = td.Cmp(t, got, td.Nil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = td.Cmp(t, got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = td.Cmp(t, got, td.Nil())
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
td.Cmp(t, got, td.Nil(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNil godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNil).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	var got fmt.Stringer // interface

	// nil value can be compared directly with nil, no need of Nil() here
	ok := td.Cmp(t, got, nil)
	fmt.Println(ok)

	// But it works with Nil() anyway
	ok = td.CmpNil(t, got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with nil fails, as the interface is not nil
	ok = td.Cmp(t, got, nil)
	fmt.Println(ok)

	// In this case Nil() succeed
	ok = td.CmpNil(t, got)
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
t.Cmp(got, td.Nil(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Nil godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Nil).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

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
