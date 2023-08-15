go-testdeep
===========

[![Build Status](https://github.com/maxatome/go-testdeep/workflows/Build/badge.svg?branch=master)](https://github.com/maxatome/go-testdeep/actions?query=workflow%3ABuild)
[![Coverage Status](https://coveralls.io/repos/github/maxatome/go-testdeep/badge.svg?branch=master)](https://coveralls.io/github/maxatome/go-testdeep?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxatome/go-testdeep)](https://goreportcard.com/report/github.com/maxatome/go-testdeep)
[![GoDoc](https://pkg.go.dev/badge/github.com/maxatome/go-testdeep)](https://pkg.go.dev/github.com/maxatome/go-testdeep/td)
[![Version](https://img.shields.io/github/tag/maxatome/go-testdeep.svg)](https://github.com/maxatome/go-testdeep/releases)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go/#testing)

![go-testdeep](https://github.com/maxatome/go-testdeep-site/raw/master/docs_src/static/images/logo.png)

**Extremely flexible golang deep comparison, extends the go testing package.**

Currently supports go 1.16 → 1.21.

- [Latest news](#latest-news)
- [Synopsis](#synopsis)
- [Description](#description)
- [Installation](#installation)
- [Functions](https://go-testdeep.zetta.rocks/functions/)
- [Available operators](https://go-testdeep.zetta.rocks/operators/)
- [Helpers](#helpers)
  - [`tdhttp` or HTTP API testing helper](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp)
  - [`tdsuite` or testing suite helper](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdsuite)
  - [`tdutil` aka the helper of helpers](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdutil)
- [See also](#see-also)
- [License](#license)
- [FAQ](https://go-testdeep.zetta.rocks/faq/)


## Latest news

- 2023/03/18: [v1.13.0 release](https://github.com/maxatome/go-testdeep/releases/tag/v1.13.0);
- 2022/08/07: [v1.12.0 release](https://github.com/maxatome/go-testdeep/releases/tag/v1.12.0);
- 2022/01/05: [v1.11.0 release](https://github.com/maxatome/go-testdeep/releases/tag/v1.11.0);
- see [commits history](https://github.com/maxatome/go-testdeep/commits/master)
  for other/older changes.


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
[TestDeep operators](https://go-testdeep.zetta.rocks/operators/):

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
flexible [operators](https://go-testdeep.zetta.rocks/operators/):

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
  "id":         $id,             // set by the API/DB
  "name":       "Alice",
  "age":        Between(40, 45), // ← ④
  "created_at": "$createdAt",    // set by the API/DB
}`,
      td.Tag("id", td.Catch(&id, td.NotZero())),        // ← ⑤
      td.Tag("createdAt", td.All(                       // ← ⑥
        td.HasSuffix("Z"),                              // ← ⑦
        td.Smuggle(func(s string) (time.Time, error) {  // ← ⑧
          return time.Parse(time.RFC3339Nano, s)
        }, td.Catch(&createdAt, td.Between(testAPI.SentAt(), time.Now()))), // ← ⑨
      )),
    ))
  if !testAPI.Failed() {
    t.Logf("The new Person ID is %d and was created at %s", id, createdAt)
  }
}
```

1. the API handler ready to be tested;
2. the POST request with automatic JSON marshalling;
3. the expected response HTTP status should be `http.StatusCreated`
   and the line just below, the body should match the [`JSON`] operator;
4. some operators can be embedded, like [`Between`] here;
5. for the `$id` placeholder, [`Catch`] its
   value: put it in `id` variable and check it is [`NotZero`];
6. for the `$createdAt` placeholder, use the [`All`]
   operator. It combines several operators like a AND;
7. check that `$createdAt` date ends with "Z" using
   [`HasSuffix`]. As we expect a RFC3339
   date, we require it in UTC time zone;
8. convert `$createdAt` date into a `time.Time` using a custom
   function thanks to the [`Smuggle`] operator;
9. then [`Catch`] the resulting value: put it in
   `createdAt` variable and check it is greater or equal than
   `testAPI.SentAt()` (the time just before the request is handled) and lesser
   or equal than `time.Now()`.

See [`tdhttp`] helper or the
[FAQ](https://go-testdeep.zetta.rocks/faq/#what-about-testing-the-response-using-my-api)
for details about HTTP API testing.

Example of produced error in case of mismatch:

![error output](https://github.com/maxatome/go-testdeep-site/raw/master/docs_src/static/images/colored-output.svg)


## Description

go-testdeep is historically a go rewrite and adaptation of wonderful
[Test::Deep perl](https://metacpan.org/pod/Test::Deep).

In golang, comparing data structure is usually done using
[reflect.DeepEqual](https://golang.org/pkg/reflect/#DeepEqual) or
using a package that uses this function behind the scene.

This function works very well, but it is not flexible. Both compared
structures must match exactly and when a difference is returned, it is
up to the caller to display it. Not easy when comparing big data
structures.

The purpose of go-testdeep, via the
[`td` package](https://pkg.go.dev/github.com/maxatome/go-testdeep/td)
and its
[helpers](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers),
is to do its best to introduce this missing flexibility using
["operators"](https://go-testdeep.zetta.rocks/operators/), when the
expected value (or one of its component) cannot be matched exactly,
mixed with some useful
[comparison functions](https://go-testdeep.zetta.rocks/functions/).

**See [go-testdeep.zetta.rocks](https://go-testdeep.zetta.rocks/) for
details.**


## Installation

```sh
$ go get github.com/maxatome/go-testdeep
```


## Helpers

The goal of helpers is to make use of `go-testdeep` even more powerful
by providing common features using
[TestDeep operators](https://go-testdeep.zetta.rocks/operators/)
behind the scene.

### `tdhttp` or HTTP API testing helper

The package `github.com/maxatome/go-testdeep/helpers/tdhttp` provides
some functions to easily test HTTP handlers.

See [`tdhttp`] documentation for details or
[FAQ](https://go-testdeep.zetta.rocks/faq/#what-about-testing-the-response-using-my-api) for an
example of use.

### `tdsuite` or testing suite helper

The package `github.com/maxatome/go-testdeep/helpers/tdsuite` adds tests
suite feature to go-testdeep in a non-intrusive way, but easily and powerfully.

A tests suite is a set of tests run sequentially that share some data.

Some hooks can be set to be automatically called before the suite is run,
before, after and/or between each test, and at the end of the suite. 

See [`tdsuite`] documentation for details.

### `tdutil` aka the helper of helpers

The package `github.com/maxatome/go-testdeep/helpers/tdutil` allows to
write unit tests for go-testdeep helpers and so provides some helpful
functions.

See
[`tdutil`](https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdutil)
for details.


## See also

- [testify](https://github.com/stretchr/testify): a toolkit with common assertions and mocks that plays nicely with the standard library
- [go-cmp](https://github.com/google/go-cmp): package for comparing Go values in tests


## License

`go-testdeep` is released under the BSD-style license found in the
[`LICENSE`](LICENSE) file in the root directory of this source tree.

Internal function `deepValueEqual` is based on `deepValueEqual` from
[`reflect` golang package](https://golang.org/pkg/reflect/) licensed
under the BSD-style license found in the [`LICENSE` file in the golang
repository](https://github.com/golang/go/blob/master/LICENSE).

Uses two files (`bypass.go` & `bypasssafe.go`) from
[Go-spew](https://github.com/davecgh/go-spew) which is licensed under
the [copyfree](http://copyfree.org) ISC License.

[Public Domain Gopher](https://github.com/egonelbre/gophers) provided
by [Egon Elbre](http://egonelbre.com/). The Go gopher was designed by
[Renee French](https://reneefrench.blogspot.com/).


## FAQ

See [FAQ](https://go-testdeep.zetta.rocks/faq/).


<!-- links:begin -->
[`T`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T
[`TestDeep`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#TestDeep
[`Cmp`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp

[`tdhttp`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp
[`tdsuite`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdsuite
[`tdutil`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdutil

[`BeLax` config flag]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ContextConfig.BeLax
[`error`]: https://pkg.go.dev/builtin#error


[`fmt.Stringer`]: https://pkg.go.dev/fmt/#Stringer
[`time.Time`]: https://pkg.go.dev/time/#Time
[`math.NaN`]: https://pkg.go.dev/math/#NaN
[`All`]: https://go-testdeep.zetta.rocks/operators/all/
[`Any`]: https://go-testdeep.zetta.rocks/operators/any/
[`Array`]: https://go-testdeep.zetta.rocks/operators/array/
[`ArrayEach`]: https://go-testdeep.zetta.rocks/operators/arrayeach/
[`Bag`]: https://go-testdeep.zetta.rocks/operators/bag/
[`Between`]: https://go-testdeep.zetta.rocks/operators/between/
[`Cap`]: https://go-testdeep.zetta.rocks/operators/cap/
[`Catch`]: https://go-testdeep.zetta.rocks/operators/catch/
[`Code`]: https://go-testdeep.zetta.rocks/operators/code/
[`Contains`]: https://go-testdeep.zetta.rocks/operators/contains/
[`ContainsKey`]: https://go-testdeep.zetta.rocks/operators/containskey/
[`Delay`]: https://go-testdeep.zetta.rocks/operators/delay/
[`Empty`]: https://go-testdeep.zetta.rocks/operators/empty/
[`ErrorIs`]: https://go-testdeep.zetta.rocks/operators/erroris/
[`First`]: https://go-testdeep.zetta.rocks/operators/first/
[`Grep`]: https://go-testdeep.zetta.rocks/operators/grep/
[`Gt`]: https://go-testdeep.zetta.rocks/operators/gt/
[`Gte`]: https://go-testdeep.zetta.rocks/operators/gte/
[`HasPrefix`]: https://go-testdeep.zetta.rocks/operators/hasprefix/
[`HasSuffix`]: https://go-testdeep.zetta.rocks/operators/hassuffix/
[`Ignore`]: https://go-testdeep.zetta.rocks/operators/ignore/
[`Isa`]: https://go-testdeep.zetta.rocks/operators/isa/
[`JSON`]: https://go-testdeep.zetta.rocks/operators/json/
[`JSONPointer`]: https://go-testdeep.zetta.rocks/operators/jsonpointer/
[`Keys`]: https://go-testdeep.zetta.rocks/operators/keys/
[`Last`]: https://go-testdeep.zetta.rocks/operators/last/
[`Lax`]: https://go-testdeep.zetta.rocks/operators/lax/
[`Len`]: https://go-testdeep.zetta.rocks/operators/len/
[`Lt`]: https://go-testdeep.zetta.rocks/operators/lt/
[`Lte`]: https://go-testdeep.zetta.rocks/operators/lte/
[`Map`]: https://go-testdeep.zetta.rocks/operators/map/
[`MapEach`]: https://go-testdeep.zetta.rocks/operators/mapeach/
[`N`]: https://go-testdeep.zetta.rocks/operators/n/
[`NaN`]: https://go-testdeep.zetta.rocks/operators/nan/
[`Nil`]: https://go-testdeep.zetta.rocks/operators/nil/
[`None`]: https://go-testdeep.zetta.rocks/operators/none/
[`Not`]: https://go-testdeep.zetta.rocks/operators/not/
[`NotAny`]: https://go-testdeep.zetta.rocks/operators/notany/
[`NotEmpty`]: https://go-testdeep.zetta.rocks/operators/notempty/
[`NotNaN`]: https://go-testdeep.zetta.rocks/operators/notnan/
[`NotNil`]: https://go-testdeep.zetta.rocks/operators/notnil/
[`NotZero`]: https://go-testdeep.zetta.rocks/operators/notzero/
[`PPtr`]: https://go-testdeep.zetta.rocks/operators/pptr/
[`Ptr`]: https://go-testdeep.zetta.rocks/operators/ptr/
[`Re`]: https://go-testdeep.zetta.rocks/operators/re/
[`ReAll`]: https://go-testdeep.zetta.rocks/operators/reall/
[`Recv`]: https://go-testdeep.zetta.rocks/operators/recv/
[`Set`]: https://go-testdeep.zetta.rocks/operators/set/
[`Shallow`]: https://go-testdeep.zetta.rocks/operators/shallow/
[`Slice`]: https://go-testdeep.zetta.rocks/operators/slice/
[`Smuggle`]: https://go-testdeep.zetta.rocks/operators/smuggle/
[`SStruct`]: https://go-testdeep.zetta.rocks/operators/sstruct/
[`String`]: https://go-testdeep.zetta.rocks/operators/string/
[`Struct`]: https://go-testdeep.zetta.rocks/operators/struct/
[`SubBagOf`]: https://go-testdeep.zetta.rocks/operators/subbagof/
[`SubJSONOf`]: https://go-testdeep.zetta.rocks/operators/subjsonof/
[`SubMapOf`]: https://go-testdeep.zetta.rocks/operators/submapof/
[`SubSetOf`]: https://go-testdeep.zetta.rocks/operators/subsetof/
[`SuperBagOf`]: https://go-testdeep.zetta.rocks/operators/superbagof/
[`SuperJSONOf`]: https://go-testdeep.zetta.rocks/operators/superjsonof/
[`SuperMapOf`]: https://go-testdeep.zetta.rocks/operators/supermapof/
[`SuperSetOf`]: https://go-testdeep.zetta.rocks/operators/supersetof/
[`SuperSliceOf`]: https://go-testdeep.zetta.rocks/operators/supersliceof/
[`Tag`]: https://go-testdeep.zetta.rocks/operators/tag/
[`TruncTime`]: https://go-testdeep.zetta.rocks/operators/trunctime/
[`Values`]: https://go-testdeep.zetta.rocks/operators/values/
[`Zero`]: https://go-testdeep.zetta.rocks/operators/zero/

[`CmpAll`]: https://go-testdeep.zetta.rocks/operators/all/#cmpall-shortcut
[`CmpAny`]: https://go-testdeep.zetta.rocks/operators/any/#cmpany-shortcut
[`CmpArray`]: https://go-testdeep.zetta.rocks/operators/array/#cmparray-shortcut
[`CmpArrayEach`]: https://go-testdeep.zetta.rocks/operators/arrayeach/#cmparrayeach-shortcut
[`CmpBag`]: https://go-testdeep.zetta.rocks/operators/bag/#cmpbag-shortcut
[`CmpBetween`]: https://go-testdeep.zetta.rocks/operators/between/#cmpbetween-shortcut
[`CmpCap`]: https://go-testdeep.zetta.rocks/operators/cap/#cmpcap-shortcut
[`CmpCode`]: https://go-testdeep.zetta.rocks/operators/code/#cmpcode-shortcut
[`CmpContains`]: https://go-testdeep.zetta.rocks/operators/contains/#cmpcontains-shortcut
[`CmpContainsKey`]: https://go-testdeep.zetta.rocks/operators/containskey/#cmpcontainskey-shortcut
[`CmpEmpty`]: https://go-testdeep.zetta.rocks/operators/empty/#cmpempty-shortcut
[`CmpErrorIs`]: https://go-testdeep.zetta.rocks/operators/erroris/#cmperroris-shortcut
[`CmpFirst`]: https://go-testdeep.zetta.rocks/operators/first/#cmpfirst-shortcut
[`CmpGrep`]: https://go-testdeep.zetta.rocks/operators/grep/#cmpgrep-shortcut
[`CmpGt`]: https://go-testdeep.zetta.rocks/operators/gt/#cmpgt-shortcut
[`CmpGte`]: https://go-testdeep.zetta.rocks/operators/gte/#cmpgte-shortcut
[`CmpHasPrefix`]: https://go-testdeep.zetta.rocks/operators/hasprefix/#cmphasprefix-shortcut
[`CmpHasSuffix`]: https://go-testdeep.zetta.rocks/operators/hassuffix/#cmphassuffix-shortcut
[`CmpIsa`]: https://go-testdeep.zetta.rocks/operators/isa/#cmpisa-shortcut
[`CmpJSON`]: https://go-testdeep.zetta.rocks/operators/json/#cmpjson-shortcut
[`CmpJSONPointer`]: https://go-testdeep.zetta.rocks/operators/jsonpointer/#cmpjsonpointer-shortcut
[`CmpKeys`]: https://go-testdeep.zetta.rocks/operators/keys/#cmpkeys-shortcut
[`CmpLast`]: https://go-testdeep.zetta.rocks/operators/last/#cmplast-shortcut
[`CmpLax`]: https://go-testdeep.zetta.rocks/operators/lax/#cmplax-shortcut
[`CmpLen`]: https://go-testdeep.zetta.rocks/operators/len/#cmplen-shortcut
[`CmpLt`]: https://go-testdeep.zetta.rocks/operators/lt/#cmplt-shortcut
[`CmpLte`]: https://go-testdeep.zetta.rocks/operators/lte/#cmplte-shortcut
[`CmpMap`]: https://go-testdeep.zetta.rocks/operators/map/#cmpmap-shortcut
[`CmpMapEach`]: https://go-testdeep.zetta.rocks/operators/mapeach/#cmpmapeach-shortcut
[`CmpN`]: https://go-testdeep.zetta.rocks/operators/n/#cmpn-shortcut
[`CmpNaN`]: https://go-testdeep.zetta.rocks/operators/nan/#cmpnan-shortcut
[`CmpNil`]: https://go-testdeep.zetta.rocks/operators/nil/#cmpnil-shortcut
[`CmpNone`]: https://go-testdeep.zetta.rocks/operators/none/#cmpnone-shortcut
[`CmpNot`]: https://go-testdeep.zetta.rocks/operators/not/#cmpnot-shortcut
[`CmpNotAny`]: https://go-testdeep.zetta.rocks/operators/notany/#cmpnotany-shortcut
[`CmpNotEmpty`]: https://go-testdeep.zetta.rocks/operators/notempty/#cmpnotempty-shortcut
[`CmpNotNaN`]: https://go-testdeep.zetta.rocks/operators/notnan/#cmpnotnan-shortcut
[`CmpNotNil`]: https://go-testdeep.zetta.rocks/operators/notnil/#cmpnotnil-shortcut
[`CmpNotZero`]: https://go-testdeep.zetta.rocks/operators/notzero/#cmpnotzero-shortcut
[`CmpPPtr`]: https://go-testdeep.zetta.rocks/operators/pptr/#cmppptr-shortcut
[`CmpPtr`]: https://go-testdeep.zetta.rocks/operators/ptr/#cmpptr-shortcut
[`CmpRe`]: https://go-testdeep.zetta.rocks/operators/re/#cmpre-shortcut
[`CmpReAll`]: https://go-testdeep.zetta.rocks/operators/reall/#cmpreall-shortcut
[`CmpRecv`]: https://go-testdeep.zetta.rocks/operators/recv/#cmprecv-shortcut
[`CmpSet`]: https://go-testdeep.zetta.rocks/operators/set/#cmpset-shortcut
[`CmpShallow`]: https://go-testdeep.zetta.rocks/operators/shallow/#cmpshallow-shortcut
[`CmpSlice`]: https://go-testdeep.zetta.rocks/operators/slice/#cmpslice-shortcut
[`CmpSmuggle`]: https://go-testdeep.zetta.rocks/operators/smuggle/#cmpsmuggle-shortcut
[`CmpSStruct`]: https://go-testdeep.zetta.rocks/operators/sstruct/#cmpsstruct-shortcut
[`CmpString`]: https://go-testdeep.zetta.rocks/operators/string/#cmpstring-shortcut
[`CmpStruct`]: https://go-testdeep.zetta.rocks/operators/struct/#cmpstruct-shortcut
[`CmpSubBagOf`]: https://go-testdeep.zetta.rocks/operators/subbagof/#cmpsubbagof-shortcut
[`CmpSubJSONOf`]: https://go-testdeep.zetta.rocks/operators/subjsonof/#cmpsubjsonof-shortcut
[`CmpSubMapOf`]: https://go-testdeep.zetta.rocks/operators/submapof/#cmpsubmapof-shortcut
[`CmpSubSetOf`]: https://go-testdeep.zetta.rocks/operators/subsetof/#cmpsubsetof-shortcut
[`CmpSuperBagOf`]: https://go-testdeep.zetta.rocks/operators/superbagof/#cmpsuperbagof-shortcut
[`CmpSuperJSONOf`]: https://go-testdeep.zetta.rocks/operators/superjsonof/#cmpsuperjsonof-shortcut
[`CmpSuperMapOf`]: https://go-testdeep.zetta.rocks/operators/supermapof/#cmpsupermapof-shortcut
[`CmpSuperSetOf`]: https://go-testdeep.zetta.rocks/operators/supersetof/#cmpsupersetof-shortcut
[`CmpSuperSliceOf`]: https://go-testdeep.zetta.rocks/operators/supersliceof/#cmpsupersliceof-shortcut
[`CmpTruncTime`]: https://go-testdeep.zetta.rocks/operators/trunctime/#cmptrunctime-shortcut
[`CmpValues`]: https://go-testdeep.zetta.rocks/operators/values/#cmpvalues-shortcut
[`CmpZero`]: https://go-testdeep.zetta.rocks/operators/zero/#cmpzero-shortcut

[`T.All`]: https://go-testdeep.zetta.rocks/operators/all/#tall-shortcut
[`T.Any`]: https://go-testdeep.zetta.rocks/operators/any/#tany-shortcut
[`T.Array`]: https://go-testdeep.zetta.rocks/operators/array/#tarray-shortcut
[`T.ArrayEach`]: https://go-testdeep.zetta.rocks/operators/arrayeach/#tarrayeach-shortcut
[`T.Bag`]: https://go-testdeep.zetta.rocks/operators/bag/#tbag-shortcut
[`T.Between`]: https://go-testdeep.zetta.rocks/operators/between/#tbetween-shortcut
[`T.Cap`]: https://go-testdeep.zetta.rocks/operators/cap/#tcap-shortcut
[`T.Code`]: https://go-testdeep.zetta.rocks/operators/code/#tcode-shortcut
[`T.Contains`]: https://go-testdeep.zetta.rocks/operators/contains/#tcontains-shortcut
[`T.ContainsKey`]: https://go-testdeep.zetta.rocks/operators/containskey/#tcontainskey-shortcut
[`T.Empty`]: https://go-testdeep.zetta.rocks/operators/empty/#tempty-shortcut
[`T.CmpErrorIs`]: https://go-testdeep.zetta.rocks/operators/erroris/#tcmperroris-shortcut
[`T.First`]: https://go-testdeep.zetta.rocks/operators/first/#tfirst-shortcut
[`T.Grep`]: https://go-testdeep.zetta.rocks/operators/grep/#tgrep-shortcut
[`T.Gt`]: https://go-testdeep.zetta.rocks/operators/gt/#tgt-shortcut
[`T.Gte`]: https://go-testdeep.zetta.rocks/operators/gte/#tgte-shortcut
[`T.HasPrefix`]: https://go-testdeep.zetta.rocks/operators/hasprefix/#thasprefix-shortcut
[`T.HasSuffix`]: https://go-testdeep.zetta.rocks/operators/hassuffix/#thassuffix-shortcut
[`T.Isa`]: https://go-testdeep.zetta.rocks/operators/isa/#tisa-shortcut
[`T.JSON`]: https://go-testdeep.zetta.rocks/operators/json/#tjson-shortcut
[`T.JSONPointer`]: https://go-testdeep.zetta.rocks/operators/jsonpointer/#tjsonpointer-shortcut
[`T.Keys`]: https://go-testdeep.zetta.rocks/operators/keys/#tkeys-shortcut
[`T.Last`]: https://go-testdeep.zetta.rocks/operators/last/#tlast-shortcut
[`T.CmpLax`]: https://go-testdeep.zetta.rocks/operators/lax/#tcmplax-shortcut
[`T.Len`]: https://go-testdeep.zetta.rocks/operators/len/#tlen-shortcut
[`T.Lt`]: https://go-testdeep.zetta.rocks/operators/lt/#tlt-shortcut
[`T.Lte`]: https://go-testdeep.zetta.rocks/operators/lte/#tlte-shortcut
[`T.Map`]: https://go-testdeep.zetta.rocks/operators/map/#tmap-shortcut
[`T.MapEach`]: https://go-testdeep.zetta.rocks/operators/mapeach/#tmapeach-shortcut
[`T.N`]: https://go-testdeep.zetta.rocks/operators/n/#tn-shortcut
[`T.NaN`]: https://go-testdeep.zetta.rocks/operators/nan/#tnan-shortcut
[`T.Nil`]: https://go-testdeep.zetta.rocks/operators/nil/#tnil-shortcut
[`T.None`]: https://go-testdeep.zetta.rocks/operators/none/#tnone-shortcut
[`T.Not`]: https://go-testdeep.zetta.rocks/operators/not/#tnot-shortcut
[`T.NotAny`]: https://go-testdeep.zetta.rocks/operators/notany/#tnotany-shortcut
[`T.NotEmpty`]: https://go-testdeep.zetta.rocks/operators/notempty/#tnotempty-shortcut
[`T.NotNaN`]: https://go-testdeep.zetta.rocks/operators/notnan/#tnotnan-shortcut
[`T.NotNil`]: https://go-testdeep.zetta.rocks/operators/notnil/#tnotnil-shortcut
[`T.NotZero`]: https://go-testdeep.zetta.rocks/operators/notzero/#tnotzero-shortcut
[`T.PPtr`]: https://go-testdeep.zetta.rocks/operators/pptr/#tpptr-shortcut
[`T.Ptr`]: https://go-testdeep.zetta.rocks/operators/ptr/#tptr-shortcut
[`T.Re`]: https://go-testdeep.zetta.rocks/operators/re/#tre-shortcut
[`T.ReAll`]: https://go-testdeep.zetta.rocks/operators/reall/#treall-shortcut
[`T.Recv`]: https://go-testdeep.zetta.rocks/operators/recv/#trecv-shortcut
[`T.Set`]: https://go-testdeep.zetta.rocks/operators/set/#tset-shortcut
[`T.Shallow`]: https://go-testdeep.zetta.rocks/operators/shallow/#tshallow-shortcut
[`T.Slice`]: https://go-testdeep.zetta.rocks/operators/slice/#tslice-shortcut
[`T.Smuggle`]: https://go-testdeep.zetta.rocks/operators/smuggle/#tsmuggle-shortcut
[`T.SStruct`]: https://go-testdeep.zetta.rocks/operators/sstruct/#tsstruct-shortcut
[`T.String`]: https://go-testdeep.zetta.rocks/operators/string/#tstring-shortcut
[`T.Struct`]: https://go-testdeep.zetta.rocks/operators/struct/#tstruct-shortcut
[`T.SubBagOf`]: https://go-testdeep.zetta.rocks/operators/subbagof/#tsubbagof-shortcut
[`T.SubJSONOf`]: https://go-testdeep.zetta.rocks/operators/subjsonof/#tsubjsonof-shortcut
[`T.SubMapOf`]: https://go-testdeep.zetta.rocks/operators/submapof/#tsubmapof-shortcut
[`T.SubSetOf`]: https://go-testdeep.zetta.rocks/operators/subsetof/#tsubsetof-shortcut
[`T.SuperBagOf`]: https://go-testdeep.zetta.rocks/operators/superbagof/#tsuperbagof-shortcut
[`T.SuperJSONOf`]: https://go-testdeep.zetta.rocks/operators/superjsonof/#tsuperjsonof-shortcut
[`T.SuperMapOf`]: https://go-testdeep.zetta.rocks/operators/supermapof/#tsupermapof-shortcut
[`T.SuperSetOf`]: https://go-testdeep.zetta.rocks/operators/supersetof/#tsupersetof-shortcut
[`T.SuperSliceOf`]: https://go-testdeep.zetta.rocks/operators/supersliceof/#tsupersliceof-shortcut
[`T.TruncTime`]: https://go-testdeep.zetta.rocks/operators/trunctime/#ttrunctime-shortcut
[`T.Values`]: https://go-testdeep.zetta.rocks/operators/values/#tvalues-shortcut
[`T.Zero`]: https://go-testdeep.zetta.rocks/operators/zero/#tzero-shortcut
<!-- links:end -->
