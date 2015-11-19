package stringutil

import "sort"

// StringCount is a struct store count of String
type StringCount struct {
	String string
	Count  int
}

// StringCountList is slice of Stringcount
type StringCountList []StringCount

func (b StringCountList) Len() int { return len(b) }
func (b StringCountList) Less(i, j int) bool {
	// return b[i].Count < b[j].Count
	// This will return unwanted result: return b[i].Count < b[j].Count || b[i].String < b[j].String
	if b[i].Count < b[j].Count {
		return true
	}
	if b[i].Count == b[j].Count {
		if b[i].String < b[j].String {
			return true
		}
		return false
	}
	return false
}
func (b StringCountList) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// ReversedStringCountList is Reversed StringCountList
type ReversedStringCountList struct {
	StringCountList
}

// Less is different from the Less of StringCountList
func (b ReversedStringCountList) Less(i, j int) bool {
	// return b.StringCountList[i].Count > b.StringCountList[j].Count
	if b.StringCountList[i].Count > b.StringCountList[j].Count {
		return true
	}
	if b.StringCountList[i].Count == b.StringCountList[j].Count {
		if b.StringCountList[i].String < b.StringCountList[j].String {
			return true
		}
		return false
	}
	return false
}

// CountOfString returns the count of String for a String slice
func CountOfString(s []string) map[string]int {
	count := make(map[string]int)
	for _, b := range s {
		count[b]++
	}
	return count
}

// SortCountOfString sorts count of String
func SortCountOfString(count map[string]int, reverse bool) StringCountList {
	countList := make(StringCountList, len(count))
	i := 0
	for b, c := range count {
		countList[i] = StringCount{b, c}
		i++
	}
	if reverse {
		sort.Sort(ReversedStringCountList{countList})
	} else {
		sort.Sort(countList)
	}
	return countList
}
