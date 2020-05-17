---
title: "NotZero"
weight: 10
---

```go
func NotZero() TestDeep
```

[`NotZero`]({{< ref "NotZero" >}}) operator checks that data is not zero regarding its type.

- `nil` is the zero value of pointers, maps, slices, channels and functions;
- 0 is the zero value of numbers;
- "" is the 0 value of strings;
- false is the zero value of booleans;
- zero value of structs is the struct with no fields initialized.


Beware that:

```go
td.Cmp(t, AnyStruct{}, td.NotZero())          // is false
td.Cmp(t, &AnyStruct{}, td.NotZero())         // is true, coz pointer ≠ nil
td.Cmp(t, &AnyStruct{}, td.Ptr(td.NotZero())) // is false
```


> See also [<i class='fas fa-book'></i> NotZero godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NotZero).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t, 0, td.NotZero()) // fails
	fmt.Println(ok)

	ok = td.Cmp(t, float64(0), td.NotZero()) // fails
	fmt.Println(ok)

	ok = td.Cmp(t, 12, td.NotZero())
	fmt.Println(ok)

	ok = td.Cmp(t, (map[string]int)(nil), td.NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = td.Cmp(t, map[string]int{}, td.NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = td.Cmp(t, ([]int)(nil), td.NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = td.Cmp(t, []int{}, td.NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{}, td.NotZero()) // fails
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{0, 1}, td.NotZero()) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = td.Cmp(t, bytes.Buffer{}, td.NotZero()) // fails
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.NotZero()) // succeeds, as pointer not nil
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.Ptr(td.NotZero())) // fails as deref by Ptr()
	fmt.Println(ok)

	// Output:
	// false
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false

```{{% /expand%}}
## CmpNotZero shortcut

```go
func CmpNotZero(t TestingT, got interface{}, args ...interface{}) bool
```

CmpNotZero is a shortcut for:

```go
td.Cmp(t, got, td.NotZero(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpNotZero godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNotZero).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.CmpNotZero(t, 0) // fails
	fmt.Println(ok)

	ok = td.CmpNotZero(t, float64(0)) // fails
	fmt.Println(ok)

	ok = td.CmpNotZero(t, 12)
	fmt.Println(ok)

	ok = td.CmpNotZero(t, (map[string]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = td.CmpNotZero(t, map[string]int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = td.CmpNotZero(t, ([]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = td.CmpNotZero(t, []int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = td.CmpNotZero(t, [3]int{}) // fails
	fmt.Println(ok)

	ok = td.CmpNotZero(t, [3]int{0, 1}) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = td.CmpNotZero(t, bytes.Buffer{}) // fails
	fmt.Println(ok)

	ok = td.CmpNotZero(t, &bytes.Buffer{}) // succeeds, as pointer not nil
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.Ptr(td.NotZero())) // fails as deref by Ptr()
	fmt.Println(ok)

	// Output:
	// false
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false

```{{% /expand%}}
## T.NotZero shortcut

```go
func (t *T) NotZero(got interface{}, args ...interface{}) bool
```

[`NotZero`]({{< ref "NotZero" >}}) is a shortcut for:

```go
t.Cmp(got, td.NotZero(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.NotZero godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.NotZero).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	ok := t.NotZero(0) // fails
	fmt.Println(ok)

	ok = t.NotZero(float64(0)) // fails
	fmt.Println(ok)

	ok = t.NotZero(12)
	fmt.Println(ok)

	ok = t.NotZero((map[string]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = t.NotZero(map[string]int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = t.NotZero(([]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = t.NotZero([]int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = t.NotZero([3]int{}) // fails
	fmt.Println(ok)

	ok = t.NotZero([3]int{0, 1}) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = t.NotZero(bytes.Buffer{}) // fails
	fmt.Println(ok)

	ok = t.NotZero(&bytes.Buffer{}) // succeeds, as pointer not nil
	fmt.Println(ok)

	ok = t.Cmp(&bytes.Buffer{}, td.Ptr(td.NotZero())) // fails as deref by Ptr()
	fmt.Println(ok)

	// Output:
	// false
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false

```{{% /expand%}}
