// Copyright 2014 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

// ParallelMap - A lock-free parallel map in go.
// ParallelMap uses a backend goroutine for the sequential excution of
// Get and Set or custom function, which was inspired by section 14.17
// in book << The Way to Go >>.
package ParallelMap

// The type of the map key is interface{}
type keyType interface{}

// The type of the map value is interface{}
type valueType interface{}

// Usage
//
//    import (
//        "fmt"
//        pmap "github.com/shenwei356/util/ParallelMap"
//        "runtime"
//    )
//
//    func main() {
//        runtime.GOMAXPROCS(runtime.NumCPU())
//
//        m := pmap.NewParallelMap()
//        m.Set("year", 2014)
//
//        if v, ok := m.Get("year"); ok {
//            fmt.Println(v)
//        }
//    }
type ParallelMap struct {
	// map
	Map map[keyType]valueType
	// backend goroutine for the sequential excution
	Op chan func()
}

// NewParallelMap
func NewParallelMap() *ParallelMap {
	m := &ParallelMap{make(map[keyType]valueType), make(chan func())}
	go m.backend()
	return m
}

// Run operation channel as backend
func (this *ParallelMap) backend() {
	var f func()
	for {
		select {
		case f = <-this.Op:
			f()
		}
	}
}

// Getting element of the map is executed sequentially
func (this *ParallelMap) Get(key keyType) (valueType, bool) {
	/*v, ok := this.Map[key]
	return v, ok*/

	c1 := make(chan valueType)
	c2 := make(chan bool)
	this.Op <- func() {
		value, ok := this.Map[key]
		c1 <- value
		c2 <- ok
	}
	return <-c1, <-c2
}

// Setting operation is executed sequentially to ensure the
// operation is atomic.
func (this *ParallelMap) Set(key keyType, value valueType) {
	c := make(chan bool)
	this.Op <- func() {
		this.Map[key] = value
		c <- true
	}
	<-c
}

// Execute a custom function.
//
// Example: An element increasing function
//
//    m.ExecuteFunc(func() {
//        if v, ok := m.Map[i]; ok {
//            m.Map[i] = v.(int) + 1
//        } else {
//            m.Map[i] = int(1)
//        }
//    })
func (this *ParallelMap) ExecuteFunc(f func()) {
	c := make(chan bool)
	this.Op <- func() {
		f()
		c <- true
	}
	<-c
}
