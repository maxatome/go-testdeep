// Copyright (c) 2019, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package util

import (
	"encoding/json"
	"fmt"
	"sync"
	"unicode"
	"unicode/utf8"
)

var jsonErrorMesg = "<NO_JSON_ERROR!>" // will be overwritten by UnmarshalJSON
var jsonErrorMesgOnce sync.Once

// `{"foo": $bar}`, 8 → `{"foo": "$bar"}`
func stringifyPlaceholder(buf []byte, dollar int64) ([]byte, error) {
	r, size := utf8.DecodeRune(buf[dollar+1:]) // just after $
	cur := dollar + 1 + int64(size)

	var end int64

	// Numeric placeholder: $1234
	if r >= '0' && r <= '9' {
		for i, c := range buf[cur:] {
			switch c {
			case ' ', '\t', '\r', '\n', ',', '}', ']':
				end = cur + int64(i)
				goto endFound
			default:
				if c < '0' || c > '9' {
					return nil,
						fmt.Errorf(`invalid numeric placeholder at offset %d`, dollar+1)
				}
			}
		}
		end = int64(len(buf))
	endFound:
	} else if unicode.IsLetter(r) || r == '_' { // Named placeholder: $pïpô12
	runes:
		for max := int64(len(buf)); cur < max; cur += int64(size) {
			r, size = utf8.DecodeRune(buf[cur:])
			switch r {
			case '_':
			case ' ', '\t', '\r', '\n', ',', '}', ']':
				break runes
			default:
				if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
					return nil,
						fmt.Errorf(`invalid named placeholder at offset %d`, dollar+1)
				}
			}
		}
		end = cur
	} else {
		return nil, fmt.Errorf(`invalid placeholder at offset %d`, dollar+1)
	}

	// put "" around $éé123 or $12345
	if cap(buf) == len(buf) {
		// allocate room for 20 extra placeholders
		buf = append(make([]byte, 0, len(buf)+40), buf...)
	}
	buf = append(buf, 0, 0)
	copy(buf[end+2:], buf[end:])
	buf[end+1] = '"'
	copy(buf[dollar+1:], buf[dollar:end])
	buf[dollar] = '"'

	return buf, nil
}

// UnmarshalJSON is a custom json.Unmarshal function allowing to
// handle placeholders not enclosed in strings. It relies on
// json.SyntaxError errors detected before any memory allocation. So
// the performance should not be too bad, avoiding to implement our
// own JSON parser...
func UnmarshalJSON(buf []byte, target interface{}) error {
	jsonErrorMesgOnce.Do(func() {
		var dummy interface{}
		err := json.Unmarshal([]byte(`$x`), &dummy)
		if jerr, ok := err.(*json.SyntaxError); ok {
			jsonErrorMesg = jerr.Error()
		}
	})

	for {
		err := json.Unmarshal(buf, target)
		if err == nil {
			return nil
		}
		jerr, ok := err.(*json.SyntaxError)
		if !ok || jerr.Error() != jsonErrorMesg || jerr.Offset >= int64(len(buf)) {
			return err
		}

		buf, err = stringifyPlaceholder(buf, jerr.Offset-1) // $ pos
		if err != nil {
			return err
		}
	}
}
