---
title: "testdeep.T"
---

### Constructing [`*testdeep.T`]

```go
import (
  "testing"
  "github.com/maxatome/go-testdeep"
)

func TestMyFunc(tt *testing.T) {
  t := testdeep.NewT(tt)
  t.Cmp(MyFunc(), 12)
}
```

- [`func NewT(t TestingFT, config ...ContextConfig) *T`](https://godoc.org/github.com/maxatome/go-testdeep#NewT)


### Configuring [`*testdeep.T`]

```go
func TestMyFunc(tt *testing.T) {
  t := testdeep.NewT(tt).UseEqual().RootName("RECORD")
  ...
}
```

- [`func (t *T) BeLax(enable ...bool) *T`](https://godoc.org/github.com/maxatome/go-testdeep#T.BeLax)
- [`func (t *T) FailureIsFatal(enable ...bool) *T`](https://godoc.org/github.com/maxatome/go-testdeep#T.FailureIsFatal)
- [`func (t *T) RootName(rootName string) *T`](https://godoc.org/github.com/maxatome/go-testdeep#T.RootName)
- [`func (t *T) UseEqual(enable ...bool) *T`](https://godoc.org/github.com/maxatome/go-testdeep#T.UseEqual)


### Main methods of [`*testdeep.T`]

```go
import (
  "testing"
  "github.com/maxatome/go-testdeep"
)

func TestMyFunc(tt *testing.T) {
  t := testdeep.NewT(tt).UseEqual()

  // Compares MyFunc() result against a fixed value
  t.Cmp(MyFunc(), 128, "MyFunc() result is 128")

  // Compares MyFunc() result using the Between Testdeep operator
  t.Cmp(MyFunc(), testdeep.Between(100, 199),
    "MyFunc() result is between 100 and 199")
}
```

- [`func (t *T) Cmp(got, expected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Cmp)
- [`func (t *T) CmpError(got error, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpError)
- [`func (t *T) CmpLax(got interface{}, expected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpLax)
  (in fact the shortcut of [`Lax` operator]({{< ref "operators/Lax" >}}))
- [`func (t *T) CmpNoError(got error, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpNoError)
- [`func (t *T) CmpNotPanic(fn func(), args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpNotPanic)
- [`func (t *T) CmpPanic(fn func(), expected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpPanic)
- [`func (t *T) False(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.False)
- [`func (t *T) Not(got interface{}, notExpected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Not)
  (in fact the shortcut of [`Not` operator]({{< ref "operators/Not" >}}))
- [`func (t *T) RunT(name string, f func(t *T)) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.RunT)
- [`func (t *T) True(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.True)

[`CmpDeeply()`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpDeeply)
method is now replaced by
[`Cmp()`](https://godoc.org/github.com/maxatome/go-testdeep#T.Cmp),
but it is still available for backward compatibility purpose.


### Shortcut methods of [`*testdeep.T`]

```go
import (
  "testing"
  "github.com/maxatome/go-testdeep"
)

func TestMyFunc(tt *testing.T) {
  t := testdeep.NewT(tt).UseEqual()
  t.Between(MyFunc(), 100, 199, testdeep.BoundsInIn,
    "MyFunc() result is between 100 and 199")
}
```

For each of these methods, it is always a shortcut on
[`T.Cmp()`](https://godoc.org/github.com/maxatome/go-testdeep#T.Cmp) and
the correponding [Testdeep operator]({{< ref "operators" >}}):

```
T.HasPrefix(got, expected, …) ⇒ T.Cmp(t, got, HasPrefix(expected), …)
  ^-------^                                   ^-------^
      +-------------------------------------------+
```

Excluding [`Lax` operator]({{< ref "operators/Lax" >}}) for which the
shortcut method stays [`CmpLax`]({{< ref "operators/Lax#cmplax-shortcut" >}}).

Each shortcut method is described in the corresponding operator
page. See [operators list]({{< ref "operators" >}}).


[`testdeep.T`]: https://godoc.org/github.com/maxatome/go-testdeep#T
[`*testdeep.T`]: https://godoc.org/github.com/maxatome/go-testdeep#T
