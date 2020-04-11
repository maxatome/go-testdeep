---
title: "Empty"
weight: 10
---

```go
func Empty() TestDeep
```

[`Empty`]({{< ref "Empty" >}}) operator checks that an array, a channel, a map, a slice or a
`string` is empty. As a special case (non-typed) `nil`, as well as `nil`
channel, map or slice are considered empty.

Note that the compared data can be a pointer (of pointer of pointer
etc.) on an array, a channel, a map, a slice or a `string`.

```go
td.Cmp(t, "", td.Empty())                // succeeds
td.Cmp(t, map[string]bool{}, td.Empty()) // succeeds
td.Cmp(t, []string{"foo"}, td.Empty())   // fails
```


> See also [<i class='fas fa-book'></i> Empty godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Empty).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t, nil, td.Empty()) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = td.Cmp(t, (*int)(nil), td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, "", td.Empty())
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = td.Cmp(t, 0, td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, (map[string]int)(nil), td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, map[string]int{}, td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, ([]int)(nil), td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, []int{}, td.Empty())
	fmt.Println(ok)

	ok = td.Cmp(t, []int{3}, td.Empty()) // fails, as not empty
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{}, td.Empty()) // fails, Empty() is not Zero()!
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
	// false
	// true
	// true
	// true
	// true
	// false
	// false

```{{% /expand%}}
{{%expand "Pointers example" %}}```go
	t := &testing.T{}

	type MySlice []int

	ok := td.Cmp(t, MySlice{}, td.Empty()) // Ptr() not needed
	fmt.Println(ok)

	ok = td.Cmp(t, &MySlice{}, td.Empty())
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = td.Cmp(t, &l3, td.Empty())
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = td.Cmp(t, &MyStruct{}, td.Empty()) // fails, use Zero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false

```{{% /expand%}}
## CmpEmpty shortcut

```go
func CmpEmpty(t TestingT, got interface{}, args ...interface{}) bool
```

CmpEmpty is a shortcut for:

```go
td.Cmp(t, got, td.Empty(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpEmpty godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpEmpty).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.CmpEmpty(t, nil) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = td.CmpEmpty(t, (*int)(nil))
	fmt.Println(ok)

	ok = td.CmpEmpty(t, "")
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = td.CmpEmpty(t, 0)
	fmt.Println(ok)

	ok = td.CmpEmpty(t, (map[string]int)(nil))
	fmt.Println(ok)

	ok = td.CmpEmpty(t, map[string]int{})
	fmt.Println(ok)

	ok = td.CmpEmpty(t, ([]int)(nil))
	fmt.Println(ok)

	ok = td.CmpEmpty(t, []int{})
	fmt.Println(ok)

	ok = td.CmpEmpty(t, []int{3}) // fails, as not empty
	fmt.Println(ok)

	ok = td.CmpEmpty(t, [3]int{}) // fails, Empty() is not Zero()!
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
	// false
	// true
	// true
	// true
	// true
	// false
	// false

```{{% /expand%}}
{{%expand "Pointers example" %}}```go
	t := &testing.T{}

	type MySlice []int

	ok := td.CmpEmpty(t, MySlice{}) // Ptr() not needed
	fmt.Println(ok)

	ok = td.CmpEmpty(t, &MySlice{})
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = td.CmpEmpty(t, &l3)
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = td.CmpEmpty(t, &MyStruct{}) // fails, use Zero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false

```{{% /expand%}}
## T.Empty shortcut

```go
func (t *T) Empty(got interface{}, args ...interface{}) bool
```

[`Empty`]({{< ref "Empty" >}}) is a shortcut for:

```go
t.Cmp(got, td.Empty(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Empty godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Empty).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	ok := t.Empty(nil) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = t.Empty((*int)(nil))
	fmt.Println(ok)

	ok = t.Empty("")
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = t.Empty(0)
	fmt.Println(ok)

	ok = t.Empty((map[string]int)(nil))
	fmt.Println(ok)

	ok = t.Empty(map[string]int{})
	fmt.Println(ok)

	ok = t.Empty(([]int)(nil))
	fmt.Println(ok)

	ok = t.Empty([]int{})
	fmt.Println(ok)

	ok = t.Empty([]int{3}) // fails, as not empty
	fmt.Println(ok)

	ok = t.Empty([3]int{}) // fails, Empty() is not Zero()!
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
	// false
	// true
	// true
	// true
	// true
	// false
	// false

```{{% /expand%}}
{{%expand "Pointers example" %}}```go
	t := td.NewT(&testing.T{})

	type MySlice []int

	ok := t.Empty(MySlice{}) // Ptr() not needed
	fmt.Println(ok)

	ok = t.Empty(&MySlice{})
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = t.Empty(&l3)
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = t.Empty(&MyStruct{}) // fails, use Zero() instead
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true
	// false

```{{% /expand%}}
