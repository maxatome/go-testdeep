package testdeep

import (
	"reflect"
)

// copyValue copies the first level of val in a new reflect.Value instance.
func copyValue(val reflect.Value) (reflect.Value, bool) {
	if val.Kind() == reflect.Ptr {
		refVal, ok := copyValue(val.Elem())
		if !ok {
			return reflect.Value{}, false
		}
		newPtrVal := reflect.New(refVal.Type())
		newPtrVal.Elem().Set(refVal)
		return newPtrVal, true
	}

	var newVal reflect.Value

	switch val.Kind() {
	case reflect.Bool:
		newPtrVal := reflect.New(val.Type())
		newVal = newPtrVal.Elem()
		newVal.SetBool(val.Bool())

	case reflect.Complex64, reflect.Complex128:
		newPtrVal := reflect.New(val.Type())
		newVal = newPtrVal.Elem()
		newVal.SetComplex(val.Complex())

	case reflect.Float32, reflect.Float64:
		newPtrVal := reflect.New(val.Type())
		newVal = newPtrVal.Elem()
		newVal.SetFloat(val.Float())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		newPtrVal := reflect.New(val.Type())
		newVal = newPtrVal.Elem()
		newVal.SetInt(val.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		newPtrVal := reflect.New(val.Type())
		newVal = newPtrVal.Elem()
		newVal.SetUint(val.Uint())

	case reflect.Array:
		newPtrVal := reflect.New(val.Type())
		newVal = newPtrVal.Elem()
		var (
			item reflect.Value
			ok   bool
		)
		for i := val.Len() - 1; i >= 0; i-- {
			item, ok = copyValue(val.Index(i))
			if !ok {
				return reflect.Value{}, false
			}
			newVal.Index(i).Set(item)
		}

	case reflect.Slice:
		newVal = reflect.MakeSlice(val.Type(), val.Len(), val.Cap())
		var (
			item reflect.Value
			ok   bool
		)
		for i := val.Len() - 1; i >= 0; i-- {
			item, ok = copyValue(val.Index(i))
			if !ok {
				return reflect.Value{}, false
			}
			newVal.Index(i).Set(item)
		}

	case reflect.Map:
		newVal = reflect.MakeMapWithSize(val.Type(), val.Len())
		var (
			key, value reflect.Value
			ok         bool
		)
		for _, keyVal := range val.MapKeys() {
			key, ok = copyValue(keyVal)
			if !ok {
				return reflect.Value{}, false
			}
			value, ok = copyValue(val.MapIndex(key))
			if !ok {
				return reflect.Value{}, false
			}
			newVal.SetMapIndex(key, value)
		}

	case reflect.Interface:
		newPtrVal := reflect.New(val.Type())
		newVal = newPtrVal.Elem()
		refVal, ok := copyValue(val.Elem())
		if !ok {
			return reflect.Value{}, false
		}
		newVal.Set(refVal)

	case reflect.String:
		newPtrVal := reflect.New(val.Type())
		newVal = newPtrVal.Elem()
		newVal.SetString(val.String())

	default:
		return reflect.Value{}, false
	}
	return newVal, true
}
