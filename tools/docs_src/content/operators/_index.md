+++
title = "Operators"
date = 2019-10-03T22:20:09+02:00
weight = 15
+++

## All operators

<!-- operators:begin -->
[`All`]({{< ref "All" >}})
: all expected values have to match

[`Any`]({{< ref "Any" >}})
: at least one expected value have to match

[`Array`]({{< ref "Array" >}})
: compares the contents of an array or a pointer on an array

[`ArrayEach`]({{< ref "ArrayEach" >}})
: compares each array or slice item

[`Bag`]({{< ref "Bag" >}})
: compares the contents of an array or a slice without taking care of the order of items

[`Between`]({{< ref "Between" >}})
: checks that a number, string or [`time.Time`] is between two bounds

[`Cap`]({{< ref "Cap" >}})
: checks an array, slice or channel capacity

[`Catch`]({{< ref "Catch" >}})
: catches data on the fly before comparing it

[`Code`]({{< ref "Code" >}})
: checks using a custom function

[`Contains`]({{< ref "Contains" >}})
: checks that a string, `[]byte`, [`error`] or [`fmt.Stringer`] interfaces contain a rune, byte or a sub-string; or a slice contains a single value or a sub-slice; or an array or map contain a single value

[`ContainsKey`]({{< ref "ContainsKey" >}})
: checks that a map contains a key

[`Delay`]({{< ref "Delay" >}})
: delays the operator construction till first use

[`Empty`]({{< ref "Empty" >}})
: checks that an array, a channel, a map, a slice or a string is empty

[`Gt`]({{< ref "Gt" >}})
: checks that a number, string or [`time.Time`] is greater than a value

[`Gte`]({{< ref "Gte" >}})
: checks that a number, string or [`time.Time`] is greater or equal than a value

[`HasPrefix`]({{< ref "HasPrefix" >}})
: checks the prefix of a string, `[]byte`, [`error`] or [`fmt.Stringer`] interfaces

[`HasSuffix`]({{< ref "HasSuffix" >}})
: checks the suffix of a string, `[]byte`, [`error`] or [`fmt.Stringer`] interfaces

[`Ignore`]({{< ref "Ignore" >}})
: allows to ignore a comparison

[`Isa`]({{< ref "Isa" >}})
: checks the data type or whether data implements an interface or not

[`JSON`]({{< ref "JSON" >}})
: compares against JSON representation

[`Keys`]({{< ref "Keys" >}})
: checks keys of a map

[`Lax`]({{< ref "Lax" >}})
: temporarily enables [`BeLax` config flag]

[`Len`]({{< ref "Len" >}})
: checks an array, slice, map, string or channel length

[`Lt`]({{< ref "Lt" >}})
: checks that a number, string or [`time.Time`] is lesser than a value

[`Lte`]({{< ref "Lte" >}})
: checks that a number, string or [`time.Time`] is lesser or equal than a value

[`Map`]({{< ref "Map" >}})
: compares the contents of a map

[`MapEach`]({{< ref "MapEach" >}})
: compares each map entry

[`N`]({{< ref "N" >}})
: compares a number with a tolerance value

[`NaN`]({{< ref "NaN" >}})
: checks a floating number is [`math.NaN`]

[`Nil`]({{< ref "Nil" >}})
: compares to `nil`

[`None`]({{< ref "None" >}})
: no values have to match

[`Not`]({{< ref "Not" >}})
: value must not match

[`NotAny`]({{< ref "NotAny" >}})
: compares the contents of an array or a slice, no values have to match

[`NotEmpty`]({{< ref "NotEmpty" >}})
: checks that an array, a channel, a map, a slice or a string is not empty

[`NotNaN`]({{< ref "NotNaN" >}})
: checks a floating number is not [`math.NaN`]

[`NotNil`]({{< ref "NotNil" >}})
: checks that data is not `nil`

[`NotZero`]({{< ref "NotZero" >}})
: checks that data is not zero regarding its type

[`PPtr`]({{< ref "PPtr" >}})
: allows to easily test a pointer of pointer value

[`Ptr`]({{< ref "Ptr" >}})
: allows to easily test a pointer value

[`Re`]({{< ref "Re" >}})
: allows to apply a regexp on a string (or convertible), `[]byte`, [`error`] or [`fmt.Stringer`] interfaces, and even test the captured groups

[`ReAll`]({{< ref "ReAll" >}})
: allows to successively apply a regexp on a string (or convertible), `[]byte`, [`error`] or [`fmt.Stringer`] interfaces, and even test the captured groups

