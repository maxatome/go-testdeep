---
title: "None"
weight: 10
---

```go
func None(notExpectedValues ...interface{}) TestDeep
```

[`None`]({{< ref "None" >}}) operator compares data against several not expected
values. During a match, none of them have to match to succeed.

```go
td.Cmp(t, 12, td.None(8, 10, 14))     // succeeds
td.Cmp(t, 12, td.None(8, 10, 12, 14)) // fails
```

Note [`Flatten`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Flatten) function can be used to group or reuse some values or
operators and so avoid boring and inefficient copies:

```go
prime := td.Flatten([]int{1, 2, 3, 5, 7, 11, 13})
even := td.Flatten([]int{2, 4, 6, 8, 10, 12, 14})
td.Cmp(t, 9, td.None(prime, even)) // succeeds
```


> See also [<i class='fas fa-book'></i> None godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#None).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 18

	ok := td.Cmp(t, got, td.None(0, 10, 20, 30, td.Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = td.Cmp(t, got, td.None(0, 10, 20, 30, td.Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = td.Cmp(t, got, td.None(0, 10, 20, 30, td.Between(100, 199)),
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	prime := td.Flatten([]int{1, 2, 3, 5, 7, 11, 13})
	even := td.Flatten([]int{2, 4, 6, 8, 10, 12, 14})
	for _, got := range [...]int{9, 3, 8, 15} {
		ok = td.Cmp(t, got, td.None(prime, even, td.Gt(14)),
			"checks %v is not prime number, nor an even number and not > 14")
		fmt.Printf("%d → %t\n", got, ok)
	}

	// Output:
	// true
	// false
	// false
	// 9 → true
	// 3 → false
	// 8 → false
	// 15 → false

```{{% /expand%}}
## CmpNone shortcut

```go
func CmpNone(t TestingT, got interface{}, notExpectedValues []interface{}, args ...interface{}) bool
```

CmpNone is a shortcut for:

```go
td.Cmp(t, got, td.None(notExpectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNone godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNone).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := 18

	ok := td.CmpNone(t, got, []interface{}{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = td.CmpNone(t, got, []interface{}{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = td.CmpNone(t, got, []interface{}{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	prime := td.Flatten([]int{1, 2, 3, 5, 7, 11, 13})
	even := td.Flatten([]int{2, 4, 6, 8, 10, 12, 14})
	for _, got := range [...]int{9, 3, 8, 15} {
		ok = td.CmpNone(t, got, []interface{}{prime, even, td.Gt(14)},
			"checks %v is not prime number, nor an even number and not > 14")
		fmt.Printf("%d → %t\n", got, ok)
	}

	// Output:
	// true
	// false
	// false
	// 9 → true
	// 3 → false
	// 8 → false
	// 15 → false

```{{% /expand%}}
## T.None shortcut

```go
func (t *T) None(got interface{}, notExpectedValues []interface{}, args ...interface{}) bool
```

[`None`]({{< ref "None" >}}) is a shortcut for:

```go
t.Cmp(got, td.None(notExpectedValues...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.None godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.None).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	got := 18

	ok := t.None(got, []interface{}{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 20

	ok = t.None(got, []interface{}{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	got = 142

	ok = t.None(got, []interface{}{0, 10, 20, 30, td.Between(100, 199)},
		"checks %v is non-null, and ≠ 10, 20 & 30, and not in [100-199]", got)
	fmt.Println(ok)

	prime := td.Flatten([]int{1, 2, 3, 5, 7, 11, 13})
	even := td.Flatten([]int{2, 4, 6, 8, 10, 12, 14})
	for _, got := range [...]int{9, 3, 8, 15} {
		ok = t.None(got, []interface{}{prime, even, td.Gt(14)},
			"checks %v is not prime number, nor an even number and not > 14")
		fmt.Printf("%d → %t\n", got, ok)
	}

	// Output:
	// true
	// false
	// false
	// 9 → true
	// 3 → false
	// 8 → false
	// 15 → false

```{{% /expand%}}
