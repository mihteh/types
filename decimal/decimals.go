package decimal

import (
	"sort"
)

type Decimals []Decimal

func (decimals Decimals) Len() int           { return len(decimals) }
func (decimals Decimals) Less(i, j int) bool { return decimals[i].Cmp(decimals[j]) < 0 }
func (decimals Decimals) Swap(i, j int)      { decimals[i], decimals[j] = decimals[j], decimals[i] }
func (decimals Decimals) Sort()              { sort.Sort(decimals) }

// Copy copies a slice of decimals
func (decimals Decimals) Copy() Decimals {
	result := make(Decimals, len(decimals))
	for i, d := range decimals {
		result[i] = d.Copy()
	}
	return result
}

// Equal compares two slices, decimals and decimals2, of Decimal
// if ordered is true, comparing order also, otherwise order may differ
// if slices are equal according to parameter conditions, it returns true, otherwise returns false
func (decimals Decimals) Equal(decimals2 Decimals, ordered bool) bool {
	if len(decimals) != len(decimals2) {
		return false
	}

	copy1, copy2 := decimals.Copy(), decimals2.Copy()
	if !ordered {
		copy1.Sort()
		copy2.Sort()
	}

	for i, _ := range copy1 {
		if copy1[i].Ne(copy2[i]) {
			return false
		}
	}

	return true
}

// RemoveDuplicates returns a slice decimals which don't contain duplicate elements from decimals slice
func (decimals Decimals) RemoveDuplicates() Decimals {
	copy := decimals.Copy()
	result := Decimals{}
	for i := 0; i < len(copy); i++ {
		exists := false
		for j := 0; j < i; j++ {
			if copy[j].Eq(copy[i]) {
				exists = true
				break
			}
		}
		if !exists { // append this one if no previous same element exists
			result = append(result, copy[i])
		}
	}
	return result
}
