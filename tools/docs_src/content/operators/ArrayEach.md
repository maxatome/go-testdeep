---
title: "ArrayEach"
weight: 10
---

```go
func ArrayEach(expectedValue interface{}) TestDeep
```

[`ArrayEach`]({{< ref "ArrayEach" >}}) operator has to be applied on arrays or slices or on
pointers on array/slice. It compares each item of data array/slice
against *expectedValue*. During a match, all items have to match to
succeed.


> See also [<i class='fas fa-book'></i> ArrayEach godoc](https://godoc.org/github.com/maxatome/go-testdeep#ArrayEach).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := Cmp(t, got, ArrayEach(Between(25, 60)),
		"checks each item of array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := Cmp(t, got, ArrayEach(Between(25, 60)),
		"checks each item of typed array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = Cmp(t, &got, ArrayEach(Between(25, 60)),
		"checks each item of typed array pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
{{%expand "Slice example" %}}```go
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := Cmp(t, got, ArrayEach(Between(25, 60)),
		"checks each item of slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedSlice example" %}}```go
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := Cmp(t, got, ArrayEach(Between(25, 60)),
		"checks each item of typed slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = Cmp(t, &got, ArrayEach(Between(25, 60)),
		"checks each item of typed slice pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpArrayEach shortcut

```go
func CmpArrayEach(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool
```

CmpArrayEach is a shortcut for:

```go
Cmp(t, got, ArrayEach(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpArrayEach godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpArrayEach).

### Examples

{{%expand "Array example" %}}```go
	t := &testing.T{}

	got := [3]int{42, 58, 26}

	ok := CmpArrayEach(t, got, Between(25, 60),
		"checks each item of array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := &testing.T{}

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := CmpArrayEach(t, got, Between(25, 60),
		"checks each item of typed array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = CmpArrayEach(t, &got, Between(25, 60),
		"checks each item of typed array pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
{{%expand "Slice example" %}}```go
	t := &testing.T{}

	got := []int{42, 58, 26}

	ok := CmpArrayEach(t, got, Between(25, 60),
		"checks each item of slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedSlice example" %}}```go
	t := &testing.T{}

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := CmpArrayEach(t, got, Between(25, 60),
		"checks each item of typed slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = CmpArrayEach(t, &got, Between(25, 60),
		"checks each item of typed slice pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.ArrayEach shortcut

```go
func (t *T) ArrayEach(got interface{}, expectedValue interface{}, args ...interface{}) bool
```

[`ArrayEach`]({{< ref "ArrayEach" >}}) is a shortcut for:

```go
t.Cmp(got, ArrayEach(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.ArrayEach godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.ArrayEach).

### Examples

{{%expand "Array example" %}}```go
	t := NewT(&testing.T{})

	got := [3]int{42, 58, 26}

	ok := t.ArrayEach(got, Between(25, 60),
		"checks each item of array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedArray example" %}}```go
	t := NewT(&testing.T{})

	type MyArray [3]int

	got := MyArray{42, 58, 26}

	ok := t.ArrayEach(got, Between(25, 60),
		"checks each item of typed array %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = t.ArrayEach(&got, Between(25, 60),
		"checks each item of typed array pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
{{%expand "Slice example" %}}```go
	t := NewT(&testing.T{})

	got := []int{42, 58, 26}

	ok := t.ArrayEach(got, Between(25, 60),
		"checks each item of slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedSlice example" %}}```go
	t := NewT(&testing.T{})

	type MySlice []int

	got := MySlice{42, 58, 26}

	ok := t.ArrayEach(got, Between(25, 60),
		"checks each item of typed slice %v is in [25 .. 60]", got)
	fmt.Println(ok)

	ok = t.ArrayEach(&got, Between(25, 60),
		"checks each item of typed slice pointer %v is in [25 .. 60]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
