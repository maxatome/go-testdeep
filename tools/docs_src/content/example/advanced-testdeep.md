+++
title = "Advanced go-testdeep technique"
weight = 30
+++

Of course we can test struct fields one by one, but with go-testdeep
and the [`td` package], the whole struct can be compared with one
[`Cmp`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp) call.

We can choose to ignore the non-guessable fields set by
`CreateRecord()`:

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
    td.Cmp(t, record,
      td.Struct(
        &Record{
          Name: "Bob",
          Age:  23,
        },
        nil),
      "Newly created record")
  }
}
```

The [`Struct`]({{< ref "operators/Struct" >}}) operator, used here,
ignores zero fields in its model parameter.

But it is better to check all fields:

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
    td.Cmp(t, record,
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

See the use of the [`Struct`]({{< ref "operators/Struct" >}})
operator. It is needed here to overcome the go static typing system
and so use other [go-testdeep operators]({{< ref "operators" >}})
for some fields, here [`NotZero`]({{< ref "operators/NotZero" >}}) for
`Id` and [`Between`]({{< ref "operators/Between" >}}) for `CreatedAt`.

Not only structs can be compared. A lot of operators can be
found to cover most (all?) needed tests. See the
[operators list]({{< ref "operators" >}}).

Say `CreateRecord()` does not set correctly `CreatedAt` field, then:
```sh
go test -run=TestCreateRecord
```

outputs for last `td.Cmp` call:

![error output](/images/colored-newly1.svg)

If `CreateRecord()` had not set correctly `Id` field, output would have
been:

![error output](/images/colored-newly2.svg)

If `CreateRecord()` had set `Name` field to "Alice" value instead of
expected "Bob", output would have been:

![error output](/images/colored-newly3.svg)
