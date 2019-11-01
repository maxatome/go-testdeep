---
title: "Struct"
weight: 10
---

```go
func Struct(model interface{}, expectedFields StructFields) TestDeep
```

[`Struct`]({{< ref "Struct" >}}) operator compares the contents of a struct or a pointer on a
struct against the non-zero values of *model* (if any) and the
values of *expectedFields*.

*model* must be the same type as compared data.

*expectedFields* can be `nil`, if no zero entries are expected and
no [TestDeep operator]({{< ref "operators" >}}) are involved.

During a match, all expected fields must be found to
succeed. Non-expected fields are ignored.

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> Struct godoc](https://godoc.org/github.com/maxatome/go-testdeep#Struct).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	// As NumChildren is zero in Struct() call, it is not checked
	ok := Cmp(t, got,
		Struct(Person{Name: "Foobar"}, StructFields{
			"Age": Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = Cmp(t, got,
		Struct(Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = Cmp(t, &got,
		Struct(&Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = Cmp(t, &got,
		Struct((*Person)(nil), StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## CmpStruct shortcut

```go
func CmpStruct(t TestingT, got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool
```

CmpStruct is a shortcut for:

```go
Cmp(t, got, Struct(model, expectedFields), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpStruct godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpStruct).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	// As NumChildren is zero in Struct() call, it is not checked
	ok := CmpStruct(t, got, Person{Name: "Foobar"}, StructFields{
		"Age": Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = CmpStruct(t, got, Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = CmpStruct(t, &got, &Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = CmpStruct(t, &got, (*Person)(nil), StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## T.Struct shortcut

```go
func (t *T) Struct(got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool
```

[`Struct`]({{< ref "Struct" >}}) is a shortcut for:

```go
t.Cmp(got, Struct(model, expectedFields), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Struct godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Struct).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	type Person struct {
		Name        string
		Age         int
		NumChildren int
	}

	got := Person{
		Name:        "Foobar",
		Age:         42,
		NumChildren: 3,
	}

	// As NumChildren is zero in Struct() call, it is not checked
	ok := t.Struct(got, Person{Name: "Foobar"}, StructFields{
		"Age": Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = t.Struct(got, Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = t.Struct(&got, &Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = t.Struct(&got, (*Person)(nil), StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
