---
title: "Ignore"
weight: 10
---

```go
func Ignore() TestDeep
```

[`Ignore`]({{< ref "Ignore" >}}) operator is always true, whatever data is. It is useful when
comparing a slice and wanting to ignore some indexes, for example.


> See also [<i class='fas fa-book'></i> Ignore godoc](https://godoc.org/github.com/maxatome/go-testdeep#Ignore).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := Cmp(t, []int{1, 2, 3},
		Slice([]int{}, ArrayEntries{
			0: 1,
			1: Ignore(), // do not care about this entry
			2: 3,
		}))
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