[`Set`]({{< ref "Set" >}})
: compares the contents of an array or a slice ignoring duplicates and without taking care of the order of items

[`Shallow`]({{< ref "Shallow" >}})
: compares pointers only, not their contents

[`Slice`]({{< ref "Slice" >}})
: compares the contents of a slice or a pointer on a slice

[`Smuggle`]({{< ref "Smuggle" >}})
: changes data contents or mutates it into another type via a custom function or a struct fields-path before stepping down in favor of generic comparison process

[`SStruct`]({{< ref "SStruct" >}})
: strictly compares the contents of a struct or a pointer on a struct

[`String`]({{< ref "String" >}})
: checks a string, `[]byte`, [`error`] or [`fmt.Stringer`] interfaces string contents

[`Struct`]({{< ref "Struct" >}})
: compares the contents of a struct or a pointer on a struct

[`SubBagOf`]({{< ref "SubBagOf" >}})
: compares the contents of an array or a slice without taking care of the order of items but with potentially some exclusions

[`SubJSONOf`]({{< ref "SubJSONOf" >}})
: compares struct or map against JSON representation but with potentially some exclusions

[`SubMapOf`]({{< ref "SubMapOf" >}})
: compares the contents of a map but with potentially some exclusions

[`SubSetOf`]({{< ref "SubSetOf" >}})
: compares the contents of an array or a slice ignoring duplicates and without taking care of the order of items but with potentially some exclusions

[`SuperBagOf`]({{< ref "SuperBagOf" >}})
: compares the contents of an array or a slice without taking care of the order of items but with potentially some extra items

[`SuperJSONOf`]({{< ref "SuperJSONOf" >}})
: compares struct or map against JSON representation but with potentially extra entries

[`SuperMapOf`]({{< ref "SuperMapOf" >}})
: compares the contents of a map but with potentially some extra entries

[`SuperSetOf`]({{< ref "SuperSetOf" >}})
: compares the contents of an array or a slice ignoring duplicates and without taking care of the order of items but with potentially some extra items

[`Tag`]({{< ref "Tag" >}})
: names an operator or a value. Only useful as a parameter of JSON operator, to name placeholders

[`TruncTime`]({{< ref "TruncTime" >}})
: compares [`time.Time`] (or assignable) values after truncating them

[`Values`]({{< ref "Values" >}})
: checks values of a map

[`Zero`]({{< ref "Zero" >}})
: checks data against its zero'ed conterpart

<!-- operators:end -->


## Smuggler operators

A smuggler operator is an operator able to transform the value (by
changing its value or even its type) before comparing it.

The following operators are smuggler ones:

<!-- smugglers:begin -->
[`Cap`]({{< ref "Cap" >}})
: checks an array, slice or channel capacity

[`Catch`]({{< ref "Catch" >}})
: catches data on the fly before comparing it

[`Contains`]({{< ref "Contains" >}})
: checks that a string, `[]byte`, [`error`] or [`fmt.Stringer`] interfaces contain a rune, byte or a sub-string; or a slice contains a single value or a sub-slice; or an array or map contain a single value

[`ContainsKey`]({{< ref "ContainsKey" >}})
: checks that a map contains a key

[`Keys`]({{< ref "Keys" >}})
: checks keys of a map

[`Lax`]({{< ref "Lax" >}})
: temporarily enables [`BeLax` config flag]

[`Len`]({{< ref "Len" >}})
: checks an array, slice, map, string or channel length

[`PPtr`]({{< ref "PPtr" >}})
: allows to easily test a pointer of pointer value

[`Ptr`]({{< ref "Ptr" >}})
: allows to easily test a pointer value

[`Smuggle`]({{< ref "Smuggle" >}})
: changes data contents or mutates it into another type via a custom function or a struct fields-path before stepping down in favor of generic comparison process

[`Tag`]({{< ref "Tag" >}})
: names an operator or a value. Only useful as a parameter of JSON operator, to name placeholders

[`Values`]({{< ref "Values" >}})
: checks values of a map

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
<!-- smugglers:end -->

[`T.Cmp`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.Cmp
[`T.CmpLax`]: https://pkg.go.dev/github.com/maxatome/go-testdeep/td#T.CmpLax


## `TypeBehind` method

```go
TypeBehind() reflect.Type
```

This method returns the type handled by the operator or `nil` if it is
not known. tdhttp helper uses it to know how to unmarshal HTTP
responses bodies before comparing them using the operator.

It is usually not used outside the
[go-testdeep repository](https://github.com/maxatome/go-testdeep).
