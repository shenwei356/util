package stringutil

import (
	"sort"
	"testing"
)

func TestMultiKeySortUnicodeLength(t *testing.T) {
	types := []SortType{{Index: 0}, {Index: 1, Length: true, Reverse: true}, {Index: 2, Number: true}}
	values := [][]string{
		{"A", "中", "3"},
		{"A", "aa", "2"},
		{"A", "a", "1"},
		{"B", "é", "5"},
		{"B", "中文", "6"},
		{"B", "b", "4"},
	}
	list := make(MultiKeyStringSliceList, len(values))
	for i, value := range values {
		list[i] = MultiKeyStringSlice{SortTypes: &types, Value: value}
	}
	sort.Sort(list)
	wantIDs := []string{"2", "1", "3", "6", "4", "5"}
	for i, want := range wantIDs {
		if got := list[i].Value[2]; got != want {
			t.Fatalf("row %d: got ID %s, want %s", i, got, want)
		}
	}
}

func TestMultiKeySortLengthThenText(t *testing.T) {
	types := []SortType{{Index: 0, Length: true}, {Index: 0}}
	list := MultiKeyStringSliceList{
		{SortTypes: &types, Value: []string{"中"}},
		{SortTypes: &types, Value: []string{"aa"}},
		{SortTypes: &types, Value: []string{"é"}},
		{SortTypes: &types, Value: []string{"a"}},
	}
	sort.Sort(list)
	want := []string{"a", "é", "中", "aa"}
	for i, value := range want {
		if got := list[i].Value[0]; got != value {
			t.Fatalf("row %d: got %q, want %q", i, got, value)
		}
	}
}
