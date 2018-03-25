package testdeep_test

import (
	"fmt"
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestSet(t *testing.T) {
	type MyArray [5]int
	type MySlice []int

	for idx, got := range []interface{}{
		[]int{1, 3, 4, 4, 5},
		[...]int{1, 3, 4, 4, 5},
		MySlice{1, 3, 4, 4, 5},
		MyArray{1, 3, 4, 4, 5},
		&MySlice{1, 3, 4, 4, 5},
		&MyArray{1, 3, 4, 4, 5},
	} {
		testName := fmt.Sprintf("Test #%d → %v", idx, got)

		//
		// Set
		checkOK(t, got, Set(5, 4, 1, 3), testName)
		checkOK(t, got, Set(5, 4, 1, 3, 3, 3, 3), testName) // duplicated fields

		checkError(t, got, Set(5, 4, 1),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra items: ((int) 3)"),
			},
			testName)

		checkError(t, got, Set(5, 4, 1, 3, 66),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66)"),
			},
			testName)

		checkError(t, got, Set(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66)"),
			},
			testName)

		checkError(t, got, Set(5, 66, 4, 1, 3, 67),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66,\n                (int) 67)"),
			},
			testName)

		checkError(t, got, Set(5, 66, 4, 3),
			expectedError{
				Message: mustBe("comparing %% as a Set"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66)\n  Extra items: ((int) 1)"),
			},
			testName)

		//
		// SubSetOf
		checkOK(t, got, SubSetOf(5, 4, 1, 3), testName)
		checkOK(t, got, SubSetOf(5, 4, 1, 3, 66), testName)

		checkError(t, got, SubSetOf(5, 66, 4, 3),
			expectedError{
				Message: mustBe("comparing %% as a SubSetOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra items: ((int) 1)"),
			},
			testName)

		//
		// SuperSetOf
		checkOK(t, got, SuperSetOf(5, 4, 1, 3), testName)
		checkOK(t, got, SuperSetOf(5, 4), testName)

		checkError(t, got, SuperSetOf(5, 66, 4, 1, 3),
			expectedError{
				Message: mustBe("comparing %% as a SuperSetOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Missing items: ((int) 66)"),
			},
			testName)

		//
		// NoneOf
		checkOK(t, got, NoneOf(10, 20, 30), testName)

		checkError(t, got, NoneOf(3, 66),
			expectedError{
				Message: mustBe("comparing %% as a NoneOf"),
				Path:    mustBe("DATA"),
				Summary: mustBe("Extra items: ((int) 3)"),
			},
			testName)
	}

	checkOK(t, []interface{}{123, "foo", nil, "bar", nil},
		Set("foo", "bar", 123, nil))

	var nilSlice MySlice
	for idx, got := range []interface{}{([]int)(nil), &nilSlice} {
		testName := fmt.Sprintf("Test #%d", idx)

		checkOK(t, got, Set(), testName)
		checkOK(t, got, SubSetOf(), testName)
		checkOK(t, got, SubSetOf(1, 2), testName)
		checkOK(t, got, SuperSetOf(), testName)
		checkOK(t, got, NoneOf(), testName)
		checkOK(t, got, NoneOf(1, 2), testName)
	}

	for idx, set := range []TestDeep{
		Set(123),
		SubSetOf(123),
		SuperSetOf(123),
		NoneOf(123),
	} {
		testName := fmt.Sprintf("Test #%d → %s", idx, set)

		checkError(t, 123, set,
			expectedError{
				Message:  mustBe("bad type"),
				Path:     mustBe("DATA"),
				Got:      mustBe("int"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)

		num := 123
		checkError(t, &num, set,
			expectedError{
				Message:  mustBe("bad type"),
				Path:     mustBe("DATA"),
				Got:      mustBe("*int"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)

		var list *MySlice
		checkError(t, list, set,
			expectedError{
				Message:  mustBe("nil pointer"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil *testdeep_test.MySlice"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)

		checkError(t, nil, set,
			expectedError{
				Message:  mustBe("bad type"),
				Path:     mustBe("DATA"),
				Got:      mustBe("nil"),
				Expected: mustBe("Slice OR Array OR *Slice OR *Array"),
			},
			testName)
	}
	//
	// String
	equalStr(t, Set(1).String(), "Set((int) 1)")
	equalStr(t, Set(1, 2).String(), "Set((int) 1,\n    (int) 2)")

	equalStr(t, SubSetOf(1).String(), "SubSetOf((int) 1)")
	equalStr(t, SubSetOf(1, 2).String(), "SubSetOf((int) 1,\n         (int) 2)")

	equalStr(t, SuperSetOf(1).String(), "SuperSetOf((int) 1)")
	equalStr(t, SuperSetOf(1, 2).String(),
		"SuperSetOf((int) 1,\n           (int) 2)")

	equalStr(t, NoneOf(1).String(), "NoneOf((int) 1)")
	equalStr(t, NoneOf(1, 2).String(), "NoneOf((int) 1,\n       (int) 2)")
}
