# Godoc table of contents

`testdeep` package provides a lot of functions and methods.
As [godoc format](https://godoc.org/github.com/maxatome/go-testdeep)
does not provide a way to tidy up all of these using sections, this
document tries to overcome.

- [Main functions](#main-functions)
- [Main shortcut functions](#main-shortcut-functions)
- [`testdeep.T`](#testdeept)
  - [Constructing `*testdeep.T`](#constructing-testdeept)
  - [Configuring `*testdeep.T`](#configuring-testdeept)
  - [Main methods of `*testdeep.T`](#main-methods-of-testdeept)
  - [Shortcut methods of `*testdeep.T`](#shortcut-methods-of-testdeept)
- [`Testdeep` operators](#testdeep-operators)


## Main functions

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
- [`func CmpNoError(t TestingT, got error, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNoError)
- [`func CmpNotPanic(t TestingT, fn func(), args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNotPanic)
- [`func CmpPanic(t TestingT, fn func(), expectedPanic interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpPanic)
- [`func CmpTrue(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpTrue)
- [`func EqDeeply(got, expected interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#EqDeeply)
- [`func EqDeeplyError(got, expected interface{}) error`](https://godoc.org/github.com/maxatome/go-testdeep#EqDeeplyError)

Note that the convenient
[`CmpNot()`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNot)
function is listed in the shortcut section below.

[`CmpDeeply()`](https://godoc.org/github.com/maxatome/go-testdeep#Cmp)
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
the correponding [Testdeep operator](#testdeep-operators):

```
CmpNot(t, got, expected, args...) ⇒ Cmp(t, got, Not(expected), args...)
   ^-^                                          ^-^
    +--------------------------------------------+
```

- [`func CmpAll(t TestingT, got interface{}, expectedValues []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpAll)
- [`func CmpAny(t TestingT, got interface{}, expectedValues []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpAny)
- [`func CmpArray(t TestingT, got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpArray)
- [`func CmpArrayEach(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpArrayEach)
- [`func CmpBag(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpBag)
- [`func CmpBetween(t TestingT, got interface{}, from interface{}, to interface{}, bounds BoundsKind, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpBetween)
- [`func CmpCap(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpCap)
- [`func CmpCode(t TestingT, got interface{}, fn interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpCode)
- [`func CmpContains(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpContains)
- [`func CmpContainsKey(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpContainsKey)
- [`func CmpEmpty(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpEmpty)
- [`func CmpGt(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpGt)
- [`func CmpGte(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpGte)
- [`func CmpHasPrefix(t TestingT, got interface{}, expected string, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpHasPrefix)
- [`func CmpHasSuffix(t TestingT, got interface{}, expected string, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpHasSuffix)
- [`func CmpIsa(t TestingT, got interface{}, model interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpIsa)
- [`func CmpKeys(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpKeys)
- [`func CmpLax(t TestingT, got interface{}, expected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpLax)
- [`func CmpLen(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpLen)
- [`func CmpLt(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpLt)
- [`func CmpLte(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpLte)
- [`func CmpMap(t TestingT, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpMap)
- [`func CmpMapEach(t TestingT, got interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpMapEach)
- [`func CmpN(t TestingT, got interface{}, num interface{}, tolerance interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpN)
- [`func CmpNaN(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNaN)
- [`func CmpNil(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNil)
- [`func CmpNone(t TestingT, got interface{}, unexpectedValues []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNone)
- [`func CmpNot(t TestingT, got interface{}, unexpected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNot)
- [`func CmpNotAny(t TestingT, got interface{}, unexpectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNotAny)
- [`func CmpNotEmpty(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNotEmpty)
- [`func CmpNotNaN(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNotNaN)
- [`func CmpNotNil(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNotNil)
- [`func CmpNotZero(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpNotZero)
- [`func CmpPPtr(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpPPtr)
- [`func CmpPtr(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpPtr)
- [`func CmpRe(t TestingT, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpRe)
- [`func CmpReAll(t TestingT, got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpReAll)
- [`func CmpSet(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSet)
- [`func CmpShallow(t TestingT, got interface{}, expectedPtr interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpShallow)
- [`func CmpSlice(t TestingT, got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSlice)
- [`func CmpSmuggle(t TestingT, got interface{}, fn interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSmuggle)
- [`func CmpString(t TestingT, got interface{}, expected string, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpString)
- [`func CmpStruct(t TestingT, got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpStruct)
- [`func CmpSubBagOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSubBagOf)
- [`func CmpSubMapOf(t TestingT, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSubMapOf)
- [`func CmpSubSetOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSubSetOf)
- [`func CmpSuperBagOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSuperBagOf)
- [`func CmpSuperMapOf(t TestingT, got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSuperMapOf)
- [`func CmpSuperSetOf(t TestingT, got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpSuperSetOf)
- [`func CmpTruncTime(t TestingT, got interface{}, expectedTime interface{}, trunc time.Duration, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpTruncTime)
- [`func CmpValues(t TestingT, got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpValues)
- [`func CmpZero(t TestingT, got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#CmpZero)

## [`testdeep.T`]

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
- [`func (t *T) CmpNoError(got error, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpNoError)
- [`func (t *T) CmpNotPanic(fn func(), args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpNotPanic)
- [`func (t *T) CmpPanic(fn func(), expected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpPanic)
- [`func (t *T) False(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.False)
- [`func (t *T) Run(name string, f func(t *T)) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Run)
- [`func (t *T) True(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.True)

Note that the convenient
[`Not()`](https://godoc.org/github.com/maxatome/go-testdeep#T.Not)
method is listed in the shortcut section below.

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
the correponding [Testdeep operator](#testdeep-operators):

```
T.Not(got, expected, args...) ⇒ T.Cmp(t, got, Not(expected), args...)
  ^-^                                         ^-^
   +-------------------------------------------+
```

- [`func (t *T) All(got interface{}, expectedValues []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.All)
- [`func (t *T) Any(got interface{}, expectedValues []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Any)
- [`func (t *T) Array(got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Array)
- [`func (t *T) ArrayEach(got interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.ArrayEach)
- [`func (t *T) Bag(got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Bag)
- [`func (t *T) Between(got interface{}, from interface{}, to interface{}, bounds BoundsKind, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Between)
- [`func (t *T) Cap(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Cap)
- [`func (t *T) Code(got interface{}, fn interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Code)
- [`func (t *T) Contains(got interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Contains)
- [`func (t *T) ContainsKey(got interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.ContainsKey)
- [`func (t *T) Empty(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Empty)
- [`func (t *T) Gt(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Gt)
- [`func (t *T) Gte(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Gte)
- [`func (t *T) HasPrefix(got interface{}, expected string, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.HasPrefix)
- [`func (t *T) HasSuffix(got interface{}, expected string, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.HasSuffix)
- [`func (t *T) Isa(got interface{}, model interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Isa)
- [`func (t *T) Keys(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Keys)
- `Lax` operator appears as more readable
  [`T.CmpLax()`](https://godoc.org/github.com/maxatome/go-testdeep#T.CmpLax)
  considering it just sets the
  [`BeLax` config flag](https://godoc.org/github.com/maxatome/go-testdeep#ContextConfig)
  but provides a first-class global feature
- [`func (t *T) Len(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Len)
- [`func (t *T) Lt(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Lt)
- [`func (t *T) Lte(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Lte)
- [`func (t *T) Map(got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Map)
- [`func (t *T) MapEach(got interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.MapEach)
- [`func (t *T) N(got interface{}, num interface{}, tolerance interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.N)
- [`func (t *T) NaN(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.NaN)
- [`func (t *T) Nil(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Nil)
- [`func (t *T) None(got interface{}, unexpectedValues []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.None)
- [`func (t *T) Not(got interface{}, unexpected interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Not)
- [`func (t *T) NotAny(got interface{}, unexpectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.NotAny)
- [`func (t *T) NotEmpty(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.NotEmpty)
- [`func (t *T) NotNaN(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.NotNaN)
- [`func (t *T) NotNil(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.NotNil)
- [`func (t *T) NotZero(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.NotZero)
- [`func (t *T) PPtr(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.PPtr)
- [`func (t *T) Ptr(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Ptr)
- [`func (t *T) Re(got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Re)
- [`func (t *T) ReAll(got interface{}, reg interface{}, capture interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.ReAll)
- [`func (t *T) Set(got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Set)
- [`func (t *T) Shallow(got interface{}, expectedPtr interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Shallow)
- [`func (t *T) Slice(got interface{}, model interface{}, expectedEntries ArrayEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Slice)
- [`func (t *T) Smuggle(got interface{}, fn interface{}, expectedValue interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Smuggle)
- [`func (t *T) String(got interface{}, expected string, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.String)
- [`func (t *T) Struct(got interface{}, model interface{}, expectedFields StructFields, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Struct)
- [`func (t *T) SubBagOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.SubBagOf)
- [`func (t *T) SubMapOf(got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.SubMapOf)
- [`func (t *T) SubSetOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.SubSetOf)
- [`func (t *T) SuperBagOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.SuperBagOf)
- [`func (t *T) SuperMapOf(got interface{}, model interface{}, expectedEntries MapEntries, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.SuperMapOf)
- [`func (t *T) SuperSetOf(got interface{}, expectedItems []interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.SuperSetOf)
- [`func (t *T) TruncTime(got interface{}, expectedTime interface{}, trunc time.Duration, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.TruncTime)
- [`func (t *T) Values(got interface{}, val interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Values)
- [`func (t *T) Zero(got interface{}, args ...interface{}) bool`](https://godoc.org/github.com/maxatome/go-testdeep#T.Zero)


## `Testdeep` operators

- [`All`] all expected values have to match;
- [`Any`] at least one expected value have to match;
- [`Array`] compares the contents of an array or a pointer on an
  array;
- [`ArrayEach`] compares each array or slice item;
- [`Bag`] compares the contents of an array or a slice without taking
  care of the order of items;
- [`Between`] checks that a number, string or [`time.Time`] is between two
  bounds;
- [`Cap`] checks an array, slice or channel capacity;
- [`Code`] allows to use a custom function;
- [`Contains`] checks that a string, [`error`] or [`fmt.Stringer`]
  interfaces contain a sub-string; or an array, slice or map contain a
  value;
- [`ContainsKey`] checks that a map contains a key;
- [`Empty`] checks that an array, a channel, a map, a slice or a
  string is empty;
- [`Gt`] checks that a number, string or [`time.Time`] is greater than a
  value;
- [`Gte`] checks that a number, string or [`time.Time`] is greater or equal
  than a value;
- [`HasPrefix`] checks the prefix of a string, [`error`] or
  [`fmt.Stringer`] interfaces;
- [`HasSuffix`] checks the suffix of a string, [`error`] or
  [`fmt.Stringer`] interfaces;
- [`Ignore`] allows to ignore a comparison;
- [`Isa`] checks the data type or whether data implements an interface
  or not;
- [`Keys`] checks keys of a map;
- [`Lax`] temporarily enables
  [`BeLax` config flag](https://godoc.org/github.com/maxatome/go-testdeep#ContextConfig);
- [`Len`] checks an array, slice, map, string or channel length;
- [`Lt`] checks that a number, string or [`time.Time`] is lesser than a value;
- [`Lte`] checks that a number, string or [`time.Time`] is lesser or equal
  than a value;
- [`Map`] compares the contents of a map;
- [`MapEach`] compares each map entry;
- [`N`] compares a number with a tolerance value;
- [`NaN`] checks a floating number is [`math.NaN`];
- [`Nil`] compares to `nil`;
- [`None`] no values have to match;
- [`NotAny`] compares the contents of an array or a slice, no values
  have to match;
- [`Not`] value must not match;
- [`NotEmpty`] checks that an array, a channel, a map, a slice or a
  string is not empty;
- [`NotNaN`] checks a floating number is not [`math.NaN`];
- [`NotNil`] checks that data is not `nil`;
- [`NotZero`] checks that data is not zero regarding its type;
- [`PPtr`] allows to easily test a pointer of pointer value,
- [`Ptr`] allows to easily test a pointer value,
- [`Re`] allows to apply a regexp on a string (or convertible),
  `[]byte`, [`error`] or [`fmt.Stringer`] interfaces, and even test
  the captured groups;
- [`ReAll`] allows to successively apply a regexp on a string (or
  convertible), `[]byte`, [`error`] or [`fmt.Stringer`] interfaces,
  and even test the captured groups;
- [`Set`] compares the contents of an array or a slice ignoring
  duplicates and without taking care of the order of items;
- [`Shallow`] compares pointers only, not their contents;
- [`Slice`] compares the contents of a slice or a pointer on a slice;
- [`Smuggle`] changes data contents or mutates it into another type
  via a custom function or a struct fields-path before stepping down
  in favor of generic comparison process;
- [`String`] checks a string, [`error`] or [`fmt.Stringer`] interfaces
  string contents;
- [`Struct`] compares the contents of a struct or a pointer on a
  struct;
- [`SubBagOf`] compares the contents of an array or a slice without
  taking care of the order of items but with potentially some
  exclusions;
- [`SubMapOf`] compares the contents of a map but with potentially
  some exclusions;
- [`SubSetOf`] compares the contents of an array or a slice ignoring
  duplicates and without taking care of the order of items but with
  potentially some exclusions;
- [`SuperBagOf`] compares the contents of an array or a slice without
  taking care of the order of items but with potentially some extra
  items;
- [`SuperMapOf`] compares the contents of a map but with potentially
  some extra entries;
- [`SuperSetOf`] compares the contents of an array or a slice ignoring
  duplicates and without taking care of the order of items but with
  potentially some extra items;
- [`TruncTime`] compares time.Time (or assignable) values after
  truncating them;
- [`Values`] checks values of a map;
- [`Zero`] checks data against its zero'ed conterpart.

[`All`]: https://godoc.org/github.com/maxatome/go-testdeep#All
[`Any`]: https://godoc.org/github.com/maxatome/go-testdeep#Any
[`Array`]: https://godoc.org/github.com/maxatome/go-testdeep#Array
[`ArrayEach`]: https://godoc.org/github.com/maxatome/go-testdeep#ArrayEach
[`Bag`]: https://godoc.org/github.com/maxatome/go-testdeep#Bag
[`Between`]: https://godoc.org/github.com/maxatome/go-testdeep#Between
[`Cap`]: https://godoc.org/github.com/maxatome/go-testdeep#Cap
[`Code`]: https://godoc.org/github.com/maxatome/go-testdeep#Code
[`Contains`]: https://godoc.org/github.com/maxatome/go-testdeep#Contains
[`ContainsKey`]: https://godoc.org/github.com/maxatome/go-testdeep#ContainsKey
[`Empty`]: https://godoc.org/github.com/maxatome/go-testdeep#Empty
[`Gt`]: https://godoc.org/github.com/maxatome/go-testdeep#Gt
[`Gte`]: https://godoc.org/github.com/maxatome/go-testdeep#Gte
[`HasPrefix`]: https://godoc.org/github.com/maxatome/go-testdeep#HasPrefix
[`HasSuffix`]: https://godoc.org/github.com/maxatome/go-testdeep#HasSuffix
[`Ignore`]: https://godoc.org/github.com/maxatome/go-testdeep#Isa
[`Isa`]: https://godoc.org/github.com/maxatome/go-testdeep#Isa
[`Keys`]: https://godoc.org/github.com/maxatome/go-testdeep#Keys
[`Lax`]: https://godoc.org/github.com/maxatome/go-testdeep#Lax
[`Len`]: https://godoc.org/github.com/maxatome/go-testdeep#Len
[`Lt`]: https://godoc.org/github.com/maxatome/go-testdeep#Lt
[`Lte`]: https://godoc.org/github.com/maxatome/go-testdeep#Lte
[`Map`]: https://godoc.org/github.com/maxatome/go-testdeep#Map
[`MapEach`]: https://godoc.org/github.com/maxatome/go-testdeep#MapEach
[`N`]: https://godoc.org/github.com/maxatome/go-testdeep#N
[`NaN`]: https://godoc.org/github.com/maxatome/go-testdeep#NaN
[`Nil`]: https://godoc.org/github.com/maxatome/go-testdeep#Nil
[`None`]: https://godoc.org/github.com/maxatome/go-testdeep#None
[`NotAny`]: https://godoc.org/github.com/maxatome/go-testdeep#NotAny
[`Not`]: https://godoc.org/github.com/maxatome/go-testdeep#Not
[`NotEmpty`]: https://godoc.org/github.com/maxatome/go-testdeep#NotEmpty
[`NotNaN`]: https://godoc.org/github.com/maxatome/go-testdeep#NotNaN
[`NotNil`]: https://godoc.org/github.com/maxatome/go-testdeep#NotNil
[`NotZero`]: https://godoc.org/github.com/maxatome/go-testdeep#NotZero
[`PPtr`]: https://godoc.org/github.com/maxatome/go-testdeep#PPtr
[`Ptr`]: https://godoc.org/github.com/maxatome/go-testdeep#Ptr
[`Re`]: https://godoc.org/github.com/maxatome/go-testdeep#Re
[`ReAll`]: https://godoc.org/github.com/maxatome/go-testdeep#ReAll
[`Set`]: https://godoc.org/github.com/maxatome/go-testdeep#Set
[`Shallow`]: https://godoc.org/github.com/maxatome/go-testdeep#Shallow
[`Slice`]: https://godoc.org/github.com/maxatome/go-testdeep#Slice
[`Smuggle`]: https://godoc.org/github.com/maxatome/go-testdeep#Smuggle
[`String`]: https://godoc.org/github.com/maxatome/go-testdeep#String
[`Struct`]: https://godoc.org/github.com/maxatome/go-testdeep#Struct
[`SubBagOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SubBagOf
[`SubMapOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SubMapOf
[`SubSetOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SubSetOf
[`SuperBagOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SuperBagOf
[`SuperMapOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SuperMapOf
[`SuperSetOf`]: https://godoc.org/github.com/maxatome/go-testdeep#SuperSetOf
[`TruncTime`]: https://godoc.org/github.com/maxatome/go-testdeep#TruncTime
[`Values`]: https://godoc.org/github.com/maxatome/go-testdeep#Values
[`Zero`]: https://godoc.org/github.com/maxatome/go-testdeep#Zero

[`testdeep.T`]: https://godoc.org/github.com/maxatome/go-testdeep#T
[`*testdeep.T`]: https://godoc.org/github.com/maxatome/go-testdeep#T

[`error`]: https://golang.org/ref/spec#Errors
[`fmt.Stringer`]: https://golang.org/pkg/fmt/#Stringer
[`time.Time`]: https://golang.org/pkg/time/
[`math.NaN`]: https://golang.org/pkg/math/#NaN
