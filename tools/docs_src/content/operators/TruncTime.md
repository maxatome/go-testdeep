---
title: "TruncTime"
weight: 10
---

```go
func TruncTime(expectedTime interface{}, trunc ...time.Duration) TestDeep
```

[`TruncTime`]({{< ref "TruncTime" >}}) operator compares [`time.Time`](https://pkg.go.dev/time/#Time) (or assignable) values after
truncating them to the optional *trunc* duration. See [`time.Truncate`](https://pkg.go.dev/time/#Truncate)
for details about the truncation.

If *trunc* is missing, it defaults to 0.

During comparison, location does not matter as [`time.Equal`](https://pkg.go.dev/time/#Equal) method is
used behind the scenes: a time instant in two different locations
is the same time instant.

Whatever the *trunc* value is, the monotonic clock is stripped
before the comparison against *expectedTime*.

```go
gotDate := time.Date(2018, time.March, 9, 1, 2, 3, 999999999, time.UTC).
  In(time.FixedZone("UTC+2", 2))

expected := time.Date(2018, time.March, 9, 1, 2, 3, 0, time.UTC)

td.Cmp(t, gotDate, td.TruncTime(expected))              // fails, ns differ
td.Cmp(t, gotDate, td.TruncTime(expected, time.Second)) // succeeds
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://pkg.go.dev/reflect/#Type) of *expectedTime*.


> See also [<i class='fas fa-book'></i> TruncTime godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#TruncTime).

### Example

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
	ok := td.Cmp(t, got, td.TruncTime(expected, time.Second),
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = td.Cmp(t, got, td.TruncTime(expected, 24*time.Hour),
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = td.Cmp(t, got, td.TruncTime(expected),
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
td.Cmp(t, got, td.TruncTime(expectedTime, trunc), args...)
```

See above for details.

[`TruncTime()`]({{< ref "TruncTime" >}}) optional parameter *trunc* is here mandatory.
0 value should be passed to mimic its absence in
original [`TruncTime()`]({{< ref "TruncTime" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpTruncTime godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpTruncTime).

### Example

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
	ok := td.CmpTruncTime(t, got, expected, time.Second,
		"checks date %v, truncated to the second", got)
	fmt.Println(ok)

	// Compare dates ignoring time and so monotonic parts
	expected = dateToTime("2018-05-01T11:22:33.444444444Z")
	ok = td.CmpTruncTime(t, got, expected, 24*time.Hour,
		"checks date %v, truncated to the day", got)
	fmt.Println(ok)

	// Compare dates exactly but ignoring monotonic part
	expected = dateToTime("2018-05-01T12:45:53.123456789Z")
	ok = td.CmpTruncTime(t, got, expected, 0,
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
t.Cmp(got, td.TruncTime(expectedTime, trunc), args...)
```

See above for details.

[`TruncTime()`]({{< ref "TruncTime" >}}) optional parameter *trunc* is here mandatory.
0 value should be passed to mimic its absence in
original [`TruncTime()`]({{< ref "TruncTime" >}}) call.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://pkg.go.dev/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://pkg.go.dev/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.TruncTime godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.TruncTime).

### Example

{{%expand "Base example" %}}```go
	t := td.NewT(&testing.T{})

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
