---
title: "Zero"
weight: 10
---

```go
func Zero() TestDeep
```

[`Zero`]({{< ref "Zero" >}}) operator checks that data is zero regarding its type.

- `nil` is the zero value of pointers, maps, slices, channels and functions;
- 0 is the zero value of numbers;
- "" is the 0 value of strings;
- false is the zero value of booleans;
- zero value of structs is the struct with no fields initialized.


Beware that:

```go
td.Cmp(t, AnyStruct{}, td.Zero())          // is true
td.Cmp(t, &AnyStruct{}, td.Zero())         // is false, coz pointer ≠ nil
td.Cmp(t, &AnyStruct{}, td.Ptr(td.Zero())) // is true
```


> See also [<i class='fas fa-book'></i> Zero godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Zero).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t, 0, td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, float64(0), td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, 12, td.Zero()) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = td.Cmp(t, (map[string]int)(nil), td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, map[string]int{}, td.Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = td.Cmp(t, ([]int)(nil), td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, []int{}, td.Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{}, td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, [3]int{0, 1}, td.Zero()) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = td.Cmp(t, bytes.Buffer{}, td.Zero())
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.Zero()) // fails, as pointer not nil
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.Ptr(td.Zero())) // OK with the help of Ptr()
	fmt.Println(ok)

	// Output:
	// true
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
	// true

```{{% /expand%}}
## CmpZero shortcut

```go
func CmpZero(t TestingT, got interface{}, args ...interface{}) bool
```

CmpZero is a shortcut for:

```go
td.Cmp(t, got, td.Zero(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpZero godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpZero).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.CmpZero(t, 0)
	fmt.Println(ok)

	ok = td.CmpZero(t, float64(0))
	fmt.Println(ok)

	ok = td.CmpZero(t, 12) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = td.CmpZero(t, (map[string]int)(nil))
	fmt.Println(ok)

	ok = td.CmpZero(t, map[string]int{}) // fails, as not nil
	fmt.Println(ok)

	ok = td.CmpZero(t, ([]int)(nil))
	fmt.Println(ok)

	ok = td.CmpZero(t, []int{}) // fails, as not nil
	fmt.Println(ok)

	ok = td.CmpZero(t, [3]int{})
	fmt.Println(ok)

	ok = td.CmpZero(t, [3]int{0, 1}) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = td.CmpZero(t, bytes.Buffer{})
	fmt.Println(ok)

	ok = td.CmpZero(t, &bytes.Buffer{}) // fails, as pointer not nil
	fmt.Println(ok)

	ok = td.Cmp(t, &bytes.Buffer{}, td.Ptr(td.Zero())) // OK with the help of Ptr()
	fmt.Println(ok)

	// Output:
	// true
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
	// true

```{{% /expand%}}
## T.Zero shortcut

```go
func (t *T) Zero(got interface{}, args ...interface{}) bool
```

[`Zero`]({{< ref "Zero" >}}) is a shortcut for:

```go
t.Cmp(got, td.Zero(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Zero godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Zero).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	ok := t.Zero(0)
	fmt.Println(ok)

	ok = t.Zero(float64(0))
	fmt.Println(ok)

	ok = t.Zero(12) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = t.Zero((map[string]int)(nil))
	fmt.Println(ok)

	ok = t.Zero(map[string]int{}) // fails, as not nil
	fmt.Println(ok)

	ok = t.Zero(([]int)(nil))
	fmt.Println(ok)

	ok = t.Zero([]int{}) // fails, as not nil
	fmt.Println(ok)

	ok = t.Zero([3]int{})
	fmt.Println(ok)

	ok = t.Zero([3]int{0, 1}) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = t.Zero(bytes.Buffer{})
	fmt.Println(ok)

	ok = t.Zero(&bytes.Buffer{}) // fails, as pointer not nil
	fmt.Println(ok)

	ok = t.Cmp(&bytes.Buffer{}, td.Ptr(td.Zero())) // OK with the help of Ptr()
	fmt.Println(ok)

	// Output:
	// true
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
	// true

```{{% /expand%}}
