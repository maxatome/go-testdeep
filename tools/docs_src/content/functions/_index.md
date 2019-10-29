+++
title = "Functions"
weight = 13
+++

```go
import (
  "testing"
  "github.com/maxatome/go-testdeep"
)

func TestMyFunc(t *testing.T) {
  // Compares MyFunc() result against a fixed value
  testdeep.Cmp(t, MyFunc(), 128, "MyFunc() result is 128")

  // Compares MyFunc() result using the Between Testdeep operator
  testdeep.Cmp(t, MyFunc(), testdeep.Between(100, 199),
    "MyFunc() result is between 100 and 199")
}
```

- [`func Cmp(t TestingT, got, expected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#Cmp)
- [`func CmpError(t TestingT, got error, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpError)
- [`func CmpFalse(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpFalse)
- [`func CmpLax(t TestingT, got interface{}, expected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpLax)
  (in fact the shortcut of [`Lax` operator]({{< ref "operators/Lax" >}}))
- [`func CmpNoError(t TestingT, got error, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNoError)
- [`func CmpNot(t TestingT, got interface{}, notExpected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNot)
  (in fact the shortcut of [`Not` operator]({{< ref "operators/Not" >}}))
- [`func CmpNotPanic(t TestingT, fn func(), args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNotPanic)
- [`func CmpPanic(t TestingT, fn func(), expectedPanic interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpPanic)
- [`func CmpTrue(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpTrue)
- [`func EqDeeply(got, expected interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#EqDeeply)
- [`func EqDeeplyError(got, expected interface{}) error`](https://godoc.org/github.com/maxatome/go-testdeep#EqDeeplyError)

[`CmpDeeply()`](https://godoc.org/github.com/maxatome/go-testdeep#CmpDeeply)
is now replaced by
[`Cmp()`](https://godoc.org/github.com/maxatome/go-testdeep#Cmp), but it
is still available for backward compatibility purpose.


## Main shortcut functions

```go
import (
  "testing"
  "github.com/maxatome/go-testdeep"
)

func TestMyFunc(t *testing.T) {
  testdeep.CmpBetween(t, MyFunc(), 100, 199, testdeep.BoundsInIn,
    "MyFunc() result is between 100 and 199")
}
```

For each of these functions, it is always a shortcut on
[`Cmp()`](https://godoc.org/github.com/maxatome/go-testdeep#Cmp) and
the correponding [Testdeep operator]({{< ref "operators" >}}):

```
CmpHasPrefix(t, got, expected, …) ⇒ Cmp(t, got, HasPrefix(expected), …)
   ^-------^                                    ^-------^
       +--------------------------------------------+
```
Each shortcut method is described in the corresponding operator
page. See [operators list]({{< ref "operators" >}}).
