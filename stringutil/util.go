package stringutil

import (
	"unsafe"

	"github.com/shenwei356/util/byteutil"
)

// Split splits a byte slice by giveen letters
func Split(slice string, letters string) []string {
	result := byteutil.Split([]byte(slice), []byte(letters))
	result2 := []string{}
	for _, s := range result {
		result2 = append(result2, string(s))
	}
	return result2
}

// Str2Bytes convert string to byte slice. Warning: it's unsafe!!!
func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// ReverseStringSlice reverses StringSlice
func ReverseStringSlice(s []string) []string {
	// make a copy of s
	l := len(s)
	t := make([]string, l)
	for i := 0; i < l; i++ {
		t[i] = s[i]
	}

	// reverse
	for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}
	return t
}
