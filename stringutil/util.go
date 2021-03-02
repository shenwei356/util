package stringutil

import (
	"bytes"
	"regexp"
	"unsafe"
)

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

// ReverseStringSliceInplace reverses StringSlice
func ReverseStringSliceInplace(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// EscapeSymbols escape custom symbols
func EscapeSymbols(s, symbols string) string {
	m := make(map[rune]struct{})
	for _, c := range symbols {
		m[c] = struct{}{}
	}
	var buf bytes.Buffer
	var ok bool
	for _, c := range s {
		if _, ok = m[c]; ok {
			buf.WriteByte('\\')
		}
		buf.WriteRune(c)
	}
	return buf.String()
}

// UnEscaper returns a function for unescaping string
func UnEscaper() func(s string) string {
	var re = regexp.MustCompile(`\\([abfnrtv'"?])`)
	var m = map[string]string{
		`\a`: "\a",
		`\b`: "\b",
		`\f`: "\f",
		`\n`: "\n",
		`\r`: "\r",
		`\t`: "\t",
		`\v`: "\v",
		`\\`: "\\",
		`\'`: "'",
		`\"`: "\"",
		`\?`: "?",
	}
	var mapping = func(key string) string {
		if v, ok := m[key]; ok {
			return v
		}
		return key
	}

	return func(s string) string {
		return re.ReplaceAllStringFunc(s, mapping)
	}
}
