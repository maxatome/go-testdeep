---
title: "Len"
weight: 10
---

```go
func Len(expectedLen interface{}) TestDeep
```

[`Len`]({{< ref "Len" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}). It takes data, applies `len()` function
on it and compares its result to *expectedLen*. Of course, the
compared value must be an array, a channel, a map, a slice or a
`string`.

*expectedLen* can be an `int` value:
```go
Cmp(t, gotSlice, Len(12))
```
as well as an other operator:
```go
Cmp(t, gotSlice, Len(Between(3, 4)))
```


### Examples

{{%expand "Slice example" %}}	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := Cmp(t, got, Len(3), "checks %v len is 3", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = Cmp(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
{{% /expand%}}
{{%expand "Map example" %}}	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := Cmp(t, got, Len(3), "checks %v len is 3", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = Cmp(t, got, Len(0), "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
{{% /expand%}}
{{%expand "OperatorSlice example" %}}	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := Cmp(t, got, Len(Between(3, 8)),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Len(Lt(5)), "checks %v len is < 5", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
{{% /expand%}}
{{%expand "OperatorMap example" %}}	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := Cmp(t, got, Len(Between(3, 8)),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Len(Gte(3)), "checks %v len is ≥ 3", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
{{% /expand%}}
## CmpLen shortcut

```go
func CmpLen(t TestingT, got interface{}, expectedLen interface{}, args ...interface{}) bool
```

CmpLen is a shortcut for:

```go
Cmp(t, got, Len(expectedLen), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Slice example" %}}	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := CmpLen(t, got, 3, "checks %v len is 3", got)
	fmt.Println(ok)

	ok = CmpLen(t, got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = CmpLen(t, got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
{{% /expand%}}
{{%expand "Map example" %}}	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := CmpLen(t, got, 3, "checks %v len is 3", got)
	fmt.Println(ok)

	ok = CmpLen(t, got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = CmpLen(t, got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
{{% /expand%}}
{{%expand "OperatorSlice example" %}}	t := &testing.T{}

	got := []int{11, 22, 33}

	ok := CmpLen(t, got, Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = CmpLen(t, got, Lt(5), "checks %v len is < 5", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
{{% /expand%}}
{{%expand "OperatorMap example" %}}	t := &testing.T{}

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := CmpLen(t, got, Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = CmpLen(t, got, Gte(3), "checks %v len is ≥ 3", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
{{% /expand%}}
## T.Len shortcut

```go
func (t *T) Len(got interface{}, expectedLen interface{}, args ...interface{}) bool
```

[`Len`]({{< ref "Len" >}}) is a shortcut for:

```go
t.Cmp(got, Len(expectedLen), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Slice example" %}}	t := NewT(&testing.T{})

	got := []int{11, 22, 33}

	ok := t.Len(got, 3, "checks %v len is 3", got)
	fmt.Println(ok)

	ok = t.Len(got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = t.Len(got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
{{% /expand%}}
{{%expand "Map example" %}}	t := NewT(&testing.T{})

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := t.Len(got, 3, "checks %v len is 3", got)
	fmt.Println(ok)

	ok = t.Len(got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	got = nil

	ok = t.Len(got, 0, "checks %v len is 0", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// true
{{% /expand%}}
{{%expand "OperatorSlice example" %}}	t := NewT(&testing.T{})

	got := []int{11, 22, 33}

	ok := t.Len(got, Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = t.Len(got, Lt(5), "checks %v len is < 5", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
{{% /expand%}}
{{%expand "OperatorMap example" %}}	t := NewT(&testing.T{})

	got := map[int]bool{11: true, 22: false, 33: false}

	ok := t.Len(got, Between(3, 8),
		"checks %v len is in [3 .. 8]", got)
	fmt.Println(ok)

	ok = t.Len(got, Gte(3), "checks %v len is ≥ 3", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
{{% /expand%}}
