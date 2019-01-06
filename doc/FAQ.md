# FAQ

## How to mix strict requirements and simple assertions?

```golang
import (
  "testing"

  td "github.com/maxatome/go-testdeep"
)

func TestAssertionsAndRequirements(t *testing.T) {
  assert := td.NewT(t)
  require := assert.FatalOnError()

  got := SomeFunction()

  require.CmpDeeply(got, expected) // if it fails: report error + abort
  assert.CmpDeeply(got, expected)  // if it fails: report error + continue
}
```

## How to test io.Reader contents, like http.Response.Body for example?

The [`Smuggle`](https://godoc.org/github.com/maxatome/go-testdeep#Smuggle)
is done for that, here with the help of
[`ReadAll`](https://golang.org/pkg/io/ioutil/#ReadAll).

```golang
import (
  "io"
  "net/http"
  "testing"
  "io/ioutil"

  td "github.com/maxatome/go-testdeep"
)

func TestResponseBody(t *testing.T) {
  // Expect this response sends "Expected Response!"
  var resp *http.Response = GetResponse()

  td.CmpDeeply(t, resp.Body,
    td.Smuggle(ioutil.ReadAll, []byte("Expected Response!")))
}
```

## OK, but I prefer comparing strings instead of bytes

No problem, [`ReadAll`](https://golang.org/pkg/io/ioutil/#ReadAll) the
body by yourself and cast returned `[]byte` contents to `string`:

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

  td.CmpDeeply(t, body, td.Smuggle(
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

  type MyOwnStruct struct {
    ID   uint64
    Name string
    Age  int
  }

  td.CmpDeeply(t, body, td.Smuggle(
    func(body io.Reader) (&MyOwnStruct, error) {
      b, err := ioutil.ReadAll(body)
      if err != nil {
        return nil, err
      }
      var s MyOwnStruct
      return &s, json.Unmarshal(b, &s)
    },
    &MyOwnStruct{
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

  type MyOwnStruct struct {
    ID   uint64
    Name string
    Age  int
  }

  td.CmpDeeply(t, body, td.Smuggle(
    func(body io.Reader) (*MyOwnStruct, error) {
      b, err := ioutil.ReadAll(body)
      if err != nil {
        return nil, err
      }
      var s MyOwnStruct
      return &s, json.Unmarshal(b, &s)
    },
    td.Struct(&MyOwnStruct{
      Name: "Bob",
      Age:  28,
    }, td.StructFields{
      "ID": td.NotZero(),
    })))
}
```

## How to add a new operator?

You want to add a new `FooBar` operator.

- check that another operator does not exist with the same meaning;
- add the operator definition in `td_foo_bar.go` file and fully
  document its usage;
- add operator tests in `td_foo_bar_test.go` file;
- in `example_test.go` file, add examples function(s) `ExampleFooBar*`
  in alphabetical order;
- automatically generate `CmpFooBar` & `T.FooBar` (+ examples) code:
  `./tools/gen_funcs.pl .`
- do not forget to run tests: `go test ./...`
- run `gometalinter` as in [`.travis.yml`](.travis.yml);
- add `FooBar` with a small description in
  [`README.md`](README.md#available-operators), respecting the
  alphabetical order.

Each time you change `example_test.go`, re-run `./tools/gen_funcs.pl .`
to update corresponding `CmpFooBar` & `T.FooBar` examples.

Test coverage must be 100%.
