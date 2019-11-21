---
title: "Catch"
weight: 10
---

```go
func Catch(target interface{}, expectedValue interface{}) TestDeep
```

[`Catch`]({{< ref "Catch" >}}) is a [smuggler operator]({{< ref "operators#smuggler-operators" >}}). It allows to copy data in *target* on
the fly before comparing it as usual against *expectedValue*.

*target* must be a non-`nil` pointer and data should be assignable to
its pointed type. If BeLax config flag is true or called under [`Lax`]({{< ref "Lax" >}})
(and so [`JSON`]({{< ref "JSON" >}})) operator, data should be convertible to its pointer
type.

```go
var id int64
if Cmp(t, CreateRecord("test"),
  JSON(`{"id": $1, "name": "test"}`, Catch(&id, NotZero()))) {
  t.Logf("Created record ID is %d", id)
}
```

It is really useful when used with [`JSON`]({{< ref "JSON" >}}) operator and/or tdhttp helper.

```go
var id int64
if tdhttp.CmpJSONResponse(t,
  tdhttp.NewRequest("POST", "/item", `{"name":"foo"}`),
  api.Handler,
  tdhttp.Response{
    Status: http.StatusCreated,
    Body: testdeep.JSON(`{"id": $id, "name": "foo"}`,
      testdeep.Tag("id", testdeep.Catch(&id, testdeep.Gt(0)))),
  }) {
  t.Logf("Created record ID is %d", id)
}
```

If you need to only catch data without comparing it, use [`Ignore`]({{< ref "Ignore" >}})
operator as *expectedValue* as in:

```go
var id int64
if Cmp(t, CreateRecord("test"),
  JSON(`{"id": $1, "name": "test"}`, Catch(&id, Ignore()))) {
  t.Logf("Created record ID is %d", id)
}
```


> See also [<i class='fas fa-book'></i> Catch godoc](https://godoc.org/github.com/maxatome/go-testdeep#Catch).

### Examples

{{%expand "Base example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	var age int
	ok := Cmp(t, got,
		JSON(`{"age":$1,"fullname":"Bob"}`,
			Catch(&age, Between(40, 45))))
	fmt.Println("check got age+fullname:", ok)
	fmt.Println("caught age:", age)

	// Output:
	// check got age+fullname: true
	// caught age: 42

```{{% /expand%}}