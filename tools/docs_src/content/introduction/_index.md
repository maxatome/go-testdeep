+++
title = "Introduction"
weight = 5
+++

## Synopsis

Make golang tests easy, from simplest usage:

```go
import (
  "testing"

  td "github.com/maxatome/go-testdeep"
)

func TestMyFunc(t *testing.T) {
  td.Cmp(t, MyFunc(), &Info{Name: "Alice", Age: 42})
}
```

To most complex one, allowing to easily test golang API routes, using
flexible [operators]({{< ref "operators" >}}):

```go
import (
  "testing"
  "time"

  td "github.com/maxatome/go-testdeep"
  "github.com/maxatome/go-testdeep/helpers/tdhttp"
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

  beforeCreate := time.Now().Truncate(0)

  tdhttp.CmpJSONResponse(td.NewT(t).FailureIsFatal(), // ← t.Fatal() if test fails
    tdhttp.PostJSON("/person", Person{Name: "Bob", Age: 42}), // ← the request
    myAPI.ServeHTTP, // ← the API handler
    tdhttp.Response{ // ← the expected response
      Status: http.StatusCreated,
      // Header can be tested too… See tdhttp doc.
      Body: td.JSON(`
{
  "id":         $id,
  "name":       "Bob",
  "age":        42,
  "created_at": "$createdAt",
}`,
        td.Tag("id", td.Catch(&id, td.NotZero())), // catch $id and check ≠ 0
        td.Tag("created_at", td.All( // ← All combines several operators like a AND
          td.HasSuffix("Z"), // check the RFC3339 $created_at date ends with "Z"
          td.Smuggle(func(s string) (time.Time, error) { // convert to time.Time
            return time.Parse(time.RFC3339Nano, s)
          }, td.Catch(&createdAt, td.Gte(beforeCreate))), // catch it and check ≥ beforeCreate
        )),
      ),
    },
    "Create a new Person")

  t.Logf("The new Person ID is %d and was created at %s", id, createdAt)
}
```

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

The purpose of testdeep package is to do its best to introduce this
missing flexibility using ["operators"]({{< ref "operators" >}}), when
the expected value (or one of its component) cannot be matched
exactly, mixed with some useful
[comparison functions]({{< ref "functions" >}}).
