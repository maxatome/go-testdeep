+++
title = "Using anchoring"
weight = 60
+++

Last, operators can directly be anchored in litterals, still using the
[`td.T` type]({{< ref "functions/td-t" >}}), avoiding the use of the
[`Struct`]({{< ref "Struct" >}}) operator:

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
    t.RootName("RECORD"). // Use RECORD instead of DATA in failure reports
      Cmp(record,
        &Record{
          Name:      "Bob",
          Age:       23,
          Id:        t.Anchor(td.NotZero(), uint64(0)).(uint64),
          CreatedAt: t.Anchor(td.Between(before, time.Now())).(time.Time),
        },
        "Newly created record")
  }
}
```

See the
[`Anchor`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Anchor)
method documentation for details. Note that
[`A`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.A) method
is also a synonym for Anchor.

```go
          Id:        t.A(td.NotZero(), uint64(0)).(uint64),
          CreatedAt: t.A(td.Between(before, time.Now())).(time.Time),
```
