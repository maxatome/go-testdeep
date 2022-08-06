// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// Package testdeep allows extremely flexible deep comparison. It is
// built for testing.
//
// It is a go rewrite and adaptation of wonderful [Test::Deep] perl
// module.
//
// In golang, comparing data structure is usually done using
// [reflect.DeepEqual] or using a package that uses this function
// behind the scene.
//
// This function works very well, but it is not flexible. Both
// compared structures must match exactly.
//
// The purpose of go-testdeep is to do its best to introduce this
// missing flexibility using ["operators"] when the expected value (or
// one of its component) cannot be matched exactly.
//
// testdeep package should not be used in new code, even if it can for
// backward compatibility reasons, but [td] package.
//
// All variables and types of testdeep package are aliases to
// respectively functions and types of [td] package. They are only
// here for compatibility purpose as
//
//	import "github.com/maxatome/go-testdeep/td"
//
// should now be used, in preference of older, but still supported:
//
//	import td "github.com/maxatome/go-testdeep"
//
// For easy HTTP API testing, see [tdhttp] package.
//
// For tests suites also just as easy, see [tdsuite] package.
//
// [Test::Deep]: https://metacpan.org/pod/Test::Deep
// ["operators"]: https://go-testdeep.zetta.rocks/operators/
// [tdhttp]: https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdhttp
// [tdsuite]: https://pkg.go.dev/github.com/maxatome/go-testdeep/helpers/tdsuite
package testdeep // import "github.com/maxatome/go-testdeep"
