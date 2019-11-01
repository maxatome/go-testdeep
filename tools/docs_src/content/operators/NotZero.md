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
Cmp(t, AnyStruct{}, NotZero())       // is false
Cmp(t, &AnyStruct{}, NotZero())      // is true, coz pointer ≠ nil
Cmp(t, &AnyStruct{}, Ptr(NotZero())) // is false
```


### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := Cmp(t, 0, NotZero()) // fails
	fmt.Println(ok)

	ok = Cmp(t, float64(0), NotZero()) // fails
	fmt.Println(ok)

	ok = Cmp(t, 12, NotZero())
	fmt.Println(ok)

	ok = Cmp(t, (map[string]int)(nil), NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = Cmp(t, map[string]int{}, NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = Cmp(t, ([]int)(nil), NotZero()) // fails, as nil
	fmt.Println(ok)

	ok = Cmp(t, []int{}, NotZero()) // succeeds, as not nil
	fmt.Println(ok)

	ok = Cmp(t, [3]int{}, NotZero()) // fails
	fmt.Println(ok)

	ok = Cmp(t, [3]int{0, 1}, NotZero()) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = Cmp(t, bytes.Buffer{}, NotZero()) // fails
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, NotZero()) // succeeds, as pointer not nil
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, Ptr(NotZero())) // fails as deref by Ptr()
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
Cmp(t, got, NotZero(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := CmpNotZero(t, 0) // fails
	fmt.Println(ok)

	ok = CmpNotZero(t, float64(0)) // fails
	fmt.Println(ok)

	ok = CmpNotZero(t, 12)
	fmt.Println(ok)

	ok = CmpNotZero(t, (map[string]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = CmpNotZero(t, map[string]int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = CmpNotZero(t, ([]int)(nil)) // fails, as nil
	fmt.Println(ok)

	ok = CmpNotZero(t, []int{}) // succeeds, as not nil
	fmt.Println(ok)

	ok = CmpNotZero(t, [3]int{}) // fails
	fmt.Println(ok)

	ok = CmpNotZero(t, [3]int{0, 1}) // succeeds, DATA[1] is not 0
	fmt.Println(ok)

	ok = CmpNotZero(t, bytes.Buffer{}) // fails
	fmt.Println(ok)

	ok = CmpNotZero(t, &bytes.Buffer{}) // succeeds, as pointer not nil
	fmt.Println(ok)

	ok = Cmp(t, &bytes.Buffer{}, Ptr(NotZero())) // fails as deref by Ptr()
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
t.Cmp(got, NotZero(), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

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

	ok = t.Cmp(&bytes.Buffer{}, Ptr(NotZero())) // fails as deref by Ptr()
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
