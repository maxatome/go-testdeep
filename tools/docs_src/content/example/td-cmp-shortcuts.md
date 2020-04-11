+++
title = "go-testdeep Cmp shortcuts"
weight = 40
+++

The [`Cmp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp)
function is the keystone of this package, but to make the writing of
tests even easier, the family of [`Cmp*`]({{< ref "functions" >}})
functions are provided and act as shortcuts. Using
[`CmpStruct`]({{< ref "operators/Struct#cmpstruct-shortcut" >}})
function, the previous example can be written as:

```go
import (
  "testing"
  "time"

  "github.com/maxatome/go-testdeep/td"
)

func TestCreateRecord(t *testing.T) {
  before := time.Now().Truncate(time.Second)
  record, err := CreateRecord()

  if td.CmpNoError(t, err) {
    td.CmpStruct(t, record,
      &Record{
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
