---
title: "Ignore"
weight: 10
---

```go
func Ignore() TestDeep
```

[`Ignore`]({{< ref "Ignore" >}}) operator is always true, whatever data is. It is useful when
comparing a slice with [`Slice`]({{< ref "Slice" >}}) and wanting to ignore some indexes,
for example. Or comparing a struct with [`SStruct`]({{< ref "SStruct" >}}) and wanting to
ignore some fields:

```go
td.Cmp(t, td.SStruct(
  Person{
    Name: "John Doe",
  },
  td.StructFields{
    Age:      td.Between(40, 45),
    Children: td.Ignore(),
  }),
)
```


> See also [<i class='fas fa-book'></i> Ignore godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Ignore).

### Example

{{%expand "Base example" %}}```go
	t := &testing.T{}

	ok := td.Cmp(t, []int{1, 2, 3},
		td.Slice([]int{}, td.ArrayEntries{
			0: 1,
			1: td.Ignore(), // do not care about this entry
			2: 3,
		}))
	fmt.Println(ok)

	// Output:
	// true

```{{% /expand%}}
