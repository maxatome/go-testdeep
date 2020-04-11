+++
title = "Functions"
weight = 13
+++

```go
import (
  "testing"
  "github.com/maxatome/go-testdeep/td"
)

func TestMyFunc(t *testing.T) {
  // Compares MyFunc() result against a fixed value
  td.Cmp(t, MyFunc(), 128, "MyFunc() result is 128")

  // Compares MyFunc() result using the Between Testdeep operator
  td.Cmp(t, MyFunc(), td.Between(100, 199),
    "MyFunc() result is between 100 and 199")
}
```

- [`func Cmp(t TestingT, got, expected interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp)
- [`func CmpError(t TestingT, got error, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpError)
- [`func CmpFalse(t TestingT, got interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpFalse)
- [`func CmpLax(t TestingT, got interface{}, expected interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpLax)
  (in fact the shortcut of [`Lax` operator]({{< ref "operators/Lax" >}}))
- [`func CmpNoError(t TestingT, got error, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNoError)
- [`func CmpNot(t TestingT, got interface{}, notExpected interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNot)
  (in fact the shortcut of [`Not` operator]({{< ref "operators/Not" >}}))
- [`func CmpNotPanic(t TestingT, fn func(), args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpNotPanic)
- [`func CmpPanic(t TestingT, fn func(), expectedPanic interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpPanic)
- [`func CmpTrue(t TestingT, got interface{}, args ...interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpTrue)
- [`func EqDeeply(got, expected interface{}) bool`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#EqDeeply)
- [`func EqDeeplyError(got, expected interface{}) error`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#EqDeeplyError)

[`CmpDeeply()`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#CmpDeeply)
is now replaced by
[`Cmp()`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp), but it
is still available for backward compatibility purpose.


## Main shortcut functions

```go
import (
  "testing"
  "github.com/maxatome/go-testdeep/td"
)

func TestMyFunc(t *testing.T) {
  td.CmpBetween(t, MyFunc(), 100, 199, td.BoundsInIn,
    "MyFunc() result is between 100 and 199")
}
```

For each of these functions, it is always a shortcut on
[`Cmp()`](https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp) and
the correponding [Testdeep operator]({{< ref "operators" >}}):

```
CmpHasPrefix(t, got, expected, …) ⇒ Cmp(t, got, HasPrefix(expected), …)
   ^-------^                                    ^-------^
       +--------------------------------------------+
```
Each shortcut method is described in the corresponding operator
page. See [operators list]({{< ref "operators" >}}).
