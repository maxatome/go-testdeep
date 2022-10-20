// Copyright (c) 2021, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package tdhttp

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/maxatome/go-testdeep/internal/color"
)

// Q allows to easily declare query parameters for use in [NewRequest]
// and related [http.Request] builders, as [Get] for example:
//
//	req := tdhttp.Get("/path", tdhttp.Q{
//	  "id":     []int64{1234, 4567},
//	  "dryrun": true,
//	})
//
// See [NewRequest] for several examples of use.
//
// Accepted types as values are:
//   - [fmt.Stringer]
//   - string
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - bool
//   - slice or array of any type above, plus any
//   - pointer on any type above, plus any or any other pointer
type Q map[string]any

var _ URLValuesEncoder = Q(nil)

// AddTo adds the q contents to qp.
func (q Q) AddTo(qp url.Values) error {
	for param, value := range q {
		// Ignore nil values
		if value == nil {
			continue
		}
		err := q.addParamTo(param, reflect.ValueOf(value), true, qp)
		if err != nil {
			return err
		}
	}
	return nil
}

// Values returns a [url.Values] instance corresponding to q. It panics
// if a value cannot be converted.
func (q Q) Values() url.Values {
	qp := make(url.Values, len(q))
	err := q.AddTo(qp)
	if err != nil {
		panic(errors.New(color.Bad(err.Error())))
	}
	return qp
}

// Encode does the same as [url.Values.Encode] does. So quoting its doc,
// it encodes the values into “URL encoded” form ("bar=baz&foo=quux")
// sorted by key.
//
// It panics if a value cannot be converted.
func (q Q) Encode() string {
	return q.Values().Encode()
}

func (q Q) addParamTo(param string, v reflect.Value, allowArray bool, qp url.Values) error {
	var str string
	for {
		if s, ok := v.Interface().(fmt.Stringer); ok {
			qp.Add(param, s.String())
			return nil
		}

		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			if !allowArray {
				return fmt.Errorf("%s is only allowed at the root level for param %q",
					v.Kind(), param)
			}
			for i, l := 0, v.Len(); i < l; i++ {
				err := q.addParamTo(param, v.Index(i), false, qp)
				if err != nil {
					return err
				}
			}
			return nil

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str = strconv.FormatInt(v.Int(), 10)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			str = strconv.FormatUint(v.Uint(), 10)

		case reflect.Float32, reflect.Float64:
			str = strconv.FormatFloat(v.Float(), 'g', -1, 64)

		case reflect.String:
			str = v.String()

		case reflect.Bool:
			str = strconv.FormatBool(v.Bool())

		case reflect.Ptr, reflect.Interface:
			if !v.IsNil() {
				v = v.Elem()
				continue
			}
			return nil // mimic url.Values behavior ⇒ ignore

		default:
			return fmt.Errorf("don't know how to add type %s (%s) to param %q",
				v.Type(), v.Kind(), param)
		}

		qp.Add(param, str)
		return nil
	}
}
