ParallelMap
===========

A lock-free parallel map in go.

ParallelMap uses a backend goroutine for the sequential excution of 
Get and Set or custom function, which was inspired by section 14.17
in book *The Way to Go*.

Install
-------
This package is "go-gettable", just:

    go get github.com/shenwei356/util/ParallelMap

Example
-------
    
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
    	m := pmap.NewParallelMap()
    	
    	// In this exmaple, the Update function will be used.
    	// To call this function, the UpdateValueFunc must be specified.
    	m.SetUpdateValueFunc(func(oldValue pmap.ValueType, newValue pmap.ValueType) pmap.ValueType {
    		return oldValue.(int) + newValue.(int)
    	})
    
    	// number of elements in map
    	var n int = 1 << 9
    
    	var wg sync.WaitGroup
    	for i := 1; i <= N; i++ {
    		wg.Add(1)
    
    		go func() {
    			defer func() {
    				wg.Done()
    			}()
    
    			for j := 0; j < n; j++ {
    				m.Update(j, 1)
    			}
    		}()
    	}
    
    	// wait for all operations to complement
    	wg.Wait()
    	// Stop the map backend
    	m.Stop()
    
    	// do something else
    	length := len(m.Map)
    	fmt.Printf("%d elements in map\n", length)
    }
    
 

Documentation
-------------

[See documentation on gowalker for more detail](http://gowalker.org/github.com/shenwei356/util/ParallelMap).

[MIT License](https://github.com/shenwei356/util/blob/master/ParallelMap/LICENSE)