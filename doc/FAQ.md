# FAQ

- [Table of contents for all these functions/methods?](#table-of-contents-for-all-these-functionsmethods)
- [How to mix strict requirements and simple assertions?](#how-to-mix-strict-requirements-and-simple-assertions)
- [How to test `io.Reader` contents, like `http.Response.Body` for example?](#how-to-test-ioreader-contents-like-httpresponsebody-for-example)
- [OK, but I prefer comparing `string`s instead of `byte`s](#ok-but-I-prefer-comparing-strings-instead-of-bytes)
- [OK, but my response is in fact a JSON marshaled struct of my own](#ok-but-my-response-is-in-fact-a-json-marshaled-struct-of-my-own)
- [OK, but you are funny, this response sends a new created object, so I don't know the ID in advance!](#ok-but-you-are-funny-this-response-sends-a-new-created-object-so-i-dont-know-the-id-in-advance)
- [What about testing the response using my API?](#what-about-testing-the-response-using-my-api)
- [Arf, I use Gin Gonic, and so no `net/http` handlers](#arf-i-use-gin-gonic-and-so-no-nethttp-handlers)
- [How to add a new operator?](#how-to-add-a-new-operator)


## Table of contents for all these functions/methods?

Of course! See the [Godoc table of contents](toc.md#godoc-table-of-contents).


## How to mix strict requirements and simple assertions?

```golang
import (
  "testing"

  td "github.com/maxatome/go-testdeep"
)

func TestAssertionsAndRequirements(t *testing.T) {
  assert := td.NewT(t)
  require := assert.FailureIsFatal()

  got := SomeFunction()

  require.Cmp(got, expected) // if it fails: report error + abort
  assert.Cmp(got, expected)  // if it fails: report error + continue
}
```

## How to test `io.Reader` contents, like http.Response.Body for example?

The [`Smuggle`](https://godoc.org/github.com/maxatome/go-testdeep#Smuggle)
operator is done for that, here with the help of
[`ReadAll`](https://golang.org/pkg/io/ioutil/#ReadAll).

```golang
import (
  "net/http"
  "testing"
  "io/ioutil"

  td "github.com/maxatome/go-testdeep"
)

func TestResponseBody(t *testing.T) {
  // Expect this response sends "Expected Response!"
  var resp *http.Response = GetResponse()

  td.Cmp(t, resp.Body,
    td.Smuggle(ioutil.ReadAll, []byte("Expected Response!")))
}
```

## OK, but I prefer comparing `string`s instead of `byte`s

No problem, [`ReadAll`](https://golang.org/pkg/io/ioutil/#ReadAll) the
body by yourself and cast returned `[]byte` contents to `string`,
still using
[`Smuggle`](https://godoc.org/github.com/maxatome/go-testdeep#Smuggle)
operator:

```golang
import (
  "io"
  "io/ioutil"
  "net/http"
  "testing"

  td "github.com/maxatome/go-testdeep"
)

func TestResponseBody(t *testing.T) {
  // Expect this response sends "Expected Response!"
  var resp *http.Response = GetResponse()

  td.Cmp(t, body, td.Smuggle(
    func(body io.Reader) (string, error) {
      b, err := ioutil.ReadAll(body)
      return string(b), err
    },
    "Expected Response!"))
}

```

## OK, but my response is in fact a JSON marshaled struct of my own

No problem, [JSON unmarshal](https://golang.org/pkg/encoding/json/#Unmarshal)
it just after reading the body:

```golang
import (
  "encoding/json"
  "io"
  "io/ioutil"
  "net/http"
  "testing"

  td "github.com/maxatome/go-testdeep"
)

func TestResponseBody(t *testing.T) {
  // Expect this response sends `{"ID":42,"Name":"Bob","Age":28}`
  var resp *http.Response = GetResponse()

  type Person struct {
    ID   uint64
    Name string
    Age  int
  }

  td.Cmp(t, body, td.Smuggle(
    func(body io.Reader) (&Person, error) {
      b, err := ioutil.ReadAll(body)
      if err != nil {
        return nil, err
      }
      var s Person
      return &s, json.Unmarshal(b, &s)
    },
    &Person{
      ID:   42,
      Name: "Bob",
      Age:  28,
    }))
}
```

## OK, but you are funny, this response sends a new created object, so I don't know the ID in advance!

No problem, use
[`Struct`](https://godoc.org/github.com/maxatome/go-testdeep#Struct)
operator to test that ID field is non-zero:

```golang
import (
  "encoding/json"
  "io"
  "io/ioutil"
  "net/http"
  "testing"

  td "github.com/maxatome/go-testdeep"
)

func TestResponseBody(t *testing.T) {
  // Expect this response sends `{"ID":42,"Name":"Bob","Age":28}`
  var resp *http.Response = GetResponse()

  type Person struct {
    ID   uint64
    Name string
    Age  int
  }

  td.Cmp(t, body, td.Smuggle(
    func(body io.Reader) (*Person, error) {
      b, err := ioutil.ReadAll(body)
      if err != nil {
        return nil, err
      }
      var s Person
      return &s, json.Unmarshal(b, &s)
    },
    td.Struct(&Person{
      Name: "Bob",
      Age:  28,
    }, td.StructFields{
      "ID": td.NotZero(),
    })))
}
```

## What about testing the response using my API?

[`tdhttp` helper](https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp)
is done for that!

```golang
import (
  "encoding/json"
  "net/http"
  "testing"

  td "github.com/maxatome/go-testdeep"
  "github.com/maxatome/go-testdeep/helpers/tdhttp"
)

type Person struct {
  ID   uint64
  Name string
  Age  int
}

// MyApi defines our API.
func MyAPI() *http.ServeMux {
  mux := http.NewServeMux()

  // GET /json
  mux.HandleFunc("/json", func(w http.ResponseWriter, req *http.Request) {
    if req.Method != "GET" {
      http.NotFound(w, req)
      return
    }

    b, err := json.Marshal(Person{
      ID:   42,
      Name: "Bob",
      Age:  28,
    })
    if err != nil {
      http.Error(w, "Internal server error", http.StatusInternalServerError)
      return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(b)
  })

  return mux
}

func TestMyApi(t *testing.T) {
  myAPI := MyAPI()

  tdhttp.CmpJSONResponse(t,
    tdhttp.NewRequest("GET", "/json", nil),
    myAPI.ServeHTTP,
    tdhttp.Response{
      Status: http.StatusOK,
      // Header can be tested too… See tdhttp doc.
      Body: td.Struct(&Person{
        Name: "Bob",
        Age:  28,
      }, td.StructFields{
        "ID": td.NotZero(),
      }),
    },
    "Testing GET /json")
}
```

## Arf, I use Gin Gonic, and so no `net/http` handlers

It is exactly the same as for `net/http` handlers as [`*gin.Engine`
implements `http.Handler`
interface](https://godoc.org/github.com/gin-gonic/gin#Engine.ServeHTTP)!
So keep using
[`tdhttp` helper](https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp):

```go
import (
  "net/http"
  "testing"

  "github.com/gin-gonic/gin"

  td "github.com/maxatome/go-testdeep"
  "github.com/maxatome/go-testdeep/helpers/tdhttp"
)

type Person struct {
  ID   uint64
  Name string
  Age  int
}

// MyGinGonicApi defines our API.
func MyGinGonicAPI() *gin.Engine {
  router := gin.Default() // or gin.New() or receive the router by param it doesn't matter

  router.GET("/json", func(c *gin.Context) {
    c.JSON(http.StatusOK, Person{
      ID:   42,
      Name: "Bob",
      Age:  28,
    })
  })

  return router
}

func TestMyGinGonicApi(t *testing.T) {
  myAPI := MyGinGonicAPI()

  tdhttp.CmpJSONResponse(t,
    tdhttp.NewRequest("GET", "/json", nil),
    myAPI.ServeHTTP,
    tdhttp.Response{
      Status: http.StatusOK,
      // Header can be tested too… See tdhttp doc.
      Body: td.Struct(&Person{
        Name: "Bob",
        Age:  28,
      }, td.StructFields{
        "ID": td.NotZero(),
      }),
    },
    "Testing GET /json")
}
```

## How to add a new operator?

You want to add a new `FooBar` operator.

- [ ] check that another operator does not exist with the same meaning;
- [ ] add the operator definition in `td_foo_bar.go` file and fully
  document its usage;
- [ ] add operator tests in `td_foo_bar_test.go` file;
- [ ] in `example_test.go` file, add examples function(s) `ExampleFooBar*`
  in alphabetical order;
- [ ] automatically generate `CmpFooBar` & `T.FooBar` (+ examples) code:
  `./tools/gen_funcs.pl .`
- [ ] do not forget to run tests: `go test ./...`
- [ ] run `golangci-lint` as in [`.travis.yml`](../.travis.yml);
- [ ] in [`README.md`](../README.md), add this new `FooBar` operator:
  - in [Available operators](../README.md#available-operators) with a
    small description, respecting the alphabetical order;
  - in the [Operators vs go types](../README.md#operators-vs-go-types)
    matrix, still respecting the alphabetical order.
- [ ] in [`toc.md#godoc-table-of-contents`](toc.md#godoc-table-of-contents),
  add this new new `FooBar` operator:
  - in [Main shortcut functions](toc.md#main-shortcut-functions);
  - in [Shortcut methods of `*testdeep.T`](toc.md#shortcut-methods-of-testdeep-t);
  - in [`Testdeep` operators](toc.md#testdeep-operators), a simple copy
    of the line inserted in [Available operators](../README.md#available-operators)
	and its corresponding link of course.

Each time you change `example_test.go`, re-run `./tools/gen_funcs.pl .`
to update corresponding `CmpFooBar` & `T.FooBar` examples.

Test coverage must be 100%.
