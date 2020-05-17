---
title: "Delay"
weight: 10
---

```go
func Delay(delayed func() TestDeep) TestDeep
```

[`Delay`]({{< ref "Delay" >}}) operator allows to delay the construction of an operator to
the time it is used for the first time. Most of the time, it is
used with helpers. See the example for a very simple use case.


> See also [<i class='fas fa-book'></i> Delay godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Delay).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	cmpNow := func(expected td.TestDeep) bool {
		time.Sleep(time.Microsecond) // imagine a DB insert returning a CreatedAt
		return td.Cmp(t, time.Now(), expected)
	}

	before := time.Now()

	ok := cmpNow(td.Between(before, time.Now()))
	fmt.Println("Between called before compare:", ok)

	ok = cmpNow(td.Delay(func() td.TestDeep {
		return td.Between(before, time.Now())
	}))
	fmt.Println("Between delayed until compare:", ok)

	// Output:
	// Between called before compare: false
	// Between delayed until compare: true

```{{% /expand%}}
