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

```go
got := map[string]string{"test": "foo", "buzz": "bar"}
td.Cmp(t, got, td.MapEach("bar"))     // fails, coz "foo" ≠ "bar"
td.Cmp(t, got, td.MapEach(td.Len(3))) // succeeds as values are 3 chars long
```


> See also [<i class='fas fa-book'></i> MapEach godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#MapEach).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got, td.MapEach(td.Between(10, 90)),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.Cmp(t, got, td.MapEach(td.Between(10, 90)),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = td.Cmp(t, &got, td.MapEach(td.Between(10, 90)),
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
td.Cmp(t, got, td.MapEach(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpMapEach godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpMapEach).

### Examples

{{%expand "Map example" %}}```go
	t := &testing.T{}

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpMapEach(t, got, td.Between(10, 90),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := &testing.T{}

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := td.CmpMapEach(t, got, td.Between(10, 90),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = td.CmpMapEach(t, &got, td.Between(10, 90),
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
t.Cmp(got, td.MapEach(expectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.MapEach godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.MapEach).

### Examples

{{%expand "Map example" %}}```go
	t := td.NewT(&testing.T{})

	got := map[string]int{"foo": 12, "bar": 42, "zip": 89}

	ok := t.MapEach(got, td.Between(10, 90),
		"checks each value of map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
{{%expand "TypedMap example" %}}```go
	t := td.NewT(&testing.T{})

	type MyMap map[string]int

	got := MyMap{"foo": 12, "bar": 42, "zip": 89}

	ok := t.MapEach(got, td.Between(10, 90),
		"checks each value of typed map %v is in [10 .. 90]", got)
	fmt.Println(ok)

	ok = t.MapEach(&got, td.Between(10, 90),
		"checks each value of typed map pointer %v is in [10 .. 90]", got)
	fmt.Println(ok)

	// Output:
	// true
	// true

```{{% /expand%}}
