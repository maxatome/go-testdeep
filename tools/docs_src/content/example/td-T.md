+++
title = "td.T type"
weight = 50
+++

[`testing.T`](https://golang.org/pkg/testing/#T) can also be
encapsulated in [`td.T` type]({{< ref "functions/td-t" >}}),
simplifying again the test:

```go
import (
  "testing"
  "time"

  "github.com/maxatome/go-testdeep/td"
)

func TestCreateRecord(tt *testing.T) {
  t := td.NewT(tt)

  before := time.Now().Truncate(time.Second)
  record, err := CreateRecord()

  if t.CmpNoError(err) {
    t := t.RootName("RECORD") // Use RECORD instead of DATA in failure reports

    // Using Struct shortcut method
    t.Struct(record,
      &Record{
        Name: "Bob",
        Age:  23,
      },
      td.StructFields{
        "Id":        td.NotZero(),
        "CreatedAt": td.Between(before, time.Now()),
      },
      "Newly created record")

    // Or using Cmp method, it's a matter of taste
    t.Cmp(record,
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

Note the use of
[`RootName`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.RootName)
method, it allows to name what we are going to test, instead of the
default "DATA".

If `CreateRecord()` had set `Name` field to "Alice" value instead of
expected "Bob", output would have been (note "RECORD" replaced default
"DATA"):

![error output](/images/colored-newly4.svg)
