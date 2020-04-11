+++
title = "Introduction"
weight = 5
+++

## Synopsis

Make golang tests easy, from simplest usage:

```go
import (
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestMyFunc(t *testing.T) {
  td.Cmp(t, MyFunc(), &Info{Name: "Alice", Age: 42})
}
```

To a bit more complex one, allowing flexible comparisons using
[TestDeep operators]({{< ref "operators" >}}):

```go
import (
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestMyFunc(t *testing.T) {
  td.Cmp(t, MyFunc(), td.Struct(
    &Info{Name: "Alice"},
    td.StructFields{
      "Age": td.Between(40, 45),
    },
  ))
}
```

Or anchoring operators directly in literals, as in:

```go
import (
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestMyFunc(tt *testing.T) {
  t := td.NewT(tt)

  t.Cmp(MyFunc(), &Info{
    Name: "Alice",
    Age:  t.Anchor(td.Between(40, 45)).(int),
  })
}
```

To most complex one, allowing to easily test HTTP API routes, using
flexible [operators]({{< ref "operators" >}}) and the
[`tdhttp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp)
helper:

```go
import (
  "testing"
  "time"

  "github.com/maxatome/go-testdeep/helpers/tdhttp"
  "github.com/maxatome/go-testdeep/td"
)

type Person struct {
  ID        uint64    `json:"id"`
  Name      string    `json:"name"`
  Age       int       `json:"age"`
  CreatedAt time.Time `json:"created_at"`
}

func TestMyApi(t *testing.T) {
  var id uint64
  var createdAt time.Time

  testAPI := tdhttp.NewTestAPI(t, myAPI) // ← ①

  testAPI.PostJSON("/person", Person{Name: "Bob", Age: 42}). // ← ②
    Name("Create a new Person").
    CmpStatus(http.StatusCreated). // ← ③
    CmpJSONBody(td.JSON(`
// Note that comments are allowed
{
  "id":         $id,          // set by the API/DB
  "name":       "Bob",
  "age":        42,
  "created_at": "$createdAt", // set by the API/DB
}`,
      td.Tag("id", td.Catch(&id, td.NotZero())),        // ← ④
      td.Tag("created_at", td.All(                      // ← ⑤
        td.HasSuffix("Z"),                              // ← ⑥
        td.Smuggle(func(s string) (time.Time, error) {  // ← ⑦
          return time.Parse(time.RFC3339Nano, s)
        }, td.Catch(&createdAt, td.Between(testAPI.SentAt(), time.Now()))), // ← ⑧
      )),
    ))
  if !testAPI.Failed() {
    t.Logf("The new Person ID is %d and was created at %s", id, createdAt)
  }
}
```

1. the API handler ready to be tested;
1. the POST request with automatic JSON marshalling;
1. the expected response HTTP status should be `http.StatusCreated`
   and the line just below, the body should match the
   [`JSON`]({{< ref "JSON" >}}) operator;
1. for the `$id` placeholder, [`Catch`]({{< ref "Catch" >}}) its
   value: put it in `id` variable and check it is
   [`NotZero`]({{< ref "NotZero" >}});
1. for the `$created_at` placeholder, use the [`All`]({{< ref "All" >}})
   operator. It combines several operators like a AND;
1. check that `$created_at` date ends with "Z" using
   [`HasSuffix`]({{< ref "HasSuffix" >}}). As we expect a RFC3339
   date, we require it in UTC time zone;
1. convert `$created_at` date into a `time.Time` using a custom
   function thanks to the [`Smuggle`]({{< ref "Smuggle" >}}) operator;
1. then [`Catch`]({{< ref "Catch" >}}) the resulting value: put it in
   `createdAt` variable and check it is greater or equal than
   `testAPI.SentAt()` (the time just before the request is handled) and lesser
   or equal than `time.Now()`.


Example of produced error in case of mismatch:

![error output](/images/colored-output.svg)


## Description

go-testdeep is a go rewrite and adaptation of wonderful
[Test::Deep perl](https://metacpan.org/pod/Test::Deep).

In golang, comparing data structure is usually done using
[reflect.DeepEqual](https://golang.org/pkg/reflect/#DeepEqual) or
using a package that uses this function behind the scene.

This function works very well, but it is not flexible. Both compared
structures must match exactly and when a difference is returned, it is
up to the caller to display it. Not easy when comparing big data
structures.

The purpose of go-testdeep, via
[`td` package](https://pkg.go.dev/github.com/maxatome/go-testdeep/td)
and its
[helpers](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers),
is to do its best to introduce this missing flexibility using
["operators"]({{< ref "operators" >}}), when the expected value (or
one of its component) cannot be matched exactly, mixed with some
useful [comparison functions]({{< ref "functions" >}}).
