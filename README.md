go-testdeep
===========

[![Build Status](https://travis-ci.org/maxatome/go-testdeep.svg?branch=master)](https://travis-ci.org/maxatome/go-testdeep)
[![Coverage Status](https://coveralls.io/repos/github/maxatome/go-testdeep/badge.svg?branch=master)](https://coveralls.io/github/maxatome/go-testdeep?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxatome/go-testdeep)](https://goreportcard.com/report/github.com/maxatome/go-testdeep)
[![GoDoc](https://godoc.org/github.com/maxatome/go-testdeep?status.svg)](https://godoc.org/github.com/maxatome/go-testdeep)
[![Version](https://img.shields.io/github/tag/maxatome/go-testdeep.svg)](https://github.com/maxatome/go-testdeep/releases)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go/#testing)

![testdeep](docs/image.png)

**Extremely flexible golang deep comparison, extends the go testing package.**

- [Latest news](#latest-news)
- [Synopsis](#synopsis)
- [Godoc table of contents](docs/toc.md#godoc-table-of-contents)
- [Installation](#installation)
- [Available operators](#available-operators)
- [Helpers](#helpers)
  - [`tdhttp` or HTTP API testing helper](https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp)
- [See also](#see-also)
- [License](#license)
- [FAQ](docs/FAQ.md)


## Latest news

- 2019/09/22: new
  [`BeLax` feature](https://godoc.org/github.com/maxatome/go-testdeep#T.BeLax)
  with its
  [`Lax`](https://godoc.org/github.com/maxatome/go-testdeep#Lax)
  operator counterpart (and its friends
  [`CmpLax`](https://godoc.org/github.com/maxatome/go-testdeep#CmpLax)
  &
  [`T.CmpLax`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpLax));
- 2019/07/07: multiple changes occurred:
  - `*T` type now implements `TestingFT`,
  - add [`UseEqual` feature](https://godoc.org/github.com/maxatome/go-testdeep#T.UseEqual)
    aka. delegates comparison to `Equal()` method of object,
  - [`tdhttp.NewRequest()`](https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp#NewRequest),
    [`tdhttp.NewJSONRequest()`](https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp#NewJSONRequest)
    and
    [`tdhttp.NewXMLRequest()`](https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp#NewXMLRequest)
    now accept headers definition,
- 2019/05/01: new
  [`Keys`](https://godoc.org/github.com/maxatome/go-testdeep#Keys) &
  [`Values`](https://godoc.org/github.com/maxatome/go-testdeep#Values)
  operators (and their friends
  [`CmpKeys`](https://godoc.org/github.com/maxatome/go-testdeep#CmpKeys),
  [`CmpValues`](https://godoc.org/github.com/maxatome/go-testdeep#CmpValues),
  [`T.Keys`](https://godoc.org/github.com/maxatome/go-testdeep#T.Keys)
  &
  [`T.Values`](https://godoc.org/github.com/maxatome/go-testdeep#T.Values));
- 2019/04/27: new
  [`Cmp`](https://godoc.org/github.com/maxatome/go-testdeep#Cmp)
  function and
  [`T.Cmp`](https://godoc.org/github.com/maxatome/go-testdeep#T.Cmp)
  method, shorter versions of
  [`CmpDeeply`](https://godoc.org/github.com/maxatome/go-testdeep#CmpDeeply)
  and [`T.CmpDeeply`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpDeeply);
- see [commits history](https://github.com/maxatome/go-testdeep/commits/master)
  for other/older changes.


## Synopsis

Simplest usage:

```go
import (
  "testing"
  td "github.com/maxatome/go-testdeep"
)

func TestMyFunc(t *testing.T) {
  td.Cmp(t, MyFunc(), &Info{Name: "Alice", Age: 42})
}
```

Example of produced error in case of mismatch:

![error output](tools/docs_src/static/images/colored-output.svg)


See [godoc table of contents](docs/toc.md#godoc-table-of-contents) for details.


## Installation

```sh
$ go get -u github.com/maxatome/go-testdeep
```


## Available operators

See functions returning [`TestDeep` interface][`TestDeep`]:

<!-- operators:begin -->
- [`All`] all expected values have to match;
- [`Any`] at least one expected value have to match;
- [`Array`] compares the contents of an array or a pointer on an array;
- [`ArrayEach`] compares each array or slice item;
- [`Bag`] compares the contents of an array or a slice without taking care of the order of items;
- [`Between`] checks that a number, string or [`time.Time`] is between two bounds;
- [`Cap`] checks an array, slice or channel capacity;
- [`Code`] checks using a custom function;
- [`Contains`] checks that a string, [`error`] or [`fmt.Stringer`] interfaces contain a sub-string; or an array, slice or map contain a value;
- [`ContainsKey`] checks that a map contains a key;
- [`Empty`] checks that an array, a channel, a map, a slice or a string is empty;
- [`Gt`] checks that a number, string or [`time.Time`] is greater than a value;
- [`Gte`] checks that a number, string or [`time.Time`] is greater or equal than a value;
- [`HasPrefix`] checks the prefix of a string, [`error`] or [`fmt.Stringer`] interfaces;
- [`HasSuffix`] checks the suffix of a string, [`error`] or [`fmt.Stringer`] interfaces;
- [`Ignore`] allows to ignore a comparison;
- [`Isa`] checks the data type or whether data implements an interface or not;
- [`Keys`] checks keys of a map;
- [`Lax`] temporarily enables [`BeLax` config flag];
- [`Len`] checks an array, slice, map, string or channel length;
- [`Lt`] checks that a number, string or [`time.Time`] is lesser than a value;
- [`Lte`] checks that a number, string or [`time.Time`] is lesser or equal than a value;
- [`Map`] compares the contents of a map;
- [`MapEach`] compares each map entry;
- [`N`] compares a number with a tolerance value;
- [`NaN`] checks a floating number is [`math.NaN`];
- [`Nil`] compares to `nil`;
- [`None`] no values have to match;
- [`Not`] value must not match;
- [`NotAny`] compares the contents of an array or a slice, no values have to match;
- [`NotEmpty`] checks that an array, a channel, a map, a slice or a string is not empty;
- [`NotNaN`] checks a floating number is not [`math.NaN`];
- [`NotNil`] checks that data is not `nil`;
- [`NotZero`] checks that data is not zero regarding its type;
- [`PPtr`] allows to easily test a pointer of pointer value;
- [`Ptr`] allows to easily test a pointer value;
- [`Re`] allows to apply a regexp on a string (or convertible), `[]byte`, [`error`] or [`fmt.Stringer`] interfaces, and even test the captured groups;
- [`ReAll`] allows to successively apply a regexp on a string (or convertible), `[]byte`, [`error`] or [`fmt.Stringer`] interfaces, and even test the captured groups;
- [`Set`] compares the contents of an array or a slice ignoring duplicates and without taking care of the order of items;
- [`Shallow`] compares pointers only, not their contents;
- [`Slice`] compares the contents of a slice or a pointer on a slice;
- [`Smuggle`] changes data contents or mutates it into another type via a custom function or a struct fields-path before stepping down in favor of generic comparison process;
- [`String`] checks a string, [`error`] or [`fmt.Stringer`] interfaces string contents;
- [`Struct`] compares the contents of a struct or a pointer on a struct;
- [`SubBagOf`] compares the contents of an array or a slice without taking care of the order of items but with potentially some exclusions;
- [`SubMapOf`] compares the contents of a map but with potentially some exclusions;
- [`SubSetOf`] compares the contents of an array or a slice ignoring duplicates and without taking care of the order of items but with potentially some exclusions;
- [`SuperBagOf`] compares the contents of an array or a slice without taking care of the order of items but with potentially some extra items;
- [`SuperMapOf`] compares the contents of a map but with potentially some extra entries;
- [`SuperSetOf`] compares the contents of an array or a slice ignoring duplicates and without taking care of the order of items but with potentially some extra items;
- [`TruncTime`] compares [`time.Time`] (or assignable) values after truncating them;
- [`Values`] checks values of a map;
- [`Zero`] checks data against its zero'ed conterpart.
<!-- operators:end -->


## Helpers

The goal of helpers is to make use of `go-testdeep` even more powerful
by providing common features using [TestDeep operators](#available-operators)
behind the scene.

### `tdhttp` or HTTP API testing helper

The package `github.com/maxatome/go-testdeep/helpers/tdhttp` provides
some functions to easily test HTTP handlers.

See [`tdhttp`] documentation for details or
[FAQ](docs/FAQ.md#what-about-testing-the-response-using-my-api) for an
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

See [FAQ](docs/FAQ.md).

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
[`All`]: https://godoc.org/github.com/maxatome/go-testdeep#All
[`Any`]: https://godoc.org/github.com/maxatome/go-testdeep#Any
[`Array`]: https://godoc.org/github.com/maxatome/go-testdeep#Array
[`ArrayEach`]: https://godoc.org/github.com/maxatome/go-testdeep#ArrayEach
[`Bag`]: https://godoc.org/github.com/maxatome/go-testdeep#Bag
[`Between`]: https://godoc.org/github.com/maxatome/go-testdeep#Between
[`Cap`]: https://godoc.org/github.com/maxatome/go-testdeep#Cap
[`Code`]: https://godoc.org/github.com/maxatome/go-testdeep#Code
[`Contains`]: https://godoc.org/github.com/maxatome/go-testdeep#Contains
[`ContainsKey`]: https://godoc.org/github.com/maxatome/go-testdeep#ContainsKey
[`Empty`]: https://godoc.org/github.com/maxatome/go-testdeep#Empty
[`Gt`]: https://godoc.org/github.com/maxatome/go-testdeep#Gt
[`Gte`]: https://godoc.org/github.com/maxatome/go-testdeep#Gte
[`HasPrefix`]: https://godoc.org/github.com/maxatome/go-testdeep#HasPrefix
[`HasSuffix`]: https://godoc.org/github.com/maxatome/go-testdeep#HasSuffix
[`Ignore`]: https://godoc.org/github.com/maxatome/go-testdeep#Ignore
[`Isa`]: https://godoc.org/github.com/maxatome/go-testdeep#Isa
[`Keys`]: https://godoc.org/github.com/maxatome/go-testdeep#Keys
[`Lax`]: https://godoc.org/github.com/maxatome/go-testdeep#Lax
[`Len`]: https://godoc.org/github.com/maxatome/go-testdeep#Len
[`Lt`]: https://godoc.org/github.com/maxatome/go-testdeep#Lt
[`Lte`]: https://godoc.org/github.com/maxatome/go-testdeep#Lte
[`Map`]: https://godoc.org/github.com/maxatome/go-testdeep#Map
[`MapEach`]: https://godoc.org/github.com/maxatome/go-testdeep#MapEach
[`N`]: https://godoc.org/github.com/maxatome/go-testdeep#N
[`NaN`]: https://godoc.org/github.com/maxatome/go-testdeep#NaN
[`Nil`]: https://godoc.org/github.com/maxatome/go-testdeep#Nil
[`None`]: https://godoc.org/github.com/maxatome/go-testdeep#None
[`Not`]: https://godoc.org/github.com/maxatome/go-testdeep#Not
[`NotAny`]: https://godoc.org/github.com/maxatome/go-testdeep#NotAny
[`NotEmpty`]: https://godoc.org/github.com/maxatome/go-testdeep#NotEmpty
[`NotNaN`]: https://godoc.org/github.com/maxatome/go-testdeep#NotNaN
[`NotNil`]: https://godoc.org/github.com/maxatome/go-testdeep#NotNil
[`NotZero`]: https://godoc.org/github.com/maxatome/go-testdeep#NotZero
[`PPtr`]: https://godoc.org/github.com/maxatome/go-testdeep#PPtr
[`Ptr`]: https://godoc.org/github.com/maxatome/go-testdeep#Ptr
[`Re`]: https://godoc.org/github.com/maxatome/go-testdeep#Re
[`ReAll`]: https://godoc.org/github.com/maxatome/go-testdeep#ReAll
[`Set`]: https://godoc.org/github.com/maxatome/go-testdeep#Set
[`Shallow`]: https://godoc.org/github.com/maxatome/go-testdeep#Shallow
[`Slice`]: https://godoc.org/github.com/maxatome/go-testdeep#Slice
[`Smuggle`]: https://godoc.org/github.com/maxatome/go-testdeep#Smuggle
[`String`]: https://godoc.org/github.com/maxatome/go-testdeep#String
[`Struct`]: https://godoc.org/github.com/maxatome/go-testdeep#Struct
[`SubBagOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SubBagOf
[`SubMapOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SubMapOf
[`SubSetOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SubSetOf
[`SuperBagOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SuperBagOf
[`SuperMapOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SuperMapOf
[`SuperSetOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SuperSetOf
[`TruncTime`]: https://godoc.org/github.com/maxatome/go-testdeep#TruncTime
[`Values`]: https://godoc.org/github.com/maxatome/go-testdeep#Values
[`Zero`]: https://godoc.org/github.com/maxatome/go-testdeep#Zero

[`CmpAll`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpAll
[`CmpAny`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpAny
[`CmpArray`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpArray
[`CmpArrayEach`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpArrayEach
[`CmpBag`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpBag
[`CmpBetween`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpBetween
[`CmpCap`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpCap
[`CmpCode`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpCode
[`CmpContains`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpContains
[`CmpContainsKey`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpContainsKey
[`CmpEmpty`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpEmpty
[`CmpGt`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpGt
[`CmpGte`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpGte
[`CmpHasPrefix`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpHasPrefix
[`CmpHasSuffix`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpHasSuffix
[`CmpIsa`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpIsa
[`CmpKeys`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpKeys
[`CmpLax`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpLax
[`CmpLen`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpLen
[`CmpLt`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpLt
[`CmpLte`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpLte
[`CmpMap`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpMap
[`CmpMapEach`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpMapEach
[`CmpN`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpN
[`CmpNaN`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpNaN
[`CmpNil`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpNil
[`CmpNone`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpNone
[`CmpNot`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpNot
[`CmpNotAny`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpNotAny
[`CmpNotEmpty`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpNotEmpty
[`CmpNotNaN`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpNotNaN
[`CmpNotNil`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpNotNil
[`CmpNotZero`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpNotZero
[`CmpPPtr`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpPPtr
[`CmpPtr`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpPtr
[`CmpRe`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpRe
[`CmpReAll`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpReAll
[`CmpSet`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpSet
[`CmpShallow`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpShallow
[`CmpSlice`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpSlice
[`CmpSmuggle`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpSmuggle
[`CmpString`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpString
[`CmpStruct`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpStruct
[`CmpSubBagOf`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpSubBagOf
[`CmpSubMapOf`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpSubMapOf
[`CmpSubSetOf`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpSubSetOf
[`CmpSuperBagOf`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpSuperBagOf
[`CmpSuperMapOf`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpSuperMapOf
[`CmpSuperSetOf`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpSuperSetOf
[`CmpTruncTime`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpTruncTime
[`CmpValues`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpValues
[`CmpZero`]:https://godoc.org/github.com/maxatome/go-testdeep#CmpZero

[`T.All`]: https://godoc.org/github.com/maxatome/go-testdeep#T.All
[`T.Any`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Any
[`T.Array`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Array
[`T.ArrayEach`]: https://godoc.org/github.com/maxatome/go-testdeep#T.ArrayEach
[`T.Bag`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Bag
[`T.Between`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Between
[`T.Cap`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Cap
[`T.Code`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Code
[`T.Contains`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Contains
[`T.ContainsKey`]: https://godoc.org/github.com/maxatome/go-testdeep#T.ContainsKey
[`T.Empty`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Empty
[`T.Gt`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Gt
[`T.Gte`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Gte
[`T.HasPrefix`]: https://godoc.org/github.com/maxatome/go-testdeep#T.HasPrefix
[`T.HasSuffix`]: https://godoc.org/github.com/maxatome/go-testdeep#T.HasSuffix
[`T.Isa`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Isa
[`T.Keys`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Keys
[`T.CmpLax`]: https://godoc.org/github.com/maxatome/go-testdeep#T.CmpLax
[`T.Len`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Len
[`T.Lt`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Lt
[`T.Lte`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Lte
[`T.Map`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Map
[`T.MapEach`]: https://godoc.org/github.com/maxatome/go-testdeep#T.MapEach
[`T.N`]: https://godoc.org/github.com/maxatome/go-testdeep#T.N
[`T.NaN`]: https://godoc.org/github.com/maxatome/go-testdeep#T.NaN
[`T.Nil`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Nil
[`T.None`]: https://godoc.org/github.com/maxatome/go-testdeep#T.None
[`T.Not`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Not
[`T.NotAny`]: https://godoc.org/github.com/maxatome/go-testdeep#T.NotAny
[`T.NotEmpty`]: https://godoc.org/github.com/maxatome/go-testdeep#T.NotEmpty
[`T.NotNaN`]: https://godoc.org/github.com/maxatome/go-testdeep#T.NotNaN
[`T.NotNil`]: https://godoc.org/github.com/maxatome/go-testdeep#T.NotNil
[`T.NotZero`]: https://godoc.org/github.com/maxatome/go-testdeep#T.NotZero
[`T.PPtr`]: https://godoc.org/github.com/maxatome/go-testdeep#T.PPtr
[`T.Ptr`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Ptr
[`T.Re`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Re
[`T.ReAll`]: https://godoc.org/github.com/maxatome/go-testdeep#T.ReAll
[`T.Set`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Set
[`T.Shallow`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Shallow
[`T.Slice`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Slice
[`T.Smuggle`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Smuggle
[`T.String`]: https://godoc.org/github.com/maxatome/go-testdeep#T.String
[`T.Struct`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Struct
[`T.SubBagOf`]: https://godoc.org/github.com/maxatome/go-testdeep#T.SubBagOf
[`T.SubMapOf`]: https://godoc.org/github.com/maxatome/go-testdeep#T.SubMapOf
[`T.SubSetOf`]: https://godoc.org/github.com/maxatome/go-testdeep#T.SubSetOf
[`T.SuperBagOf`]: https://godoc.org/github.com/maxatome/go-testdeep#T.SuperBagOf
[`T.SuperMapOf`]: https://godoc.org/github.com/maxatome/go-testdeep#T.SuperMapOf
[`T.SuperSetOf`]: https://godoc.org/github.com/maxatome/go-testdeep#T.SuperSetOf
[`T.TruncTime`]: https://godoc.org/github.com/maxatome/go-testdeep#T.TruncTime
[`T.Values`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Values
[`T.Zero`]: https://godoc.org/github.com/maxatome/go-testdeep#T.Zero
<!-- links:end -->
