package numbers

import (
	"math/rand"
	"testing"
)

var data []*[]uint64
var data2 []*[]uint64

func init() {
	N := 100000 // N lists
	n := 300    // n elements for a list
	data = make([]*[]uint64, N)
	data2 = make([]*[]uint64, N)
	for i := 0; i < N; i++ {
		_n := rand.Intn(n)
		if _n < 0 {
			_n = -_n
		}
		_data := make([]uint64, _n)
		for j := 0; j < _n; j++ {
			_data[j] = uint64(float64(rand.Intn(_n)) / float64(1))
		}
		_data2 := make([]uint64, _n)
		copy(_data2, _data)

		data[i] = &_data
		data2[i] = &_data2
	}
}

func TestUniq(t *testing.T) {
	for _, _data := range data {
		u1 := Uniq(_data)

		c := make([]uint64, len(*_data))
		copy(c, *_data)

		UniqInplace(&c)

		// fmt.Printf("original: %v\n", _data)
		// fmt.Printf(" Inplace: %v\n", c)
		// fmt.Printf("    uniq: %v\n", *u1)
		if !Equal(*u1, c) {
			// fmt.Printf("original: %v\n", _data)
			// fmt.Printf(" Inplace: %v\n", c)
			// fmt.Printf("    uniq: %v\n", *u1)
			t.Error("error")
		}
	}
}

var result *[]uint64

func BenchmarkUniq(b *testing.B) {
	var _result *[]uint64
	// for i := 0; i < b.N; i++ {
	for _, _data := range data {
		_result = Uniq(_data)
	}
	// }
	result = _result
}

func BenchmarkUniqInplace(b *testing.B) {
	var _result *[]uint64
	// for i := 0; i < b.N; i++ {
	for _, _data := range data2 {
		UniqInplace(_data)
	}
	// }
	result = _result
}
