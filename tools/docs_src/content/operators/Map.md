---
title: "Map"
weight: 10
---

```go
func Map(model interface{}, expectedEntries MapEntries) TestDeep
```

[`Map`]({{< ref "Map" >}}) operator compares the contents of a map against the non-zero
values of *model* (if any) and the values of *expectedEntries*.

*model* must be the same type as compared data.

*expectedEntries* can be `nil`, if no zero entries are expected and
no [TestDeep operator]({{< ref "operators" >}}) are involved.

During a match, all expected entries must be found and all data
entries must be expected to succeed.

```go
got := map[string]string{
  "foo": "test",
  "bar": "wizz",
  "zip": "buzz",
}
td.Cmp(t, got, td.Map(
  map[string]string{
    "foo": "test",
    "bar": "wizz",
  },
  td.MapEntries{
    "zip": td.HasSuffix("zz"),
  }),
) // succeeds
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> Map godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Map).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got,
		td.Map(map[string]int{"bar": 42},
			td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got,
		td.Map(map[string]int{},
			td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got,
		td.Map((map[string]int)(nil),
			td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got,
		td.Map(MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()}),
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Map(&MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Map(&MyMap{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Map((*MyMap)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## CmpMap shortcut

```go
func CmpMap(t TestingT, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool
```

CmpMap is a shortcut for:

```go
td.Cmp(t, got, td.Map(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpMap godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpMap).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpMap(t, got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, got, map[string]int{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, got, (map[string]int)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpMap(t, got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, &got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, &got, &MyMap{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = td.CmpMap(t, &got, (*MyMap)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## T.Map shortcut

```go
func (t *T) Map(got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool
```

[`Map`]({{< ref "Map" >}}) is a shortcut for:

```go
t.Cmp(got, td.Map(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Map godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Map).

### Examples

{{%expand "Map example" %}}```go
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.Map(got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = t.Map(got, map[string]int{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = t.Map(got, (map[string]int)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := td.NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.Map(got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, &MyMap{}, td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, (*MyMap)(nil), td.MapEntries{"bar": 42, "foo": td.Lt(15), "zip": td.Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
