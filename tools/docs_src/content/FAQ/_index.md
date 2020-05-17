+++
title = "FAQ"
date = 2019-10-08T21:28:21+02:00
weight = 30
+++

## How to mix strict requirements and simple assertions?

```golang
import (
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestAssertionsAndRequirements(t *testing.T) {
  assert, require := td.AssertRequire(t)

  got := SomeFunction()

  require.Cmp(got, expected) // if it fails: report error + abort
  assert.Cmp(got, expected)  // if it fails: report error + continue
}
```


## Why `nil` is handled so specifically?

```golang
var pn *int
td.Cmp(t, pn, nil)
```

fails with the error ![error output](/images/faq-nil-colored-output.svg)

And, yes, it is normal. (TL;DR use
[`CmpNil`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNil)
instead, safer, or use
[`CmpLax`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpLax),
but be careful of edge cases.)

To understand why, look at the following examples:

```golang
var err error
td.Cmp(t, err, nil)
```

works (and you want it works), but

```golang
var err error = (*MyError)(nil)
td.Cmp(t, err, nil)
```

fails with the error
![error output](/images/faq-error-nil-colored-output.svg)

and in most cases you want it fails, because `err` is not nil! The
pointer stored in the interface is nil, but not the interface itself.

As [`Cmp`] `got` parameter type is `interface{}`, when you pass an
interface variable in it (whatever the interface is), [`Cmp`] always
receives an `interface{}`. So here, [`Cmp`] receives `(*MyError)(nil)`
in the `got` interface, and not `error((*MyError)(nil))` ⇒ the `error`
interface information is lost at the compilation time.

In other words, [`Cmp`] has no abilities to tell the difference
between `error((*MyError)(nil))` and `(*MyError)(nil)` when passed in
`got` parameter.

That is why [`Cmp`] is strict by default, and requires that nil be
strongly typed, to be able to detect when a non-nil interface contains
a nil pointer.

So to recap:
```golang
var pn *int
td.Cmp(t, pn, nil)               // fails as nil is not strongly typed
td.Cmp(t, pn, (*int)(nil))       // succeeds
td.Cmp(t, pn, td.Nil())          // succeeds
td.CmpNil(t, pn)                 // succeeds
td.Cmp(t, pn, td.Lax(nil))       // succeeds
td.CmpLax(t, pn, nil)            // succeeds

var err error
td.Cmp(t, err, nil)              // succeeds
td.Cmp(t, err, (*MyError)(nil))  // fails as err does not contain any value
td.Cmp(t, err, td.Nil())         // succeeds
td.CmpNil(t, err)                // succeeds
td.Cmp(t, err, td.Lax(nil))      // succeeds
td.CmpLax(t, err, nil)           // succeeds
td.CmpError(t, err)              // fails as err is nil
td.CmpNoError(t, err)            // succeeds

err = (*MyError)(nil)
td.Cmp(t, err, nil)              // fails as err contains a value
td.Cmp(t, err, (*MyError)(nil))  // succeeds
td.Cmp(t, err, td.Nil())         // succeeds
td.CmpNil(t, err)                // succeeds
td.Cmp(t, err, td.Lax(nil))      // succeeds *** /!\ be careful here! ***
td.CmpLax(t, err, nil)           // succeeds *** /!\ be careful here! ***
td.CmpError(t, err)              // succeeds
td.CmpNoError(t, err)            // fails as err contains a value
```

Morality:
- to compare a pointer against nil, use
  [`CmpNil`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNil)
  or strongly type nil (e.g. `(*int)(nil)`) in expected parameter of [`Cmp`];
- to compare an error against nil, use
  [`CmpNoError`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNoError)
  or nil direcly in expected parameter of [`Cmp`].


## How does operator anchoring work?

Take this struct, returned by a `GetPerson()` function:

```golang
type Person struct {
  ID   int64
  Name string
  Age  uint8
}
```

For the `Person` returned by `GetPerson()`, we expect that:

- `ID` field should be ≠ 0;
- `Name` field should always be "Bob";
- `Age` field should be ≥ 40 and ≤ 45.

Without operator anchoring:

```golang
func TestPerson(tt *testing.T) {
  t := td.NewT(tt)

  t.Cmp(GetPerson(),               // ← ①
    td.Struct(Person{Name: "Bob"}, // ← ②
      td.StructFields{             // ← ③
        "ID":  td.NotZero(),       // ← ④
        "Age": td.Between(uint8(40), uint8(45)), // ← ⑤
      }))
}
```

