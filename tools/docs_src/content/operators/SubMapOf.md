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
Cmp(t, map[string]int{"a": 1},
  SubMapOf(map[string]int{"a": 1, "b": 2}, nil)) // succeeds

Cmp(t, map[string]int{"a": 1, "c": 3},
  SubMapOf(map[string]int{"a": 1, "b": 2}, nil)) // fails, extra {"c": 3}
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> SubMapOf godoc](https://godoc.org/github.com/maxatome/go-testdeep#SubMapOf).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42}

	ok := Cmp(t, got,
		SubMapOf(map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666}),
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := Cmp(t, got,
		SubMapOf(MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666}),
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		SubMapOf(&MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666}),
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
Cmp(t, got, SubMapOf(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSubMapOf godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpSubMapOf).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42}

	ok := CmpSubMapOf(t, got, map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666},
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := CmpSubMapOf(t, got, MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666},
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = CmpSubMapOf(t, &got, &MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666},
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
t.Cmp(got, SubMapOf(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.SubMapOf godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.SubMapOf).

### Examples

{{%expand "Map example" %}}```go
	t := NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42}

	ok := t.SubMapOf(got, map[string]int{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666},
		"checks map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42}

	ok := t.SubMapOf(got, MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666},
		"checks typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	ok = t.SubMapOf(&got, &MyMap{"bar": 42}, MapEntries{"foo": Lt(15), "zip": 666},
		"checks pointed typed map %v is included in expected keys/values", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
