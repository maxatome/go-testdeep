---
title: "Keys"
weight: 10
---

```go
func Keys(val interface{}) TestDeep
```

[`Keys`]({{< ref "Keys" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}). It takes a map and compares its
ordered keys to *val*.

*val* can be a slice of items of the same type as the map keys:

```go
got := map[string]bool{"c": true, "a": false, "b": true}
td.Cmp(t, got, td.Keys([]string{"a", "b", "c"})) // succeeds, keys sorted
td.Cmp(t, got, td.Keys([]string{"c", "a", "b"})) // fails as not sorted
```

as well as an other operator as [`Bag`]({{< ref "Bag" >}}), for example, to test keys in
an unsorted manner:

```go
got := map[string]bool{"c": true, "a": false, "b": true}
td.Cmp(t, got, td.Keys(td.Bag("c", "a", "b"))) // succeeds
```


> See also [<i class='fas fa-book'></i> Keys godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Keys).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := td.Cmp(t, got, td.Keys([]string{"bar", "foo", "zip"}))
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = td.Cmp(t, got, td.Keys([]string{"zip", "bar", "foo"}))
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = td.Cmp(t, got, td.Keys(td.Bag("zip", "bar", "foo")))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = td.Cmp(t, got, td.Keys(td.ArrayEach(td.Len(3))))
	fmt.Println("Each key is 3 bytes long:", ok)

	// Output:
	// All sorted keys are found: true
	// All unsorted keys are found: false
	// All unsorted keys are found, with the help of Bag operator: true
	// Each key is 3 bytes long: true

```{{% /expand%}}
## CmpKeys shortcut

```go
func CmpKeys(t TestingT, got interface{}, val interface{}, args ...interface{}) bool
```

CmpKeys is a shortcut for:

```go
td.Cmp(t, got, td.Keys(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpKeys godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpKeys).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := td.CmpKeys(t, got, []string{"bar", "foo", "zip"})
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = td.CmpKeys(t, got, []string{"zip", "bar", "foo"})
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = td.CmpKeys(t, got, td.Bag("zip", "bar", "foo"))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = td.CmpKeys(t, got, td.ArrayEach(td.Len(3)))
	fmt.Println("Each key is 3 bytes long:", ok)

	// Output:
	// All sorted keys are found: true
	// All unsorted keys are found: false
	// All unsorted keys are found, with the help of Bag operator: true
	// Each key is 3 bytes long: true

```{{% /expand%}}
## T.Keys shortcut

```go
func (t *T) Keys(got interface{}, val interface{}, args ...interface{}) bool
```

[`Keys`]({{< ref "Keys" >}}) is a shortcut for:

```go
t.Cmp(got, td.Keys(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Keys godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Keys).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := t.Keys(got, []string{"bar", "foo", "zip"})
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = t.Keys(got, []string{"zip", "bar", "foo"})
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = t.Keys(got, td.Bag("zip", "bar", "foo"))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = t.Keys(got, td.ArrayEach(td.Len(3)))
	fmt.Println("Each key is 3 bytes long:", ok)

	// Output:
	// All sorted keys are found: true
	// All unsorted keys are found: false
	// All unsorted keys are found, with the help of Bag operator: true
	// Each key is 3 bytes long: true

```{{% /expand%}}
