---
title: "SStruct"
weight: 10
---

```go
func SStruct(model interface{}, expectedFields StructFields) TestDeep
```

[`SStruct`]({{< ref "SStruct" >}}) operator (aka strict-[`Struct`]({{< ref "Struct" >}})) compares the contents of a
struct or a pointer on a struct against values of *model* (if any)
and the values of *expectedFields*. The zero values are compared
too even if they are omitted from *expectedFields*: that is the
difference with [`Struct`]({{< ref "Struct" >}}) operator.

*model* must be the same type as compared data.

*expectedFields* can be `nil`, if no [TestDeep operators]({{< ref "operators" >}}) are involved.

To ignore a field, one has to specify it in *expectedFields* and
use the [`Ignore`]({{< ref "Ignore" >}}) operator.

```go
td.Cmp(t, td.SStruct(
  Person{
    Name: "John Doe",
  },
  td.StructFields{
    Age:      td.Between(40, 45),
    Children: td.Ignore(),
  }),
)
```

During a match, all expected and zero fields must be found to
succeed.

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> SStruct godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SStruct).

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
		NumChildren: 0,
	}

	// NumChildren is not listed in expected fields so it must be zero
	ok := td.Cmp(t, got,
		td.SStruct(Person{Name: "Foobar"}, td.StructFields{
			"Age": td.Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	got.NumChildren = 3
	ok = td.Cmp(t, got,
		td.SStruct(Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = td.Cmp(t, &got,
		td.SStruct(&Person{}, td.StructFields{
			"Name":        "Foobar",
			"Age":         td.Between(40, 50),
			"NumChildren": td.Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = td.Cmp(t, &got,
		td.SStruct((*Person)(nil), td.StructFields{
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
## CmpSStruct shortcut

```go
func CmpSStruct(t TestingT, got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool
```

CmpSStruct is a shortcut for:

```go
td.Cmp(t, got, td.SStruct(model, expectedFields), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSStruct godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSStruct).

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
		NumChildren: 0,
	}

	// NumChildren is not listed in expected fields so it must be zero
	ok := td.CmpSStruct(t, got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	got.NumChildren = 3
	ok = td.CmpSStruct(t, got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = td.CmpSStruct(t, &got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = td.CmpSStruct(t, &got, (*Person)(nil), td.StructFields{
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
## T.SStruct shortcut

```go
func (t *T) SStruct(got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool
```

[`SStruct`]({{< ref "SStruct" >}}) is a shortcut for:

```go
t.Cmp(got, td.SStruct(model, expectedFields), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SStruct godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SStruct).

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
		NumChildren: 0,
	}

	// NumChildren is not listed in expected fields so it must be zero
	ok := t.SStruct(got, Person{Name: "Foobar"}, td.StructFields{
		"Age": td.Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	got.NumChildren = 3
	ok = t.SStruct(got, Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = t.SStruct(&got, &Person{}, td.StructFields{
		"Name":        "Foobar",
		"Age":         td.Between(40, 50),
		"NumChildren": td.Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = t.SStruct(&got, (*Person)(nil), td.StructFields{
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
