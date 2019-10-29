+++
title = "Using testing package"
weight = 10
+++

To check the freshly created record contents using standard
[`testing` package](https://golang.org/pkg/testing/), we have to do
something like that:

```go
import (
  "testing"
  "time"
)

func TestCreateRecord(t *testing.T) {
  before := time.Now().Truncate(time.Second)
  record, err := CreateRecord()

  if err != nil {
    t.Errorf("An error occurred: %s", err)
  } else {
    expected := Record{Name: "Bob", Age: 23}

    if record.Id == 0 {
      t.Error("Id probably not initialized")
    }
    if record.Name != expected.Name {
      t.Errorf("Name field differs, got=%s, expected=%s",
        record.Name, expected.Name)
    }
    if record.Age != expected.Age {
      t.Errorf("Age field differs, got=%s, expected=%s",
        record.Age, expected.Age)
    }
    if before.After(record.CreatedAt) ||
      time.Now().Before(record.CreatedAt) {
      t.Errorf("CreatedAt field not expected: %s", record.CreatedAt)
    }
  }
}
```
