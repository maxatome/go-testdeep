---
title: "Values"
weight: 10
---

```go
func Values(val interface{}) TestDeep
```

[`Values`]({{< ref "Values" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}). It takes a map and compares its
ordered values to *val*.

*val* can be a slice of items of the same type as the map values:

```go
got := map[int]string{3: "c", 1: "a", 2: "b"}
td.Cmp(t, got, td.Values([]string{"a", "b", "c"})) // succeeds, values sorted
td.Cmp(t, got, td.Values([]string{"c", "a", "b"})) // fails as not sorted
```

as well as an other operator as [`Bag`]({{< ref "Bag" >}}), for example, to test values in
an unsorted manner:

```go
got := map[int]string{3: "c", 1: "a", 2: "b"}
td.Cmp(t, got, td.Values(td.Bag("c", "a", "b"))) // succeeds
```


> See also [<i class='fas fa-book'></i> Values godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Values).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Values tests values in an ordered manner
	ok := td.Cmp(t, got, td.Values([]int{1, 2, 3}))
	fmt.Println("All sorted values are found:", ok)

	// If the expected values are not ordered, it fails
	ok = td.Cmp(t, got, td.Values([]int{3, 1, 2}))
	fmt.Println("All unsorted values are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = td.Cmp(t, got, td.Values(td.Bag(3, 1, 2)))
	fmt.Println("All unsorted values are found, with the help of Bag operator:", ok)

	// Check that each value is between 1 and 3
	ok = td.Cmp(t, got, td.Values(td.ArrayEach(td.Between(1, 3))))
	fmt.Println("Each value is between 1 and 3:", ok)

	// Output:
	// All sorted values are found: true
	// All unsorted values are found: false
	// All unsorted values are found, with the help of Bag operator: true
	// Each value is between 1 and 3: true

```{{% /expand%}}
## CmpValues shortcut

```go
func CmpValues(t TestingT, got interface{}, val interface{}, args ...interface{}) bool
```

CmpValues is a shortcut for:

```go
td.Cmp(t, got, td.Values(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpValues godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpValues).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Values tests values in an ordered manner
	ok := td.CmpValues(t, got, []int{1, 2, 3})
	fmt.Println("All sorted values are found:", ok)

	// If the expected values are not ordered, it fails
	ok = td.CmpValues(t, got, []int{3, 1, 2})
	fmt.Println("All unsorted values are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = td.CmpValues(t, got, td.Bag(3, 1, 2))
	fmt.Println("All unsorted values are found, with the help of Bag operator:", ok)

	// Check that each value is between 1 and 3
	ok = td.CmpValues(t, got, td.ArrayEach(td.Between(1, 3)))
	fmt.Println("Each value is between 1 and 3:", ok)

	// Output:
	// All sorted values are found: true
	// All unsorted values are found: false
	// All unsorted values are found, with the help of Bag operator: true
	// Each value is between 1 and 3: true

```{{% /expand%}}
## T.Values shortcut

```go
func (t *T) Values(got interface{}, val interface{}, args ...interface{}) bool
```

[`Values`]({{< ref "Values" >}}) is a shortcut for:

```go
t.Cmp(got, td.Values(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Values godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Values).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Values tests values in an ordered manner
	ok := t.Values(got, []int{1, 2, 3})
	fmt.Println("All sorted values are found:", ok)

	// If the expected values are not ordered, it fails
	ok = t.Values(got, []int{3, 1, 2})
	fmt.Println("All unsorted values are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = t.Values(got, td.Bag(3, 1, 2))
	fmt.Println("All unsorted values are found, with the help of Bag operator:", ok)

	// Check that each value is between 1 and 3
	ok = t.Values(got, td.ArrayEach(td.Between(1, 3)))
	fmt.Println("Each value is between 1 and 3:", ok)

	// Output:
	// All sorted values are found: true
	// All unsorted values are found: false
	// All unsorted values are found, with the help of Bag operator: true
	// Each value is between 1 and 3: true

```{{% /expand%}}
