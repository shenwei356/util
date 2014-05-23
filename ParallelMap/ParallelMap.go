// Copyright 2014 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.
//
// ParallelMap - A lock-free parallel map in go.
// ParallelMap uses a backend goroutine for sequential excution of
// Get and Set or custom function, which was inspired by section 14.17
// in book << The Way to Go >>.
//
// Example:
//
//    N := runtime.NumCPU() * 30
//    runtime.GOMAXPROCS(N)
//
//    m := NewParallelMap(func(v1 ValueType, v2 ValueType) ValueType {
//        return v1.(int) + v2.(int)
//    }, N)
//
//    var length int = 1 << 10
//
//    var wg sync.WaitGroup
//    for i := 1; i <= N; i++ {
//        wg.Add(1)
//
//        go func() {
//            defer func() {
//                wg.Done()
//                m.UnboundAGoroutine()
//            }()
//
//            for j := 0; j < length; j++ {
//                m.Update(j, 1)
//            }
//        }()
//    }
//
//    // wait for all operations to complement
//    wg.Wait()
//    m.Wait()
//
package ParallelMap

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// The type of the map key is interface{}
type KeyType interface{}

// The type of the map value is interface{}
type ValueType interface{}

// ParallelMap
type ParallelMap struct {
	// map
	Map map[KeyType]ValueType

	// backend goroutine for sequential operations
	Op chan func() error
	// counter of operations in Op
	opCounter int64

	// function to update value
	UpdateValueFunc func(ValueType, ValueType) ValueType

	// stop signal
	signalNum  int
	signalRsv  int
	signal     chan int
	signalExit chan int
	// status checking interval
	checkInterval time.Duration
}

// Constructor of ParallelMap
func NewParallelMap(
	updateValueFunc func(ValueType, ValueType) ValueType,
	parallelNum int) *ParallelMap {

	this := new(ParallelMap)

	this.Map = make(map[KeyType]ValueType)

	this.Op = make(chan func() error)
	this.opCounter = int64(0)

	this.UpdateValueFunc = updateValueFunc

	this.signalNum = parallelNum
	this.signal = make(chan int)
	this.signalExit = make(chan int)
	this.signalRsv = 0
	this.checkInterval = time.Millisecond * 50

	go this.backend()
	return this
}

// Run operation channel as backend
func (this *ParallelMap) backend() {
	ticker := time.NewTicker(this.checkInterval)

	var (
		f   func() error
		err error
		sig int
	)
	for {
		select {
		case sig = <-this.signal:
			if sig == 1 {
				this.signalRsv += sig
			} else {
				// bad signal
			}
		case f = <-this.Op:
			err = f()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		case <-ticker.C:
			if this.signalRsv == this.signalNum &&
				this.opCounter == 0 {
				// all operations excuted
				this.signalExit <- 0
			}
		}
	}
}

// UnboundAGoroutine
func (this *ParallelMap) UnboundAGoroutine() {
	this.signal <- 1
}

// Wait for all operations to complete
func (this *ParallelMap) Wait() {
	<-this.signalExit
}

// Getting element of the map is executed sequentially
func (this *ParallelMap) Get(key KeyType) (ValueType, bool) {
	c1 := make(chan ValueType)
	c2 := make(chan bool)
	this.OpCounterPlusOne()
	this.Op <- func() error {
		value, ok := this.Map[key]
		c1 <- value
		c2 <- ok
		this.OpCounterMinusOne()
		return nil
	}
	return <-c1, <-c2
}

// Setting operation is executed sequentially to ensure the
// operation is atomic.
func (this *ParallelMap) Set(key KeyType, value ValueType) {
	c := make(chan bool)
	this.OpCounterPlusOne()
	this.Op <- func() error {
		this.Map[key] = value
		c <- true
		this.OpCounterMinusOne()
		return nil
	}
	<-c
}

// Update function
func (this *ParallelMap) Update(key KeyType, value ValueType) {
	c := make(chan bool)
	this.OpCounterPlusOne()
	this.Op <- func() error {
		val, ok := this.Map[key]
		if ok {
			this.Map[key] = this.UpdateValueFunc(val, value)
		} else {
			this.Map[key] = value
		}

		c <- true
		this.OpCounterMinusOne()
		return nil
	}
	<-c
}

// Execute a custom function.
//
// Example: An element increasing function
//
//    m.ExecuteFunc(func() error {
//        if v, ok := m.Map[i]; ok {
//            m.Map[i] = v.(int) + 1
//        } else {
//            m.Map[i] = int(1)
//        }
//        return nil
//    })
func (this *ParallelMap) ExecuteFunc(f func() error) {
	c := make(chan bool)
	this.OpCounterPlusOne()
	this.Op <- func() error {
		err := f()
		c <- true
		this.OpCounterMinusOne()
		return err
	}
	<-c
}

// OpCounterPlusOne
func (this *ParallelMap) OpCounterPlusOne() {
	atomic.AddInt64(&this.opCounter, 1)
}

// OpCounterMinusOne
func (this *ParallelMap) OpCounterMinusOne() {
	atomic.AddInt64(&this.opCounter, -1)
}
