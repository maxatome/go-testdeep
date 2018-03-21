package testdeep

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

type boundCmp uint8

const (
	boundNone boundCmp = iota
	boundIn
	boundOut
)

type tdBetween struct {
	TestDeepBase
	expectedMin reflect.Value
	expectedMax reflect.Value

	minBound boundCmp
	maxBound boundCmp
}

var _ TestDeep = &tdBetween{}

type BoundsKind uint8

const (
	BoundsInIn BoundsKind = iota
	BoundsInOut
	BoundsOutIn
	BoundsOutOut
)

type tdBetweenTime struct {
	tdBetween
	expectedType reflect.Type
	mustConvert  bool
}

var _ TestDeep = &tdBetweenTime{}

func Between(from interface{}, to interface{}, bounds ...BoundsKind) TestDeep {
	b := tdBetween{
		expectedMin: reflect.ValueOf(from),
		expectedMax: reflect.ValueOf(to),
	}

	const usage = "usage: Between(NUM|TIME, NUM|TIME[, BOUNDS_KIND])"

	if len(bounds) > 0 {
		if len(bounds) > 1 {
			panic(usage)
		}

		if bounds[0] == BoundsInIn || bounds[0] == BoundsInOut {
			b.minBound = boundIn
		} else {
			b.minBound = boundOut
		}

		if bounds[0] == BoundsInIn || bounds[0] == BoundsOutIn {
			b.maxBound = boundIn
		} else {
			b.maxBound = boundOut
		}
	} else {
		b.minBound = boundIn
		b.maxBound = boundIn
	}

	if b.expectedMax.Type() != b.expectedMin.Type() {
		panic("from and to params must have the same type")
	}

	return b.initBetween(usage)
}

func (b *tdBetween) initBetween(usage string) TestDeep {
	b.TestDeepBase = NewTestDeepBase(4)

	if !b.expectedMax.IsValid() {
		b.expectedMax = b.expectedMin
	}

	switch b.expectedMin.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if b.expectedMin.Int() > b.expectedMax.Int() {
			b.expectedMin, b.expectedMax = b.expectedMax, b.expectedMin
		}
		return b

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if b.expectedMin.Uint() > b.expectedMax.Uint() {
			b.expectedMin, b.expectedMax = b.expectedMax, b.expectedMin
		}
		return b

	case reflect.Float32, reflect.Float64:
		if b.expectedMin.Float() > b.expectedMax.Float() {
			b.expectedMin, b.expectedMax = b.expectedMax, b.expectedMin
		}
		return b

	case reflect.Struct:
		var bt tdBetweenTime
		if b.expectedMin.Type() == timeType {
			bt = tdBetweenTime{
				tdBetween:    *b,
				expectedType: timeType,
			}
		} else if b.expectedMin.Type().ConvertibleTo(timeType) {
			bt = tdBetweenTime{
				tdBetween:    *b,
				expectedType: b.expectedMin.Type(),
				mustConvert:  true,
			}
			bt.expectedMin = b.expectedMin.Convert(timeType)
			bt.expectedMax = b.expectedMax.Convert(timeType)
		}

		if bt.expectedMin.Interface().(time.Time).After(
			bt.expectedMax.Interface().(time.Time)) {
			bt.expectedMin, bt.expectedMax = bt.expectedMax, bt.expectedMin
		}

		return &bt
	}
	panic(usage)
}

func N(num interface{}, tolerance ...interface{}) TestDeep {
	n := tdBetween{
		TestDeepBase: NewTestDeepBase(3),
		expectedMin:  reflect.ValueOf(num),
		minBound:     boundIn,
		maxBound:     boundIn,
	}

	const usage = "usage: N({,U}INT{,8,16,32,64}|FLOAT{32,64}[, TOLERANCE])"

	switch n.expectedMin.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
	default:
		panic(usage)
	}

	n.expectedMax = n.expectedMin

	if len(tolerance) > 0 {
		if len(tolerance) > 1 {
			panic(usage)
		}

		tol := reflect.ValueOf(tolerance[0])
		if tol.Type() != n.expectedMin.Type() {
			panic("tolerance param must have the same type as num one")
		}

		switch tol.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if diff := tol.Int(); diff != 0 {
				base := n.expectedMin.Int()

				max := base + diff
				if max < base {
					max = math.MaxInt64
				}

				min := base - diff
				if min > base {
					min = math.MinInt64
				}

				n.expectedMin = reflect.New(tol.Type()).Elem()
				n.expectedMin.SetInt(min)

				n.expectedMax = reflect.New(tol.Type()).Elem()
				n.expectedMax.SetInt(max)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
			reflect.Uint64:
			if diff := tol.Uint(); diff != 0 {
				base := n.expectedMin.Uint()

				max := base + diff
				if max < base {
					max = math.MaxUint64
				}

				min := base - diff
				if min > base {
					min = 0
				}

				n.expectedMin = reflect.New(tol.Type()).Elem()
				n.expectedMin.SetUint(min)

				n.expectedMax = reflect.New(tol.Type()).Elem()
				n.expectedMax.SetUint(max)
			}

		default: // case reflect.Float32, reflect.Float64:
			if diff := tol.Float(); diff != 0 {
				base := n.expectedMin.Float()

				n.expectedMin = reflect.New(tol.Type()).Elem()
				n.expectedMin.SetFloat(base - diff)

				n.expectedMax = reflect.New(tol.Type()).Elem()
				n.expectedMax.SetFloat(base + diff)
			}
		}
	}

	return &n
}

