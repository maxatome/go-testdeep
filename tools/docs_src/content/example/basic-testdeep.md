+++
title = "Basic go-testdeep approach"
weight = 20
+++

[`td` package](https://pkg.go.dev/github.com/maxatome/go-testdeep/td),
via its [`Cmp*`]({{< ref "functions" >}}) functions, handles the tests
and all the error message boiler plate. Let's do it:

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
    td.Cmp(t, record.Id, td.NotZero(), "Id initialized")
    td.Cmp(t, record.Name, "Bob")
    td.Cmp(t, record.Age, 23)
    td.Cmp(t, record.CreatedAt, td.Between(before, time.Now()))
  }
}
```

As we cannot guess the `Id` field value before its creation, we use the
[`NotZero`]({{< ref "operators/NotZero" >}}) operator to check it is
set by  `CreateRecord()` call. The same is true for the creation date
field `CreatedAt`. Thanks to the [`Between`]({{< ref "operators/Between" >}})
operator we can check it is set with a value included between
the date before `CreateRecord()` call and the date just after.

Note that if `Id` and `CreateAt` could be known in advance, we could
simply do:

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
    td.Cmp(t, record, &Record{
        Id:        1234,
        Name:      "Bob",
        Age:       23,
        CreatedAt: time.Date(2019, time.May, 1, 12, 13, 14, 0, time.UTC),
      })
  }
}
```

But unfortunately, it is common to not know exactly the value of some
fields…
