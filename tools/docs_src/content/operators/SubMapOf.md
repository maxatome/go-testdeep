---
title: "SubMapOf"
weight: 10
---

```go
func SubMapOf(model interface{}, expectedEntries MapEntries) TestDeep
```

[`SubMapOf`]({{< ref "SubMapOf" >}}) operator compares the contents of a map against the non-zero
values of *model* (if any) and the values of *expectedEntries*.

*model* must be the same type as compared data.

*expectedEntries* can be `nil`, if no zero entries are expected and
no [TestDeep operator]({{< ref "operators" >}}) are involved.

During a match, each map entry should be matched by an expected
entry to succeed. But some expected entries can be missing from the
compared map.

```go
got := map[string]string{
  "foo": "test",
  "zip": "buzz",
}
td.Cmp(t, got, td.SubMapOf(
  map[string]string{
    "foo": "test",
    "bar": "wizz",
  },
  td.MapEntries{
    "zip": td.HasSuffix("zz"),
  }),
) // succeeds

td.Cmp(t, got, td.SubMapOf(
  map[string]string{
    "bar": "wizz",
  },
  td.MapEntries{
    "zip": td.HasSuffix("zz"),
  }),
) // fails, extra {"foo": "test"} in got
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> SubMapOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#SubMapOf).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42}

	ok := td.Cmp(t, got,
		td.SubMapOf(map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666}),
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := td.Cmp(t, got,
		td.SubMapOf(MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666}),
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.SubMapOf(&MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666}),
		"checks pointed typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpSubMapOf shortcut

```go
func CmpSubMapOf(t TestingT, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool
```

CmpSubMapOf is a shortcut for:

```go
td.Cmp(t, got, td.SubMapOf(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSubMapOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSubMapOf).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42}

	ok := td.CmpSubMapOf(t, got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := td.CmpSubMapOf(t, got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = td.CmpSubMapOf(t, &got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks pointed typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.SubMapOf shortcut

```go
func (t *T) SubMapOf(got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool
```

[`SubMapOf`]({{< ref "SubMapOf" >}}) is a shortcut for:

```go
t.Cmp(got, td.SubMapOf(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SubMapOf godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SubMapOf).

### Examples

{{%expand "Map example" %}}```go
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42}

	ok := t.SubMapOf(got, map[string]int{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := td.NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := t.SubMapOf(got, MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = t.SubMapOf(&got, &MyMap{"bar": 42}, td.MapEntries{"foo": td.Lt(15), "zip": 666},
		"checks pointed typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
