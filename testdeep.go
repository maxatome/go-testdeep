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
// check the record content:
//
//   before := time.Now()
//   record, err := CreateRecord()
//
//   if err != nil {
//     t.Errorf("An error occurred: %s", err)
//   } else {
//     expected := Record{Name: "Bob", Age: 23}
//
//     if record.Id == 0 {
//       t.Error("Id probably not initialized")
//     }
//     if before.After(record.CreatedAt) ||
//       time.Now().Before(record.CreatedAt) {
//       t.Errorf("CreatedAt field not expected: %s", record.CreatedAt)
//     }
//     if record.Name != expected.Name {
//       t.Errorf("Name field differ, got=%s, expected=%s",
//         record.Name, expected.Name)
//     }
//     if record.Age != expected.Age {
//       t.Errorf("Age field differ, got=%s, expected=%s",
//         record.Age, expected.Age)
//     }
//   }
//
// With testdeep, it is a way simple, thanks to CmpDeeply function:
//
//   import (
//     td "github.com/maxatome/go-testdeep"
//   )
//
//   ...
//
//   before := time.Now()
//   record, err := CreateRecord()
//
//   if td.CmpDeeply(t, err, nil) {
//     td.CmpDeeply(t, record,
//       Struct(
//         Record{
//           Name: "Bob",
//           Age:  23,
//         },
//         StructFields{
//           Id:        td.Not(0),
//           CreatedAt: td.Between(before, time.Now()),
//         }),
//       "Newly created record")
//   }
//
// Of course not only structs can be compared. A lot of operators can
// be found below to cover most (all?) needed tests.
//
// The CmpDeeply function is the keystone of this package, but to make
// the writing of tests even easier, the family of Cmp* functions are
// provided and act as shortcuts. Using CmpStruct function, the
// previous example can be written as:
//
//   before := time.Now()
//   record, err := CreateRecord()
//
//   if td.CmpDeeply(t, err, nil) {
//     td.CmpStruct(t, record,
//       Record{
//         Name: "Bob",
//         Age:  23,
//       },
//       StructFields{
//         Id:        td.Not(0),
//         CreatedAt: td.Between(before, time.Now()),
//       },
//       "Newly created record")
//   }
package testdeep
