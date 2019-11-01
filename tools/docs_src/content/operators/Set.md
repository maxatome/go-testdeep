---
title: "Set"
weight: 10
---

```go
func Set(expectedItems ...interface{}) TestDeep
```

[`Set`]({{< ref "Set" >}}) operator compares the contents of an array or a slice (or a
pointer on array/slice) ignoring duplicates and without taking care
of the order of items.

During a match, each expected item should match in the compared
array/slice, and each array/slice item should be matched by an
expected item to succeed.

```go
Cmp(t, []int{1, 1, 2}, Set(1, 2))    // succeeds
Cmp(t, []int{1, 1, 2}, Set(2, 1))    // succeeds
Cmp(t, []int{1, 1, 2}, Set(1, 2, 3)) // fails, 3 is missing
```


### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present, ignoring duplicates
	ok := Cmp(t, got, Set(1, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = Cmp(t, got, Set(1, 2, 2, 2, 2, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = Cmp(t, got, Set(Between(1, 4), 3, Between(2, 10)),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
## CmpSet shortcut

```go
func CmpSet(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

CmpSet is a shortcut for:

```go
Cmp(t, got, Set(expectedItems...), args...)
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

	// Matches as all items are present, ignoring duplicates
	ok := CmpSet(t, got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = CmpSet(t, got, []interface{}{1, 2, 2, 2, 2, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = CmpSet(t, got, []interface{}{Between(1, 4), 3, Between(2, 10)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
## T.Set shortcut

```go
func (t *T) Set(got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

[`Set`]({{< ref "Set" >}}) is a shortcut for:

```go
t.Cmp(got, Set(expectedItems...), args...)
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

	// Matches as all items are present, ignoring duplicates
	ok := t.Set(got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Duplicates are ignored in a Set
	ok = t.Set(got, []interface{}{1, 2, 2, 2, 2, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Tries its best to not raise an error when a value can be matched
	// by several Set entries
	ok = t.Set(got, []interface{}{Between(1, 4), 3, Between(2, 10)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
