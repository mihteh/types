package decimal

import (
	"testing"
)

// TestSort checks sorting of Decimals
func TestSort(t *testing.T) {
	testCase := Decimals{F(-1.25), Zero, F(1.25), F(-1.26), F(-1.259), F(1), F(-1.261), F(1.249), F(1.251), F(1)}
	testCase.Sort()
	expectedResult := Decimals{F(-1.261), F(-1.26), F(-1.259), F(-1.25), Zero, F(1), F(1), F(1.249), F(1.25), F(1.251)}
	if !testCase.Equal(expectedResult, true) {
		t.Fatalf("Sorting don't work as expected: %v", testCase)
	}
}

// TestEqual checks two Decimals equality
func TestEqual(t *testing.T) {
	type testCase struct {
		input1         Decimals
		input2         Decimals
		ordered        bool
		expectedResult bool
	}

	testCases := []testCase{
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.25)}, true, true},
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.251)}, true, false},
		testCase{Decimals{F(-1.251), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.25)}, true, false},
		testCase{Decimals{F(-1.249), Zero, F(1.25)}, Decimals{F(-1.25), Zero, F(1.25)}, true, false},
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), F(1.25), Zero}, true, false},
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), F(1.25), Zero}, false, true},
		testCase{Decimals{F(-1.25), Zero, F(1.25)}, Decimals{F(-1.25), F(1.25), Zero, Zero}, false, false},
	}

	for i, tc := range testCases {
		result := tc.input1.Equal(tc.input2, tc.ordered)
		if result != tc.expectedResult {
			t.Fatalf("Failed at index %d, ordered=%v, result: %v, expected: %v",
				i, tc.ordered, result, tc.expectedResult)
		}
	}
}

// TestRemoveDuplicates checks removing duplicates from Decimals
func TestRemoveDuplicates(t *testing.T) {
	testCase := Decimals{F(-1.25), Zero, F(1.25), F(-1.26), F(-1.259), F(1), F(-1.261), F(1.249), F(1.251), F(1),
		F(-1.26), F(-1.261), F(-1.259)}
	expectedResult := Decimals{F(-1.25), Zero, F(1.25), F(-1.26), F(-1.259), F(1), F(-1.261), F(1.249), F(1.251)}
	result := testCase.RemoveDuplicates()
	if !result.Equal(expectedResult, false) {
		t.Fatalf("RemoveDuplicates don't work as expected: %v", result)
	}
}
