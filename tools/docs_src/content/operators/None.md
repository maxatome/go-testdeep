---
title: "None"
weight: 10
---

```go
func None(notExpectedValues ...interface{}) TestDeep
```

[`None`]({{< ref "None" >}}) operator compares data against several not expected
values. During a match, none of them have to match to succeed.


> See also [<i class='fas fa-book'></i> None godoc](https://godoc.org/github.com/maxatome/go-testdeep#None).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 18

	ok := Cmp(t, got, None(0, 10, 20, 30, Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = Cmp(t, got, None(0, 10, 20, 30, Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = Cmp(t, got, None(0, 10, 20, 30, Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false

```{{% /expand%}}
## CmpNone shortcut

```go
func CmpNone(t TestingT, got interface{}, notExpectedValues []interface{}, args ...interface{}) bool
```

CmpNone is a shortcut for:

```go
Cmp(t, got, None(notExpectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNone godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpNone).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 18

	ok := CmpNone(t, got, []interface{}{0, 10, 20, 30, Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = CmpNone(t, got, []interface{}{0, 10, 20, 30, Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = CmpNone(t, got, []interface{}{0, 10, 20, 30, Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false

```{{% /expand%}}
## T.None shortcut

```go
func (t *T) None(got interface{}, notExpectedValues []interface{}, args ...interface{}) bool
```

[`None`]({{< ref "None" >}}) is a shortcut for:

```go
t.Cmp(got, None(notExpectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.None godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.None).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	got := 18

	ok := t.None(got, []interface{}{0, 10, 20, 30, Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = t.None(got, []interface{}{0, 10, 20, 30, Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = t.None(got, []interface{}{0, 10, 20, 30, Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	// Output:
	// true
	// false
	// false

```{{% /expand%}}
