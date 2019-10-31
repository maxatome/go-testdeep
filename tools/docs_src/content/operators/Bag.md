---
title: "Bag"
weight: 10
---

```go
func Bag(expectedItems ...interface{}) TestDeep
```

[`Bag`]({{< ref "Bag" >}}) operator compares the contents of an array or a slice (or a
pointer on array/slice) without taking care of the order of items.

During a match, each expected item should match in the compared
array/slice, and each array/slice item should be matched by an
expected item to succeed.

```go
Cmp(t, []int{1, 1, 2}, Bag(1, 1, 2))    // succeeds
Cmp(t, []int{1, 1, 2}, Bag(1, 2, 1))    // succeeds
Cmp(t, []int{1, 1, 2}, Bag(2, 1, 1))    // succeeds
Cmp(t, []int{1, 1, 2}, Bag(1, 2))       // fails, one 1 is missing
Cmp(t, []int{1, 1, 2}, Bag(1, 2, 1, 3)) // fails, 3 is missing
```


> See also [<i class='fas fa-book'></i> Bag godoc](https://godoc.org/github.com/maxatome/go-testdeep#Bag).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present
	ok := Cmp(t, got, Bag(1, 1, 2, 3, 5, 8, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Does not match as got contains 2 times 1 and 8, and these
	// duplicates are not expected
	ok = Cmp(t, got, Bag(1, 2, 3, 5, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 8, 2}

	// Duplicates of 1 and 8 are expected but not present in got
	ok = Cmp(t, got, Bag(1, 1, 2, 3, 5, 8, 8),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Matches as all items are present
	ok = Cmp(t, got, Bag(1, 2, 3, 5, Gt(7)),
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
	// true

```{{% /expand%}}
## CmpBag shortcut

```go
func CmpBag(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

CmpBag is a shortcut for:

```go
Cmp(t, got, Bag(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpBag godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpBag).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present
	ok := CmpBag(t, got, []interface{}{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Does not match as got contains 2 times 1 and 8, and these
	// duplicates are not expected
	ok = CmpBag(t, got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 8, 2}

	// Duplicates of 1 and 8 are expected but not present in got
	ok = CmpBag(t, got, []interface{}{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Matches as all items are present
	ok = CmpBag(t, got, []interface{}{1, 2, 3, 5, Gt(7)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
	// true

```{{% /expand%}}
## T.Bag shortcut

```go
func (t *T) Bag(got interface{}, expectedItems []interface{}, args ...interface{}) bool
```

[`Bag`]({{< ref "Bag" >}}) is a shortcut for:

```go
t.Cmp(got, Bag(expectedItems...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Bag godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Bag).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	got := []int{1, 3, 5, 8, 8, 1, 2}

	// Matches as all items are present
	ok := t.Bag(got, []interface{}{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Does not match as got contains 2 times 1 and 8, and these
	// duplicates are not expected
	ok = t.Bag(got, []interface{}{1, 2, 3, 5, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	got = []int{1, 3, 5, 8, 2}

	// Duplicates of 1 and 8 are expected but not present in got
	ok = t.Bag(got, []interface{}{1, 1, 2, 3, 5, 8, 8},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Matches as all items are present
	ok = t.Bag(got, []interface{}{1, 2, 3, 5, Gt(7)},
		"checks all items are present, in any order")
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false
	// true

```{{% /expand%}}
