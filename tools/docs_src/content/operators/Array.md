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

[TypeBehind]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> Array godoc](https://godoc.org/github.com/maxatome/go-testdeep#Array).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := Cmp(t, got, Array([3]int{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := Cmp(t, got, Array(MyArray{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got, Array(&MyArray{42}, ArrayEntries{1: 58, 2: Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Array(&MyArray{}, ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = Cmp(t, &got,
		Array((*MyArray)(nil), ArrayEntries{0: 42, 1: 58, 2: Ignore()}),
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
Cmp(t, got, Array(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpArray godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpArray).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := CmpArray(t, got, [3]int{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := CmpArray(t, got, MyArray{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = CmpArray(t, &got, &MyArray{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = CmpArray(t, &got, &MyArray{}, ArrayEntries{0: 42, 1: 58, 2: Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = CmpArray(t, &got, (*MyArray)(nil), ArrayEntries{0: 42, 1: 58, 2: Ignore()},
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
t.Cmp(got, Array(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Array godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Array).

### Examples

{{%expand "Array example" %}}```go
	t := NewT(&testing.T{})

	got := [3]int{42, 58, 26}

	ok := t.Array(got, [3]int{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks array %v", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := NewT(&testing.T{})

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := t.Array(got, MyArray{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks typed array %v", got)
	fmt.Println(ok)

	ok = t.Array(&got, &MyArray{42}, ArrayEntries{1: 58, 2: Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = t.Array(&got, &MyArray{}, ArrayEntries{0: 42, 1: 58, 2: Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	ok = t.Array(&got, (*MyArray)(nil), ArrayEntries{0: 42, 1: 58, 2: Ignore()},
		"checks pointer on typed array %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
