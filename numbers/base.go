package numbers

type Uint64Slice []uint64

func (s Uint64Slice) Len() int           { return len(s) }
func (s Uint64Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s Uint64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *Uint64Slice) Push(x interface{}) {
	*s = append(*s, x.(uint64))
}

func (s *Uint64Slice) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

func Equal(s, t []uint64) bool {
	if len(s) != len(t) {
		return false
	}

	for i, v := range s {
		if v != t[i] {
			return false
		}
	}

	return true
}
