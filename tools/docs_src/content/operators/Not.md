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

```go
td.Cmp(t, 12, td.Not(10)) // succeeds
td.Cmp(t, 12, td.Not(12)) // fails
```


> See also [<i class='fas fa-book'></i> Not godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Not).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 42

	ok := td.Cmp(t, got, td.Not(0), "checks %v is non-null", got)
	fmt.Println(ok)

	ok = td.Cmp(t, got, td.Not(td.Between(10, 30)),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = td.Cmp(t, got, td.Not(0), "checks %v is non-null", got)
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
td.Cmp(t, got, td.Not(notExpected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNot godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNot).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 42

	ok := td.CmpNot(t, got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	ok = td.CmpNot(t, got, td.Between(10, 30),
		"checks %v is not in [10 .. 30]", got)
	fmt.Println(ok)

	got = 0

	ok = td.CmpNot(t, got, 0, "checks %v is non-null", got)
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
t.Cmp(got, td.Not(notExpected), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Not godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Not).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := 42

	ok := t.Not(got, 0, "checks %v is non-null", got)
	fmt.Println(ok)

	ok = t.Not(got, td.Between(10, 30),
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
