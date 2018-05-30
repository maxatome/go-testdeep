// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"testing"
)

// T is a type that encapsulates testing.T allowing to easily use
// testing.T methods as well as T ones.
type T struct {
	*testing.T
}

// NewT returns a new T instance. Typically used as:
//
//   type Record struct {
//     Id        uint64
//     Name      string
//     Age       int
//     CreatedAt time.Time
//   }
//
//   func TestCreateRecord(tt *testing.T) {
//     t := NewT(tt)
//
//     before := time.Now()
//     record, err := CreateRecord()
//
//     if t.Nil(err) {
//       t.Log("No error, can now check struct contents")
//
//       ok := t.Struct(record,
//         Record{
//           Name: "Bob",
//           Age:  23,
//         },
//         StructFields{
//           Id:        td.Not(uint64(0)),
//           CreatedAt: td.Between(before, time.Now()),
//         },
//         "Newly created record")
//       if ok {
//         t.Log(Record created successfully!")
//       }
//     }
//   }
func NewT(t *testing.T) *T {
	return &T{
		T: t,
	}
}

// CmpDeeply is shortcut for:
//
//    CmpDeeply(t.T, got, expected, args...)
func (t *T) CmpDeeply(got, expected interface{}, args ...interface{}) bool {
	t.Helper()
	return CmpDeeply(t.T, got, expected, args...)
}
