go-testdeep
===========

[![Build Status](https://travis-ci.org/maxatome/go-testdeep.svg)](https://travis-ci.org/maxatome/go-testdeep)
[![Coverage Status](https://coveralls.io/repos/github/maxatome/go-testdeep/badge.svg?branch=master)](https://coveralls.io/github/maxatome/go-testdeep?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxatome/go-testdeep)](https://goreportcard.com/report/github.com/maxatome/go-testdeep)
[![GoDoc](https://godoc.org/github.com/maxatome/go-testdeep?status.svg)](https://godoc.org/github.com/maxatome/go-testdeep)


## Synopsis

```go
import (
  "testing"
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

  if td.CmpNil(t, err) {
    td.CmpStruct(t, record,
      Record{
        Name: "Bob",
        Age:  23,
      },
      StructFields{
        "Id":        td.Not(uint64(0)),
        "CreatedAt": td.Between(before, time.Now()),
      },
      "Newly created record")
  }
}
```

Imagine `CreateRecord` does not set correctly `CreatedAt` field, then:
```sh
go test -run=TestCreateRecord
```

outputs:
```
--- FAIL: TestCreateRecord (0.00s)
  test_test.go:22: Failed test 'Newly created record'
    DATA.CreatedAt: values differ
           got: 2018-05-27 10:55:50.788166932 +0200 CEST m=-2.998149554
      expected: 2018-05-27 10:55:53.788163509 +0200 CEST m=+0.001848002 ≤ got ≤ 2018-05-27 10:55:53.788464176 +0200 CEST m=+0.002148179
    [under TestDeep operator Between at test_test.go:29]
FAIL
exit status 1
FAIL  github.com/maxatome/go-testdeep  0.006s
```

If `CreateRecord` did not set correctly `Id` field, output would be:
```
--- FAIL: TestCreateRecord (0.00s)
  test_test.go:22: Failed test 'Newly created record'
    DATA.Id: comparing with Not
           got: (uint64) 0
      expected: Not((uint64) 0)
    [under TestDeep operator Not at test_test.go:28]
FAIL
exit status 1
FAIL  github.com/maxatome/go-testdeep  0.006s
```

If `CreateRecord` set `Name` field to "Alice" value instead of
expected "Bob", output would be:
```
--- FAIL: TestCreateRecord (0.00s)
  test_test.go:22: Failed test 'Newly created record'
    DATA.Name: values differ
           got: (string) (len=5) "Alice"
      expected: (string) (len=3) "Bob"
    [called by CmpStruct at td_between_test.go:37]
FAIL
exit status 1
FAIL  github.com/maxatome/go-testdeep  0.006s
```

## Installation

```sh
$ go get github.com/maxatome/go-testdeep
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
check the record content:

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

With `testdeep`, it is a way simple, thanks to `CmpDeeply` function:

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
      Struct(
        Record{
          Name: "Bob",
          Age:  23,
        },
        StructFields{
          "Id":        td.Not(uint64(0)),
          "CreatedAt": td.Between(before, time.Now()),
        }),
      "Newly created record")
  }
}
```

Of course not only structs can be compared. A lot of operators can
be found below to cover most (all?) needed tests.

The `CmpDeeply` function is the keystone of this package, but to make
the writing of tests even easier, the family of `Cmp*` functions are
provided and act as shortcuts. Using `CmpNil` and `CmpStruct`
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

  if td.CmpNil(t, err) {
    td.CmpStruct(t, record,
      Record{
        Name: "Bob",
        Age:  23,
      },
      StructFields{
        "Id":        td.Not(uint64(0)),
        "CreatedAt": td.Between(before, time.Now()),
      },
      "Newly created record")
  }
}
```

## License

`go-testdeep` is released under the BSD-style license found in the
`LICENSE` file in the root directory of this source tree.

Internal function `deepValueEqual` is based on `deepValueEqual` from
[`reflect` golang package](https://golang.org/pkg/reflect/) licensed
under the BSD-style license found in the [`LICENSE` file in the golang
repository](https://github.com/golang/go/blob/master/LICENSE).

Uses two files (`bypass.go` & `bypasssafe.go`) from
[Go-spew](https://github.com/davecgh/go-spew) which is licensed under
the [copyfree](http://copyfree.org) ISC License.
