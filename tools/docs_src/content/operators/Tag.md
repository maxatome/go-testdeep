---
title: "Tag"
weight: 10
---

```go
func Tag(tag string, expectedValue interface{}) TestDeep
```

[`Tag`]({{< ref "Tag" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}). It only allows to name *expectedValue*,
which can be an operator or a value. The data is then compared
against *expectedValue* as if [`Tag`]({{< ref "Tag" >}}) was never called. It is only
useful as [`JSON`]({{< ref "JSON" >}}) operator parameter, to name placeholders. See [`JSON`]({{< ref "JSON" >}})
operator for more details.

```go
td.Cmp(t, gotValue,
  td.JSON(`{"fullname": $name, "age": $age, "gender": $gender}`,
    td.Tag("name", td.HasPrefix("Foo")), // matches $name
    td.Tag("age", td.Between(41, 43)),   // matches $age
    td.Tag("gender", "male")))           // matches $gender
```

[`TypeBehind`]({{< ref "operators#typebehind-method" >}}) method is delegated to *expectedValue* one if
*expectedValue* is a [TestDeep operator]({{< ref "operators" >}}), otherwise it returns the
type of *expectedValue* (or `nil` if it is originally untyped `nil`).


> See also [<i class='fas fa-book'></i> Tag godoc](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Tag).

