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

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> Map godoc](https://godoc.org/github.com/maxatome/go-testdeep#Map).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := Cmp(t, got,
		Map(map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, got,
		Map(map[string]int{},
			MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
		"checks map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, got,
		Map((map[string]int)(nil),
			MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
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

	ok := Cmp(t, got,
		Map(MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()}),
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Map(&MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Map(&MyMap{}, MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Map((*MyMap)(nil), MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()}),
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
Cmp(t, got, Map(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpMap godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpMap).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := CmpMap(t, got, map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = CmpMap(t, got, map[string]int{}, MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = CmpMap(t, got, (map[string]int)(nil), MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()},
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

	ok := CmpMap(t, got, MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()},
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = CmpMap(t, &got, &MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = CmpMap(t, &got, &MyMap{}, MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = CmpMap(t, &got, (*MyMap)(nil), MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()},
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
t.Cmp(got, Map(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Map godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Map).

### Examples

{{%expand "Map example" %}}```go
	t := NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.Map(got, map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = t.Map(got, map[string]int{}, MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	ok = t.Map(got, (map[string]int)(nil), MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()},
		"checks map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.Map(got, MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()},
		"checks typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, &MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, &MyMap{}, MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	ok = t.Map(&got, (*MyMap)(nil), MapEntries{"bar": 42, "foo": Lt(15), "zip": Ignore()},
		"checks pointer on typed map %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
