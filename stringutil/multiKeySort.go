package stringutil

import (
	"strconv"
	"strings"
)

// SortType defines the sort type
type SortType struct {
	Index   int
	Number  bool
	Reverse bool
}

// MultiKeyStringSlice sort [][]string by multiple keys
type MultiKeyStringSlice struct {
	SortTypes *[]SortType
	Value     []string
}

// MultiKeyStringSliceList is slice of MultiKeyStringSlice
type MultiKeyStringSliceList []MultiKeyStringSlice

func (list MultiKeyStringSliceList) Len() int { return len(list) }
func (list MultiKeyStringSliceList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list MultiKeyStringSliceList) Less(i, j int) bool {
	for _, t := range *list[i].SortTypes {
		var v int
		if t.Number {
			a, err := strconv.ParseFloat(list[i].Value[t.Index], 64)
			if err != nil {
				a = 0
			}
			b, err := strconv.ParseFloat(list[j].Value[t.Index], 64)
			if err != nil {
				b = 0
			}
			if a < b {
				v = -1
			} else if a == b {
				v = 0
			} else {
				v = 1
			}
		} else {
			v = strings.Compare(list[i].Value[t.Index], list[j].Value[t.Index])
		}

		if v == 0 {
		} else if v < 0 {
			if t.Reverse {
				return false
			}
			return true
		} else {
			if t.Reverse {
				return true
			}
			return false
		}
	}
	return true
}
