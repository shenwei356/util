// Copyright 2013 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

// Package bytesize provides a way to show readable values of byte sizes
// by reediting the code from http://golang.org/doc/effective_go.html.
package bytesize

import (
	"fmt"
)

type ByteSize float64

const (
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

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
