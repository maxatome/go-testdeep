---
title: "Operators matrices"
weight: 1
---

## Operator → go type matrix

<!-- op-go-matrix:begin -->

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`All`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`All`] |
| [`Any`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Any`] |
| [`Array`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ptr on array | ✓ | ✗ | ✗ | [`Array`] |
| [`ArrayEach`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`ArrayEach`] |
| [`Bag`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`Bag`] |
| [`Between`] | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗ | ✓ | ✗ | ✗ | [`Between`] |
| [`Cap`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | [`Cap`] |
| [`Catch`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Catch`] |
| [`Code`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Code`] |
| [`Contains`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`Contains`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`ContainsKey`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ | ✗ | ✗ | [`ContainsKey`] |
| [`Delay`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Delay`] |
| [`Empty`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ | ✗ | ptr on array/slice/map/string | ✓ | ✓ | ✗ | [`Empty`] |
| [`Gt`] | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗ | ✓ | ✗ | ✗ | [`Gt`] |
| [`Gte`] | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗ | ✓ | ✗ | ✗ | [`Gte`] |
| [`HasPrefix`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | `[]byte` | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`HasPrefix`] |
| [`HasSuffix`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | `[]byte` | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`HasSuffix`] |
| [`Ignore`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Ignore`] |
| [`Isa`] | ✗ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Isa`] |
| [`JSON`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✗ | ✗ | [`JSON`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`Keys`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ | ✗ | ✗ | [`Keys`] |
| [`Lax`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Lax`] |
| [`Len`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ | ✗ | ✗ | ✓ | ✓ | ✗ | [`Len`] |
| [`Lt`] | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗ | ✓ | ✗ | ✗ | [`Lt`] |
| [`Lte`] | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗ | ✓ | ✗ | ✗ | [`Lte`] |
| [`Map`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ptr on map | ✓ | ✗ | ✗ | [`Map`] |
| [`MapEach`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ptr on map | ✓ | ✗ | ✗ | [`MapEach`] |
| [`N`] | ✗ | ✗ | ✗ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | [`N`] |
| [`NaN`] | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | [`NaN`] |
| [`Nil`] | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | [`Nil`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`None`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`None`] |
| [`Not`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Not`] |
| [`NotAny`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`NotAny`] |
| [`NotEmpty`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ | ✗ | ptr on array/slice/map/string | ✓ | ✓ | ✗ | [`NotEmpty`] |
| [`NotNaN`] | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | [`NotNaN`] |
| [`NotNil`] | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | [`NotNil`] |
| [`NotZero`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`NotZero`] |
| [`PPtr`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | [`PPtr`] |
| [`Ptr`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | [`Ptr`] |
| [`Re`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | `[]byte` | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`Re`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`ReAll`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | `[]byte` | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`ReAll`] |
| [`Set`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`Set`] |
| [`Shallow`] | ✓ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | [`Shallow`] |
| [`Slice`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ptr on slice | ✓ | ✗ | ✗ | [`Slice`] |
| [`Smuggle`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Smuggle`] |
| [`SStruct`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ptr on struct | ✓ | ✗ | ✗ | [`SStruct`] |
| [`String`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | `[]byte` | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`String`] |
| [`Struct`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ptr on struct | ✓ | ✗ | ✗ | [`Struct`] |
| [`SubBagOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`SubBagOf`] |
| [`SubJSONOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ptr on map/struct | ✓ | ✗ | ✗ | [`SubJSONOf`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`SubMapOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ptr on map | ✓ | ✗ | ✗ | [`SubMapOf`] |
| [`SubSetOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`SubSetOf`] |
| [`SuperBagOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`SuperBagOf`] |
| [`SuperJSONOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ptr on map/struct | ✓ | ✗ | ✗ | [`SuperJSONOf`] |
| [`SuperMapOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ptr on map | ✓ | ✗ | ✗ | [`SuperMapOf`] |
| [`SuperSetOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`SuperSetOf`] |
| [`Tag`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Tag`] |
| [`TruncTime`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | [`time.Time`] | todo | ✓ | ✗ | ✗ | [`TruncTime`] |
| [`Values`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ | ✗ | ✗ | [`Values`] |
| [`Zero`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Zero`] |

[`T`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T
[`TestDeep`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#TestDeep
[`Cmp`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#Cmp

[`tdhttp`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp

[`BeLax` config flag]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#ContextConfig
[`error`]: https://pkg.go.dev/builtin/#error


[`fmt.Stringer`]: https://pkg.go.dev/fmt/#Stringer
[`time.Time`]: https://pkg.go.dev/time/#Time
[`math.NaN`]: https://pkg.go.dev/math/#NaN
[`All`]: {{< ref "All" >}}
[`Any`]: {{< ref "Any" >}}
[`Array`]: {{< ref "Array" >}}
[`ArrayEach`]: {{< ref "ArrayEach" >}}
[`Bag`]: {{< ref "Bag" >}}
[`Between`]: {{< ref "Between" >}}
[`Cap`]: {{< ref "Cap" >}}
[`Catch`]: {{< ref "Catch" >}}
[`Code`]: {{< ref "Code" >}}
[`Contains`]: {{< ref "Contains" >}}
[`ContainsKey`]: {{< ref "ContainsKey" >}}
[`Delay`]: {{< ref "Delay" >}}
[`Empty`]: {{< ref "Empty" >}}
[`Gt`]: {{< ref "Gt" >}}
[`Gte`]: {{< ref "Gte" >}}
[`HasPrefix`]: {{< ref "HasPrefix" >}}
[`HasSuffix`]: {{< ref "HasSuffix" >}}
[`Ignore`]: {{< ref "Ignore" >}}
[`Isa`]: {{< ref "Isa" >}}
[`JSON`]: {{< ref "JSON" >}}
[`Keys`]: {{< ref "Keys" >}}
[`Lax`]: {{< ref "Lax" >}}
[`Len`]: {{< ref "Len" >}}
[`Lt`]: {{< ref "Lt" >}}
[`Lte`]: {{< ref "Lte" >}}
[`Map`]: {{< ref "Map" >}}
[`MapEach`]: {{< ref "MapEach" >}}
[`N`]: {{< ref "N" >}}
[`NaN`]: {{< ref "NaN" >}}
[`Nil`]: {{< ref "Nil" >}}
[`None`]: {{< ref "None" >}}
[`Not`]: {{< ref "Not" >}}
[`NotAny`]: {{< ref "NotAny" >}}
[`NotEmpty`]: {{< ref "NotEmpty" >}}
[`NotNaN`]: {{< ref "NotNaN" >}}
[`NotNil`]: {{< ref "NotNil" >}}
[`NotZero`]: {{< ref "NotZero" >}}
[`PPtr`]: {{< ref "PPtr" >}}
[`Ptr`]: {{< ref "Ptr" >}}
[`Re`]: {{< ref "Re" >}}
[`ReAll`]: {{< ref "ReAll" >}}
[`Set`]: {{< ref "Set" >}}
[`Shallow`]: {{< ref "Shallow" >}}
[`Slice`]: {{< ref "Slice" >}}
[`Smuggle`]: {{< ref "Smuggle" >}}
[`SStruct`]: {{< ref "SStruct" >}}
[`String`]: {{< ref "String" >}}
[`Struct`]: {{< ref "Struct" >}}
[`SubBagOf`]: {{< ref "SubBagOf" >}}
[`SubJSONOf`]: {{< ref "SubJSONOf" >}}
[`SubMapOf`]: {{< ref "SubMapOf" >}}
[`SubSetOf`]: {{< ref "SubSetOf" >}}
[`SuperBagOf`]: {{< ref "SuperBagOf" >}}
[`SuperJSONOf`]: {{< ref "SuperJSONOf" >}}
[`SuperMapOf`]: {{< ref "SuperMapOf" >}}
[`SuperSetOf`]: {{< ref "SuperSetOf" >}}
[`Tag`]: {{< ref "Tag" >}}
[`TruncTime`]: {{< ref "TruncTime" >}}
[`Values`]: {{< ref "Values" >}}
[`Zero`]: {{< ref "Zero" >}}

[`CmpAll`]: {{< ref "All#cmpall-shortcut" >}}
[`CmpAny`]: {{< ref "Any#cmpany-shortcut" >}}
[`CmpArray`]: {{< ref "Array#cmparray-shortcut" >}}
[`CmpArrayEach`]: {{< ref "ArrayEach#cmparrayeach-shortcut" >}}
[`CmpBag`]: {{< ref "Bag#cmpbag-shortcut" >}}
[`CmpBetween`]: {{< ref "Between#cmpbetween-shortcut" >}}
[`CmpCap`]: {{< ref "Cap#cmpcap-shortcut" >}}
[`CmpCode`]: {{< ref "Code#cmpcode-shortcut" >}}
[`CmpContains`]: {{< ref "Contains#cmpcontains-shortcut" >}}
[`CmpContainsKey`]: {{< ref "ContainsKey#cmpcontainskey-shortcut" >}}
[`CmpEmpty`]: {{< ref "Empty#cmpempty-shortcut" >}}
[`CmpGt`]: {{< ref "Gt#cmpgt-shortcut" >}}
[`CmpGte`]: {{< ref "Gte#cmpgte-shortcut" >}}
[`CmpHasPrefix`]: {{< ref "HasPrefix#cmphasprefix-shortcut" >}}
[`CmpHasSuffix`]: {{< ref "HasSuffix#cmphassuffix-shortcut" >}}
[`CmpIsa`]: {{< ref "Isa#cmpisa-shortcut" >}}
[`CmpJSON`]: {{< ref "JSON#cmpjson-shortcut" >}}
[`CmpKeys`]: {{< ref "Keys#cmpkeys-shortcut" >}}
[`CmpLax`]: {{< ref "Lax#cmplax-shortcut" >}}
[`CmpLen`]: {{< ref "Len#cmplen-shortcut" >}}
[`CmpLt`]: {{< ref "Lt#cmplt-shortcut" >}}
[`CmpLte`]: {{< ref "Lte#cmplte-shortcut" >}}
[`CmpMap`]: {{< ref "Map#cmpmap-shortcut" >}}
[`CmpMapEach`]: {{< ref "MapEach#cmpmapeach-shortcut" >}}
[`CmpN`]: {{< ref "N#cmpn-shortcut" >}}
[`CmpNaN`]: {{< ref "NaN#cmpnan-shortcut" >}}
[`CmpNil`]: {{< ref "Nil#cmpnil-shortcut" >}}
[`CmpNone`]: {{< ref "None#cmpnone-shortcut" >}}
[`CmpNot`]: {{< ref "Not#cmpnot-shortcut" >}}
[`CmpNotAny`]: {{< ref "NotAny#cmpnotany-shortcut" >}}
[`CmpNotEmpty`]: {{< ref "NotEmpty#cmpnotempty-shortcut" >}}
[`CmpNotNaN`]: {{< ref "NotNaN#cmpnotnan-shortcut" >}}
[`CmpNotNil`]: {{< ref "NotNil#cmpnotnil-shortcut" >}}
[`CmpNotZero`]: {{< ref "NotZero#cmpnotzero-shortcut" >}}
[`CmpPPtr`]: {{< ref "PPtr#cmppptr-shortcut" >}}
[`CmpPtr`]: {{< ref "Ptr#cmpptr-shortcut" >}}
[`CmpRe`]: {{< ref "Re#cmpre-shortcut" >}}
[`CmpReAll`]: {{< ref "ReAll#cmpreall-shortcut" >}}
[`CmpSet`]: {{< ref "Set#cmpset-shortcut" >}}
[`CmpShallow`]: {{< ref "Shallow#cmpshallow-shortcut" >}}
[`CmpSlice`]: {{< ref "Slice#cmpslice-shortcut" >}}
[`CmpSmuggle`]: {{< ref "Smuggle#cmpsmuggle-shortcut" >}}
[`CmpSStruct`]: {{< ref "SStruct#cmpsstruct-shortcut" >}}
[`CmpString`]: {{< ref "String#cmpstring-shortcut" >}}
[`CmpStruct`]: {{< ref "Struct#cmpstruct-shortcut" >}}
[`CmpSubBagOf`]: {{< ref "SubBagOf#cmpsubbagof-shortcut" >}}
[`CmpSubJSONOf`]: {{< ref "SubJSONOf#cmpsubjsonof-shortcut" >}}
[`CmpSubMapOf`]: {{< ref "SubMapOf#cmpsubmapof-shortcut" >}}
[`CmpSubSetOf`]: {{< ref "SubSetOf#cmpsubsetof-shortcut" >}}
[`CmpSuperBagOf`]: {{< ref "SuperBagOf#cmpsuperbagof-shortcut" >}}
[`CmpSuperJSONOf`]: {{< ref "SuperJSONOf#cmpsuperjsonof-shortcut" >}}
[`CmpSuperMapOf`]: {{< ref "SuperMapOf#cmpsupermapof-shortcut" >}}
[`CmpSuperSetOf`]: {{< ref "SuperSetOf#cmpsupersetof-shortcut" >}}
[`CmpTruncTime`]: {{< ref "TruncTime#cmptrunctime-shortcut" >}}
[`CmpValues`]: {{< ref "Values#cmpvalues-shortcut" >}}
[`CmpZero`]: {{< ref "Zero#cmpzero-shortcut" >}}

[`T.All`]: {{< ref "All#tall-shortcut" >}}
[`T.Any`]: {{< ref "Any#tany-shortcut" >}}
[`T.Array`]: {{< ref "Array#tarray-shortcut" >}}
[`T.ArrayEach`]: {{< ref "ArrayEach#tarrayeach-shortcut" >}}
[`T.Bag`]: {{< ref "Bag#tbag-shortcut" >}}
[`T.Between`]: {{< ref "Between#tbetween-shortcut" >}}
[`T.Cap`]: {{< ref "Cap#tcap-shortcut" >}}
[`T.Code`]: {{< ref "Code#tcode-shortcut" >}}
[`T.Contains`]: {{< ref "Contains#tcontains-shortcut" >}}
[`T.ContainsKey`]: {{< ref "ContainsKey#tcontainskey-shortcut" >}}
[`T.Empty`]: {{< ref "Empty#tempty-shortcut" >}}
[`T.Gt`]: {{< ref "Gt#tgt-shortcut" >}}
[`T.Gte`]: {{< ref "Gte#tgte-shortcut" >}}
[`T.HasPrefix`]: {{< ref "HasPrefix#thasprefix-shortcut" >}}
[`T.HasSuffix`]: {{< ref "HasSuffix#thassuffix-shortcut" >}}
[`T.Isa`]: {{< ref "Isa#tisa-shortcut" >}}
[`T.JSON`]: {{< ref "JSON#tjson-shortcut" >}}
[`T.Keys`]: {{< ref "Keys#tkeys-shortcut" >}}
[`T.CmpLax`]: {{< ref "Lax#tcmplax-shortcut" >}}
[`T.Len`]: {{< ref "Len#tlen-shortcut" >}}
[`T.Lt`]: {{< ref "Lt#tlt-shortcut" >}}
[`T.Lte`]: {{< ref "Lte#tlte-shortcut" >}}
[`T.Map`]: {{< ref "Map#tmap-shortcut" >}}
[`T.MapEach`]: {{< ref "MapEach#tmapeach-shortcut" >}}
[`T.N`]: {{< ref "N#tn-shortcut" >}}
[`T.NaN`]: {{< ref "NaN#tnan-shortcut" >}}
[`T.Nil`]: {{< ref "Nil#tnil-shortcut" >}}
[`T.None`]: {{< ref "None#tnone-shortcut" >}}
[`T.Not`]: {{< ref "Not#tnot-shortcut" >}}
[`T.NotAny`]: {{< ref "NotAny#tnotany-shortcut" >}}
[`T.NotEmpty`]: {{< ref "NotEmpty#tnotempty-shortcut" >}}
[`T.NotNaN`]: {{< ref "NotNaN#tnotnan-shortcut" >}}
[`T.NotNil`]: {{< ref "NotNil#tnotnil-shortcut" >}}
[`T.NotZero`]: {{< ref "NotZero#tnotzero-shortcut" >}}
[`T.PPtr`]: {{< ref "PPtr#tpptr-shortcut" >}}
[`T.Ptr`]: {{< ref "Ptr#tptr-shortcut" >}}
[`T.Re`]: {{< ref "Re#tre-shortcut" >}}
[`T.ReAll`]: {{< ref "ReAll#treall-shortcut" >}}
[`T.Set`]: {{< ref "Set#tset-shortcut" >}}
[`T.Shallow`]: {{< ref "Shallow#tshallow-shortcut" >}}
[`T.Slice`]: {{< ref "Slice#tslice-shortcut" >}}
[`T.Smuggle`]: {{< ref "Smuggle#tsmuggle-shortcut" >}}
[`T.SStruct`]: {{< ref "SStruct#tsstruct-shortcut" >}}
[`T.String`]: {{< ref "String#tstring-shortcut" >}}
[`T.Struct`]: {{< ref "Struct#tstruct-shortcut" >}}
[`T.SubBagOf`]: {{< ref "SubBagOf#tsubbagof-shortcut" >}}
[`T.SubJSONOf`]: {{< ref "SubJSONOf#tsubjsonof-shortcut" >}}
[`T.SubMapOf`]: {{< ref "SubMapOf#tsubmapof-shortcut" >}}
[`T.SubSetOf`]: {{< ref "SubSetOf#tsubsetof-shortcut" >}}
[`T.SuperBagOf`]: {{< ref "SuperBagOf#tsuperbagof-shortcut" >}}
[`T.SuperJSONOf`]: {{< ref "SuperJSONOf#tsuperjsonof-shortcut" >}}
[`T.SuperMapOf`]: {{< ref "SuperMapOf#tsupermapof-shortcut" >}}
[`T.SuperSetOf`]: {{< ref "SuperSetOf#tsupersetof-shortcut" >}}
[`T.TruncTime`]: {{< ref "TruncTime#ttrunctime-shortcut" >}}
[`T.Values`]: {{< ref "Values#tvalues-shortcut" >}}
[`T.Zero`]: {{< ref "Zero#tzero-shortcut" >}}
<!-- op-go-matrix:end -->

Legend:

- ✗ means using this operator with a value type of this kind will always fail
- ✓ means using this operator with a value type of this kind can succeed
- `[]byte`, [`time.Time`], ptr on X, [`fmt.Stringer`], [`error`] means
  using this operator with this go type can succeed
- todo means should be implemented in future (PRs welcome :) )
- ¹ + ✓ means using this operator with the data behind the interface can succeed


## go type → operator matrix

Operators likely to succeed for each go type:

### Untyped `nil` value

<!-- go-nil-matrix:begin -->
- [`All`]
- [`Any`]
- [`Catch`]
- [`Code`]
- [`Delay`]
- [`Ignore`]
- [`JSON`]
- [`Lax`]
- [`Nil`]
- [`None`]
- [`Not`]
- [`NotNil`]
- [`NotZero`]
- [`Shallow`]
- [`Smuggle`]
- [`Tag`]
- [`Zero`]
<!-- go-nil-matrix:end -->

### `bool` values (or any type based on `bool`)

<!-- go-bool-matrix:begin -->
- [`All`]
- [`Any`]
- [`Catch`]
- [`Code`]
- [`Delay`]
- [`Ignore`]
- [`Isa`]
- [`JSON`]
- [`Lax`]
- [`None`]
- [`Not`]
- [`NotZero`]
- [`Smuggle`]
- [`Tag`]
- [`Zero`]
<!-- go-bool-matrix:end -->

### `string` values (or any type based on `string`)

<!-- go-str-matrix:begin -->
- [`All`]
- [`Any`]
- [`Between`]
- [`Catch`]
- [`Code`]
- [`Contains`]
- [`Delay`]
- [`Empty`]
- [`Gt`]
- [`Gte`]
- [`HasPrefix`]
- [`HasSuffix`]
- [`Ignore`]
- [`Isa`]
- [`JSON`]
- [`Lax`]
- [`Lt`]
- [`Lte`]
- [`None`]
- [`Not`]
- [`NotEmpty`]
- [`NotZero`]
- [`Re`]
- [`ReAll`]
- [`Shallow`]
- [`Smuggle`]
- [`String`]
- [`Tag`]
- [`Zero`]
<!-- go-str-matrix:end -->

### Integer values (`uint*`, `int*` or any type based on them)

<!-- go-int-matrix:begin -->
- [`All`]
- [`Any`]
- [`Between`]
- [`Catch`]
- [`Code`]
- [`Delay`]
- [`Gt`]
- [`Gte`]
- [`Ignore`]
- [`Isa`]
- [`JSON`]
- [`Lax`]
- [`Lt`]
- [`Lte`]
- [`N`]
- [`None`]
- [`Not`]
- [`NotZero`]
- [`Smuggle`]
- [`Tag`]
- [`Zero`]
<!-- go-int-matrix:end -->

### Float values (`float32`, `float64` or any type based on them)

<!-- go-float-matrix:begin -->
- [`All`]
- [`Any`]
- [`Between`]
- [`Catch`]
- [`Code`]
- [`Delay`]
- [`Gt`]
- [`Gte`]
- [`Ignore`]
- [`Isa`]
- [`JSON`]
- [`Lax`]
- [`Lt`]
- [`Lte`]
- [`N`]
- [`NaN`]
- [`None`]
- [`Not`]
- [`NotNaN`]
- [`NotZero`]
- [`Smuggle`]
- [`Tag`]
- [`Zero`]
<!-- go-float-matrix:end -->

### Complex values (`complex64`, `complex128` or any type based on them)

<!-- go-cplx-matrix:begin -->
- [`All`]
- [`Any`]
- [`Catch`]
- [`Code`]
- [`Delay`]
- [`Ignore`]
- [`Isa`]
- [`Lax`]
- [`None`]
- [`Not`]
- [`NotZero`]
- [`Smuggle`]
- [`Tag`]
- [`Zero`]
<!-- go-cplx-matrix:end -->

### Arrays

<!-- go-array-matrix:begin -->
- [`All`]
- [`Any`]
- [`Array`]
- [`ArrayEach`]
- [`Bag`]
- [`Cap`]
- [`Catch`]
- [`Code`]
- [`Contains`]
- [`Delay`]
- [`Empty`]
- [`Ignore`]
- [`Isa`]
- [`JSON`]
- [`Lax`]
- [`Len`]
- [`None`]
- [`Not`]
- [`NotAny`]
- [`NotEmpty`]
- [`NotZero`]
- [`Set`]
- [`Smuggle`]
- [`SubBagOf`]
- [`SubSetOf`]
- [`SuperBagOf`]
- [`SuperSetOf`]
- [`Tag`]
- [`Zero`]
<!-- go-array-matrix:end -->

### Slices

<!-- go-slice-matrix:begin -->
- [`All`]
- [`Any`]
- [`ArrayEach`]
- [`Bag`]
- [`Cap`]
- [`Catch`]
- [`Code`]
- [`Contains`]
- [`Delay`]
- [`Empty`]
- [`HasPrefix`] only `[]byte`
- [`HasSuffix`] only `[]byte`
- [`Ignore`]
- [`Isa`]
- [`JSON`]
- [`Lax`]
- [`Len`]
- [`Nil`]
- [`None`]
- [`Not`]
- [`NotAny`]
- [`NotEmpty`]
- [`NotNil`]
- [`NotZero`]
- [`Re`] only `[]byte`
- [`ReAll`] only `[]byte`
- [`Set`]
- [`Shallow`]
- [`Slice`]
- [`Smuggle`]
- [`String`] only `[]byte`
- [`SubBagOf`]
- [`SubSetOf`]
- [`SuperBagOf`]
- [`SuperSetOf`]
- [`Tag`]
- [`Zero`]
<!-- go-slice-matrix:end -->

### Maps

<!-- go-map-matrix:begin -->
- [`All`]
- [`Any`]
- [`Catch`]
- [`Code`]
- [`Contains`]
- [`ContainsKey`]
- [`Delay`]
- [`Empty`]
- [`Ignore`]
- [`Isa`]
- [`JSON`]
- [`Keys`]
- [`Lax`]
- [`Len`]
- [`Map`]
- [`MapEach`]
- [`Nil`]
- [`None`]
- [`Not`]
- [`NotEmpty`]
- [`NotNil`]
- [`NotZero`]
- [`Shallow`]
- [`Smuggle`]
- [`SubJSONOf`]
- [`SubMapOf`]
- [`SuperJSONOf`]
- [`SuperMapOf`]
- [`Tag`]
- [`Values`]
- [`Zero`]
<!-- go-map-matrix:end -->

### Structs

<!-- go-struct-matrix:begin -->
- [`All`]
- [`Any`]
- [`Between`] only [`time.Time`]
- [`Catch`]
- [`Code`]
- [`Delay`]
- [`Gt`] only [`time.Time`]
- [`Gte`] only [`time.Time`]
- [`Ignore`]
- [`Isa`]
- [`JSON`]
- [`Lax`]
- [`Lt`] only [`time.Time`]
- [`Lte`] only [`time.Time`]
- [`None`]
- [`Not`]
- [`NotZero`]
- [`SStruct`]
- [`Smuggle`]
- [`Struct`]
- [`SubJSONOf`]
- [`SuperJSONOf`]
- [`Tag`]
- [`TruncTime`] only [`time.Time`]
- [`Zero`]
<!-- go-struct-matrix:end -->

### Interface values

As all operators accept interface values, only specific interfaces are
listed below:

<!-- go-if-matrix:begin -->
- [`Contains`] → [`fmt.Stringer`]/[`error`]
- [`HasPrefix`] → [`fmt.Stringer`]/[`error`]
- [`HasSuffix`] → [`fmt.Stringer`]/[`error`]
- [`Re`] → [`fmt.Stringer`]/[`error`]
- [`ReAll`] → [`fmt.Stringer`]/[`error`]
- [`String`] → [`fmt.Stringer`]/[`error`]
<!-- go-if-matrix:end -->

### Any pointer

<!-- go-ptr-matrix:begin -->
- [`All`]
- [`Any`]
- [`Array`] only ptr on array
- [`ArrayEach`] only ptr on array/slice
- [`Bag`] only ptr on array/slice
- [`Catch`]
- [`Code`]
- [`Delay`]
- [`Empty`] only ptr on array/slice/map/string
- [`Ignore`]
- [`Isa`]
- [`JSON`]
- [`Lax`]
- [`Map`] only ptr on map
- [`MapEach`] only ptr on map
- [`Nil`]
- [`None`]
- [`Not`]
- [`NotAny`] only ptr on array/slice
- [`NotEmpty`] only ptr on array/slice/map/string
- [`NotNil`]
- [`NotZero`]
- [`PPtr`]
- [`Ptr`]
- [`SStruct`] only ptr on struct
- [`Set`] only ptr on array/slice
- [`Shallow`]
- [`Slice`] only ptr on slice
- [`Smuggle`]
- [`Struct`] only ptr on struct
- [`SubBagOf`] only ptr on array/slice
- [`SubJSONOf`] only ptr on map/struct
- [`SubMapOf`] only ptr on map
- [`SubSetOf`] only ptr on array/slice
- [`SuperBagOf`] only ptr on array/slice
- [`SuperJSONOf`] only ptr on map/struct
- [`SuperMapOf`] only ptr on map
- [`SuperSetOf`] only ptr on array/slice
- [`Tag`]
- [`Zero`]
<!-- go-ptr-matrix:end -->

### Channels

<!-- go-chan-matrix:begin -->
- [`All`]
- [`Any`]
- [`Cap`]
- [`Catch`]
- [`Code`]
- [`Delay`]
- [`Empty`]
- [`Ignore`]
- [`Isa`]
- [`Lax`]
- [`Len`]
- [`Nil`]
- [`None`]
- [`Not`]
- [`NotEmpty`]
- [`NotNil`]
- [`NotZero`]
- [`Shallow`]
- [`Smuggle`]
- [`Tag`]
- [`Zero`]
<!-- go-chan-matrix:end -->

### Functions

<!-- go-func-matrix:begin -->
- [`All`]
- [`Any`]
- [`Catch`]
- [`Code`]
- [`Delay`]
- [`Ignore`]
- [`Isa`]
- [`Lax`]
- [`Nil`]
- [`None`]
- [`Not`]
- [`NotNil`]
- [`NotZero`]
- [`Shallow`]
- [`Smuggle`]
- [`Tag`]
- [`Zero`]
<!-- go-func-matrix:end -->
