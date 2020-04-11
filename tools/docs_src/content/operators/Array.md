---
title: "Array"
weight: 10
---

```go
func Array(model interface{}, expectedEntries ArrayEntries) TestDeep
```

[`Array`]({{< ref "Array" >}}) operator compares the contents of an array or a pointer on an
array against the non-zero values of *model* (if any) and the
values of *expectedEntries*.

*model* must be the same type as compared data.

*expectedEntries* can be `nil`, if no zero entries are expected and
no [TestDeep operator]({{< ref "operators" >}}) are involved.

```go
got := [3]int{12, 14, 17}
td.Cmp(t, got, td.Array([3]int{0, 14}, td.ArrayEntries{0: 12, 2: 17})) // succeeds
td.Cmp(t, got,
  td.Array([3]int{0, 14}, td.ArrayEntries{0: td.Gt(10), 2: td.Gt(15)})) // succeeds
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> Array godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Array).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := td.Cmp(t, got,
		td.Array([3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := td.Cmp(t, got,
		td.Array(MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Array(&MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Array(&MyArray{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Array((*MyArray)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## CmpArray shortcut

```go
func CmpArray(t TestingT, got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool
```

CmpArray is a shortcut for:

```go
td.Cmp(t, got, td.Array(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpArray godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpArray).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := td.CmpArray(t, got, [3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := td.CmpArray(t, got, MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = td.CmpArray(t, &got, &MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = td.CmpArray(t, &got, &MyArray{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = td.CmpArray(t, &got, (*MyArray)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## T.Array shortcut

```go
func (t *T) Array(got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool
```

[`Array`]({{< ref "Array" >}}) is a shortcut for:

```go
t.Cmp(got, td.Array(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Array godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Array).

### Examples

{{%expand "Array example" %}}```go
	t := td.NewT(&testing.T{})

	got := [3]int{42, 58, 26}

	ok := t.Array(got, [3]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := td.NewT(&testing.T{})

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := t.Array(got, MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = t.Array(&got, &MyArray{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = t.Array(&got, &MyArray{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = t.Array(&got, (*MyArray)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
