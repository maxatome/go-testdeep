// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// Package testdeep allows extremely flexible deep comparison, built
// for testing.
//
// It is a go rewrite and adaptation of wonderful Test::Deep perl
// module (see https://metacpan.org/pod/Test::Deep).
//
// In golang, comparing data structure is usually done using
// reflect.DeepEqual or using a package that uses this function behind
// the scene.
//
// This function works very well, but it is not flexible. Both
// compared structures must match exactly.
//
// The purpose of testdeep package is to do its best to introduce this
// missing flexibility using "operators" when the expected value (or
// one of its component) cannot be matched exactly.
//
// Imagine a function returning a struct containing a newly created
// database record. The Id and the CreatedAt fields are set by the
// database layer. In this case we have to do something like that to
// check the record contents:
//
//   import (
//     "testing"
//   )
//
//   type Record struct {
//     Id        uint64
//     Name      string
//     Age       int
//     CreatedAt time.Time
//   }
//
//   func CreateRecord(name string, age int) (*Record, error) {
//     ...
//   }
//
//   func TestCreateRecord(t *testing.T) {
//     before := time.Now()
//     record, err := CreateRecord()
//
//     if err != nil {
//       t.Errorf("An error occurred: %s", err)
//     } else {
//       expected := Record{Name: "Bob", Age: 23}
//
//       if record.Id == 0 {
//         t.Error("Id probably not initialized")
//       }
//       if before.After(record.CreatedAt) ||
//         time.Now().Before(record.CreatedAt) {
//         t.Errorf("CreatedAt field not expected: %s", record.CreatedAt)
//       }
//       if record.Name != expected.Name {
//         t.Errorf("Name field differ, got=%s, expected=%s",
//           record.Name, expected.Name)
//       }
//       if record.Age != expected.Age {
//         t.Errorf("Age field differ, got=%s, expected=%s",
//           record.Age, expected.Age)
//       }
//     }
//   }
//
// With testdeep, it is a way simple, thanks to CmpDeeply and
// CmpNoError functions:
//
//   import (
//     "testing"
//     td "github.com/maxatome/go-testdeep"
//   )
//
//   ...
//
//   func TestCreateRecord(t *testing.T) {
//     before := time.Now()
//     record, err := CreateRecord()
//
//     if td.CmpNoError(t, err) {
//       td.CmpDeeply(t, record,
//         Struct(
//           &Record{
//             Name: "Bob",
//             Age:  23,
//           },
//           td.StructFields{
//             Id:        td.NotZero(),
//             CreatedAt: td.Between(before, time.Now()),
//           }),
//         "Newly created record")
//     }
//   }
//
// Of course not only structs can be compared. A lot of operators can
// be found below to cover most (all?) needed tests. See
// https://godoc.org/github.com/maxatome/go-testdeep#TestDeep
//
// The CmpDeeply function is the keystone of this package, but to make
// the writing of tests even easier, the family of Cmp* functions are
// provided and act as shortcuts. Using CmpStruct function, the
// previous example can be written as:
//
//   import (
//     "testing"
//     td "github.com/maxatome/go-testdeep"
//   )
//
//   ...
//
//   func TestCreateRecord(t *testing.T) {
//     before := time.Now()
//     record, err := CreateRecord()
//
//     if td.CmpNoError(t, err) {
//       td.CmpStruct(t, record,
//         &Record{
//           Name: "Bob",
//           Age:  23,
//         },
//         td.StructFields{
//           Id:        td.NotZero(),
//           CreatedAt: td.Between(before, time.Now()),
//         },
//         "Newly created record")
//     }
//   }
//
// Last, testing.T can be encapsulated in testdeep T type, simplifying
// again the test:
//
//   import (
//     "testing"
//     td "github.com/maxatome/go-testdeep"
//   )
//
//   ...
//
//   func TestCreateRecord(tt *testing.T) {
//     t := td.NewT(tt)
//
//     before := time.Now()
//     record, err := CreateRecord()
//
//     if t.CmpNoError(err) {
//       t.RootName("RECORD").Struct(record,
//         &Record{
//           Name: "Bob",
//           Age:  23,
//         },
//         td.StructFields{
//           Id:        td.NotZero(),
//           CreatedAt: td.Between(before, time.Now()),
//         },
//         "Newly created record")
//     }
//   }
package testdeep // import "github.com/maxatome/go-testdeep"
