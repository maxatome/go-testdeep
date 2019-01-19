go-testdeep
===========

[![Build Status](https://travis-ci.org/maxatome/go-testdeep.svg?branch=master)](https://travis-ci.org/maxatome/go-testdeep)
[![Coverage Status](https://coveralls.io/repos/github/maxatome/go-testdeep/badge.svg?branch=master)](https://coveralls.io/github/maxatome/go-testdeep?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxatome/go-testdeep)](https://goreportcard.com/report/github.com/maxatome/go-testdeep)
[![GoDoc](https://godoc.org/github.com/maxatome/go-testdeep?status.svg)](https://godoc.org/github.com/maxatome/go-testdeep)
[![Version](https://img.shields.io/github/tag/maxatome/go-testdeep.svg)](https://github.com/maxatome/go-testdeep/releases)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go/#testing)

![testdeep](doc/image.png)

**Extremely flexible golang deep comparison, extends the go testing package.**

- [Latest new](#latest-news)
- [Synopsis](#synopsis)
- [Installation](#installation)
- [Presentation](#presentation)
- [Available operators](#available-operators)
- [Helpers](#helpers)
  - [`tdhttp` or HTTP API testing helper](#tdhttp-or-http-api-testing-helper)
- [Environment variables](#environment-variables)
- [Operators vs go types](#operators-vs-go-types)
- [See also](#see-also)
- [License](#license)
- [FAQ](doc/FAQ.md)


## Latest news

- 2019/01/13: test failures output is now colored by default. See
  [Environment variables](#environment-variables) to configure it;
- 2019/01/07: introducing TestDeep helpers. First one is
  [`tdhttp` or HTTP API testing helper](#tdhttp-or-http-api-testing-helper);
- 2018/12/16: [`Between`], [`Gt`], [`Gte`], [`Lt`] & [`Lte`] operators
  now handle strings as well;
- 2018/12/07: [`Smuggle`] operator and its friends
  [`CmpSmuggle`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSmuggle)
  &
  [`T.Smuggle`](https://godoc.org/github.com/maxatome/go-testdeep#T.Smuggle)
  are now more lax and gain the ability to work on struct
  fields-paths.  [`Shallow`] operator and its friends
  [`CmpShallow`](https://godoc.org/github.com/maxatome/go-testdeep#CmpShallow)
  &
  [`T.Shallow`](https://godoc.org/github.com/maxatome/go-testdeep#T.Shallow)
  can now work on strings;
- see [commits history](https://github.com/maxatome/go-testdeep/commits/master)
  for other/older changes.


## Synopsis

![error output](doc/colored-output.svg)

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

Using [`testdeep.T`][`T`]
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

The purpose of `testdeep` package is to do its best to introduce this
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

With `testdeep`, it is a way simple, thanks to [`CmpDeeply`] function:

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

The [`CmpDeeply`] function is the keystone of this package,
but to make the writing of tests even easier, the family of `Cmp*`
functions are provided and act as shortcuts. Using
[`CmpNoError`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNoError)
and
[`CmpStruct`](https://godoc.org/github.com/maxatome/go-testdeep#CmpStruct)
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

Last, [`testing.T`](https://golang.org/pkg/testing/#T) can be
encapsulated in [`T`] type, simplifying again the test:

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

See functions returning [`TestDeep` interface][`TestDeep`]:

- [`All`] all expected values have to match;
- [`Any`] at least one expected value have to match;
- [`Array`] compares the contents of an array or a pointer on an
  array;
- [`ArrayEach`] compares each array or slice item;
- [`Bag`] compares the contents of an array or a slice without taking
  care of the order of items;
- [`Between`] checks that a number, string or [`time.Time`] is between two
  bounds;
- [`Cap`] checks an array, slice or channel capacity;
- [`Code`] allows to use a custom function;
- [`Contains`] checks that a string, [`error`] or [`fmt.Stringer`]
  interfaces contain a sub-string; or an array, slice or map contain a
  value;
- [`ContainsKey`] checks that a map contains a key;
- [`Empty`] checks that an array, a channel, a map, a slice or a
  string is empty;
- [`Gt`] checks that a number, string or [`time.Time`] is greater than a
  value;
- [`Gte`] checks that a number, string or [`time.Time`] is greater or equal
  than a value;
- [`HasPrefix`] checks the prefix of a string, [`error`] or
  [`fmt.Stringer`] interfaces;
- [`HasSuffix`] checks the suffix of a string, [`error`] or
  [`fmt.Stringer`] interfaces;
- [`Ignore`] allows to ignore a comparison;
- [`Isa`] checks the data type or whether data implements an interface
  or not;
- [`Len`] checks an array, slice, map, string or channel length;
- [`Lt`] checks that a number, string or [`time.Time`] is lesser than a value;
- [`Lte`] checks that a number, string or [`time.Time`] is lesser or equal
  than a value;
- [`Map`] compares the contents of a map;
- [`MapEach`] compares each map entry;
- [`N`] compares a number with a tolerance value;
- [`NaN`] checks a floating number is [`math.NaN`];
- [`Nil`] compares to `nil`;
- [`None`] no values have to match;
- [`NotAny`] compares the contents of an array or a slice, no values
  have to match;
- [`Not`] value must not match;
- [`NotEmpty`] checks that an array, a channel, a map, a slice or a
  string is not empty;
- [`NotNaN`] checks a floating number is not [`math.NaN`];
- [`NotNil`] checks that data is not `nil`;
- [`NotZero`] checks that data is not zero regarding its type;
- [`PPtr`] allows to easily test a pointer of pointer value,
- [`Ptr`] allows to easily test a pointer value,
- [`Re`] allows to apply a regexp on a string (or convertible),
  `[]byte`, [`error`] or [`fmt.Stringer`] interfaces, and even test
  the captured groups;
- [`ReAll`] allows to successively apply a regexp on a string (or
  convertible), `[]byte`, [`error`] or [`fmt.Stringer`] interfaces,
  and even test the captured groups;
- [`Set`] compares the contents of an array or a slice ignoring
  duplicates and without taking care of the order of items;
- [`Shallow`] compares pointers only, not their contents;
- [`Slice`] compares the contents of a slice or a pointer on a slice;
- [`Smuggle`] changes data contents or mutates it into another type
  via a custom function or a struct fields-path before stepping down
  in favor of generic comparison process;
- [`String`] checks a string, [`error`] or [`fmt.Stringer`] interfaces
  string contents;
- [`Struct`] compares the contents of a struct or a pointer on a
  struct;
- [`SubBagOf`] compares the contents of an array or a slice without
  taking care of the order of items but with potentially some
  exclusions;
- [`SubMapOf`] compares the contents of a map but with potentially
  some exclusions;
- [`SubSetOf`] compares the contents of an array or a slice ignoring
  duplicates and without taking care of the order of items but with
  potentially some exclusions;
- [`SuperBagOf`] compares the contents of an array or a slice without
  taking care of the order of items but with potentially some extra
  items;
- [`SuperMapOf`] compares the contents of a map but with potentially
  some extra entries;
- [`SuperSetOf`] compares the contents of an array or a slice ignoring
  duplicates and without taking care of the order of items but with
  potentially some extra items;
- [`TruncTime`] compares time.Time (or assignable) values after
  truncating them;
- [`Zero`] checks data against its zero'ed conterpart.


## Helpers

The goal of helpers is to make use of `go-testdeep` even more powerful
by providing common features using [TestDeep operators](#available-operators)
behind the scene.

### `tdhttp` or HTTP API testing helper

The package `github.com/maxatome/go-testdeep/helpers/tdhttp` provides
some functions to easily test HTTP handlers.

See [`tdhttp`] documentation for details or
[FAQ](doc/FAQ.md#what-about-testing-the-response-using-my-api) for an
example of use.


## Environment variables

- `TESTDEEP_MAX_ERRORS` maximum number of errors to report before
  stopping during one comparison (one [`CmpDeeply`] execution for
  example). It defaults to 10;
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

```
TESTDEEP_COLOR_OK=black:green \
    TESTDEEP_COLOR_BAD=white:red \
    TESTDEEP_COLOR_TITLE=yellow \
    go test
```


## Operators vs go types

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`All`]             | ✓ | ✓ | ✓ | ✓ | ✓ | ✓    | ✓ | ✓ | ✓ | ✓             | ✓                  | ✓ | ✓ | ✓ | [`All`] |
| [`Any`]             | ✓ | ✓ | ✓ | ✓ | ✓ | ✓    | ✓ | ✓ | ✓ | ✓             | ✓                  | ✓ | ✓ | ✓ | [`Any`] |
| [`Array`]           | ✗ | ✗ | ✗ | ✗ | ✗ | ✗    | ✓ | ✗ | ✗ | ✗             | ptr on array       | ✓ | ✗ | ✗ | [`Array`] |
| [`ArrayEach`]       | ✗ | ✗ | ✗ | ✗ | ✗ | ✗    | ✓ | ✓ | ✗ | ✗             | ptr on array/slice | ✓ | ✗ | ✗ | [`ArrayEach`] |
| [`Bag`]             | ✗ | ✗ | ✗ | ✗ | ✗ | ✗    | ✓ | ✓ | ✗ | ✗             | ptr on array/slice | ✓ | ✗ | ✗ | [`Bag`] |
| [`Between`]         | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗                  | ✓ | ✗ | ✗ | [`Between`] |
| [`Cap`]             | ✗ | ✗ | ✗ | ✗ | ✗ | ✗    | ✓ | ✓ | ✗ | ✗             | ✗                  | ✓ | ✓ | ✗ | [`Cap`] |
| [`Code`]            | ✓ | ✓ | ✓ | ✓ | ✓ | ✓    | ✓ | ✓ | ✓ | ✓             | ✓                  | ✓ | ✓ | ✓ | [`Code`] |
| [`Contains`]        | ✗ | ✗ | ✓ | ✗ | ✗ | ✗    | ✓ | ✓ | ✓ | ✗             | ✗                  | ✓ | ✗ | ✗ | [`Contains`] |
| [`ContainsKey`]     | ✗ | ✗ | ✗ | ✗ | ✗ | ✗    | ✗ | ✗ | ✓ | ✗             | ✗                  | ✓ | ✗ | ✗ | [`ContainsKey`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`Empty`]           | ✗ | ✗ | ✓ | ✗ | ✗ | ✗    | ✓ | ✓ | ✓ | ✗             | ptr or/array/slice/map/string | ✓ | ✓ | ✗ | [`Empty`] |
| [`Gt`]              | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗                             | ✓ | ✗ | ✗ | [`Gt`] |
| [`Gte`]             | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗                             | ✓ | ✗ | ✗ | [`Gte`] |
| [`HasPrefix`]       | ✗ | ✗ | ✓ | ✗ | ✗ | ✗    | ✗ | ✗ | ✗ | ✗             | ✗                             | ✓ + [`fmt.Stringer`], [`error`] | ✗ | ✗ | [`HasPrefix`] |
| [`HasSuffix`]       | ✗ | ✗ | ✓ | ✗ | ✗ | ✗    | ✗ | ✗ | ✗ | ✗             | ✗                             | ✓ + [`fmt.Stringer`], [`error`] | ✗ | ✗ | [`HasSuffix`] |
| [`Ignore`]          | ✓ | ✓ | ✓ | ✓ | ✓ | ✓    | ✓ | ✓ | ✓ | ✓             | ✓                             | ✓ | ✓ | ✓ | [`Ignore`] |
| [`Isa`]             | ✗ | ✓ | ✓ | ✓ | ✓ | ✓    | ✓ | ✓ | ✓ | ✓             | ✓                             | ✓ | ✓ | ✓ | [`Isa`] |
| [`Len`]             | ✗ | ✗ | ✓ | ✗ | ✗ | ✗    | ✓ | ✓ | ✓ | ✗             | ✗                             | ✓ | ✓ | ✗ | [`Len`] |
| [`Lt`]              | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗                             | ✓ | ✗ | ✗ | [`Lt`] |
| [`Lte`]             | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗                             | ✓ | ✗ | ✗ | [`Lte`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`Map`]             | ✗ | ✗ | ✗ | ✗ | ✗ | ✗    | ✗ | ✗ | ✓ | ✗ | ptr on map                    | ✓ | ✗ | ✗ | [`Map`] |
| [`MapEach`]         | ✗ | ✗ | ✗ | ✗ | ✗ | ✗    | ✗ | ✗ | ✓ | ✗ | ptr on map                    | ✓ | ✗ | ✗ | [`MapEach`] |
| [`N`]               | ✗ | ✗ | ✗ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | ✗ | ✗                             | ✓ | ✗ | ✗ | [`N`] |
| [`NaN`]             | ✗ | ✗ | ✗ | ✗ | ✓ | ✗    | ✗ | ✗ | ✗ | ✗ | ✗                             | ✓ | ✗ | ✗ | [`NaN`] |
| [`Nil`]             | ✓ | ✗ | ✗ | ✗ | ✗ | ✗    | ✗ | ✓ | ✓ | ✗ | ✓                             | ✓ | ✓ | ✓ | [`Nil`] |
| [`None`]            | ✓ | ✓ | ✓ | ✓ | ✓ | ✓    | ✓ | ✓ | ✓ | ✓ | ✓                             | ✓ | ✓ | ✓ | [`None`] |
| [`NotAny`]          | ✗ | ✗ | ✗ | ✗ | ✗ | ✗    | ✓ | ✓ | ✗ | ✗ | ptr on array/slice            | ✓ | ✗ | ✗ | [`NotAny`] |
| [`Not`]             | ✓ | ✓ | ✓ | ✓ | ✓ | ✓    | ✓ | ✓ | ✓ | ✓ | ✓                             | ✓ | ✓ | ✓ | [`Not`] |
| [`NotEmpty`]        | ✗ | ✗ | ✓ | ✗ | ✗ | ✗    | ✓ | ✓ | ✓ | ✗ | ptr on array/slice/map/string | ✓ | ✓ | ✗ | [`NotEmpty`] |
| [`NotNaN`]          | ✗ | ✗ | ✗ | ✗ | ✓ | ✗    | ✗ | ✗ | ✗ | ✗ | ✗                             | ✓ | ✗ | ✗ | [`NotNaN`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`NotNil`]          | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓        | ✓ | ✗ | ✓                  | ✓ | ✓ | ✓ | [`NotNil`] |
| [`NotZero`]         | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓        | ✓ | ✓ | ✓                  | ✓ | ✓ | ✓ | [`NotZero`] |
| [`PPtr`]            | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗        | ✗ | ✗ | ✓                  | ✓ | ✗ | ✗ | [`PPtr`] |
| [`Ptr`]             | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗        | ✗ | ✗ | ✓                  | ✓ | ✗ | ✗ | [`Ptr`] |
| [`Re`]              | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | `[]byte` | ✗ | ✗ | ✗                  | ✓ + [`fmt.Stringer`], [`error`] | ✗ | ✗ | [`Re`] |
| [`ReAll`]           | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | `[]byte` | ✗ | ✗ | ✗                  | ✓ + [`fmt.Stringer`], [`error`] | ✗ | ✗ | [`ReAll`] |
| [`Set`]             | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓        | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`Set`] |
| [`Shallow`]         | ✓ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✓        | ✓ | ✗ | ✓                  | ✓ | ✓ | ✓ | [`Shallow`] |
| [`Slice`]           | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓        | ✗ | ✗ | ptr on slice       | ✓ | ✗ | ✗ | [`Slice`] |
| [`Smuggle`]         | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓        | ✓ | ✓ | ✓                  | ✓ | ✓ | ✓ | [`Smuggle`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`String`]          | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗             | ✗                  | ✓ + [`fmt.Stringer`], [`error`] | ✗ | ✗ | [`String`] |
| [`Struct`]          | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✓             | ptr on struct      | ✓ | ✗ | ✗ | [`Struct`] |
| [`SubBagOf`]        | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗             | ptr on array/slice | ✓ | ✗ | ✗ | [`SubBagOf`] |
| [`SubMapOf`]        | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗             | ptr on map         | ✓ | ✗ | ✗ | [`SubMapOf`] |
| [`SubSetOf`]        | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗             | ptr on array/slice | ✓ | ✗ | ✗ | [`SubSetOf`] |
| [`SuperBagOf`]      | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗             | ptr on array/slice | ✓ | ✗ | ✗ | [`SuperBagOf`] |
| [`SuperMapOf`]      | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗             | ptr on map         | ✓ | ✗ | ✗ | [`SuperMapOf`] |
| [`SuperSetOf`]      | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗             | ptr on array/slice | ✓ | ✗ | ✗ | [`SuperSetOf`] |
| [`TruncTime`]       | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | [`time.Time`] | todo               | ✓ | ✗ | ✗ | [`TruncTime`] |
| [`Zero`]            | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓             | ✓                  | ✓ | ✓ | ✓ | [`Zero`] |

Legend:
- ✗ means using this operator with this go type will always fail
- ✓ means using this operator with this go type can succeed
- `[]byte`, [`time.Time`], ptr on X, [`fmt.Stringer`], [`error`] means using
  this operator with this go type can succeed
- todo means should be implemented in future (PRs welcome :) )
- ¹ + ✓ means using this operator with the data behind the interface can succeed

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

See [FAQ](doc/FAQ.md).

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
[`Ignore`]: https://godoc.org/github.com/maxatome/go-testdeep#Isa
[`Isa`]: https://godoc.org/github.com/maxatome/go-testdeep#Isa
[`Len`]: https://godoc.org/github.com/maxatome/go-testdeep#Len
[`Lt`]: https://godoc.org/github.com/maxatome/go-testdeep#Lt
[`Lte`]: https://godoc.org/github.com/maxatome/go-testdeep#Lte
[`Map`]: https://godoc.org/github.com/maxatome/go-testdeep#Map
[`MapEach`]: https://godoc.org/github.com/maxatome/go-testdeep#MapEach
[`N`]: https://godoc.org/github.com/maxatome/go-testdeep#N
[`NaN`]: https://godoc.org/github.com/maxatome/go-testdeep#NaN
[`Nil`]: https://godoc.org/github.com/maxatome/go-testdeep#Nil
[`None`]: https://godoc.org/github.com/maxatome/go-testdeep#None
[`NotAny`]: https://godoc.org/github.com/maxatome/go-testdeep#NotAny
[`Not`]: https://godoc.org/github.com/maxatome/go-testdeep#Not
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
[`Zero`]: https://godoc.org/github.com/maxatome/go-testdeep#Zero

[`T`]: https://godoc.org/github.com/maxatome/go-testdeep#T
[`TestDeep`]: https://godoc.org/github.com/maxatome/go-testdeep#TestDeep
[`CmpDeeply`]: https://godoc.org/github.com/maxatome/go-testdeep#CmpDeeply

[`error`]: https://golang.org/ref/spec#Errors
[`fmt.Stringer`]: https://golang.org/pkg/fmt/#Stringer
[`time.Time`]: https://golang.org/pkg/time/
[`math.NaN`]: https://golang.org/pkg/math/#NaN

[`tdhttp`]: https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp
