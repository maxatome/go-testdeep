go-testdeep
===========

[![Build Status](https://travis-ci.org/maxatome/go-testdeep.svg)](https://travis-ci.org/maxatome/go-testdeep)
[![Coverage Status](https://coveralls.io/repos/github/maxatome/go-testdeep/badge.svg?branch=master)](https://coveralls.io/github/maxatome/go-testdeep?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxatome/go-testdeep)](https://goreportcard.com/report/github.com/maxatome/go-testdeep)
[![GoDoc](https://godoc.org/github.com/maxatome/go-testdeep?status.svg)](https://godoc.org/github.com/maxatome/go-testdeep)
[![Version](https://img.shields.io/github/tag/maxatome/go-testdeep.svg)](https://github.com/maxatome/go-testdeep/releases)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go/#testing)

Extremely flexible golang deep comparison, extends the go testing package.

- [Latest new](#latest-news)
- [Synopsis](#synopsis)
- [Installation](#installation)
- [Presentation](#presentation)
- [Available operators](#available-operators)
- [See also](#see-also)
- [License](#license)


## Latest news

- 2018/07/15: new
  [`NaN`](https://godoc.org/github.com/maxatome/go-testdeep#NaN) &
  [`NotNaN`](https://godoc.org/github.com/maxatome/go-testdeep#NotNaN) &
  operators (and their friends
  [`CmpNaN`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNaN),
  [`CmpNotNaN`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNotNaN),
  [`T.NaN`](https://godoc.org/github.com/maxatome/go-testdeep#T.NaN)
  &
  [`T.NotNaN`](https://godoc.org/github.com/maxatome/go-testdeep#T.NotNaN));
- 2018/06/24: [`Contains`](https://godoc.org/github.com/maxatome/go-testdeep#Contains)
  (and its friends
  [`CmpContains`](https://godoc.org/github.com/maxatome/go-testdeep#CmpContains)
  &
  [`T.Contains`](https://godoc.org/github.com/maxatome/go-testdeep#T.Contains))
  reworked to handle arrays, slices and maps;
- 2018/06/19: new
  [ContextConfig](https://godoc.org/github.com/maxatome/go-testdeep#ContextConfig)
  feature `FailureIsFatal` available. See
  [DefaultContextConfig](https://godoc.org/github.com/maxatome/go-testdeep#pkg-variables)
  for default global value and
  [`T.FailureIsFatal`](https://godoc.org/github.com/maxatome/go-testdeep#T.FailureIsFatal)
  method;
- see [commits history](https://github.com/maxatome/go-testdeep/commits/master)
  for other/older changes.


## Synopsis

```go
import (
  "testing"
  "time"

  td "github.com/maxatome/go-testdeep"
)

type Record struct {
  Id        uint64
  Name      string
  Age       int
  CreatedAt time.Time
}

func CreateRecord(name string, age int) (*Record, error) {
  ...
}

func TestCreateRecord(t *testing.T) {
  before := time.Now()
  record, err := CreateRecord("Bob", 23)

  if td.CmpNoError(t, err) {
    // If you know the exact value of all fields of newly created record
    td.CmpDeeply(t, record,
      &Record{
        Id:        245,
        Name:      "Bob",
        Age:       23,
        CreatedAt: time.Date(2018, time.May, 1, 11, 12, 13, 0, time.UTC),
      },
      "Fixed value for each field of newly created record")

    // But as often you cannot guess the values of DB generated fields,
    // you can choose to ignore them and only test the non-zero ones
    td.CmpDeeply(t, record,
      td.Struct(
        &Record{
          Name: "Bob",
          Age:  23,
        },
        nil),
      "Name & Age fields of newly created record")

    // Anyway, it is better to be able to test all fields!
    td.CmpDeeply(t, record,
      td.Struct(
        &Record{
          Name: "Bob",
          Age:  23,
        },
        td.StructFields{
          "Id":        td.NotZero(),
          "CreatedAt": td.Between(before, time.Now()),
        }),
      "Newly created record")
  }
}
```

Imagine `CreateRecord` does not set correctly `CreatedAt` field, then:
```sh
go test -run=TestCreateRecord
```

outputs for last `td.CmpDeeply` call:
```
--- FAIL: TestCreateRecord (0.00s)
  test_test.go:46: Failed test 'Newly created record'
    DATA.CreatedAt: values differ
           got: 2018-05-27 10:55:50.788166932 +0200 CEST m=-2.998149554
      expected: 2018-05-27 10:55:53.788163509 +0200 CEST m=+0.001848002 ≤ got ≤ 2018-05-27 10:55:53.788464176 +0200 CEST m=+0.002148179
    [under TestDeep operator Between at test_test.go:54]
FAIL
exit status 1
FAIL  github.com/maxatome/go-testdeep  0.006s
```

If `CreateRecord` had not set correctly `Id` field, output would have
been:
```
--- FAIL: TestCreateRecord (0.00s)
  test_test.go:46: Failed test 'Newly created record'
    DATA.Id: zero value
           got: (uint64) 0
      expected: NotZero()
    [under TestDeep operator Not at test_test.go:53]
FAIL
exit status 1
FAIL  github.com/maxatome/go-testdeep  0.006s
```

If `CreateRecord` had not set `Name` field to "Alice" value instead of
expected "Bob", output would have been:
```
--- FAIL: TestCreateRecord (0.00s)
  test_test.go:46: Failed test 'Newly created record'
    DATA.Name: values differ
           got: "Alice"
      expected: "Bob"
FAIL
exit status 1
FAIL  github.com/maxatome/go-testdeep  0.006s
```

Using [`testdeep.T`](https://godoc.org/github.com/maxatome/go-testdeep#T)
type, `TestCreateRecord` can also be written as:

```go
import (
  "testing"
  "time"

  td "github.com/maxatome/go-testdeep"
)

type Record struct {
  Id        uint64
  Name      string
  Age       int
  CreatedAt time.Time
}

func CreateRecord(name string, age int) (*Record, error) {
  ...
}

func TestCreateRecord(tt *testing.T) {
  t := td.NewT(tt)

  before := time.Now()
  record, err := CreateRecord("Bob", 23)

  if t.CmpNoError(err) {
    t := t.RootName("RECORD") // Use RECORD instead of DATA in failure reports

    // If you know the exact value of all fields of newly created record
    t.CmpDeeply(record,
      &Record{
        Id:        245,
        Name:      "Bob",
        Age:       23,
        CreatedAt: time.Date(2018, time.May, 1, 11, 12, 13, 0, time.UTC),
      },
      "Fixed value for each field of newly created record")

    // Anyway, it is better to be able to test all fields in a generic way!
    // Using Struct method
    t.Struct(record,
      Record{
        Name: "Bob",
        Age:  23,
      },
      td.StructFields{
        "Id":        td.NotZero(),
        "CreatedAt": td.Between(before, time.Now()),
      },
      "Newly created record")

    // Or using CmpDeeply method, it's a matter of taste
    t.CmpDeeply(record,
      td.Struct(
        Record{
          Name: "Bob",
          Age:  23,
        },
        td.StructFields{
          "Id":        td.NotZero(),
          "CreatedAt": td.Between(before, time.Now()),
        }),
      "Newly created record")
  }
}
```


## Installation

```sh
$ go get -u github.com/maxatome/go-testdeep
```


## Presentation

Package `testdeep` allows extremely flexible deep comparison, built
for testing.

It is a go rewrite and adaptation of wonderful
[`Test::Deep`](https://metacpan.org/pod/Test::Deep) perl module.

In golang, comparing data structure is usually done using
[`reflect.DeepEqual`](https://golang.org/pkg/reflect/#DeepEqual) or
using a package that uses this function behind the scene.

This function works very well, but it is not flexible. Both
compared structures must match exactly.

The purpose of testdeep package is to do its best to introduce this
missing flexibility using *operators* when the expected value (or
one of its component) cannot be matched exactly.

Imagine a function returning a struct containing a newly created
database record. The `Id` and the `CreatedAt` fields are set by the
database layer. In this case we have to do something like that to
check the record contents:

```go
import (
  "testing"
)

...

func TestCreateRecord(t *testing.T) {
  before := time.Now()
  record, err := CreateRecord("Bob", 23)

  if err != nil {
    t.Errorf("An error occurred: %s", err)
  } else {
    expected := Record{Name: "Bob", Age: 23}

    if record.Id == 0 {
      t.Error("Id probably not initialized")
    }
    if before.After(record.CreatedAt) ||
      time.Now().Before(record.CreatedAt) {
      t.Errorf("CreatedAt field not expected: %s", record.CreatedAt)
    }
    if record.Name != expected.Name {
      t.Errorf("Name field differ, got=%s, expected=%s",
        record.Name, expected.Name)
    }
    if record.Age != expected.Age {
      t.Errorf("Age field differ, got=%s, expected=%s",
        record.Age, expected.Age)
    }
  }
}
```

With `testdeep`, it is a way simple, thanks to [`CmpDeeply`](https://godoc.org/github.com/maxatome/go-testdeep#CmpDeeply) function:

```go
import (
  "testing"
  td "github.com/maxatome/go-testdeep"
)

...

func TestCreateRecord(t *testing.T) {
  before := time.Now()
  record, err := CreateRecord("Bob", 23)

  if td.CmpDeeply(t, err, nil) {
    td.CmpDeeply(t, record,
      td.Struct(
        Record{
          Name: "Bob",
          Age:  23,
        },
        td.StructFields{
          "Id":        td.NotZero(),
          "CreatedAt": td.Between(before, time.Now()),
        }),
      "Newly created record")
  }
}
```

Of course not only structs can be compared. A lot of
[operators](#available-operators) can be found below to cover most
(all?) needed tests.

The [`CmpDeeply`](https://godoc.org/github.com/maxatome/go-testdeep#CmpDeeply)
function is the keystone of this package, but to make the writing of
tests even easier, the family of `Cmp*` functions are provided and act
as shortcuts. Using
[`CmpNoError`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNoError)
and [`CmpStruct`](https://godoc.org/github.com/maxatome/go-testdeep#CmpStruct)
function, the previous example can be written as:

```go
import (
  "testing"
  td "github.com/maxatome/go-testdeep"
)

...

func TestCreateRecord(t *testing.T) {
  before := time.Now()
  record, err := CreateRecord("Bob", 23)

  if td.CmpNoError(t, err) {
    td.CmpStruct(t, record,
      Record{
        Name: "Bob",
        Age:  23,
      },
      td.StructFields{
        "Id":        td.NotZero(),
        "CreatedAt": td.Between(before, time.Now()),
      },
      "Newly created record")
  }
}
```

Last, [`testing.T`](https://golang.org/pkg/testing/#T) can be encapsulated in
[`T`](https://godoc.org/github.com/maxatome/go-testdeep#T) type,
simplifying again the test:

```go
import (
  "testing"
  td "github.com/maxatome/go-testdeep"
)

...

func TestCreateRecord(tt *testing.T) {
  t := td.NewT(tt)

  before := time.Now()
  record, err := CreateRecord()

  if t.CmpNoError(err) {
    t.Struct(record,
      Record{
        Name: "Bob",
        Age:  23,
      },
      td.StructFields{
        "Id":        td.NotZero(),
        "CreatedAt": td.Between(before, time.Now()),
      },
      "Newly created record")
  }
}
```


## Available operators

See functions returning [`TestDeep` interface](https://godoc.org/github.com/maxatome/go-testdeep#TestDeep):

- [`All`](https://godoc.org/github.com/maxatome/go-testdeep#All)
  all expected values have to match;
- [`Any`](https://godoc.org/github.com/maxatome/go-testdeep#Any)
  at least one expected value have to match;
- [`Array`](https://godoc.org/github.com/maxatome/go-testdeep#Array)
  compares the contents of an array or a pointer on an array;
- [`ArrayEach`](https://godoc.org/github.com/maxatome/go-testdeep#ArrayEach)
  compares each array or slice item;
- [`Bag`](https://godoc.org/github.com/maxatome/go-testdeep#Bag)
  compares the contents of an array or a slice without taking care of the order
  of items;
- [`Between`](https://godoc.org/github.com/maxatome/go-testdeep#Between)
  checks that a number or [`time.Time`](https://golang.org/pkg/time/)) is
  between two bounds;
- [`Cap`](https://godoc.org/github.com/maxatome/go-testdeep#Cap)
  checks an array, slice or channel capacity;
- [`Code`](https://godoc.org/github.com/maxatome/go-testdeep#Code)
  allows to use a custom function;
- [`Contains`](https://godoc.org/github.com/maxatome/go-testdeep#Contains)
  checks that a string, [`error`](https://golang.org/ref/spec#Errors) or
  [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interfaces contain
  a sub-string;
- [`Empty`](https://godoc.org/github.com/maxatome/go-testdeep#Empty)
  checks that an array, a channel, a map, a slice or a string is empty;
- [`Gt`](https://godoc.org/github.com/maxatome/go-testdeep#Gt)
  checks that a number or [`time.Time`](https://golang.org/pkg/time/)) is
  greater than a value;
- [`Gte`](https://godoc.org/github.com/maxatome/go-testdeep#Gte)
  checks that a number or [`time.Time`](https://golang.org/pkg/time/)) is
  greater or equal than a value;
- [`HasPrefix`](https://godoc.org/github.com/maxatome/go-testdeep#HasPrefix)
  checks the prefix of a string, [`error`](https://golang.org/ref/spec#Errors)
  or [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interfaces;
- [`HasSuffix`](https://godoc.org/github.com/maxatome/go-testdeep#HasSuffix)
  checks the suffix of a string, [`error`](https://golang.org/ref/spec#Errors)
  or [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interfaces;
- [`Ignore`](https://godoc.org/github.com/maxatome/go-testdeep#Isa)
  allows to ignore a comparison;
- [`Isa`](https://godoc.org/github.com/maxatome/go-testdeep#Isa)
  checks the data type or whether data implements an interface or not;
- [`Len`](https://godoc.org/github.com/maxatome/go-testdeep#Len)
  checks an array, slice, map, string or channel length;
- [`Lt`](https://godoc.org/github.com/maxatome/go-testdeep#Lt)
  checks that a number or [`time.Time`](https://golang.org/pkg/time/)) is
  lesser than a value;
- [`Lte`](https://godoc.org/github.com/maxatome/go-testdeep#Lte)
  checks that a number or [`time.Time`](https://golang.org/pkg/time/)) is
  lesser or equal than a value;
- [`Map`](https://godoc.org/github.com/maxatome/go-testdeep#Map)
  compares the contents of a map;
- [`MapEach`](https://godoc.org/github.com/maxatome/go-testdeep#MapEach)
  compares each map entry;
- [`N`](https://godoc.org/github.com/maxatome/go-testdeep#N)
  compares a number with a tolerance value;
- [`NaN`](https://godoc.org/github.com/maxatome/go-testdeep#NaN)
  checks a floating number is [`math.NaN`](https://golang.org/pkg/math/#NaN);
- [`Nil`](https://godoc.org/github.com/maxatome/go-testdeep#Nil)
  compares to `nil`;
- [`None`](https://godoc.org/github.com/maxatome/go-testdeep#None)
  no values have to match;
- [`NotAny`](https://godoc.org/github.com/maxatome/go-testdeep#NotAny)
  compares the contents of an array or a slice, no values have to match;
- [`Not`](https://godoc.org/github.com/maxatome/go-testdeep#Not)
  value must not match;
- [`NotEmpty`](https://godoc.org/github.com/maxatome/go-testdeep#NotEmpty)
  checks that an array, a channel, a map, a slice or a string is not empty;
- [`NotNaN`](https://godoc.org/github.com/maxatome/go-testdeep#NotNaN)
  checks a floating number is not
  [`math.NaN`](https://golang.org/pkg/math/#NaN);
- [`NotNil`](https://godoc.org/github.com/maxatome/go-testdeep#NotNil)
  checks that data is not `nil`;
- [`NotZero`](https://godoc.org/github.com/maxatome/go-testdeep#NotZero)
  checks that data is not zero regarding its type;
- [`PPtr`](https://godoc.org/github.com/maxatome/go-testdeep#PPtr)
  allows to easily test a pointer of pointer value,
- [`Ptr`](https://godoc.org/github.com/maxatome/go-testdeep#Ptr)
  allows to easily test a pointer value,
- [`Re`](https://godoc.org/github.com/maxatome/go-testdeep#Re) allows
  to apply a regexp on a string (or convertible), `[]byte`,
  [`error`](https://golang.org/ref/spec#Errors) or
  [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interfaces, and even
  test the captured groups;
- [`ReAll`](https://godoc.org/github.com/maxatome/go-testdeep#ReAll) allows
  to successively apply a regexp on a string (or convertible), `[]byte`,
  [`error`](https://golang.org/ref/spec#Errors) or
  [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interfaces, and even
  test the captured groups;
- [`Set`](https://godoc.org/github.com/maxatome/go-testdeep#Set)
  compares the contents of an array or a slice ignoring duplicates and
  without taking care of the order of items;
- [`Shallow`](https://godoc.org/github.com/maxatome/go-testdeep#Shallow)
  compares pointers only, not their contents;
- [`Slice`](https://godoc.org/github.com/maxatome/go-testdeep#Slice)
  compares the contents of a slice or a pointer on a slice;
- [`Smuggle`](https://godoc.org/github.com/maxatome/go-testdeep#Smuggle)
  changes data contents or mutates it into another type via a custom
  function before stepping down in favor of generic comparison process;
- [`String`](https://godoc.org/github.com/maxatome/go-testdeep#String)
  checks a string, [`error`](https://golang.org/ref/spec#Errors) or
  [`fmt.Stringer`](https://golang.org/pkg/fmt/#Stringer) interfaces
  string contents;
- [`Struct`](https://godoc.org/github.com/maxatome/go-testdeep#Struct)
  compares the contents of a struct or a pointer on a struct;
- [`SubBagOf`](https://godoc.org/github.com/maxatome/go-testdeep#SubBagOf)
  compares the contents of an array or a slice without taking care of the order
  of items but with potentially some exclusions;
- [`SubMapOf`](https://godoc.org/github.com/maxatome/go-testdeep#SubMapOf)
  compares the contents of a map but with potentially some exclusions;
- [`SubSetOf`](https://godoc.org/github.com/maxatome/go-testdeep#SubSetOf)
  compares the contents of an array or a slice ignoring duplicates and
  without taking care of the order of items but with potentially some exclusions;
- [`SuperBagOf`](https://godoc.org/github.com/maxatome/go-testdeep#SuperBagOf)
  compares the contents of an array or a slice without taking care of the order
  of items but with potentially some extra items;
- [`SuperMapOf`](https://godoc.org/github.com/maxatome/go-testdeep#SuperMapOf)
  compares the contents of a map but with potentially some extra entries;
- [`SuperSetOf`](https://godoc.org/github.com/maxatome/go-testdeep#SuperSetOf)
  compares the contents of an array or a slice ignoring duplicates and
  without taking care of the order of items but with potentially some extra
  items;
- [`TruncTime`](https://godoc.org/github.com/maxatome/go-testdeep#TruncTime)
  compares time.Time (or assignable) values after truncating them;
- [`Zero`](https://godoc.org/github.com/maxatome/go-testdeep#Zero)
  checks data against its zero'ed conterpart.


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
