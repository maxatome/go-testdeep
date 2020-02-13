---
title: "SStruct"
weight: 10
---

```go
func SStruct(model interface{}, expectedFields StructFields) TestDeep
```

[`SStruct`]({{< ref "SStruct" >}}) operator (a.k.a. strict-[`Struct`]({{< ref "Struct" >}})) compares the contents of a
struct or a pointer on a struct against values of *model* (if any)
and the values of *expectedFields*. The zero values are compared
too even if they are omitted from *expectedFields*: that is the
difference with [`Struct`]({{< ref "Struct" >}}) operator.

To ignore a field, one has to specify it in *expectedFields* and
use the [`Ignore`]({{< ref "Ignore" >}}) operator.

```go
td.SStruct(Person{Name: "Bob"},
  td.StructFields{
    "Age": td.Ignore(),
  })
```

*model* must be the same type as compared data.

*expectedFields* can be `nil`, if no [TestDeep operators]({{< ref "operators" >}}) are involved.

During a match, all expected and zero fields must be found to
succeed.

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> SStruct godoc](https://godoc.org/github.com/maxatome/go-testdeep#SStruct).

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
		NumChildren: 0,
	}

	// NumChildren is not listed in expected fields so it must be zero
	ok := Cmp(t, got,
		SStruct(Person{Name: "Foobar"}, StructFields{
			"Age": Between(40, 50),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	got.NumChildren = 3
	ok = Cmp(t, got,
		SStruct(Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = Cmp(t, &got,
		SStruct(&Person{}, StructFields{
			"Name":        "Foobar",
			"Age":         Between(40, 50),
			"NumChildren": Not(0),
		}),
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = Cmp(t, &got,
		SStruct((*Person)(nil), StructFields{
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
## CmpSStruct shortcut

```go
func CmpSStruct(t TestingT, got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool
```

CmpSStruct is a shortcut for:

```go
Cmp(t, got, SStruct(model, expectedFields), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSStruct godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpSStruct).

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
		NumChildren: 0,
	}

	// NumChildren is not listed in expected fields so it must be zero
	ok := CmpSStruct(t, got, Person{Name: "Foobar"}, StructFields{
		"Age": Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	got.NumChildren = 3
	ok = CmpSStruct(t, got, Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = CmpSStruct(t, &got, &Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = CmpSStruct(t, &got, (*Person)(nil), StructFields{
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
## T.SStruct shortcut

```go
func (t *T) SStruct(got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool
```

[`SStruct`]({{< ref "SStruct" >}}) is a shortcut for:

```go
t.Cmp(got, SStruct(model, expectedFields), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SStruct godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.SStruct).

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
		NumChildren: 0,
	}

	// NumChildren is not listed in expected fields so it must be zero
	ok := t.SStruct(got, Person{Name: "Foobar"}, StructFields{
		"Age": Between(40, 50),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model can be empty
	got.NumChildren = 3
	ok = t.SStruct(got, Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Works with pointers too
	ok = t.SStruct(&got, &Person{}, StructFields{
		"Name":        "Foobar",
		"Age":         Between(40, 50),
		"NumChildren": Not(0),
	},
		"checks %v is the right Person")
	fmt.Println(ok)

	// Model does not need to be instanciated
	ok = t.SStruct(&got, (*Person)(nil), StructFields{
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
