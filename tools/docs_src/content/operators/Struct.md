---
title: "Struct"
weight: 10
---

```go
func Struct(model interface{}, expectedFields StructFields) TestDeep
```

[`Struct`]({{< ref "Struct" >}}) operator compares the contents of a struct or a pointer on a
struct against the non-zero values of *model* (if any) and the
values of *expectedFields*. See [`SStruct`]({{< ref "SStruct" >}}) to compares against zero
fields without specifying them in *expectedFields*.

*model* must be the same type as compared data.

*expectedFields* can be `nil`, if no zero entries are expected and
no [TestDeep operators]({{< ref "operators" >}}) are involved.

```go
td.Cmp(t, td.Struct(
  Person{
    Name: "John Doe",
  },
  td.StructFields{
    Age:      td.Between(40, 45),
    Children: 0,
  }),
)
```

During a match, all expected fields must be found to
succeed. Non-expected fields are ignored.

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> Struct godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Struct).

### Example

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
	ok := td.Cmp(t, got,
		td.Struct(Person{Name: "Foobar"}, td.StructFields{
			"Age": td.Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = td.Cmp(t, got,
		td.Struct(Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = td.Cmp(t, &got,
		td.Struct(&Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = td.Cmp(t, &got,
		td.Struct((*Person)(nil), td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
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
td.Cmp(t, got, td.Struct(model, expectedFields), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpStruct godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpStruct).

### Example

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
	ok := td.CmpStruct(t, got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = td.CmpStruct(t, got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = td.CmpStruct(t, &got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = td.CmpStruct(t, &got, (*Person)(nil), td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
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
t.Cmp(got, td.Struct(model, expectedFields), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Struct godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Struct).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

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
	ok := t.Struct(got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	ok = t.Struct(got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = t.Struct(&got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = t.Struct(&got, (*Person)(nil), td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
