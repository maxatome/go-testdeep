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
Keys([]string{"a", "b", "c"})
```
as well as an other operator:
```go
Keys(Bag("c", "a", "b"))
```


> See also [<i class='fas fa-book'></i> Keys godoc](https://godoc.org/github.com/maxatome/go-testdeep#Keys).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := Cmp(t, got, Keys([]string{"bar", "foo", "zip"}))
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = Cmp(t, got, Keys([]string{"zip", "bar", "foo"}))
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = Cmp(t, got, Keys(Bag("zip", "bar", "foo")))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = Cmp(t, got, Keys(ArrayEach(Len(3))))
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
Cmp(t, got, Keys(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpKeys godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpKeys).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := CmpKeys(t, got, []string{"bar", "foo", "zip"})
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = CmpKeys(t, got, []string{"zip", "bar", "foo"})
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = CmpKeys(t, got, Bag("zip", "bar", "foo"))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = CmpKeys(t, got, ArrayEach(Len(3)))
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
t.Cmp(got, Keys(val), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Keys godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Keys).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	got := map[string]int{"foo": 1, "bar": 2, "zip": 3}

	// Keys tests keys in an ordered manner
	ok := t.Keys(got, []string{"bar", "foo", "zip"})
	fmt.Println("All sorted keys are found:", ok)

	// If the expected keys are not ordered, it fails
	ok = t.Keys(got, []string{"zip", "bar", "foo"})
	fmt.Println("All unsorted keys are found:", ok)

	// To circumvent that, one can use Bag operator
	ok = t.Keys(got, Bag("zip", "bar", "foo"))
	fmt.Println("All unsorted keys are found, with the help of Bag operator:", ok)

	// Check that each key is 3 bytes long
	ok = t.Keys(got, ArrayEach(Len(3)))
	fmt.Println("Each key is 3 bytes long:", ok)

	// Output:
	// All sorted keys are found: true
	// All unsorted keys are found: false
	// All unsorted keys are found, with the help of Bag operator: true
	// Each key is 3 bytes long: true

```{{% /expand%}}
