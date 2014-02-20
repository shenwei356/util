// Copyright 2014 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

// Package bytesize provides a way to show readable values of byte size
// by reediting the code from http://golang.org/doc/effective_go.html.
// It could also parseng byte size text to ByteSize object.
package bytesize

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ByteSize float64

const (
	B ByteSize = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

// Print readable values of byte size
func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%7.2f YB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%7.2f ZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%7.2f EB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%7.2f PB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%7.2f TB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%7.2f GB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%7.2f MB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%7.2f KB", b/KB)
	}
	return fmt.Sprintf("%7.2f  B", b)
}

// Regexp object for ByteSize Text
var BytesizeRegexp = regexp.MustCompile(`(?i)^\s*([\-\d\.]+)\s*(B|KB|MB|GB|TB|PB|EB|ZB|YB)\s*$`)

// Error information for Illegal byte size text
var ErrText = "Illegal bytesize text"

// Parse ByteSize Text to ByteSize object
//
// Example
//
//     size, err := bytesize.Parse([]byte("1.5 KB"))
//     if err != nil {
//         fmt.Println(err)
//     }
//     fmt.Printf("%.0f bytes\n", size)
//
func Parse(sizeText []byte) (ByteSize, error) {
	if !BytesizeRegexp.Match(sizeText) {
		return 0, errors.New(ErrText)
	}

	// parse value and unit
	subs := BytesizeRegexp.FindSubmatch(sizeText)

	// no need to check ParseFloat error. BytesizeRegexp could ensure this
	size, _ := strconv.ParseFloat(string(subs[1]), 64)
	unit := strings.ToUpper(string(subs[2]))

	switch unit {
	case "B":
		size = size * float64(B)
	case "KB":
		size = size * float64(KB)
	case "MB":
		size = size * float64(MB)
	case "GB":
		size = size * float64(GB)
	case "TB":
		size = size * float64(TB)
	case "PB":
		size = size * float64(PB)
	case "EB":
		size = size * float64(EB)
	case "ZB":
		size = size * float64(ZB)
	case "YB":
		size = size * float64(YB)
	}

	return ByteSize(size), nil
}
