package testdeep_test

import (
	"testing"

	. "github.com/maxatome/go-testdeep"
)

func TestIsa(t *testing.T) {
	checkOK(t, &gotStruct, Isa(&MyStruct{}))
	checkOK(t, (*MyStruct)(nil), Isa(&MyStruct{}))
	checkOK(t, (*MyStruct)(nil), Isa((*MyStruct)(nil)))
	checkOK(t, gotStruct, Isa(MyStruct{}))

	checkError(t, &gotStruct, Isa(&MyStructBase{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("*testdeep_test.MyStruct"),
		Expected: mustContain("*testdeep_test.MyStructBase"),
	})

	checkError(t, (*MyStruct)(nil), Isa(&MyStructBase{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("*testdeep_test.MyStruct"),
		Expected: mustContain("*testdeep_test.MyStructBase"),
	})

	checkError(t, gotStruct, Isa(&MyStruct{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("testdeep_test.MyStruct"),
		Expected: mustContain("*testdeep_test.MyStruct"),
	})

	checkError(t, &gotStruct, Isa(MyStruct{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("*testdeep_test.MyStruct"),
		Expected: mustContain("testdeep_test.MyStruct"),
	})

	gotSlice := []int{1, 2, 3}
	checkOK(t, gotSlice, Isa([]int{}))
	checkOK(t, &gotSlice, Isa(((*[]int)(nil))))

	checkError(t, &gotSlice, Isa([]int{}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("*[]int"),
		Expected: mustContain("[]int"),
	})

	checkError(t, gotSlice, Isa((*[]int)(nil)), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("[]int"),
		Expected: mustContain("*[]int"),
	})

	checkError(t, gotSlice, Isa([1]int{2}), expectedError{
		Message:  mustBe("type mismatch"),
		Path:     mustBe("DATA"),
		Got:      mustContain("[]int"),
		Expected: mustContain("[1]int"),
	})

	//
	// String
	equalStr(t, Isa((*MyStruct)(nil)).String(), "*testdeep_test.MyStruct")
}
