---
title: "NotEmpty"
weight: 10
---

```go
func NotEmpty() TestDeep
```

[`NotEmpty`]({{< ref "NotEmpty" >}}) operator checks that an array, a channel, a map, a slice
or a `string` is not empty. As a special case (non-typed) `nil`, as
well as `nil` channel, map or slice are considered empty.

Note that the compared data can be a pointer (of pointer of pointer
etc.) on an array, a channel, a map, a slice or a `string`.

```go
td.Cmp(t, "", td.NotEmpty())                // fails
td.Cmp(t, map[string]bool{}, td.NotEmpty()) // fails
td.Cmp(t, []string{"foo"}, td.NotEmpty())   // succeeds
```


> See also [<i class='fas fa-book'></i> NotEmpty godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotEmpty).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t, nil, td.NotEmpty()) // fails, as nil is considered empty
	fmt.Println(ok)

	ok = td.Cmp(t, "foobar", td.NotEmpty())
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use NotZero() instead
	ok = td.Cmp(t, 0, td.NotEmpty())
	fmt.Println(ok)

	ok = td.Cmp(t, map[string]int{"foobar": 42}, td.NotEmpty())
	fmt.Println(ok)

	ok = td.Cmp(t, []int{1}, td.NotEmpty())
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{}, td.NotEmpty()) // succeeds, NotEmpty() is not NotZero()!
	fmt.Println(ok)

	// Output:
	// false
	// true
	// false
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Pointers example" %}}```go
	t := &testing.T{}

	type MySlice []int

	ok := td.Cmp(t, MySlice{12}, td.NotEmpty())
	fmt.Println(ok)

	ok = td.Cmp(t, &MySlice{12}, td.NotEmpty()) // Ptr() not needed
	fmt.Println(ok)

	l1 := &MySlice{12}
	l2 := &l1
	l3 := &l2
	ok = td.Cmp(t, &l3, td.NotEmpty())
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = td.Cmp(t, &MyStruct{}, td.NotEmpty()) // fails, use NotZero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false

```{{% /expand%}}
## CmpNotEmpty shortcut

```go
func CmpNotEmpty(t TestingT, got interface{}, args ...interface{}) bool
```

CmpNotEmpty is a shortcut for:

```go
td.Cmp(t, got, td.NotEmpty(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNotEmpty godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNotEmpty).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.CmpNotEmpty(t, nil) // fails, as nil is considered empty
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, "foobar")
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use NotZero() instead
	ok = td.CmpNotEmpty(t, 0)
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, map[string]int{"foobar": 42})
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, []int{1})
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, [3]int{}) // succeeds, NotEmpty() is not NotZero()!
	fmt.Println(ok)

	// Output:
	// false
	// true
	// false
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Pointers example" %}}```go
	t := &testing.T{}

	type MySlice []int

	ok := td.CmpNotEmpty(t, MySlice{12})
	fmt.Println(ok)

	ok = td.CmpNotEmpty(t, &MySlice{12}) // Ptr() not needed
	fmt.Println(ok)

	l1 := &MySlice{12}
	l2 := &l1
	l3 := &l2
	ok = td.CmpNotEmpty(t, &l3)
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = td.CmpNotEmpty(t, &MyStruct{}) // fails, use NotZero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false

```{{% /expand%}}
## T.NotEmpty shortcut

```go
func (t *T) NotEmpty(got interface{}, args ...interface{}) bool
```

[`NotEmpty`]({{< ref "NotEmpty" >}}) is a shortcut for:

```go
t.Cmp(got, td.NotEmpty(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.NotEmpty godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.NotEmpty).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	ok := t.NotEmpty(nil) // fails, as nil is considered empty
	fmt.Println(ok)

	ok = t.NotEmpty("foobar")
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use NotZero() instead
	ok = t.NotEmpty(0)
	fmt.Println(ok)

	ok = t.NotEmpty(map[string]int{"foobar": 42})
	fmt.Println(ok)

	ok = t.NotEmpty([]int{1})
	fmt.Println(ok)

	ok = t.NotEmpty([3]int{}) // succeeds, NotEmpty() is not NotZero()!
	fmt.Println(ok)

	// Output:
	// false
	// true
	// false
	// true
	// true
	// true

```{{% /expand%}}
{{%expand "Pointers example" %}}```go
	t := td.NewT(&testing.T{})

	type MySlice []int

	ok := t.NotEmpty(MySlice{12})
	fmt.Println(ok)

	ok = t.NotEmpty(&MySlice{12}) // Ptr() not needed
	fmt.Println(ok)

	l1 := &MySlice{12}
	l2 := &l1
	l3 := &l2
	ok = t.NotEmpty(&l3)
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = t.NotEmpty(&MyStruct{}) // fails, use NotZero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false

```{{% /expand%}}
