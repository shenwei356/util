// Copyright 2014 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.
package ParallelMap

import (
	"fmt"
	"runtime"
	"testing"
)

func TestParallelMap(t *testing.T) {
	N := runtime.NumCPU() * 10
	runtime.GOMAXPROCS(N)

	m := NewParallelMap()

	var length int = 1 << 10
	c := make(chan int)
	for i := 1; i <= N; i++ {
		go IntValuePlusOne(m, length, c)
	}

	// wait for finish
	for i := 1; i <= N; i++ {
		<-c
	}

	// check length of map
	if len(m.Map) != length {
		t.Error("length error")
	}

	// check values
	for _, v := range m.Map {
		if v.(int) != int(N) {
			t.Error(fmt.Sprintf("value error: %d != %d", v.(int), int(N)))
		}
	}
}

func IntValuePlusOne(m *ParallelMap, length int, c chan int) {
	for i := 0; i < length; i++ {
		m.ExecuteFunc(func() {
			if v, ok := m.Map[i]; ok {
				m.Map[i] = v.(int) + 1
			} else {
				m.Map[i] = int(1)
			}
		})
	}
	c <- 1
}
