---
title: "Shallow"
weight: 10
---

```go
func Shallow(expectedPtr interface{}) TestDeep
```

[`Shallow`]({{< ref "Shallow" >}}) operator compares pointers only, not their contents. It
applies on channels, functions (with some restrictions), maps,
pointers, slices and strings.

During a match, the compared data must be the same as *expectedPtr*
to succeed.

```go
a, b := 123, 123
td.Cmp(t, &a, td.Shallow(&a)) // succeeds
td.Cmp(t, &a, td.Shallow(&b)) // fails even if a == b as &a != &b

back := "foobarfoobar"
a, b := back[:6], back[6:]
// a == b but...
td.Cmp(t, &a, td.Shallow(&b)) // fails
```

Be careful for slices and strings! [`Shallow`]({{< ref "Shallow" >}}) can succeed but the
slices/strings not be identical because of their different
lengths. For example:

```go
a := "foobar yes!"
b := a[:1]                    // aka "f"
td.Cmp(t, &a, td.Shallow(&b)) // succeeds as both strings point to the same area, even if len() differ
```

The same behavior occurs for slices:

```go
a := []int{1, 2, 3, 4, 5, 6}
b := a[:2]                    // aka []int{1, 2}
td.Cmp(t, &a, td.Shallow(&b)) // succeeds as both slices point to the same area, even if len() differ
```


> See also [<i class='fas fa-book'></i> Shallow godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Shallow).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type MyStruct struct {
		Value int
	}
	data := MyStruct{Value: 12}
	got := &data

	ok := td.Cmp(t, got, td.Shallow(&data),
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Same contents, but not same pointer
	ok = td.Cmp(t, got, td.Shallow(&MyStruct{Value: 12}),
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "Slice example" %}}```go
	t := &testing.T{}

	back := []int{1, 2, 3, 1, 2, 3}
	a := back[:3]
	b := back[3:]

	ok := td.Cmp(t, a, td.Shallow(back))
	fmt.Println("are ≠ but share the same area:", ok)

	ok = td.Cmp(t, b, td.Shallow(back))
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	back := "foobarfoobar"
	a := back[:6]
	b := back[6:]

	ok := td.Cmp(t, a, td.Shallow(back))
	fmt.Println("are ≠ but share the same area:", ok)

	ok = td.Cmp(t, b, td.Shallow(a))
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false

```{{% /expand%}}
## CmpShallow shortcut

```go
func CmpShallow(t TestingT, got interface{}, expectedPtr interface{}, args ...interface{}) bool
```

CmpShallow is a shortcut for:

```go
td.Cmp(t, got, td.Shallow(expectedPtr), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpShallow godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpShallow).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	type MyStruct struct {
		Value int
	}
	data := MyStruct{Value: 12}
	got := &data

	ok := td.CmpShallow(t, got, &data,
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Same contents, but not same pointer
	ok = td.CmpShallow(t, got, &MyStruct{Value: 12},
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "Slice example" %}}```go
	t := &testing.T{}

	back := []int{1, 2, 3, 1, 2, 3}
	a := back[:3]
	b := back[3:]

	ok := td.CmpShallow(t, a, back)
	fmt.Println("are ≠ but share the same area:", ok)

	ok = td.CmpShallow(t, b, back)
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := &testing.T{}

	back := "foobarfoobar"
	a := back[:6]
	b := back[6:]

	ok := td.CmpShallow(t, a, back)
	fmt.Println("are ≠ but share the same area:", ok)

	ok = td.CmpShallow(t, b, a)
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false

```{{% /expand%}}
## T.Shallow shortcut

```go
func (t *T) Shallow(got interface{}, expectedPtr interface{}, args ...interface{}) bool
```

[`Shallow`]({{< ref "Shallow" >}}) is a shortcut for:

```go
t.Cmp(got, td.Shallow(expectedPtr), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.Shallow godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Shallow).

### Examples

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

	type MyStruct struct {
		Value int
	}
	data := MyStruct{Value: 12}
	got := &data

	ok := t.Shallow(got, &data,
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Same contents, but not same pointer
	ok = t.Shallow(got, &MyStruct{Value: 12},
		"checks pointers only, not contents")
	fmt.Println(ok)

	// Output:
	// true
	// false

```{{% /expand%}}
{{%expand "Slice example" %}}```go
	t := td.NewT(&testing.T{})

	back := []int{1, 2, 3, 1, 2, 3}
	a := back[:3]
	b := back[3:]

	ok := t.Shallow(a, back)
	fmt.Println("are ≠ but share the same area:", ok)

	ok = t.Shallow(b, back)
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false

```{{% /expand%}}
{{%expand "String example" %}}```go
	t := td.NewT(&testing.T{})

	back := "foobarfoobar"
	a := back[:6]
	b := back[6:]

	ok := t.Shallow(a, back)
	fmt.Println("are ≠ but share the same area:", ok)

	ok = t.Shallow(b, a)
	fmt.Println("are = but do not point to same area:", ok)

	// Output:
	// are ≠ but share the same area: true
	// are = but do not point to same area: false

```{{% /expand%}}
