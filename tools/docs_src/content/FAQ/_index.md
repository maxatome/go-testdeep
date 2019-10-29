+++
title = "FAQ"
date = 2019-10-08T21:28:21+02:00
weight = 30
+++

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

The [`Smuggle`]({{< ref "Smuggle" >}}) operator is done for that,
here with the help of [`ReadAll`](https://golang.org/pkg/io/ioutil/#ReadAll).

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
still using [`Smuggle`]({{< ref "Smuggle" >}}) operator:

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

No problem, use [`Struct`]({{< ref "Struct" >}}) operator to test
that ID field is non-zero:

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


## go-testdeep dumps only 10 errors, how to have more (or less)?

Using the environment variable `TESTDEEP_MAX_ERRORS`.

`TESTDEEP_MAX_ERRORS` contains the maximum number of errors to report
before stopping during one comparison (one
[`Cmp`](https://godoc.org/github.com/maxatome/go-testdeep#Cmp)
execution for example). It defaults to `10`.

Example:
```shell
TESTDEEP_MAX_ERRORS=30 go test
```

Setting it to `-1` means no limit:
```shell
TESTDEEP_MAX_ERRORS=-1 go test
```


## How do I change these crappy colors?

Using some environment variables:

- `TESTDEEP_COLOR` enable (`on`) or disable (`off`) the color
  output. It defaults to `on`;
- `TESTDEEP_COLOR_TEST_NAME` color of the test name. See below
  for color format, it defaults to `yellow`;
- `TESTDEEP_COLOR_TITLE` color of the test failure title. See below
  for color format, it defaults to `cyan`;
- `TESTDEEP_COLOR_OK` color of the test expected value. See below
  for color format, it defaults to `green`;
- `TESTDEEP_COLOR_BAD` color of the test got value. See below
  for color format, it defaults to `red`;

### Color format

A color in `TESTDEEP_COLOR_*` environment variables has the following
format:

```
foreground_color                    # set foreground color, background one untouched
foreground_color:background_color   # set foreground AND background color
:background_color                   # set background color, foreground one untouched
```

`foreground_color` and `background_color` can be:

- `black`
- `red`
- `green`
- `yellow`
- `blue`
- `magenta`
- `cyan`
- `white`
- `gray`

For example:

```shell
TESTDEEP_COLOR_OK=black:green \
    TESTDEEP_COLOR_BAD=white:red \
    TESTDEEP_COLOR_TITLE=yellow \
    go test
```


## How to add a new operator?

You want to add a new `FooBar` operator.

- [ ] check that another operator does not exist with the same meaning;
- [ ] add the operator definition in `td_foo_bar.go` file and fully
  document its usage:
  - add a `// summary(FooBar): small description` line, before
    operator comment,
  - add a `// input(FooBar): …` line, just aftezr `summary(FooBar)`
    line. This one lists all inputs accepted by the operator;
- [ ] add operator tests in `td_foo_bar_test.go` file;
- [ ] in `example_test.go` file, add examples function(s) `ExampleFooBar*`
  in alphabetical order;
- [ ] automatically generate `CmpFooBar` & `T.FooBar` (+ examples) code:
  `./tools/gen_funcs.pl .`
- [ ] do not forget to run tests: `go test ./...`
- [ ] run `golangci-lint` as in [`.travis.yml`](https://github.com/maxatome/go-testdeep/blob/master/.travis.yml);

Each time you change `example_test.go`, re-run `./tools/gen_funcs.pl .`
to update corresponding `CmpFooBar` & `T.FooBar` examples.

Test coverage must be 100%.
