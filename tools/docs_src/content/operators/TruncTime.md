---
title: "TruncTime"
weight: 10
---

```go
func TruncTime(expectedTime interface{}, trunc ...time.Duration) TestDeep
```

[`TruncTime`]({{< ref "TruncTime" >}}) operator compares [`time.Time`](https://golang.org/pkg/time/#Time) (or assignable) values after
truncating them to the optional *trunc* duration. See [`time.Truncate`](https://golang.org/pkg/time/#Truncate)
for details about the truncation.

If *trunc* is missing, it defaults to 0.

During comparison, location does not matter as [`time.Equal`](https://golang.org/pkg/time/#Equal) method is
used behind the scenes: a time instant in two different locations
is the same time instant.

Whatever the *trunc* value is, the monotonic clock is stripped
before the comparison against *expectedTime*.

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *expectedTime*.


> See also [<i class='fas fa-book'></i> TruncTime godoc](https://godoc.org/github.com/maxatome/go-testdeep#TruncTime).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	dateToTime := func(str string) time.Time {
		t, err := time.Parse(time.RFC3339Nano, str)
		if err != nil {
			panic(err)
		}
		return t
	}

	got := dateToTime("2018-05-01T12:45:53.123456789Z")

	// Compare dates ignoring nanoseconds and monotonic parts
	expected := dateToTime("2018-05-01T12:45:53Z")
	ok := Cmp(t, got, TruncTime(expected, time.Second),
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = Cmp(t, got, TruncTime(expected, 24*time.Hour),
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = Cmp(t, got, TruncTime(expected),
		"checks date %v ignoring monotonic part", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
## CmpTruncTime shortcut

```go
func CmpTruncTime(t TestingT, got interface{}, expectedTime interface{}, trunc time.Duration, args ...interface{}) bool
```

CmpTruncTime is a shortcut for:

```go
Cmp(t, got, TruncTime(expectedTime, trunc), args...)
```

See above for details.

[`TruncTime()`]({{< ref "TruncTime" >}}) optional parameter *trunc* is here mandatory.
0 value should be passed to mimic its absence in
original [`TruncTime()`]({{< ref "TruncTime" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpTruncTime godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpTruncTime).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	dateToTime := func(str string) time.Time {
		t, err := time.Parse(time.RFC3339Nano, str)
		if err != nil {
			panic(err)
		}
		return t
	}

	got := dateToTime("2018-05-01T12:45:53.123456789Z")

	// Compare dates ignoring nanoseconds and monotonic parts
	expected := dateToTime("2018-05-01T12:45:53Z")
	ok := CmpTruncTime(t, got, expected, time.Second,
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = CmpTruncTime(t, got, expected, 24*time.Hour,
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = CmpTruncTime(t, got, expected, 0,
		"checks date %v ignoring monotonic part", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
## T.TruncTime shortcut

```go
func (t *T) TruncTime(got interface{}, expectedTime interface{}, trunc time.Duration, args ...interface{}) bool
```

[`TruncTime`]({{< ref "TruncTime" >}}) is a shortcut for:

```go
t.Cmp(got, TruncTime(expectedTime, trunc), args...)
```

See above for details.

[`TruncTime()`]({{< ref "TruncTime" >}}) optional parameter *trunc* is here mandatory.
0 value should be passed to mimic its absence in
original [`TruncTime()`]({{< ref "TruncTime" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.TruncTime godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.TruncTime).

### Examples

{{%expand "Base example" %}}```go
	t := NewT(&testing.T{})

	dateToTime := func(str string) time.Time {
		t, err := time.Parse(time.RFC3339Nano, str)
		if err != nil {
			panic(err)
		}
		return t
	}

	got := dateToTime("2018-05-01T12:45:53.123456789Z")

	// Compare dates ignoring nanoseconds and monotonic parts
	expected := dateToTime("2018-05-01T12:45:53Z")
	ok := t.TruncTime(got, expected, time.Second,
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = t.TruncTime(got, expected, 24*time.Hour,
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = t.TruncTime(got, expected, 0,
		"checks date %v ignoring monotonic part", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// true

```{{% /expand%}}
