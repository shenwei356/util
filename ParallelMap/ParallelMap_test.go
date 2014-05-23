// Copyright 2014 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.
package ParallelMap

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestParallelMap(t *testing.T) {
	N := runtime.NumCPU() * 30
	runtime.GOMAXPROCS(N)

	m := NewParallelMap(func(v1 ValueType, v2 ValueType) ValueType {
		return v1.(int) + v2.(int)
	}, N)

	var length int = 1 << 10

	var wg sync.WaitGroup
	for i := 1; i <= N; i++ {
		wg.Add(1)

		go func() {
			defer func() {
				wg.Done()
				m.UnboundAGoroutine()
			}()

			for j := 0; j < length; j++ {
				m.Update(j, 1)
			}
		}()
	}

	// wait for all operations to complement
	wg.Wait()
	m.Wait()

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