1. `GetPerson()` returns a `Person`;
2. as some fields of the returned `Person` are not exactly known in
   advance, we use the [`Struct`]({{< ref "Struct" >}}) operator as
   expected parameter. It allows to match exactly some fields, and use
   [TestDeep operators]({{< ref "operators" >}}) on others. Here we
   know that `Name` field should always be "Bob";
3. [`StructFields`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#StructFields)
   is a map allowing to use [TestDeep operators]({{< ref "operators" >}})
   for any field;
4. `ID` field should be ≠ 0. See [`NotZero`]({{< ref "NotZero" >}})
   operator for details;
5. `Age` field should be ≥ 40 and ≤ 45. See [`Between`]({{< ref "between" >}})
   operator for details.

With operator anchoring, the use of [`Struct`]({{< ref "Struct" >}})
operator is no longer needed:

```golang
func TestPerson(tt *testing.T) {
  t := td.NewT(tt)

  t.Cmp(GetPerson(), // ← ①
    Person{          // ← ②
      Name: "Bob",   // ← ③
      ID:   t.A(td.NotZero(), int64(0)).(int64), // ← ④
      Age:  t.A(td.Between(uint8(40), uint8(45))).(uint8), // ← ⑤
    })
}
```

1. `GetPerson()` still returns a `Person`;
2. expected parameter is directly a `Person`. No operator needed here;
3. `Name` field should always be "Bob", no change here;
4. `ID` field should be ≠ 0: anchor the [`NotZero`]({{< ref "NotZero" >}})
   operator using the [`A`] method. Break this line down:
   ```golang
   t.A(            // ← ①
     td.NotZero(), // ← ②
     int64(0),     // ← ③
   ).(int64)       // ← ④
   ```
   1. the [`A`] method is the key of the anchoring system. It saves
      the operator globally, so it can be retrieved during the
      comparison of the next [`Cmp`] call,
   2. the operator we want to anchor,
   3. this optional parameter is needed to tell [`A`] that the returned
      value must be a `int64`. Sometimes, this type can be deduced
      from the operator, but as [`NotZero`]({{< ref "NotZero" >}}) can
      handle any kind of number, it is not the case here. So we have
      to pass it,
   4. as [`A`] method returns an `interface{}`, we need to assert the
      `int64` type to bypass the golang static typing system,
5. `Age `field should be ≥ 40 and ≤ 45: anchor the
   [`Between`]({{< ref "between" >}}) operator using the [`A`]
   method. Break this line down:
   ```golang
   t.A(                                // ← ①
     td.Between(uint8(40), uint8(45)), // ← ②
   ).(uint8)                           // ← ③
   ```
   1. the [`A`] method saves the operator globally, so it can be
      retrieved during the comparison of the next [`Cmp`] call,
   2. the operator we want to anchor. As [`Between`]({{< ref "between" >}})
      knows the type of its operands (here `uint8`), there is no need
      to tell [`A`] the returned type must be `uint8`. It can be deduced
      from [`Between`]({{< ref "between" >}}),
   3. as [`A`] method returns an `interface{}`, we need to assert the
      `uint8` type to bypass the golang static typing system.

Note the [`A`] method is a shortcut of [`Anchor`] method.

Some rules have to be kept in mind:
- never cast a value returned by [`A`] or [`Anchor`] methods:
  ```golang
  t := td.NewT(tt) // tt is a *testing.T
  t.A(td.NotZero(), uint(8)).(uint8)         // OK
  uint16(t.A(td.NotZero(), uint(8)).(uint8)) // Not OK!
  t.A(td.NotZero(), uint16(0)).(uint16)      // OK
  ```
- anchored operators disappear once the next [`Cmp`] call done. To
  share them between [`Cmp`] calls, use the [`SetAnchorsPersist`]
  method as in:
  ```golang
  t := td.NewT(tt) // tt is a *testing.T
  age := t.A(td.Between(uint8(40), uint8(45))).(uint8)

  t.SetAnchorsPersist(true) // ← Don't reset anchors after next Cmp() call

  t.Cmp(GetPerson(1), Person{
    Name: "Bob",
    Age:  age,
  })

  t.Cmp(GetPerson(2), Person{
    Name: "Bob",
    Age:  age, // ← OK
  })
  ```


## How to test `io.Reader` contents, like `net/http.Response.Body` for example?

The [`Smuggle`]({{< ref "Smuggle" >}}) operator is done for that,
here with the help of [`ReadAll`](https://golang.org/pkg/io/ioutil/#ReadAll).

```golang
import (
  "net/http"
  "testing"
  "io/ioutil"

  "github.com/maxatome/go-testdeep/td"
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

  "github.com/maxatome/go-testdeep/td"
)

func TestResponseBody(t *testing.T) {
  // Expect this response sends "Expected Response!"
  var resp *http.Response = GetResponse()

  td.Cmp(t, resp.Body, td.Smuggle( // ← transform a io.Reader to a string
    func(body io.Reader) (string, error) {
      b, err := ioutil.ReadAll(body)
      return string(b), err
    },
    "Expected Response!"))
}

```


## OK, but my response is in fact a JSON marshaled struct of my own

No problem, [JSON decode](https://golang.org/pkg/encoding/json/#Decoder)
while reading the body:

```golang
import (
  "encoding/json"
  "io"
  "net/http"
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestResponseBody(t *testing.T) {
  // Expect this response sends `{"ID":42,"Name":"Bob","Age":28}`
  var resp *http.Response = GetResponse()

  type Person struct {
    ID   uint64
    Name string
    Age  int
  }

  td.Cmp(t, resp.Body, td.Smuggle( // ← transform a io.Reader in *Person
    func(body io.Reader) (*Person, error) {
      var s Person
      return &s, json.NewDecoder(body).Decode(&s)
    },
    &Person{ // ← check Person content
      ID:   42,
      Name: "Bob",
      Age:  28,
    }))
}
```


## OK, but you are funny, this response sends a new created object, so I don't know the ID in advance!

No problem, use [`Struct`]({{< ref "Struct" >}}) operator to test
that `ID` field is non-zero (as a bonus, add a `CreatedAt` field):

```golang
import (
  "encoding/json"
  "io"
  "net/http"
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestResponseBody(t *testing.T) {
  // Expect this response sends:
  //   `{"ID":42,"Name":"Bob","Age":28,"CreatedAt":"2019-01-02T11:22:33Z"}`
  var resp *http.Response = GetResponse()

  type Person struct {
    ID        uint64
    Name      string
    Age       int
    CreatedAt time.Time
  }

  y2019, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")

  td.Cmp(t, resp.Body, td.Smuggle( // ← transform a io.Reader in *Person
    func(body io.Reader) (*Person, error) {
      var s Person
      return &s, json.NewDecoder(body).Decode(&s)
    },
    td.Struct(&Person{ // ← check Person content
      Name: "Bob",
      Age:  28,
    }, td.StructFields{
      "ID":        td.NotZero(),  // check ID ≠ 0
      "CreatedAt": td.Gte(y2019), // check CreatedAt ≥ 2019/01/01
    })))

  tt := td.newT(t)
  tt.Cmp(resp.Body, td.Smuggle( // ← transform a io.Reader in *Person
    func(body io.Reader) (*Person, error) {
      var s Person
      return &s, json.NewDecoder(body).Decode(&s)
    },
    &Person{ // ← check Person content
      Name:      "Bob",
      Age:       28,
      ID:        tt.A(td.NotZero(), uint64(0)).(uint64), // check ID ≠ 0
      CreatedAt: tt.A(td.Gte(y2019)).(time.Time),        // check CreatedAt ≥ 2019/01/01
    }))
}
```


## What about testing the response using my API?

[`tdhttp` helper](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp)
is done for that!

```golang
import (
  "encoding/json"
  "net/http"
  "testing"
  "time"

  "github.com/maxatome/go-testdeep/helpers/tdhttp"
  "github.com/maxatome/go-testdeep/td"
)

type Person struct {
  ID        uint64
  Name      string
  Age       int
  CreatedAt time.Time
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
      ID:        42,
      Name:      "Bob",
      Age:       28,
      CreatedAt: time.Now().UTC(),
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

  y2019, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")

  testAPI := tdhttp.NewTestAPI(t, myAPI) // ← ①

  testAPI.Get("/json").             // ← ②
    Name("Testing GET /json").
    CmpStatus(http.StatusOK).       // ← ③
    CmpJSONBody(td.SStruct(&Person{ // ← ④
      Name: "Bob",
      Age:  28,
    }, td.StructFields{
      "ID":        td.NotZero(),    // ← ⑤
      "CreatedAt": td.Gte(y2019),   // ← ⑥
    }))

  // testAPI can be used to test another route…
}
```

1. the API handler ready to be tested;
1. the GET request;
1. the expected HTTP status should be `http.StatusOK`;
1. the expected body should match the [`SStruct`]({{< ref "SStruct" >}})
   operator;
1. check the `ID` field is [`NotZero`]({{< ref "NotZero" >}});
1. check the `CreatedAt` field is greater or equal than `y2019` variable
   (set just before `tdhttp.NewTestAPI` call).

If you prefer to do one function call instead of chaining methods as
above, you can try [CmpJSONResponse](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp#CmpJSONResponse).


## Arf, I use Gin Gonic, and so no `net/http` handlers

It is exactly the same as for `net/http` handlers as [`*gin.Engine`
implements `http.Handler`
interface](https://pkg.go.dev/github.com/gin-gonic/gin#Engine.ServeHTTP)!
So keep using
[`tdhttp` helper](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp):

```go
import (
  "net/http"
  "testing"
  "time"

  "github.com/gin-gonic/gin"

  "github.com/maxatome/go-testdeep/helpers/tdhttp"
  "github.com/maxatome/go-testdeep/td"
)

type Person struct {
  ID        uint64
  Name      string
  Age       int
  CreatedAt time.Time
}

// MyGinGonicApi defines our API.
func MyGinGonicAPI() *gin.Engine {
  router := gin.Default() // or gin.New() or receive the router by param it doesn't matter

  router.GET("/json", func(c *gin.Context) {
    c.JSON(http.StatusOK, Person{
      ID:        42,
      Name:      "Bob",
      Age:       28,
      CreatedAt: time.Now().UTC(),
    })
  })

  return router
}

func TestMyGinGonicApi(t *testing.T) {
  myAPI := MyGinGonicAPI()

  y2019, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")

  testAPI := tdhttp.NewTestAPI(t, myAPI) // ← ①

  testAPI.Get("/json").             // ← ②
    Name("Testing GET /json").
    CmpStatus(http.StatusOK).       // ← ③
    CmpJSONBody(td.SStruct(&Person{ // ← ④
      Name: "Bob",
      Age:  28,
    }, td.StructFields{
      "ID":        td.NotZero(),    // ← ⑤
      "CreatedAt": td.Gte(y2019),   // ← ⑥
    }))

  // testAPI can be used to test another route…
}
```

1. the API handler ready to be tested;
1. the GET request;
1. the expected HTTP status should be `http.StatusOK`;
1. the expected body should match the [`SStruct`]({{< ref "SStruct" >}})
   operator;
1. check the `ID` field is [`NotZero`]({{< ref "NotZero" >}});
1. check the `CreatedAt` field is greater or equal than `y2019` variable
   (set just before `tdhttp.NewTestAPI` call).

If you prefer to do one function call instead of chaining methods as
above, you can try [CmpJSONResponse](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp#CmpJSONResponse).


## Fine, the request succeeds and the ID is not 0, but what is the ID real value?

Stay with [`tdhttp` helper](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp)!

In fact you can [`Catch`]({{< ref "Catch" >}}) the `ID` before comparing
it to 0 (as well as `CreatedAt` in fact). Try:

```go
func TestMyGinGonicApi(t *testing.T) {
  myAPI := MyGinGonicAPI()

  var id uint64
  var createdAt time.Time

  y2019, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")

  testAPI := tdhttp.NewTestAPI(t, myAPI) // ← ①

  testAPI.Get("/json").             // ← ②
    Name("Testing GET /json").
    CmpStatus(http.StatusOK).       // ← ③
    CmpJSONBody(td.SStruct(&Person{ // ← ④
      Name: "Bob",
      Age:  28,
    }, td.StructFields{
      "ID":        td.Catch(&id, td.NotZero()),         // ← ⑤
      "CreatedAt": td.Catch(&createdAt, td.Gte(y2019)), // ← ⑥
    }))
  if !testAPI.Failed() {
    t.Logf("The ID is %d and was created at %s", id, createdAt)
  }

  // testAPI can be used to test another route…
}
```

1. the API handler ready to be tested;
1. the GET request;
1. the expected HTTP status should be `http.StatusOK`;
1. the expected body should match the [`SStruct`]({{< ref "SStruct" >}})
   operator;
1. [`Catch`]({{< ref "Catch" >}}) the `ID` field: put it in `id`
   variable and check it is [`NotZero`]({{< ref "NotZero" >}});
1. [`Catch`]({{< ref "Catch" >}}) the `CreatedAt` field: put it in `createdAt`
   variable and check it is greater or equal than `y2019` variable
   (set just before `tdhttp.NewTestAPI` call).

If you prefer to do one function call instead of chaining methods as
above, you can try [CmpJSONResponse](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp#CmpJSONResponse).


## OK, but how to be sure the response content is well JSONified?

Again, [`tdhttp` helper](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp)
is your friend!

With the help of [`JSON`]({{< ref "JSON" >}}) operator of course! See
it below, used with [`Catch`]({{< ref "Catch" >}}) (note it can be used
without), for a `POST` example:

```go
type Person struct {
  ID        uint64    `json:"id"`
  Name      string    `json:"name"`
  Age       int       `json:"age"`
  CreatedAt time.Time `json:"created_at"`
}

func TestMyGinGonicApi(t *testing.T) {
  myAPI := MyGinGonicAPI()

  var id uint64
  var createdAt time.Time

  testAPI := tdhttp.NewTestAPI(t, myAPI) // ← ①

  testAPI.PostJSON("/person", Person{Name: "Bob", Age: 42}), // ← ②
    Name("Create a new Person").
	CmpStatus(http.StatusCreated). // ← ③
	CmpJSONBody(td.JSON(`
{
  "id":         $id,
  "name":       "Bob",
  "age":        42,
  "created_at": "$createdAt",
}`,
      td.Tag("id", td.Catch(&id, td.NotZero())),        // ← ④
      td.Tag("created_at", td.All(                      // ← ⑤
        td.HasSuffix("Z"),                              // ← ⑥
        td.Smuggle(func(s string) (time.Time, error) {  // ← ⑦
          return time.Parse(time.RFC3339Nano, s)
        }, td.Catch(&createdAt, td.Gte(testAPI.SentAt()))), // ← ⑧
      )),
    ))
  if !testAPI.Failed() {
    t.Logf("The new Person ID is %d and was created at %s", id, createdAt)
  }

  // testAPI can be used to test another route…
}
```

1. the API handler ready to be tested;
1. the POST request with automatic JSON marshalling;
1. the expected HTTP status should be `http.StatusCreated`
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
   `testAPI.SentAt()` (the time just before the request is handled).

If you prefer to do one function call instead of chaining methods as
above, you can try [CmpJSONResponse](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp#CmpJSONResponse).


## My API use XML not JSON!

[`tdhttp`
helper](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp)
provides the same functions and methods for XML it does for JSON.

[RTFM](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp)
:)

Note that the [`JSON`]({{< ref "JSON" >}}) operator have not its `XML`
counterpart yet.
But [PRs are welcome](https://github.com/maxatome/go-testdeep/pulls)!


## Should I import `github.com/maxatome/go-testdeep` or `github.com/maxatome/go-testdeep/td`?

Historically the main package of go-testdeep was `testdeep` as in:

```go
import (
  "testing"

  "github.com/maxatome/go-testdeep"
)

func TestMyFunc(t *testing.T) {
  testdeep.Cmp(t, GetPerson(), Person{Name: "Bob", Age: 42})
}
```

As `testdeep` was boring to type, renaming it to `td` became a habit as in:

```go
import (
  "testing"

  td "github.com/maxatome/go-testdeep"
)

func TestMyFunc(t *testing.T) {
  td.Cmp(t, GetPerson(), Person{Name: "Bob", Age: 42})
}
```

Forcing the developer to systematically rename `testdeep` package to
`td` in all its tests is not very friendly. That is why a decision was
taken to create a new package `github.com/maxatome/go-testdeep/td`
while keeping `github.com/maxatome/go-testdeep` working thanks to go
type aliases.

So the previous examples (that are still working) can now be written as:

```go
import (
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestMyFunc(t *testing.T) {
  td.Cmp(t, GetPerson(), Person{Name: "Bob", Age: 42})
}
```

There is no package renaming anymore. Switching to import
`github.com/maxatome/go-testdeep/td` is advised for new code.


## What does the error `undefined: testdeep.DefaultContextConfig` mean?

Since release `v1.3.0`, this variable moved to the new
`github.com/maxatome/go-testdeep/td` package.

* If you rename the `testdeep` package to `td` as in:
  ```
  import td "github.com/maxatome/go-testdeep"
  …
    td.DefaultContextConfig = td.ContextConfig{…}
  ```
  then just change the import line to:
  ```
  import "github.com/maxatome/go-testdeep/td"
  ```

* Otherwise, you have two choices:
  1. either add a new import line:
     ```
     import "github.com/maxatome/go-testdeep/td"
     ```
     then use `td.DefaultContextConfig` instead of
     `testdeep.DefaultContextConfig`, and continue to use `testdeep`
     package elsewhere.
  2. or replace the import line:
     ```
     import "github.com/maxatome/go-testdeep"
     ```
     by
     ```
     import "github.com/maxatome/go-testdeep/td"
     ```
     then rename all occurrences of `testdeep` package to `td`.


## go-testdeep dumps only 10 errors, how to have more (or less)?

Using the environment variable `TESTDEEP_MAX_ERRORS`.

`TESTDEEP_MAX_ERRORS` contains the maximum number of errors to report
before stopping during one comparison (one [`Cmp`] execution for
example). It defaults to `10`.

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


## play.golang.org does not handle colors, error output is nasty

Just add this single line in playground:

```go
func init() { os.Setenv("TESTDEEP_COLOR", "off") }
```

Until playground supports [ANSI color escape
sequences](https://en.wikipedia.org/wiki/ANSI_escape_code#Colors).


## The `X` testing framework allows to test/do `Y` while go-testdeep not

The [`Code`]({{< ref "Code" >}}) and [`Smuggle`]({{< ref "Smuggle" >}})
operators should allow to cover all cases not handled by [other
operators]({{< ref "operators" >}}).

If you think this missing feature deserves a specific operator,
because it is frequently or widely used, file an issue and let's
discuss about it.

We plan to add a new `github.com/maxatome/go-testdeep/helpers/tdcombo`
helper package, bringing together all what we can call
combo-operators. Combo-operators are operators using any number of
[already existing operators]({{< ref "operators" >}}).

As an example of such combo-operators, the following one. It allows to
check that a string contains a RFC3339 formatted time, in UTC time
zone ("Z" suffix) and then to compare it as a `time.Time` against
`expectedValue` (which can be another [operator]({{< ref "operators" >}})
or, of course, a `time.Time` value).

```go
func RFC3339ZToTime(expectedValue interface{}) td.TestDeep {
  return td.All(
    td.HasSuffix("Z"),
    td.Smuggle(func(s string) (time.Time, error) {
      return time.Parse(time.RFC3339Nano, s)
    }, expectedValue),
  )
}
```

It could be used as:

```go
before := time.Now()
record := NewRecord()
td.Cmp(t, record,
  td.SuperJSONOf(`{"created_at": $1}`,
    tdcombo.RFC3339ZToTime(td.Between(before, time.Now()),
  )),
  "The JSONified record.created_at is UTC-RFC3339",
)
```


## How to add a new operator?

You want to add a new `FooBar` operator.

- [ ] check that another operator does not exist with the same meaning;
- [ ] add the operator definition in `td_foo_bar.go` file and fully
  document its usage:
  - add a `// summary(FooBar): small description` line, before
    operator comment,
  - add a `// input(FooBar): …` line, just after `summary(FooBar)`
    line. This one lists all inputs accepted by the operator;
- [ ] add operator tests in `td_foo_bar_test.go` file;
- [ ] in `example_test.go` file, add examples function(s) `ExampleFooBar*`
  in alphabetical order;
- [ ] automatically generate `CmpFooBar` & `T.FooBar` (+ examples) code:
  `./tools/gen_funcs.pl`
- [ ] do not forget to run tests: `go test ./...`
- [ ] run `golangci-lint` as in [`.travis.yml`](https://github.com/maxatome/go-testdeep/blob/master/.travis.yml);

Each time you change `example_test.go`, re-run `./tools/gen_funcs.pl`
to update corresponding `CmpFooBar` & `T.FooBar` examples.

Test coverage must be 100%.


[`Cmp`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp

[`A`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.A
[`Anchor`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Anchor
[`SetAnchorsPersist`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SetAnchorsPersist
