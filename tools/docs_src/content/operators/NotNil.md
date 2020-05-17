---
title: "NotNil"
weight: 10
---

```go
func NotNil() TestDeep
```

[`NotNil`]({{< ref "NotNil" >}}) operator checks that data is not `nil` (or is a non-`nil`
interface, containing a non-`nil` pointer.)

```go
got := &Person{}
td.Cmp(t, got, td.NotNil()) // succeeds
td.Cmp(t, got, td.Not(nil)) // succeeds too, but be careful it is first
// because of got type *Person ≠ untyped nil so prefer NotNil()
```

but:

```go
var got fmt.Stringer = (*bytes.Buffer)(nil)
td.Cmp(t, got, td.NotNil()) // fails
td.Cmp(t, got, td.Not(nil)) // succeeds, as the interface is not nil
```


> See also [<i class='fas fa-book'></i> NotNil godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotNil).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := td.Cmp(t, got, td.Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = td.Cmp(t, got, td.NotNil())
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = td.Cmp(t, got, td.Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = td.Cmp(t, got, td.NotNil())
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
td.Cmp(t, got, td.NotNil(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNotNil godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNotNil).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := td.Cmp(t, got, td.Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = td.CmpNotNil(t, got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = td.Cmp(t, got, td.Not(nil))
	fmt.Println(ok)

	// In this case NotNil() fails
	ok = td.CmpNotNil(t, got)
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
t.Cmp(got, td.NotNil(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.NotNil godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.NotNil).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	var got fmt.Stringer = &bytes.Buffer{}

	// nil value can be compared directly with Not(nil), no need of NotNil() here
	ok := t.Cmp(got, td.Not(nil))
	fmt.Println(ok)

	// But it works with NotNil() anyway
	ok = t.NotNil(got)
	fmt.Println(ok)

	got = (*bytes.Buffer)(nil)

	// In the case of an interface containing a nil pointer, comparing
	// with Not(nil) succeeds, as the interface is not nil
	ok = t.Cmp(got, td.Not(nil))
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
