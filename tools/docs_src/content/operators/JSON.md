---
title: "JSON"
weight: 10
---

```go
func JSON(expectedJSON interface{}, params ...interface{}) TestDeep
```

[`JSON`]({{< ref "JSON" >}}) operator allows to compare the JSON representation of data
against *expectedJSON*. *expectedJSON* can be a:

- `string` containing JSON data like `{"fullname":"Bob","age":42}`
- `string` containing a JSON filename, ending with ".json" (its
  content is [`ioutil.ReadFile`](https://golang.org/pkg/ioutil/#ReadFile) before unmarshaling)
- `[]byte` containing JSON data
- [`io.Reader`](https://golang.org/pkg/io/#Reader) stream containing JSON data (is [`ioutil.ReadAll`](https://golang.org/pkg/ioutil/#ReadAll) before
  unmarshaling)


*expectedJSON* JSON value can contain placeholders. The *params*
are for any placeholder parameters in *expectedJSON*. *params* can
contain [TestDeep operators]({{< ref "operators" >}}) as well as raw values. A placeholder can
be numeric like `$2` or named like `$name` and always references an
item in *params*.

Numeric placeholders reference the n'th "operators" item (starting
at 1). Named placeholders are used with [`Tag`]({{< ref "Tag" >}}) operator as follows:

```go
Cmp(t, gotValue,
  JSON(`{"fullname": $name, "age": $2, "gender": $3}`,
    Tag("name", HasPrefix("Foo")), // matches $1 and $name
    Between(41, 43),               // matches only $2
    "male"))                       // matches only $3
```

Note that placeholders can be double-quoted as in:

```go
Cmp(t, gotValue,
  JSON(`{"fullname": "$name", "age": "$2", "gender": "$3"}`,
    Tag("name", HasPrefix("Foo")), // matches $1 and $name
    Between(41, 43),               // matches only $2
    "male"))                       // matches only $3
```

It makes no difference whatever the underlying type of the replaced
item is (= double quoting a placeholder matching a number is not a
problem). It is just a matter of taste, double-quoting placeholders
can be preferred when the JSON data has to conform to the [`JSON`]({{< ref "JSON" >}})
specification, like when used in a ".json" file.

Note *expectedJSON* can be a `[]byte`, JSON filename or [`io.Reader`](https://golang.org/pkg/io/#Reader):

```go
Cmp(t, gotValue, JSON("file.json", Between(12, 34)))
Cmp(t, gotValue, JSON([]byte(`[1, $1, 3]`), Between(12, 34)))
Cmp(t, gotValue, JSON(osFile, Between(12, 34)))
```

A JSON filename ends with ".json".

To avoid a legit "$" `string` prefix cause a bad placeholder [`error`](https://golang.org/pkg/builtin/#error),
just double it to escape it. Note it is only needed when the "$" is
the first character of a `string`:

```go
Cmp(t, gotValue,
  JSON(`{"fullname": "$name", "details": "$$info", "age": $2}`,
    Tag("name", HasPrefix("Foo")), // matches $1 and $name
    Between(41, 43)))              // matches only $2
```

For the "details" key, the raw value "`$info`" is expected, no
placeholders are involved here.

Note that [`Lax`]({{< ref "Lax" >}}) mode is automatically enabled by [`JSON`]({{< ref "JSON" >}}) operator to
simplify numeric tests.

Last but not least, comments can be embedded in JSON data:

```go
Cmp(t, gotValue,
  JSON(`
{
  // A guy properties:
  "fullname": "$name",  // The full name of the guy
  "details":  "$$info", // Literally "$info", thanks to "$" escape
  "age":      $2        /* The age of the guy:
                           - placeholder unquoted, but could be without
                             any change
                           - to demonstrate a multi-lines comment */
}`,
    Tag("name", HasPrefix("Foo")), // matches $1 and $name
    Between(41, 43)))              // matches only $2
```

Comments, like in go, have 2 forms. To quote the Go language specification:

- line comments start with the character sequence // and stop at the
  end of the line.
- multi-lines comments start with the character sequence /* and stop
  with the first subsequent character sequence */.


[TypeBehind]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of the *expectedJSON*
[`json.Unmarshal`](https://golang.org/pkg/json/#Unmarshal)'ed. So it can be `bool`, `string`, `float64`,
`[]interface{}`, `map[string]interface{}` or `interface{}` in case
*expectedJSON* is "null".


> See also [<i class='fas fa-book'></i> JSON godoc](https://godoc.org/github.com/maxatome/go-testdeep#JSON).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := Cmp(t, got, JSON(`{"age":42,"fullname":"Bob"}`))
	fmt.Println("check got with age then fullname:", ok)

	ok = Cmp(t, got, JSON(`{"fullname":"Bob","age":42}`))
	fmt.Println("check got with fullname then age:", ok)

	ok = Cmp(t, got, JSON(`
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`))
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = Cmp(t, got, JSON(`{"fullname":"Bob","age":42,"gender":"male"}`))
	fmt.Println("check got with gender field:", ok)

	ok = Cmp(t, got, JSON(`{"fullname":"Bob"}`))
	fmt.Println("check got with fullname only:", ok)

	ok = Cmp(t, true, JSON(`true`))
	fmt.Println("check boolean got is true:", ok)

	ok = Cmp(t, 42, JSON(`42`))
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = Cmp(t, got, JSON(`null`))
	fmt.Println("check nil got is null:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with gender field: false
	// check got with fullname only: false
	// check boolean got is true: true
	// check numeric got is 42: true
	// check nil got is null: true

```{{% /expand%}}
{{%expand "Placeholders example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := Cmp(t, got, JSON(`{"age": $1, "fullname": $2}`, 42, "Bob Foobar"))
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = Cmp(t, got,
		JSON(`{"age": $1, "fullname": $2}`,
			Between(40, 45),
			HasSuffix("Foobar")))
	fmt.Println("check got with numeric placeholders:", ok)

	ok = Cmp(t, got,
		JSON(`{"age": "$1", "fullname": "$2"}`,
			Between(40, 45),
			HasSuffix("Foobar")))
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = Cmp(t, got,
		JSON(`{"age": $age, "fullname": $name}`,
			Tag("age", Between(40, 45)),
			Tag("name", HasSuffix("Foobar"))))
	fmt.Println("check got with named placeholders:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true

```{{% /expand%}}
{{%expand "File example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = ioutil.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := Cmp(t, got,
		JSON(filename,
			Tag("name", HasPrefix("Bob")),
			Tag("age", Between(40, 45)),
			Tag("gender", Re(`^(male|female)\z`))))
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = Cmp(t, got,
		JSON(file,
			Tag("name", HasPrefix("Bob")),
			Tag("age", Between(40, 45)),
			Tag("gender", Re(`^(male|female)\z`))))
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true

```{{% /expand%}}
## CmpJSON shortcut

```go
func CmpJSON(t TestingT, got interface{}, expectedJSON interface{}, params []interface{}, args ...interface{}) bool
```

CmpJSON is a shortcut for:

```go
Cmp(t, got, JSON(expectedJSON, params...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> CmpJSON godoc](https://godoc.org/github.com/maxatome/go-testdeep#CmpJSON).

### Examples

{{%expand "Basic example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := CmpJSON(t, got, `{"age":42,"fullname":"Bob"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = CmpJSON(t, got, `{"fullname":"Bob","age":42}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = CmpJSON(t, got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = CmpJSON(t, got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with gender field:", ok)

	ok = CmpJSON(t, got, `{"fullname":"Bob"}`, nil)
	fmt.Println("check got with fullname only:", ok)

	ok = CmpJSON(t, true, `true`, nil)
	fmt.Println("check boolean got is true:", ok)

	ok = CmpJSON(t, 42, `42`, nil)
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = CmpJSON(t, got, `null`, nil)
	fmt.Println("check nil got is null:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with gender field: false
	// check got with fullname only: false
	// check boolean got is true: true
	// check numeric got is 42: true
	// check nil got is null: true

```{{% /expand%}}
{{%expand "Placeholders example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := CmpJSON(t, got, `{"age": $1, "fullname": $2}`, []interface{}{42, "Bob Foobar"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = CmpJSON(t, got, `{"age": $1, "fullname": $2}`, []interface{}{Between(40, 45), HasSuffix("Foobar")})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = CmpJSON(t, got, `{"age": "$1", "fullname": "$2"}`, []interface{}{Between(40, 45), HasSuffix("Foobar")})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = CmpJSON(t, got, `{"age": $age, "fullname": $name}`, []interface{}{Tag("age", Between(40, 45)), Tag("name", HasSuffix("Foobar"))})
	fmt.Println("check got with named placeholders:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true

```{{% /expand%}}
{{%expand "File example" %}}```go
	t := &testing.T{}

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = ioutil.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := CmpJSON(t, got, filename, []interface{}{Tag("name", HasPrefix("Bob")), Tag("age", Between(40, 45)), Tag("gender", Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = CmpJSON(t, got, file, []interface{}{Tag("name", HasPrefix("Bob")), Tag("age", Between(40, 45)), Tag("gender", Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true

```{{% /expand%}}
## T.JSON shortcut

```go
func (t *T) JSON(got interface{}, expectedJSON interface{}, params []interface{}, args ...interface{}) bool
```

[`JSON`]({{< ref "JSON" >}}) is a shortcut for:

```go
t.Cmp(got, JSON(expectedJSON, params...), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


> See also [<i class='fas fa-book'></i> T.JSON godoc](https://godoc.org/github.com/maxatome/go-testdeep#T.JSON).

### Examples

{{%expand "Basic example" %}}```go
	t := NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob",
		Age:      42,
	}

	ok := t.JSON(got, `{"age":42,"fullname":"Bob"}`, nil)
	fmt.Println("check got with age then fullname:", ok)

	ok = t.JSON(got, `{"fullname":"Bob","age":42}`, nil)
	fmt.Println("check got with fullname then age:", ok)

	ok = t.JSON(got, `
// This should be the JSON representation of a struct
{
  // A person:
  "fullname": "Bob", // The name of this person
  "age":      42     /* The age of this person:
                        - 42 of course
                        - to demonstrate a multi-lines comment */
}`, nil)
	fmt.Println("check got with nicely formatted and commented JSON:", ok)

	ok = t.JSON(got, `{"fullname":"Bob","age":42,"gender":"male"}`, nil)
	fmt.Println("check got with gender field:", ok)

	ok = t.JSON(got, `{"fullname":"Bob"}`, nil)
	fmt.Println("check got with fullname only:", ok)

	ok = t.JSON(true, `true`, nil)
	fmt.Println("check boolean got is true:", ok)

	ok = t.JSON(42, `42`, nil)
	fmt.Println("check numeric got is 42:", ok)

	got = nil
	ok = t.JSON(got, `null`, nil)
	fmt.Println("check nil got is null:", ok)

	// Output:
	// check got with age then fullname: true
	// check got with fullname then age: true
	// check got with nicely formatted and commented JSON: true
	// check got with gender field: false
	// check got with fullname only: false
	// check boolean got is true: true
	// check numeric got is 42: true
	// check nil got is null: true

```{{% /expand%}}
{{%expand "Placeholders example" %}}```go
	t := NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
	}

	ok := t.JSON(got, `{"age": $1, "fullname": $2}`, []interface{}{42, "Bob Foobar"})
	fmt.Println("check got with numeric placeholders without operators:", ok)

	ok = t.JSON(got, `{"age": $1, "fullname": $2}`, []interface{}{Between(40, 45), HasSuffix("Foobar")})
	fmt.Println("check got with numeric placeholders:", ok)

	ok = t.JSON(got, `{"age": "$1", "fullname": "$2"}`, []interface{}{Between(40, 45), HasSuffix("Foobar")})
	fmt.Println("check got with double-quoted numeric placeholders:", ok)

	ok = t.JSON(got, `{"age": $age, "fullname": $name}`, []interface{}{Tag("age", Between(40, 45)), Tag("name", HasSuffix("Foobar"))})
	fmt.Println("check got with named placeholders:", ok)

	// Output:
	// check got with numeric placeholders without operators: true
	// check got with numeric placeholders: true
	// check got with double-quoted numeric placeholders: true
	// check got with named placeholders: true

```{{% /expand%}}
{{%expand "File example" %}}```go
	t := NewT(&testing.T{})

	got := &struct {
		Fullname string `json:"fullname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
	}{
		Fullname: "Bob Foobar",
		Age:      42,
		Gender:   "male",
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // clean up

	filename := tmpDir + "/test.json"
	if err = ioutil.WriteFile(filename, []byte(`
{
  "fullname": "$name",
  "age":      "$age",
  "gender":   "$gender"
}`), 0644); err != nil {
		t.Fatal(err)
	}

	// OK let's test with this file
	ok := t.JSON(got, filename, []interface{}{Tag("name", HasPrefix("Bob")), Tag("age", Between(40, 45)), Tag("gender", Re(`^(male|female)\z`))})
	fmt.Println("Full match from file name:", ok)

	// When the file is already open
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	ok = t.JSON(got, file, []interface{}{Tag("name", HasPrefix("Bob")), Tag("age", Between(40, 45)), Tag("gender", Re(`^(male|female)\z`))})
	fmt.Println("Full match from io.Reader:", ok)

	// Output:
	// Full match from file name: true
	// Full match from io.Reader: true

```{{% /expand%}}
