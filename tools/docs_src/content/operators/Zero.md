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
Cmp(t, AnyStruct{}, Zero())       // is true
Cmp(t, &AnyStruct{}, Zero())      // is false, coz pointer ≠ nil
Cmp(t, &AnyStruct{}, Ptr(Zero())) // is true
```


> See also [<i class='fas fa-book'></i> Zero godoc](https://godoc.org/github.com/maxatome/go-testdeep#Zero).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := Cmp(t, 0, Zero())
	fmt.Println(ok)

	ok = Cmp(t, float64(0), Zero())
	fmt.Println(ok)

	ok = Cmp(t, 12, Zero()) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = Cmp(t, (map[string]int)(nil), Zero())
	fmt.Println(ok)

	ok = Cmp(t, map[string]int{}, Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = Cmp(t, ([]int)(nil), Zero())
	fmt.Println(ok)

	ok = Cmp(t, []int{}, Zero()) // fails, as not nil
	fmt.Println(ok)

	ok = Cmp(t, [3]int{}, Zero())
	fmt.Println(ok)

	ok = Cmp(t, [3]int{0, 1}, Zero()) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = Cmp(t, bytes.Buffer{}, Zero())
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, Zero()) // fails, as pointer not nil
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, Ptr(Zero())) // OK with the help of Ptr()
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
Cmp(t, got, Zero(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpZero godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpZero).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := CmpZero(t, 0)
	fmt.Println(ok)

	ok = CmpZero(t, float64(0))
	fmt.Println(ok)

	ok = CmpZero(t, 12) // fails, as 12 is not 0 :)
	fmt.Println(ok)

	ok = CmpZero(t, (map[string]int)(nil))
	fmt.Println(ok)

	ok = CmpZero(t, map[string]int{}) // fails, as not nil
	fmt.Println(ok)

	ok = CmpZero(t, ([]int)(nil))
	fmt.Println(ok)

	ok = CmpZero(t, []int{}) // fails, as not nil
	fmt.Println(ok)

	ok = CmpZero(t, [3]int{})
	fmt.Println(ok)

	ok = CmpZero(t, [3]int{0, 1}) // fails, DATA[1] is not 0
	fmt.Println(ok)

	ok = CmpZero(t, bytes.Buffer{})
	fmt.Println(ok)

	ok = CmpZero(t, &bytes.Buffer{}) // fails, as pointer not nil
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, Ptr(Zero())) // OK with the help of Ptr()
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
t.Cmp(got, Zero(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Zero godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.Zero).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

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

	ok = t.Cmp(&bytes.Buffer{}, Ptr(Zero())) // OK with the help of Ptr()
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
