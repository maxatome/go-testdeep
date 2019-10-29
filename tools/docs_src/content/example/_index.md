+++
title = "Example"
weight = 11
+++

Imagine a function returning a struct containing a newly created
database record. The `Id` and the `CreatedAt` fields are set by the
database layer:

```go
type Record struct {
  Id        uint64
  Name      string
  Age       int
  CreatedAt time.Time
}

func CreateRecord(name string, age int) (*Record, error) {
  // Do INSERT INTO … and return newly created record or error if it failed
}
```
