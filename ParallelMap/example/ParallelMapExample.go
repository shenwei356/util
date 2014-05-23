package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"

	pmap "github.com/shenwei356/util/ParallelMap"
)

func main() {
	// number of goroutines that will operate on ParallelMap
	N := runtime.NumCPU() * 30
	runtime.GOMAXPROCS(N)

	// constructor
	m := pmap.NewParallelMap(func(v1 pmap.ValueType, v2 pmap.ValueType) pmap.ValueType {
		return v1.(int) + v2.(int)
	}, N)

	// number of elements in map
	var n int = 1 << 9

	var wg sync.WaitGroup
	for i := 1; i <= N; i++ {
		wg.Add(1)

		go func() {
			defer func() {
				wg.Done()
				m.UnboundAGoroutine()
			}()

			for j := 0; j < n; j++ {
				// update data in map
				m.Update(j, 1)
			}
		}()
	}

	// wait for all operations to complement
	wg.Wait()
	m.Wait()

	// do something else
	length := len(m.Map)
	fmt.Printf("%d elements in map\n", length)

	keys := make([]int, length)
	i := 0
	for k, _ := range m.Map {
		keys[i] = k.(int)
		i++
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Printf("%d => %d\n", k, m.Map[k])
	}
}
