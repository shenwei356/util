package byteutil

import (
	"bytes"
	// "fmt"
	"unsafe"
)

// ReverseByteSlice reverses a byte slice
func ReverseByteSlice(s []byte) []byte {
	// make a copy of s
	l := len(s)
	t := make([]byte, l)
	for i := 0; i < l; i++ {
		t[i] = s[i]
	}

	// reverse
	for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}
	return t
}

// ReverseByteSliceInplace reverses a byte slice
func ReverseByteSliceInplace(s []byte) {
	// reverse
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// WrapByteSlice wraps byte slice
func WrapByteSlice(s []byte, width int) []byte {
	if width < 1 {
		return s
	}
	l := len(s)
	if l == 0 {
		return []byte{}
	}
	var lines int
	if l%width == 0 {
		lines = l/width - 1
	} else {
		lines = int(l / width)
	}
	// var buffer bytes.Buffer
	buffer := bytes.NewBuffer(make([]byte, 0, l+lines))
	var start, end int
	for i := 0; i <= lines; i++ {
		start = i * width
		end = (i + 1) * width
		if end > l {
			end = l
		}

		buffer.Write(s[start:end])
		if i < lines {
			buffer.WriteString("\n")
		}
	}
	return buffer.Bytes()
}

// WrapByteSliceInplace wraps byte slice in place.
// Sadly, it's too slow. Never use this!
func WrapByteSliceInplace(s []byte, width int) []byte {
	if width < 1 {
		return s
	}
	var l, lines int

	l = len(s)
	if l%width == 0 {
		lines = l/width - 1
	} else {
		lines = int(l / width)
	}

	var end int
	j := 0
	for i := 0; i <= lines; i++ {
		end = (i+1)*width + j
		if end >= l {
			break
		}
		// fmt.Printf("len:%d, lines:%d, i:%d, j:%d, end:%d\n", l, lines, i, j, end)
		if i < lines {

			// https://github.com/golang/go/wiki/SliceTricks
			// Sadly, it's too slow
			// s = append(s, []byte(" ")[0])
			// copy(s[end+1:], s[end:])
			// s[end] = []byte("\n")[0]

			// slow too
			s = append(s[:end], append([]byte("\n"), s[end:]...)...)

			l = len(s)
			if l%width == 0 {
				lines = l/width - 1
			} else {
				lines = int(l / width)
			}

			j++
		}
	}
	return s
}

// SubSlice provides similar slice indexing as python with one exception
// that end could be equal to 0.
// So we could get the last element by SubSlice(s, -1, 0)
// or get the whole element by SubSlice(s, 0, 0)
func SubSlice(slice []byte, start int, end int) []byte {
	if start == 0 && end == 0 {
		return slice
	}
	if start == end || (start < 0 && end > 0) {
		return []byte{}
	}
	l := len(slice)
	s, e := start, end

	if s < 0 {
		s = l + s
		if s < 1 {
			s = 0
		}
	}
	if e < 0 {
		e = l + e
		if e < 0 {
			e = 0
		}
	}
	if e == 0 || e > l {
		e = l
	}
	return slice[s:e]
}

// ByteToLower lowers a byte
func ByteToLower(b byte) byte {
	if b <= '\u007F' {
		if 'A' <= b && b <= 'Z' {
			b += 'a' - 'A'
		}
		return b
	}
	return b
}

// ByteToUpper upper a byte
func ByteToUpper(b byte) byte {
	if b <= '\u007F' {
		if 'a' <= b && b <= 'z' {
			b -= 'a' - 'A'
		}
		return b
	}
	return b
}

// MakeQuerySlice is used to replace map.
// see: http://blog.shenwei.me/map-is-not-the-fastest-in-go/
func MakeQuerySlice(letters []byte) []byte {
	max := -1
	for i := 0; i < len(letters); i++ {
		j := int(letters[i])
		if max < j {
			max = j
		}
	}
	querySlice := make([]byte, max+1)
	for i := 0; i < len(letters); i++ {
		querySlice[int(letters[i])] = letters[i]
	}
	return querySlice
}

// Split splits a byte slice by giveen letters.
// It's much faster than regexp.Split
func Split(slice []byte, letters []byte) [][]byte {
	querySlice := MakeQuerySlice(letters)
	results := [][]byte{}
	tmp := []byte{}

	var j int
	var value byte
	var sliceSize = len(querySlice)
	for _, b := range slice {
		j = int(b)
		if j >= sliceSize { // not delimiter byte
			tmp = append(tmp, b)
			continue
		}
		value = querySlice[j]
		if value == 0 { // not delimiter byte
			tmp = append(tmp, b)
			continue
		} else {
			if len(tmp) > 0 {
				results = append(results, tmp)
				tmp = []byte{}
			}
		}
	}
	if len(tmp) > 0 {
		results = append(results, tmp)
	}
	return results
}

// Bytes2Str convert byte slice to string without GC
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
