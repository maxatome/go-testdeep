---
title: "Not"
weight: 10
---

```go
func Not(notExpected interface{}) TestDeep
```

[`Not`]({{< ref "Not" >}}) operator compares data against the not expected value. During a
match, it must not match to succeed.

[`Not`]({{< ref "Not" >}}) is the same operator as [`None()`]({{< ref "None" >}}) with only one argument. It is
provided as a more readable function when only one argument is
needed.


> See also [<i class='fas fa-book'></i> Not godoc](https://godoc.org/github.com/maxatome/go-testdeep#Not).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 42

	ok := Cmp(t, got, Not(0), "checks %v is non-null", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Not(Between(10, 30)),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = Cmp(t, got, Not(0), "checks %v is non-null", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false

```{{% /expand%}}
## CmpNot shortcut

```go
func CmpNot(t TestingT, got interface{}, notExpected interface{}, args ...interface{}) bool
```

CmpNot is a shortcut for:

```go
Cmp(t, got, Not(notExpected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNot godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpNot).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 42

	ok := CmpNot(t, got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	ok = CmpNot(t, got, Between(10, 30),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = CmpNot(t, got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false

```{{% /expand%}}
## T.Not shortcut

```go
func (t *T) Not(got interface{}, notExpected interface{}, args ...interface{}) bool
```

[`Not`]({{< ref "Not" >}}) is a shortcut for:

```go
t.Cmp(got, Not(notExpected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Not godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Not).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	got := 42

	ok := t.Not(got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	ok = t.Not(got, Between(10, 30),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = t.Not(got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false

```{{% /expand%}}
