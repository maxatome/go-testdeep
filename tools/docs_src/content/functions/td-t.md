---
title: "td.T"
---

### Constructing [`*td.T`]

```go
import (
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestMyFunc(tt *testing.T) {
  t := td.NewT(tt)
  t.Cmp(MyFunc(), 12)
}
```

- [`func NewT(t TestingFT, config ...ContextConfig) *T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#NewT)
- [`func Assert(t TestingFT, config ...ContextConfig) *T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Assert)
- [`func Require(t TestingFT, config ...ContextConfig) *T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Require)
- [`func AssertRequire(t TestingFT, config ...ContextConfig) (*T, *T)`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#AssertRequire)


### Configuring [`*td.T`]

```go
func TestMyFunc(tt *testing.T) {
  t := td.NewT(tt).UseEqual().RootName("RECORD")
  ...
}
```

- [`func (t *T) BeLax(enable ...bool) *T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.BeLax)
- [`func (t *T) FailureIsFatal(enable ...bool) *T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.FailureIsFatal)
- [`func (t *T) RootName(rootName string) *T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.RootName)
- [`func (t *T) UseEqual(enable ...bool) *T`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.UseEqual)


### Main methods of [`*td.T`]

```go
import (
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestMyFunc(tt *testing.T) {
  t := td.NewT(tt).UseEqual()

  // Compares MyFunc() result against a fixed value
  t.Cmp(MyFunc(), 128, "MyFunc() result is 128")

  // Compares MyFunc() result using the Between Testdeep operator
  t.Cmp(MyFunc(), td.Between(100, 199),
    "MyFunc() result is between 100 and 199")
}
```

- [`func (t *T) Cmp(got, expected interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Cmp)
- [`func (t *T) CmpError(got error, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.CmpError)
- [`func (t *T) CmpLax(got interface{}, expected interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.CmpLax)
  (in fact the shortcut of [`Lax` operator]({{< ref "operators/Lax" >}}))
- [`func (t *T) CmpNoError(got error, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.CmpNoError)
- [`func (t *T) CmpNotPanic(fn func(), args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.CmpNotPanic)
- [`func (t *T) CmpPanic(fn func(), expected interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.CmpPanic)
- [`func (t *T) False(got interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.False)
- [`func (t *T) Not(got interface{}, notExpected interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Not)
  (in fact the shortcut of [`Not` operator]({{< ref "operators/Not" >}}))
- [`func (t *T) RunT(name string, f func(t *T)) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.RunT)
- [`func (t *T) True(got interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.True)

[`CmpDeeply()`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.CmpDeeply)
method is now replaced by
[`Cmp()`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Cmp),
but it is still available for backward compatibility purpose.


### Anchoring methods of [`*td.T`]

- [`func (t *T) A(operator TestDeep, model ...interface{}) interface{}`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.A)
- [`func (t *T) Anchor(operator TestDeep, model ...interface{}) interface{}`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Anchor)
- [`func (t *T) AnchorsPersistTemporarily() func()`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.AnchorsPersistTemporarily)
- [`func (t *T) DoAnchorsPersist() bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.DoAnchorsPersist)
- [`func (t *T) ResetAnchors()`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.ResetAnchors)
- [`func (t *T) SetAnchorsPersist(persist bool)`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.SetAnchorsPersist)


### Shortcut methods of [`*td.T`]

```go
import (
  "testing"

  "github.com/maxatome/go-testdeep/td"
)

func TestMyFunc(tt *testing.T) {
  t := td.NewT(tt).UseEqual()
  t.Between(MyFunc(), 100, 199, td.BoundsInIn,
    "MyFunc() result is between 100 and 199")
}
```

For each of these methods, it is always a shortcut on
[`T.Cmp()`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Cmp) and
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


[`td.T`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T
[`*td.T`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T
