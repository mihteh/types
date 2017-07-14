package decimal

import (
	"testing"
)

// TestSort checks sorting of Decimals
func TestSort(t *testing.T) {
	testCase := Decimals{F(-1.25), Zero, F(1.25), F(-1.26), F(-1.259), F(1), F(-1.261), F(1.249), F(1.251), F(1)}
	testCase.Sort()
	expectedResult := Decimals{F(-1.261), F(-1.26), F(-1.259), F(-1.25), Zero, F(1), F(1), F(1.249), F(1.25), F(1.251)}
	if !testCase.Equal(expectedResult, true, false) {
		t.Fatalf("Sorting don't work as expected: %v", testCase)
	}
}

// TestEqual checks two Decimals equality
func TestEqual(t *testing.T) {
	oldComparePrecision := comparePrecision.Float64f()
	defer SetComparePrecision(oldComparePrecision)
	SetComparePrecision(0.00999999999)

	type testCase struct {
		input1         Decimals
		input2         Decimals
		ordered        bool
		precise        bool
		expectedResult bool
	}

	testCases := []testCase{
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.25)}, true, true, true},
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.251)}, true, false, true},
		testCase{Decimals{F(-1.251), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.25)}, true, false, true},
		testCase{Decimals{F(-1.249), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.25)}, true, false, true},
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.251)}, true, true, false},
		testCase{Decimals{F(-1.251), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.25)}, true, true, false},
		testCase{Decimals{F(-1.249), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.25)}, true, true, false},
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), F(1.25), Zero}, true, true, false},
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), F(1.25), Zero}, false, true, true},
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), F(1.25), Zero, Zero}, false, false, false},
	}

	for i, tc := range testCases {
		result := tc.input1.Equal(tc.input2, tc.ordered, tc.precise)
		if result != tc.expectedResult {
			t.Fatalf("Failed at index %d, ordered=%v, precise=%v, result: %v, expected: %v",
				i, tc.ordered, tc.precise, result, tc.expectedResult)
		}
	}
}

// TestRemoveDuplicates checks removing duplicates from Decimals
func TestRemoveDuplicates(t *testing.T) {
	oldComparePrecision := comparePrecision.Float64f()
	defer SetComparePrecision(oldComparePrecision)
	SetComparePrecision(0.00999999999)

	testCase := Decimals{F(-1.25), Zero, F(1.25), F(-1.26), F(-1.259), F(1), F(-1.261), F(1.249), F(1.251), F(1),
		F(-1.26), F(-1.261), F(-1.259)}
	expectedResult := Decimals{F(-1.25), Zero, F(1.25), F(-1.26), F(1)}
	result := testCase.RemoveDuplicates()
	if !result.Equal(expectedResult, false, true) {
		t.Fatalf("RemoveDuplicates don't work as expected: %v", result)
	}
}
