---
title: "MapEach"
weight: 10
---

```go
func MapEach(expectedValue interface{}) TestDeep
```

[`MapEach`]({{< ref "MapEach" >}}) operator has to be applied on maps. It compares each value
of data map against expected value. During a match, all values have
to match to succeed.


> See also [<i class='fas fa-book'></i> MapEach godoc](https://godoc.org/github.com/maxatome/go-testdeep#MapEach).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := Cmp(t, got, MapEach(Between(10, 90)),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := Cmp(t, got, MapEach(Between(10, 90)),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = Cmp(t, &got, MapEach(Between(10, 90)),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## CmpMapEach shortcut

```go
func CmpMapEach(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool
```

CmpMapEach is a shortcut for:

```go
Cmp(t, got, MapEach(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpMapEach godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpMapEach).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := CmpMapEach(t, got, Between(10, 90),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := CmpMapEach(t, got, Between(10, 90),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = CmpMapEach(t, &got, Between(10, 90),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
## T.MapEach shortcut

```go
func (t *T) MapEach(got interface{}, expectedValue interface{}, args ...interface{}) bool
```

[`MapEach`]({{< ref "MapEach" >}}) is a shortcut for:

```go
t.Cmp(got, MapEach(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.MapEach godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.MapEach).

### Examples

{{%expand "Map example" %}}```go
	t := NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.MapEach(got, Between(10, 90),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.MapEach(got, Between(10, 90),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = t.MapEach(&got, Between(10, 90),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
