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
| [`Code`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Code`] |
| [`Contains`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ | ✗ | ✗ | ✓ | ✗ | ✗ | [`Contains`] |
| [`ContainsKey`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ | ✗ | ✗ | [`ContainsKey`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`Empty`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ | ✗ | ptr on array/slice/map/string | ✓ | ✓ | ✗ | [`Empty`] |
| [`Gt`] | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗ | ✓ | ✗ | ✗ | [`Gt`] |
| [`Gte`] | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗ | ✓ | ✗ | ✗ | [`Gte`] |
| [`HasPrefix`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`HasPrefix`] |
| [`HasSuffix`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`HasSuffix`] |
| [`Ignore`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Ignore`] |
| [`Isa`] | ✗ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Isa`] |
| [`JSON`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✗ | ✗ | [`JSON`] |
| [`Keys`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ | ✗ | ✗ | [`Keys`] |
| [`Lax`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Lax`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`Len`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ | ✗ | ✗ | ✓ | ✓ | ✗ | [`Len`] |
| [`Lt`] | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗ | ✓ | ✗ | ✗ | [`Lt`] |
| [`Lte`] | ✗ | ✗ | ✓ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | [`time.Time`] | ✗ | ✓ | ✗ | ✗ | [`Lte`] |
| [`Map`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ptr on map | ✓ | ✗ | ✗ | [`Map`] |
| [`MapEach`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ptr on map | ✓ | ✗ | ✗ | [`MapEach`] |
| [`N`] | ✗ | ✗ | ✗ | ✓ | ✓ | todo | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | [`N`] |
| [`NaN`] | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | [`NaN`] |
| [`Nil`] | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | [`Nil`] |
| [`None`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`None`] |
| [`Not`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Not`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`NotAny`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`NotAny`] |
| [`NotEmpty`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✓ | ✓ | ✓ | ✗ | ptr on array/slice/map/string | ✓ | ✓ | ✗ | [`NotEmpty`] |
| [`NotNaN`] | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | [`NotNaN`] |
| [`NotNil`] | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | [`NotNil`] |
| [`NotZero`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`NotZero`] |
| [`PPtr`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | [`PPtr`] |
| [`Ptr`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | [`Ptr`] |
| [`Re`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | `[]byte` | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`Re`] |
| [`ReAll`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | `[]byte` | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`ReAll`] |
| [`Set`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`Set`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`Shallow`] | ✓ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | [`Shallow`] |
| [`Slice`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ptr on slice | ✓ | ✗ | ✗ | [`Slice`] |
| [`Smuggle`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Smuggle`] |
| [`String`] | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ + [`fmt.Stringer`]/[`error`] | ✗ | ✗ | [`String`] |
| [`Struct`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ptr on struct | ✓ | ✗ | ✗ | [`Struct`] |
| [`SubBagOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`SubBagOf`] |
| [`SubMapOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ptr on map | ✓ | ✗ | ✗ | [`SubMapOf`] |
| [`SubSetOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`SubSetOf`] |
| [`SuperBagOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`SuperBagOf`] |
| [`SuperMapOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ptr on map | ✓ | ✗ | ✗ | [`SuperMapOf`] |

| Operator vs go type | nil | bool | string | {u,}int* | float* | complex* | array | slice | map | struct | pointer | interface¹ | chan | func | operator |
| ------------------- | --- | ---- | ------ | -------- | ------ | -------- | ----- | ----- | --- | ------ | ------- | ---------- | ---- | ---- | -------- |
| [`SuperSetOf`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✓ | ✗ | ✗ | ptr on array/slice | ✓ | ✗ | ✗ | [`SuperSetOf`] |
| [`Tag`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Tag`] |
| [`TruncTime`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | [`time.Time`] | todo | ✓ | ✗ | ✗ | [`TruncTime`] |
| [`Values`] | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✓ | ✗ | ✗ | [`Values`] |
| [`Zero`] | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | [`Zero`] |
[`T`]: https://godoc.org/github.com/maxatome/go-testdeep#T
[`TestDeep`]: https://godoc.org/github.com/maxatome/go-testdeep#TestDeep
[`Cmp`]: https://godoc.org/github.com/maxatome/go-testdeep#Cmp

[`tdhttp`]: https://godoc.org/github.com/maxatome/go-testdeep/helpers/tdhttp

[`BeLax` config flag]: https://godoc.org/github.com/maxatome/go-testdeep#ContextConfig
[`error`]: https://golang.org/pkg/builtin/#error


[`fmt.Stringer`]: https://godoc.org/pkg/fmt/#Stringer
[`time.Time`]: https://godoc.org/pkg/time/#Time
[`math.NaN`]: https://godoc.org/pkg/math/#NaN
[`All`]: {{< ref "All" >}}
[`Any`]: {{< ref "Any" >}}
[`Array`]: {{< ref "Array" >}}
[`ArrayEach`]: {{< ref "ArrayEach" >}}
[`Bag`]: {{< ref "Bag" >}}
[`Between`]: {{< ref "Between" >}}
[`Cap`]: {{< ref "Cap" >}}
[`Code`]: {{< ref "Code" >}}
[`Contains`]: {{< ref "Contains" >}}
[`ContainsKey`]: {{< ref "ContainsKey" >}}
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
[`String`]: {{< ref "String" >}}
[`Struct`]: {{< ref "Struct" >}}
[`SubBagOf`]: {{< ref "SubBagOf" >}}
[`SubMapOf`]: {{< ref "SubMapOf" >}}
[`SubSetOf`]: {{< ref "SubSetOf" >}}
[`SuperBagOf`]: {{< ref "SuperBagOf" >}}
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
[`CmpString`]: {{< ref "String#cmpstring-shortcut" >}}
[`CmpStruct`]: {{< ref "Struct#cmpstruct-shortcut" >}}
[`CmpSubBagOf`]: {{< ref "SubBagOf#cmpsubbagof-shortcut" >}}
[`CmpSubMapOf`]: {{< ref "SubMapOf#cmpsubmapof-shortcut" >}}
[`CmpSubSetOf`]: {{< ref "SubSetOf#cmpsubsetof-shortcut" >}}
[`CmpSuperBagOf`]: {{< ref "SuperBagOf#cmpsuperbagof-shortcut" >}}
[`CmpSuperMapOf`]: {{< ref "SuperMapOf#cmpsupermapof-shortcut" >}}
[`CmpSuperSetOf`]: {{< ref "SuperSetOf#cmpsupersetof-shortcut" >}}
[`CmpTruncTime`]: {{< ref "TruncTime#cmptrunctime-shortcut" >}}
[`CmpValues`]: {{< ref "Values#cmpvalues-shortcut" >}}
[`CmpZero`]: {{< ref "Zero#cmpzero-shortcut" >}}

[`T.All`]: {{< ref "All#t-all-shortcut" >}}
[`T.Any`]: {{< ref "Any#t-any-shortcut" >}}
[`T.Array`]: {{< ref "Array#t-array-shortcut" >}}
[`T.ArrayEach`]: {{< ref "ArrayEach#t-arrayeach-shortcut" >}}
[`T.Bag`]: {{< ref "Bag#t-bag-shortcut" >}}
[`T.Between`]: {{< ref "Between#t-between-shortcut" >}}
[`T.Cap`]: {{< ref "Cap#t-cap-shortcut" >}}
[`T.Code`]: {{< ref "Code#t-code-shortcut" >}}
[`T.Contains`]: {{< ref "Contains#t-contains-shortcut" >}}
[`T.ContainsKey`]: {{< ref "ContainsKey#t-containskey-shortcut" >}}
[`T.Empty`]: {{< ref "Empty#t-empty-shortcut" >}}
[`T.Gt`]: {{< ref "Gt#t-gt-shortcut" >}}
[`T.Gte`]: {{< ref "Gte#t-gte-shortcut" >}}
[`T.HasPrefix`]: {{< ref "HasPrefix#t-hasprefix-shortcut" >}}
[`T.HasSuffix`]: {{< ref "HasSuffix#t-hassuffix-shortcut" >}}
[`T.Isa`]: {{< ref "Isa#t-isa-shortcut" >}}
[`T.JSON`]: {{< ref "JSON#t-json-shortcut" >}}
[`T.Keys`]: {{< ref "Keys#t-keys-shortcut" >}}
[`T.CmpLax`]: {{< ref "Lax#t-cmplax-shortcut" >}}
[`T.Len`]: {{< ref "Len#t-len-shortcut" >}}
[`T.Lt`]: {{< ref "Lt#t-lt-shortcut" >}}
[`T.Lte`]: {{< ref "Lte#t-lte-shortcut" >}}
[`T.Map`]: {{< ref "Map#t-map-shortcut" >}}
[`T.MapEach`]: {{< ref "MapEach#t-mapeach-shortcut" >}}
[`T.N`]: {{< ref "N#t-n-shortcut" >}}
[`T.NaN`]: {{< ref "NaN#t-nan-shortcut" >}}
[`T.Nil`]: {{< ref "Nil#t-nil-shortcut" >}}
[`T.None`]: {{< ref "None#t-none-shortcut" >}}
[`T.Not`]: {{< ref "Not#t-not-shortcut" >}}
[`T.NotAny`]: {{< ref "NotAny#t-notany-shortcut" >}}
[`T.NotEmpty`]: {{< ref "NotEmpty#t-notempty-shortcut" >}}
[`T.NotNaN`]: {{< ref "NotNaN#t-notnan-shortcut" >}}
[`T.NotNil`]: {{< ref "NotNil#t-notnil-shortcut" >}}
[`T.NotZero`]: {{< ref "NotZero#t-notzero-shortcut" >}}
[`T.PPtr`]: {{< ref "PPtr#t-pptr-shortcut" >}}
[`T.Ptr`]: {{< ref "Ptr#t-ptr-shortcut" >}}
[`T.Re`]: {{< ref "Re#t-re-shortcut" >}}
[`T.ReAll`]: {{< ref "ReAll#t-reall-shortcut" >}}
[`T.Set`]: {{< ref "Set#t-set-shortcut" >}}
[`T.Shallow`]: {{< ref "Shallow#t-shallow-shortcut" >}}
[`T.Slice`]: {{< ref "Slice#t-slice-shortcut" >}}
[`T.Smuggle`]: {{< ref "Smuggle#t-smuggle-shortcut" >}}
[`T.String`]: {{< ref "String#t-string-shortcut" >}}
[`T.Struct`]: {{< ref "Struct#t-struct-shortcut" >}}
[`T.SubBagOf`]: {{< ref "SubBagOf#t-subbagof-shortcut" >}}
[`T.SubMapOf`]: {{< ref "SubMapOf#t-submapof-shortcut" >}}
[`T.SubSetOf`]: {{< ref "SubSetOf#t-subsetof-shortcut" >}}
[`T.SuperBagOf`]: {{< ref "SuperBagOf#t-superbagof-shortcut" >}}
[`T.SuperMapOf`]: {{< ref "SuperMapOf#t-supermapof-shortcut" >}}
[`T.SuperSetOf`]: {{< ref "SuperSetOf#t-supersetof-shortcut" >}}
[`T.TruncTime`]: {{< ref "TruncTime#t-trunctime-shortcut" >}}
[`T.Values`]: {{< ref "Values#t-values-shortcut" >}}
[`T.Zero`]: {{< ref "Zero#t-zero-shortcut" >}}
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
- [`Code`]
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
- [`Code`]
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
- [`Code`]
- [`Contains`]
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
- [`Code`]
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
- [`Code`]
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
- [`Code`]
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
- [`Code`]
- [`Contains`]
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
- [`Code`]
- [`Contains`]
- [`Empty`]
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
- [`Code`]
- [`Contains`]
- [`ContainsKey`]
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
- [`SubMapOf`]
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
- [`Code`]
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
- [`Smuggle`]
- [`Struct`]
- [`Tag`]
- [`TruncTime`] only [`time.Time`]
- [`Zero`]
<!-- go-struct-matrix:end -->

### Interface values

As all operators accept interface values, only specific interfaces are
listed below:

<!-- go-if-matrix:begin -->
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
- [`Code`]
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
- [`Set`] only ptr on array/slice
- [`Shallow`]
- [`Slice`] only ptr on slice
- [`Smuggle`]
- [`Struct`] only ptr on struct
- [`SubBagOf`] only ptr on array/slice
- [`SubMapOf`] only ptr on map
- [`SubSetOf`] only ptr on array/slice
- [`SuperBagOf`] only ptr on array/slice
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
- [`Code`]
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
- [`Code`]
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
