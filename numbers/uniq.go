package numbers

import "sort"

// Uniq removes duplicated elements in the list, and returns a new one.
func Uniq(list *[]uint64) *[]uint64 {
	if len(*list) == 0 {
		return &[]uint64{}
	} else if len(*list) == 1 {
		return &[]uint64{(*list)[0]}
	}

	sort.Sort(Uint64Slice(*list))

	s := make([]uint64, 0, len(*list))
	p := (*list)[0]
	s = append(s, p)
	for _, v := range (*list)[1:] {
		if v != p {
			s = append(s, v)
		}
		p = v
	}
	return &s
}

// UniqInplace is faster than Uniq for short slice (<1000).
func UniqInplace(list *[]uint64) {
	if len(*list) == 0 || len(*list) == 1 {
		return
	}

	sort.Sort(Uint64Slice(*list))

	var i, j int
	var p, v uint64
	var flag bool
	p = (*list)[0]
	for i = 1; i < len(*list); i++ {
		v = (*list)[i]
		if v == p {
			if !flag {
				j = i // mark insertion position
				flag = true
			}
			continue
		}

		if flag { // need to insert to previous position
			(*list)[j] = v
			j++
		}
		p = v
	}
	if j > 0 {
		*list = (*list)[:j]
	}
}
