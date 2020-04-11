---
title: "Isa"
weight: 10
---

```go
func Isa(model interface{}) TestDeep
```

[`Isa`]({{< ref "Isa" >}}) operator checks the data type or whether data implements an
interface or not.

Typical type checks:

```go
td.Cmp(t, time.Now(), td.Isa(time.Time{}))  // succeeds
td.Cmp(t, time.Now(), td.Isa(&time.Time{})) // fails, as not a *time.Time
td.Cmp(t, got, td.Isa(map[string]time.Time{}))
```

For interfaces, it is a bit more complicated, as:

```go
fmt.Stringer(nil)
```

is not an interface, but just `nil`… To bypass this golang
limitation, [`Isa`]({{< ref "Isa" >}}) accepts pointers on interfaces. So checking that
data implements [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer) interface should be written as:

```go
td.Cmp(t, bytes.Buffer{}, td.Isa((*fmt.Stringer)(nil))) // succeeds
```

Of course, in the latter case, if checked data type is
[`*fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer), [`Isa`]({{< ref "Isa" >}}) will match too (in fact before checking whether
it implements [`fmt.Stringer`](https://pkg.go.dev/fmt/#Stringer) or not).

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> Isa godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Isa).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type TstStruct struct {
		Field int
	}

	got := TstStruct{Field: 1}

	ok := td.Cmp(t, got, td.Isa(TstStruct{}), "checks got is a TstStruct")
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Isa(&TstStruct{}),
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.Isa(&TstStruct{}),
		"checks &got is a pointer on a TstStruct")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
{{%expand "Interface example" %}}```go
	t := &testing.T{}

	got := bytes.NewBufferString("foobar")

	ok := td.Cmp(t, got, td.Isa((*fmt.Stringer)(nil)),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = td.Cmp(t, errGot, td.Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil… So it
	// does not match
	errGot = nil

	ok = td.Cmp(t, errGot, td.Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = td.Cmp(t, &errGot, td.Isa((*error)(nil)),
		"checks &errGot is a *error or implements error interface")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true

```{{% /expand%}}
## CmpIsa shortcut

```go
func CmpIsa(t TestingT, got interface{}, model interface{}, args ...interface{}) bool
```

CmpIsa is a shortcut for:

```go
td.Cmp(t, got, td.Isa(model), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpIsa godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpIsa).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type TstStruct struct {
		Field int
	}

	got := TstStruct{Field: 1}

	ok := td.CmpIsa(t, got, TstStruct{}, "checks got is a TstStruct")
	fmt.Println(ok)

	ok = td.CmpIsa(t, got, &TstStruct{},
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = td.CmpIsa(t, &got, &TstStruct{},
		"checks &got is a pointer on a TstStruct")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
{{%expand "Interface example" %}}```go
	t := &testing.T{}

	got := bytes.NewBufferString("foobar")

	ok := td.CmpIsa(t, got, (*fmt.Stringer)(nil),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = td.CmpIsa(t, errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil… So it
	// does not match
	errGot = nil

	ok = td.CmpIsa(t, errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = td.CmpIsa(t, &errGot, (*error)(nil),
		"checks &errGot is a *error or implements error interface")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true

```{{% /expand%}}
## T.Isa shortcut

```go
func (t *T) Isa(got interface{}, model interface{}, args ...interface{}) bool
```

[`Isa`]({{< ref "Isa" >}}) is a shortcut for:

```go
t.Cmp(got, td.Isa(model), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Isa godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Isa).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	type TstStruct struct {
		Field int
	}

	got := TstStruct{Field: 1}

	ok := t.Isa(got, TstStruct{}, "checks got is a TstStruct")
	fmt.Println(ok)

	ok = t.Isa(got, &TstStruct{},
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = t.Isa(&got, &TstStruct{},
		"checks &got is a pointer on a TstStruct")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true

```{{% /expand%}}
{{%expand "Interface example" %}}```go
	t := td.NewT(&testing.T{})

	got := bytes.NewBufferString("foobar")

	ok := t.Isa(got, (*fmt.Stringer)(nil),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = t.Isa(errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil… So it
	// does not match
	errGot = nil

	ok = t.Isa(errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = t.Isa(&errGot, (*error)(nil),
		"checks &errGot is a *error or implements error interface")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
	// true

```{{% /expand%}}
