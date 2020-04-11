---
title: "Slice"
weight: 10
---

```go
func Slice(model interface{}, expectedEntries ArrayEntries) TestDeep
```

[`Slice`]({{< ref "Slice" >}}) operator compares the contents of a slice or a pointer on a
slice against the non-zero values of *model* (if any) and the
values of *expectedEntries*.

*model* must be the same type as compared data.

*expectedEntries* can be `nil`, if no zero entries are expected and
no [TestDeep operator]({{< ref "operators" >}}) are involved.

```go
got := []int{12, 14, 17}
td.Cmp(t, got, td.Slice([]int{0, 14}, td.ArrayEntries{0: 12, 2: 17})) // succeeds
td.Cmp(t, got,
  td.Slice([]int{0, 14}, td.ArrayEntries{0: td.Gt(10), 2: td.Gt(15)})) // succeeds
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *model*.


> See also [<i class='fas fa-book'></i> Slice godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Slice).

### Examples

{{%expand "Slice example" %}}```go
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := td.Cmp(t, got, td.Slice([]int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got,
		td.Slice([]int{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got,
		td.Slice(([]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "TypedSlice example" %}}```go
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := td.Cmp(t, got, td.Slice(MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.Slice(&MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Slice(&MySlice{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got,
		td.Slice((*MySlice)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()}),
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## CmpSlice shortcut

```go
func CmpSlice(t TestingT, got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool
```

CmpSlice is a shortcut for:

```go
td.Cmp(t, got, td.Slice(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpSlice godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpSlice).

### Examples

{{%expand "Slice example" %}}```go
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := td.CmpSlice(t, got, []int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, got, []int{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, got, ([]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "TypedSlice example" %}}```go
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := td.CmpSlice(t, got, MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, &got, &MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, &got, &MySlice{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = td.CmpSlice(t, &got, (*MySlice)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
## T.Slice shortcut

```go
func (t *T) Slice(got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool
```

[`Slice`]({{< ref "Slice" >}}) is a shortcut for:

```go
t.Cmp(got, td.Slice(model, expectedEntries), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Slice godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Slice).

### Examples

{{%expand "Slice example" %}}```go
	t := td.NewT(&testing.T{})

	got := []int{42, 58, 26}

	ok := t.Slice(got, []int{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(got, []int{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(got, ([]int)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "TypedSlice example" %}}```go
	t := td.NewT(&testing.T{})

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := t.Slice(got, MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks typed slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(&got, &MySlice{42}, td.ArrayEntries{1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(&got, &MySlice{}, td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	ok = t.Slice(&got, (*MySlice)(nil), td.ArrayEntries{0: 42, 1: 58, 2: td.Ignore()},
		"checks pointer on typed slice %v", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// true

```{{% /expand%}}
