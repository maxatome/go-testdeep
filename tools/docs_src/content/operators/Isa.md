---
title: "Isa"
weight: 10
---

```go
func Isa(model interface{}) TestDeep
```

[`Isa`]({{< ref "Isa" >}}) operator checks the data type or whether data implements an
interface or not.

Typically type checks:
```go
Isa(time.Time{})
Isa(&time.Time{})
Isa(map[string]time.Time{})
```

For interfaces it is a bit more complicated, as:
```go
fmt.Stringer(nil)
```
is not an interface, but just `nil`... To bypass this golang
limitation, [`Isa`]({{< ref "Isa" >}}) accepts pointers on interfaces. So checking that
data implements [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interface should be written as:
```go
Isa((*fmt.Stringer)(nil))
```

Of course, in the latter case, if data type is [`*fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer), [`Isa`]({{< ref "Isa" >}})
will match too (in fact before checking whether it implements
[`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) or not.)

[TypeBehind]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> Isa godoc](https://godoc.org/github.com/maxatome/go-testdeep#Isa).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type TstStruct struct {
		Field int
	}

	got := TstStruct{Field: 1}

	ok := Cmp(t, got, Isa(TstStruct{}), "checks got is a TstStruct")
	fmt.Println(ok)

	ok = Cmp(t, got, Isa(&TstStruct{}),
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = Cmp(t, &got, Isa(&TstStruct{}),
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

	ok := Cmp(t, got, Isa((*fmt.Stringer)(nil)),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = Cmp(t, errGot, Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil... So it
	// does not match
	errGot = nil

	ok = Cmp(t, errGot, Isa((*error)(nil)),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = Cmp(t, &errGot, Isa((*error)(nil)),
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
Cmp(t, got, Isa(model), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpIsa godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpIsa).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type TstStruct struct {
		Field int
	}

	got := TstStruct{Field: 1}

	ok := CmpIsa(t, got, TstStruct{}, "checks got is a TstStruct")
	fmt.Println(ok)

	ok = CmpIsa(t, got, &TstStruct{},
		"checks got is a pointer on a TstStruct")
	fmt.Println(ok)

	ok = CmpIsa(t, &got, &TstStruct{},
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

	ok := CmpIsa(t, got, (*fmt.Stringer)(nil),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = CmpIsa(t, errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil... So it
	// does not match
	errGot = nil

	ok = CmpIsa(t, errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// BUT if its address is passed, now it is OK as the types match
	ok = CmpIsa(t, &errGot, (*error)(nil),
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
t.Cmp(got, Isa(model), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Isa godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Isa).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

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
	t := NewT(&testing.T{})

	got := bytes.NewBufferString("foobar")

	ok := t.Isa(got, (*fmt.Stringer)(nil),
		"checks got implements fmt.Stringer interface")
	fmt.Println(ok)

	errGot := fmt.Errorf("An error #%d occurred", 123)

	ok = t.Isa(errGot, (*error)(nil),
		"checks errGot is a *error or implements error interface")
	fmt.Println(ok)

	// As nil, is passed below, it is not an interface but nil... So it
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
