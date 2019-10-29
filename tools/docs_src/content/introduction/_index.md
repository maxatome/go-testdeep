+++
title = "Introduction"
weight = 5
+++

## Synopsis

Simplest usage:

```go
import (
  "testing"
  td "github.com/maxatome/go-testdeep"
)

func TestMyFunc(t *testing.T) {
  td.Cmp(t, MyFunc(), &Info{Name: "Alice", Age: 42})
}
```

Example of produced error in case of mismatch:

![error output](/images/colored-output.svg)


## Description

go-testdeep is a go rewrite and adaptation of wonderful
[Test::Deep perl](https://metacpan.org/pod/Test::Deep).

In golang, comparing data structure is usually done using
[reflect.DeepEqual](https://golang.org/pkg/reflect/#DeepEqual) or
using a package that uses this function behind the scene.

This function works very well, but it is not flexible. Both compared
structures must match exactly and when a difference is returned, it is
up to the caller to display it. Not easy when comparing big data
structures.

The purpose of testdeep package is to do its best to introduce this
missing flexibility using ["operators"]({{< ref "operators" >}}), when
the expected value (or one of its component) cannot be matched
exactly, mixed with some useful
[comparison functions]({{< ref "functions" >}}).
