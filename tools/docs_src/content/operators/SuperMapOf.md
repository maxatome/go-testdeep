---
title: "SuperMapOf"
weight: 10
---

```go
func SuperMapOf(model interface{}, expectedEntries MapEntries) TestDeep
```

[`SuperMapOf`]({{< ref "SuperMapOf" >}}) operator compares the contents of a map against the non-zero
values of *model* (if any) and the values of *expectedEntries*.

*model* must be the same type as compared data.

*expectedEntries* can be `nil`, if no zero entries are expected and
no [TestDeep operator]({{< ref "operators" >}}) are involved.

During a match, each expected entry should match in the compared
map. But some entries in the compared map may not be expected.

```go
got := map[string]string{
  "foo": "test",
  "bar": "wizz",
  "zip": "buzz",
}
td.Cmp(t, got, td.SuperMapOf(
  map[string]string{
    "foo": "test",
  },
  td.MapEntries{
    "zip": td.HasSuffix("zz"),
  }),
) // succeeds

td.Cmp(t, got, td.SuperMapOf(
  map[string]string{
    "foo": "test",
  },
  td.MapEntries{
    "biz": td.HasSuffix("zz"),
  }),
) // fails, missing {"biz": …} in got
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> SuperMapOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SuperMapOf).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got,
		td.SuperMapOf(map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15)}),
		"checks map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got,
		td.SuperMapOf(MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)}),
		"checks typed map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.SuperMapOf(&MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)}),
		"checks pointed typed map %v contains at leat all expected keys/values",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpSuperMapOf shortcut

```go
func CmpSuperMapOf(t TestingT, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool
```

CmpSuperMapOf is a shortcut for:

```go
td.Cmp(t, got, td.SuperMapOf(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSuperMapOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSuperMapOf).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpSuperMapOf(t, got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpSuperMapOf(t, got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks typed map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	ok = td.CmpSuperMapOf(t, &got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks pointed typed map %v contains at leat all expected keys/values",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.SuperMapOf shortcut

```go
func (t *T) SuperMapOf(got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool
```

[`SuperMapOf`]({{< ref "SuperMapOf" >}}) is a shortcut for:

```go
t.Cmp(got, td.SuperMapOf(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SuperMapOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SuperMapOf).

### Examples

{{%expand "Map example" %}}```go
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.SuperMapOf(got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := td.NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.SuperMapOf(got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks typed map %v contains at leat all expected keys/values", got)
	fmt.Println(ok)

	ok = t.SuperMapOf(&got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15)},
		"checks pointed typed map %v contains at leat all expected keys/values",
		got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
