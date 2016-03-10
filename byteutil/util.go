package byteutil

import "bytes"

// ReverseByteSlice reverses a byte slice
func ReverseByteSlice(s []byte) []byte {
	// make a copy of s
	l := len(s)
	t := make([]byte, l)
	for i := 0; i < l; i++ {
		t[i] = s[i]
	}

	// reverse
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}
	return t
}

// WrapByteSlice wraps byte slice
func WrapByteSlice(s []byte, width int) []byte {
	if width < 1 {
		return s
	}
	var buffer bytes.Buffer
	l := len(s)
	var lines int
	if l%width == 0 {
		lines = l/width - 1
	} else {
		lines = int(l / width)
	}
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
