# FAQ

## How to mix strict requirements and simple assertions

```golang
func TestAssertionsAndRequirements(t *testing.T) {
  assert := td.NewT(t)
  require := assert.FatalOnError()

  require.CmpDeeply(got, expected) // if it fails: report error + abort
  assert.CmpDeeply(got, expected)  // if it fails: report error + continue
}
```

## How to add a new operator

You want to add a new `FooBar` operator.

- check that another operator does not exist with the same meaning;
- add the operator definition in `td_foo_bar.go` file and fully
  document its usage;
- add operator tests in `td_foo_bar_test.go` file;
- in `example_test.go` file, add examples function(s) `ExampleFooBar*`
  in alphabetical order;
- automatically generate `CmpFooBar` & `T.FooBar` (+ examples) code:
  `./tools/gen_funcs.pl .`
- do not forget to run tests: `go test ./...`
- run `gometalinter` as in [`.travis.yml`](.travis.yml);
- add `FooBar` with a small description in
  [`README.md`](README.md#available-operators), respecting the
  alphabetical order.

Each time you change `example_test.go`, re-run `./tools/gen_funcs.pl .`
to update corresponding `CmpFooBar` & `T.FooBar` examples.

Test coverage must be 100%.