func Gt(val interface{}) TestDeep {
	b := &tdBetween{
		expectedMin: reflect.ValueOf(val),
		minBound:    boundOut,
	}
	return b.initBetween("usage: Gt(NUM|TIME)")
}

func Gte(val interface{}) TestDeep {
	b := &tdBetween{
		expectedMin: reflect.ValueOf(val),
		minBound:    boundIn,
	}
	return b.initBetween("usage: Gte(NUM|TIME)")
}

func Lt(val interface{}) TestDeep {
	b := &tdBetween{
		expectedMin: reflect.ValueOf(val),
		maxBound:    boundOut,
	}
	return b.initBetween("usage: Lt(NUM|TIME)")
}

func Lte(val interface{}) TestDeep {
	b := &tdBetween{
		expectedMin: reflect.ValueOf(val),
		maxBound:    boundIn,
	}
	return b.initBetween("usage: Lte(NUM|TIME)")
}

func (b *tdBetween) Match(ctx Context, got reflect.Value) *Error {
	if got.Type() != b.expectedMin.Type() {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "type mismatch",
			Got:      rawString(got.Type().String()),
			Expected: rawString(b.expectedMin.Type().String()),
			Location: b.GetLocation(),
		}
	}

	var ok bool

	switch got.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch b.minBound {
		case boundIn:
			ok = got.Int() >= b.expectedMin.Int()
		case boundOut:
			ok = got.Int() > b.expectedMin.Int()
		default:
			ok = true
		}
		if ok {
			switch b.maxBound {
			case boundIn:
				ok = got.Int() <= b.expectedMax.Int()
			case boundOut:
				ok = got.Int() < b.expectedMax.Int()
			default:
				ok = true
			}
		}

	case reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch b.minBound {
		case boundIn:
			ok = got.Uint() >= b.expectedMin.Uint()
		case boundOut:
			ok = got.Uint() > b.expectedMin.Uint()
		default:
			ok = true
		}
		if ok {
			switch b.maxBound {
			case boundIn:
				ok = got.Uint() <= b.expectedMax.Uint()
			case boundOut:
				ok = got.Uint() < b.expectedMax.Uint()
			default:
				ok = true
			}
		}

	default: // == case reflect.Float32, reflect.Float64:
		switch b.minBound {
		case boundIn:
			ok = got.Float() >= b.expectedMin.Float()
		case boundOut:
			ok = got.Float() > b.expectedMin.Float()
		default:
			ok = true
		}
		if ok {
			switch b.maxBound {
			case boundIn:
				ok = got.Float() <= b.expectedMax.Float()
			case boundOut:
				ok = got.Float() < b.expectedMax.Float()
			default:
				ok = true
			}
		}
	}

	if ok {
		return nil
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "values differ",
		Got:      rawString(fmt.Sprintf("%v", got.Interface())),
		Expected: rawString(b.String()),
		Location: b.GetLocation(),
	}
}

func (b *tdBetween) String() string {
	var min, max interface{}

	if b.minBound != boundNone {
		min = b.expectedMin.Interface()
	}
	if b.maxBound != boundNone {
		max = b.expectedMax.Interface()
	}

	if min == max {
		return fmt.Sprintf("%v", min)
	}

	if min != nil {
		if max != nil {
			return fmt.Sprintf("%v %c got %c %v",
				min, ternRune(b.minBound == boundIn, '≤', '<'),
				ternRune(b.maxBound == boundIn, '≤', '<'), max)
		}

		return fmt.Sprintf("%c %v",
			ternRune(b.minBound == boundIn, '≥', '>'), min)
	}

	return fmt.Sprintf("%c %v",
		ternRune(b.maxBound == boundIn, '≤', '<'), max)
}

var _ TestDeep = &tdBetweenTime{}

func (b *tdBetweenTime) Match(ctx Context, got reflect.Value) *Error {
	if got.Type() != b.expectedType {
		if ctx.booleanError {
			return booleanError
		}
		return &Error{
			Context:  ctx,
			Message:  "type mismatch",
			Got:      rawString(got.Type().String()),
			Expected: rawString(b.expectedType.String()),
			Location: b.GetLocation(),
		}
	}

	var cmpGot time.Time
	if b.mustConvert {
		cmpGot = got.Convert(timeType).Interface().(time.Time)
	} else {
		cmpGot = got.Interface().(time.Time)
	}

	var ok bool
	if b.minBound != boundNone {
		min := b.expectedMin.Interface().(time.Time)

		if b.minBound == boundIn {
			ok = min.Before(cmpGot)
		} else {
			ok = cmpGot.After(min)
		}
	} else {
		ok = true
	}

	if ok && b.maxBound != boundNone {
		max := b.expectedMax.Interface().(time.Time)

		if b.maxBound == boundIn {
			ok = max.After(cmpGot)
		} else {
			ok = cmpGot.Before(max)
		}
	}

	if ok {
		return nil
	}

	if ctx.booleanError {
		return booleanError
	}
	return &Error{
		Context:  ctx,
		Message:  "values differ",
		Got:      rawString(fmt.Sprintf("%v", got.Interface())),
		Expected: rawString(b.String()),
		Location: b.GetLocation(),
	}
}
