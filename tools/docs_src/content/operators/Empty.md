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


> See also [<i class='fas fa-book'></i> Empty godoc](https://godoc.org/github.com/maxatome/go-testdeep#Empty).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := Cmp(t, nil, Empty()) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = Cmp(t, (*int)(nil), Empty())
	fmt.Println(ok)

	ok = Cmp(t, "", Empty())
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = Cmp(t, 0, Empty())
	fmt.Println(ok)

	ok = Cmp(t, (map[string]int)(nil), Empty())
	fmt.Println(ok)

	ok = Cmp(t, map[string]int{}, Empty())
	fmt.Println(ok)

	ok = Cmp(t, ([]int)(nil), Empty())
	fmt.Println(ok)

	ok = Cmp(t, []int{}, Empty())
	fmt.Println(ok)

	ok = Cmp(t, []int{3}, Empty()) // fails, as not empty
	fmt.Println(ok)

	ok = Cmp(t, [3]int{}, Empty()) // fails, Empty() is not Zero()!
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

	ok := Cmp(t, MySlice{}, Empty()) // Ptr() not needed
	fmt.Println(ok)

	ok = Cmp(t, &MySlice{}, Empty())
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = Cmp(t, &l3, Empty())
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = Cmp(t, &MyStruct{}, Empty()) // fails, use Zero() instead
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
Cmp(t, got, Empty(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpEmpty godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpEmpty).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := CmpEmpty(t, nil) // special case: nil is considered empty
	fmt.Println(ok)

	// fails, typed nil is not empty (expect for channel, map, slice or
	// pointers on array, channel, map slice and strings)
	ok = CmpEmpty(t, (*int)(nil))
	fmt.Println(ok)

	ok = CmpEmpty(t, "")
	fmt.Println(ok)

	// Fails as 0 is a number, so not empty. Use Zero() instead
	ok = CmpEmpty(t, 0)
	fmt.Println(ok)

	ok = CmpEmpty(t, (map[string]int)(nil))
	fmt.Println(ok)

	ok = CmpEmpty(t, map[string]int{})
	fmt.Println(ok)

	ok = CmpEmpty(t, ([]int)(nil))
	fmt.Println(ok)

	ok = CmpEmpty(t, []int{})
	fmt.Println(ok)

	ok = CmpEmpty(t, []int{3}) // fails, as not empty
	fmt.Println(ok)

	ok = CmpEmpty(t, [3]int{}) // fails, Empty() is not Zero()!
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

	ok := CmpEmpty(t, MySlice{}) // Ptr() not needed
	fmt.Println(ok)

	ok = CmpEmpty(t, &MySlice{})
	fmt.Println(ok)

	l1 := &MySlice{}
	l2 := &l1
	l3 := &l2
	ok = CmpEmpty(t, &l3)
	fmt.Println(ok)

	// Works the same for array, map, channel and string

	// But not for others types as:
	type MyStruct struct {
		Value int
	}

	ok = CmpEmpty(t, &MyStruct{}) // fails, use Zero() instead
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
t.Cmp(got, Empty(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Empty godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Empty).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

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
	t := NewT(&testing.T{})

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
