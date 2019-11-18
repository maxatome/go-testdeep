go-testdeep
===========

[![Build Status](https://travis-ci.org/maxatome/go-testdeep.svg?branch=master)](https://travis-ci.org/maxatome/go-testdeep)
[![Coverage Status](https://coveralls.io/repos/github/maxatome/go-testdeep/badge.svg?branch=master)](https://coveralls.io/github/maxatome/go-testdeep?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxatome/go-testdeep)](https://goreportcard.com/report/github.com/maxatome/go-testdeep)
[![GoDoc](https://godoc.org/github.com/maxatome/go-testdeep?status.svg)](https://godoc.org/github.com/maxatome/go-testdeep)
[![Version](https://img.shields.io/github/tag/maxatome/go-testdeep.svg)](https://github.com/maxatome/go-testdeep/releases)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go/#testing)

![testdeep](tools/docs_src/static/images/logo.png)

**Extremely flexible golang deep comparison, extends the go testing package.**

- [Latest news](#latest-news)
- [Synopsis](#synopsis)
- [Description](#description)
- [Installation](#installation)
- [Functions](https://go-testdeep.zetta.rocks/functions/)
- [Available operators](https://go-testdeep.zetta.rocks/operators/)
- [Helpers](#helpers)
  - [`tdhttp` or HTTP API testing helper](https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp)
- [See also](#see-also)
- [License](#license)
- [FAQ](https://go-testdeep.zetta.rocks/faq/)


## Latest news

- 2019/11/18:
  - new [`SubJSONOf`] & [`SuperJSONOf`] operators (and their
    friends [`CmpSubJSONOf`], [`CmpSuperJSONOf`], [`T.SubJSONOf`] &
    [`T.SuperJSONOf`]),
  - JSON data can now contain comments and some operator shortcuts;
- 2019/11/01: new [`Catch`] operator;
- 2019/10/31: new [`JSON`] operator (and its friends [`CmpJSON`]
  & [`T.JSON`] along with new fully dedicated [`Tag`] operator;
- 2019/10/29: new web site
  [go-testdeep.zetta.rocks](https://go-testdeep.zetta.rocks/)
- see [commits history](https://github.com/maxatome/go-testdeep/commits/master)
  for other/older changes.


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
flexible [operators](https://go-testdeep.zetta.rocks/operators/):

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
// Note that comments are allowed
{
  "id":         $id,          // set by the API/DB
  "name":       "Bob",
  "age":        42,
  "created_at": "$createdAt", // set by the API/DB
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

![error output](tools/docs_src/static/images/colored-output.svg)


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
missing flexibility using
["operators"](https://go-testdeep.zetta.rocks/operators/), when the
expected value (or one of its component) cannot be matched exactly,
mixed with some useful
[comparison functions](https://go-testdeep.zetta.rocks/functions/).

**See [go-testdeep.zetta.rocks](https://go-testdeep.zetta.rocks/) for
details.**


## Installation

```sh
$ go get -u github.com/maxatome/go-testdeep
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
[`T`]: https://godoc.org/github.com/maxatome/go-testdeep#T
[`TestDeep`]: https://godoc.org/github.com/maxatome/go-testdeep#TestDeep
[`Cmp`]: https://godoc.org/github.com/maxatome/go-testdeep#Cmp

[`tdhttp`]: https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp

[`BeLax` config flag]: https://godoc.org/github.com/maxatome/go-testdeep#ContextConfig
[`error`]: https://golang.org/pkg/builtin/#error


[`fmt.Stringer`]: https://godoc.org/pkg/fmt/#Stringer
[`time.Time`]: https://godoc.org/pkg/time/#Time
[`math.NaN`]: https://godoc.org/pkg/math/#NaN
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
[`Empty`]: https://go-testdeep.zetta.rocks/operators/empty/
[`Gt`]: https://go-testdeep.zetta.rocks/operators/gt/
[`Gte`]: https://go-testdeep.zetta.rocks/operators/gte/
[`HasPrefix`]: https://go-testdeep.zetta.rocks/operators/hasprefix/
[`HasSuffix`]: https://go-testdeep.zetta.rocks/operators/hassuffix/
[`Ignore`]: https://go-testdeep.zetta.rocks/operators/ignore/
[`Isa`]: https://go-testdeep.zetta.rocks/operators/isa/
[`JSON`]: https://go-testdeep.zetta.rocks/operators/json/
[`Keys`]: https://go-testdeep.zetta.rocks/operators/keys/
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
[`Set`]: https://go-testdeep.zetta.rocks/operators/set/
[`Shallow`]: https://go-testdeep.zetta.rocks/operators/shallow/
[`Slice`]: https://go-testdeep.zetta.rocks/operators/slice/
[`Smuggle`]: https://go-testdeep.zetta.rocks/operators/smuggle/
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
[`Tag`]: https://go-testdeep.zetta.rocks/operators/tag/
[`TruncTime`]: https://go-testdeep.zetta.rocks/operators/trunctime/
[`Values`]: https://go-testdeep.zetta.rocks/operators/values/
[`Zero`]: https://go-testdeep.zetta.rocks/operators/zero/

[`CmpAll`]:https://go-testdeep.zetta.rocks/operators/all/#cmpall-shortcut
[`CmpAny`]:https://go-testdeep.zetta.rocks/operators/any/#cmpany-shortcut
[`CmpArray`]:https://go-testdeep.zetta.rocks/operators/array/#cmparray-shortcut
[`CmpArrayEach`]:https://go-testdeep.zetta.rocks/operators/arrayeach/#cmparrayeach-shortcut
[`CmpBag`]:https://go-testdeep.zetta.rocks/operators/bag/#cmpbag-shortcut
[`CmpBetween`]:https://go-testdeep.zetta.rocks/operators/between/#cmpbetween-shortcut
[`CmpCap`]:https://go-testdeep.zetta.rocks/operators/cap/#cmpcap-shortcut
[`CmpCode`]:https://go-testdeep.zetta.rocks/operators/code/#cmpcode-shortcut
[`CmpContains`]:https://go-testdeep.zetta.rocks/operators/contains/#cmpcontains-shortcut
[`CmpContainsKey`]:https://go-testdeep.zetta.rocks/operators/containskey/#cmpcontainskey-shortcut
[`CmpEmpty`]:https://go-testdeep.zetta.rocks/operators/empty/#cmpempty-shortcut
[`CmpGt`]:https://go-testdeep.zetta.rocks/operators/gt/#cmpgt-shortcut
[`CmpGte`]:https://go-testdeep.zetta.rocks/operators/gte/#cmpgte-shortcut
[`CmpHasPrefix`]:https://go-testdeep.zetta.rocks/operators/hasprefix/#cmphasprefix-shortcut
[`CmpHasSuffix`]:https://go-testdeep.zetta.rocks/operators/hassuffix/#cmphassuffix-shortcut
[`CmpIsa`]:https://go-testdeep.zetta.rocks/operators/isa/#cmpisa-shortcut
[`CmpJSON`]:https://go-testdeep.zetta.rocks/operators/json/#cmpjson-shortcut
[`CmpKeys`]:https://go-testdeep.zetta.rocks/operators/keys/#cmpkeys-shortcut
[`CmpLax`]:https://go-testdeep.zetta.rocks/operators/lax/#cmplax-shortcut
[`CmpLen`]:https://go-testdeep.zetta.rocks/operators/len/#cmplen-shortcut
[`CmpLt`]:https://go-testdeep.zetta.rocks/operators/lt/#cmplt-shortcut
[`CmpLte`]:https://go-testdeep.zetta.rocks/operators/lte/#cmplte-shortcut
[`CmpMap`]:https://go-testdeep.zetta.rocks/operators/map/#cmpmap-shortcut
[`CmpMapEach`]:https://go-testdeep.zetta.rocks/operators/mapeach/#cmpmapeach-shortcut
[`CmpN`]:https://go-testdeep.zetta.rocks/operators/n/#cmpn-shortcut
[`CmpNaN`]:https://go-testdeep.zetta.rocks/operators/nan/#cmpnan-shortcut
[`CmpNil`]:https://go-testdeep.zetta.rocks/operators/nil/#cmpnil-shortcut
[`CmpNone`]:https://go-testdeep.zetta.rocks/operators/none/#cmpnone-shortcut
[`CmpNot`]:https://go-testdeep.zetta.rocks/operators/not/#cmpnot-shortcut
[`CmpNotAny`]:https://go-testdeep.zetta.rocks/operators/notany/#cmpnotany-shortcut
[`CmpNotEmpty`]:https://go-testdeep.zetta.rocks/operators/notempty/#cmpnotempty-shortcut
[`CmpNotNaN`]:https://go-testdeep.zetta.rocks/operators/notnan/#cmpnotnan-shortcut
[`CmpNotNil`]:https://go-testdeep.zetta.rocks/operators/notnil/#cmpnotnil-shortcut
[`CmpNotZero`]:https://go-testdeep.zetta.rocks/operators/notzero/#cmpnotzero-shortcut
[`CmpPPtr`]:https://go-testdeep.zetta.rocks/operators/pptr/#cmppptr-shortcut
[`CmpPtr`]:https://go-testdeep.zetta.rocks/operators/ptr/#cmpptr-shortcut
[`CmpRe`]:https://go-testdeep.zetta.rocks/operators/re/#cmpre-shortcut
[`CmpReAll`]:https://go-testdeep.zetta.rocks/operators/reall/#cmpreall-shortcut
[`CmpSet`]:https://go-testdeep.zetta.rocks/operators/set/#cmpset-shortcut
[`CmpShallow`]:https://go-testdeep.zetta.rocks/operators/shallow/#cmpshallow-shortcut
[`CmpSlice`]:https://go-testdeep.zetta.rocks/operators/slice/#cmpslice-shortcut
[`CmpSmuggle`]:https://go-testdeep.zetta.rocks/operators/smuggle/#cmpsmuggle-shortcut
[`CmpString`]:https://go-testdeep.zetta.rocks/operators/string/#cmpstring-shortcut
[`CmpStruct`]:https://go-testdeep.zetta.rocks/operators/struct/#cmpstruct-shortcut
[`CmpSubBagOf`]:https://go-testdeep.zetta.rocks/operators/subbagof/#cmpsubbagof-shortcut
[`CmpSubJSONOf`]:https://go-testdeep.zetta.rocks/operators/subjsonof/#cmpsubjsonof-shortcut
[`CmpSubMapOf`]:https://go-testdeep.zetta.rocks/operators/submapof/#cmpsubmapof-shortcut
[`CmpSubSetOf`]:https://go-testdeep.zetta.rocks/operators/subsetof/#cmpsubsetof-shortcut
[`CmpSuperBagOf`]:https://go-testdeep.zetta.rocks/operators/superbagof/#cmpsuperbagof-shortcut
[`CmpSuperJSONOf`]:https://go-testdeep.zetta.rocks/operators/superjsonof/#cmpsuperjsonof-shortcut
[`CmpSuperMapOf`]:https://go-testdeep.zetta.rocks/operators/supermapof/#cmpsupermapof-shortcut
[`CmpSuperSetOf`]:https://go-testdeep.zetta.rocks/operators/supersetof/#cmpsupersetof-shortcut
[`CmpTruncTime`]:https://go-testdeep.zetta.rocks/operators/trunctime/#cmptrunctime-shortcut
[`CmpValues`]:https://go-testdeep.zetta.rocks/operators/values/#cmpvalues-shortcut
[`CmpZero`]:https://go-testdeep.zetta.rocks/operators/zero/#cmpzero-shortcut

[`T.All`]: https://go-testdeep.zetta.rocks/operators/all/#t-all-shortcut
[`T.Any`]: https://go-testdeep.zetta.rocks/operators/any/#t-any-shortcut
[`T.Array`]: https://go-testdeep.zetta.rocks/operators/array/#t-array-shortcut
[`T.ArrayEach`]: https://go-testdeep.zetta.rocks/operators/arrayeach/#t-arrayeach-shortcut
[`T.Bag`]: https://go-testdeep.zetta.rocks/operators/bag/#t-bag-shortcut
[`T.Between`]: https://go-testdeep.zetta.rocks/operators/between/#t-between-shortcut
[`T.Cap`]: https://go-testdeep.zetta.rocks/operators/cap/#t-cap-shortcut
[`T.Code`]: https://go-testdeep.zetta.rocks/operators/code/#t-code-shortcut
[`T.Contains`]: https://go-testdeep.zetta.rocks/operators/contains/#t-contains-shortcut
[`T.ContainsKey`]: https://go-testdeep.zetta.rocks/operators/containskey/#t-containskey-shortcut
[`T.Empty`]: https://go-testdeep.zetta.rocks/operators/empty/#t-empty-shortcut
[`T.Gt`]: https://go-testdeep.zetta.rocks/operators/gt/#t-gt-shortcut
[`T.Gte`]: https://go-testdeep.zetta.rocks/operators/gte/#t-gte-shortcut
[`T.HasPrefix`]: https://go-testdeep.zetta.rocks/operators/hasprefix/#t-hasprefix-shortcut
[`T.HasSuffix`]: https://go-testdeep.zetta.rocks/operators/hassuffix/#t-hassuffix-shortcut
[`T.Isa`]: https://go-testdeep.zetta.rocks/operators/isa/#t-isa-shortcut
[`T.JSON`]: https://go-testdeep.zetta.rocks/operators/json/#t-json-shortcut
[`T.Keys`]: https://go-testdeep.zetta.rocks/operators/keys/#t-keys-shortcut
[`T.CmpLax`]: https://go-testdeep.zetta.rocks/operators/lax/#t-cmplax-shortcut
[`T.Len`]: https://go-testdeep.zetta.rocks/operators/len/#t-len-shortcut
[`T.Lt`]: https://go-testdeep.zetta.rocks/operators/lt/#t-lt-shortcut
[`T.Lte`]: https://go-testdeep.zetta.rocks/operators/lte/#t-lte-shortcut
[`T.Map`]: https://go-testdeep.zetta.rocks/operators/map/#t-map-shortcut
[`T.MapEach`]: https://go-testdeep.zetta.rocks/operators/mapeach/#t-mapeach-shortcut
[`T.N`]: https://go-testdeep.zetta.rocks/operators/n/#t-n-shortcut
[`T.NaN`]: https://go-testdeep.zetta.rocks/operators/nan/#t-nan-shortcut
[`T.Nil`]: https://go-testdeep.zetta.rocks/operators/nil/#t-nil-shortcut
[`T.None`]: https://go-testdeep.zetta.rocks/operators/none/#t-none-shortcut
[`T.Not`]: https://go-testdeep.zetta.rocks/operators/not/#t-not-shortcut
[`T.NotAny`]: https://go-testdeep.zetta.rocks/operators/notany/#t-notany-shortcut
[`T.NotEmpty`]: https://go-testdeep.zetta.rocks/operators/notempty/#t-notempty-shortcut
[`T.NotNaN`]: https://go-testdeep.zetta.rocks/operators/notnan/#t-notnan-shortcut
[`T.NotNil`]: https://go-testdeep.zetta.rocks/operators/notnil/#t-notnil-shortcut
[`T.NotZero`]: https://go-testdeep.zetta.rocks/operators/notzero/#t-notzero-shortcut
[`T.PPtr`]: https://go-testdeep.zetta.rocks/operators/pptr/#t-pptr-shortcut
[`T.Ptr`]: https://go-testdeep.zetta.rocks/operators/ptr/#t-ptr-shortcut
[`T.Re`]: https://go-testdeep.zetta.rocks/operators/re/#t-re-shortcut
[`T.ReAll`]: https://go-testdeep.zetta.rocks/operators/reall/#t-reall-shortcut
[`T.Set`]: https://go-testdeep.zetta.rocks/operators/set/#t-set-shortcut
[`T.Shallow`]: https://go-testdeep.zetta.rocks/operators/shallow/#t-shallow-shortcut
[`T.Slice`]: https://go-testdeep.zetta.rocks/operators/slice/#t-slice-shortcut
[`T.Smuggle`]: https://go-testdeep.zetta.rocks/operators/smuggle/#t-smuggle-shortcut
[`T.String`]: https://go-testdeep.zetta.rocks/operators/string/#t-string-shortcut
[`T.Struct`]: https://go-testdeep.zetta.rocks/operators/struct/#t-struct-shortcut
[`T.SubBagOf`]: https://go-testdeep.zetta.rocks/operators/subbagof/#t-subbagof-shortcut
[`T.SubJSONOf`]: https://go-testdeep.zetta.rocks/operators/subjsonof/#t-subjsonof-shortcut
[`T.SubMapOf`]: https://go-testdeep.zetta.rocks/operators/submapof/#t-submapof-shortcut
[`T.SubSetOf`]: https://go-testdeep.zetta.rocks/operators/subsetof/#t-subsetof-shortcut
[`T.SuperBagOf`]: https://go-testdeep.zetta.rocks/operators/superbagof/#t-superbagof-shortcut
[`T.SuperJSONOf`]: https://go-testdeep.zetta.rocks/operators/superjsonof/#t-superjsonof-shortcut
[`T.SuperMapOf`]: https://go-testdeep.zetta.rocks/operators/supermapof/#t-supermapof-shortcut
[`T.SuperSetOf`]: https://go-testdeep.zetta.rocks/operators/supersetof/#t-supersetof-shortcut
[`T.TruncTime`]: https://go-testdeep.zetta.rocks/operators/trunctime/#t-trunctime-shortcut
[`T.Values`]: https://go-testdeep.zetta.rocks/operators/values/#t-values-shortcut
[`T.Zero`]: https://go-testdeep.zetta.rocks/operators/zero/#t-zero-shortcut
<!-- links:end -->
