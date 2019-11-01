---
title: "SuperSetOf"
weight: 10
---

```go
func SuperSetOf(expectedItems ...interface{}) TestDeep
```

[`SuperSetOf`]({{< ref "SuperSetOf" >}}) operator compares the contents of an array or a slice (or
a pointer on array/slice) ignoring duplicates and without taking
care of the order of items.

During a match, each expected item should match in the compared
array/slice. But some items in the compared array/slice may not be
expected.

```go
Cmp(t, []int{1, 1, 2}, SuperSetOf(1))    // succeeds
Cmp(t, []int{1, 1, 2}, SuperSetOf(1, 3)) // fails, 3 is missing
```


### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := Cmp(t, got, SuperSetOf(1, 2, 3),
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = Cmp(t, got, SuperSetOf(Gt(5), Lte(2)),
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpSuperSetOf shortcut

```go
func CmpSuperSetOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

CmpSuperSetOf is a shortcut for:

```go
Cmp(t, got, SuperSetOf(expectedItems...), args...)
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

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := CmpSuperSetOf(t, got, []interface{}{1, 2, 3},
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = CmpSuperSetOf(t, got, []interface{}{Gt(5), Lte(2)},
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.SuperSetOf shortcut

```go
func (t *T) SuperSetOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

[`SuperSetOf`]({{< ref "SuperSetOf" >}}) is a shortcut for:

```go
t.Cmp(got, SuperSetOf(expectedItems...), args...)
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

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	ok := t.SuperSetOf(got, []interface{}{1, 2, 3},
		"checks the items are present, in any order and ignoring duplicates")
	fmt.Println(ok)

	ok = t.SuperSetOf(got, []interface{}{Gt(5), Lte(2)},
		"checks at least 2 items of %v match ignoring duplicates", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
